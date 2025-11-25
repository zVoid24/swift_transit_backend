package ticket

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"swift_transit/domain"
	"swift_transit/infra/payment"
	"swift_transit/user"
	"time"

	"github.com/go-pdf/fpdf"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/skip2/go-qrcode"
)

type service struct {
	repo       TicketRepo
	userRepo   user.UserRepo
	redis      *redis.Client
	sslCommerz *payment.SSLCommerz
	ctx        context.Context
}

func NewService(repo TicketRepo, userRepo user.UserRepo, redis *redis.Client, sslCommerz *payment.SSLCommerz, ctx context.Context) Service {
	return &service{
		repo:       repo,
		userRepo:   userRepo,
		redis:      redis,
		sslCommerz: sslCommerz,
		ctx:        ctx,
	}
}

func (s *service) BuyTicket(req BuyTicketRequest) (*BuyTicketResponse, error) {
	qrCode := uuid.New().String()
	now := time.Now().Format(time.RFC3339)

	ticket := domain.Ticket{
		UserId:           req.UserId,
		RouteId:          req.RouteId,
		BusName:          req.BusName,
		StartDestination: req.StartDestination,
		EndDestination:   req.EndDestination,
		Fare:             req.Fare,
		QRCode:           qrCode,
		CreatedAt:        now,
	}

	if req.PaymentMethod == "wallet" {
		// Deduct balance
		err := s.userRepo.DeductBalance(req.UserId, req.Fare)
		if err != nil {
			return nil, fmt.Errorf("payment failed: %w", err)
		}
		ticket.PaidStatus = true
	} else {
		ticket.PaidStatus = false
	}

	// Create ticket in DB
	createdTicket, err := s.repo.Create(ticket)
	if err != nil {
		return nil, err
	}

	if req.PaymentMethod == "wallet" {
		// Store in Redis (valid for 5 hours)
		ticketJSON, _ := json.Marshal(createdTicket)
		key := fmt.Sprintf("ticket:%d", createdTicket.Id)
		fmt.Printf("DEBUG: Attempting to set Redis key: %s with value: %s\n", key, string(ticketJSON))

		err = s.redis.Set(s.ctx, key, ticketJSON, 5*time.Hour).Err()
		if err != nil {
			fmt.Printf("DEBUG: failed to cache ticket: %v\n", err)
		} else {
			fmt.Println("DEBUG: Redis Set successful")
		}

		val, err := s.redis.Get(s.ctx, key).Result()
		if err != nil {
			fmt.Printf("DEBUG: Redis Get failed immediately after Set: %v\n", err)
		} else {
			fmt.Printf("DEBUG: Redis Get value: %s\n", val)
		}
		return &BuyTicketResponse{
			Ticket:  createdTicket,
			Message: "Ticket purchased successfully",
		}, nil
	} else {
		// Init SSLCommerz
		tranID := fmt.Sprintf("TICKET-%d-%s", createdTicket.Id, uuid.New().String()[:8])
		// URLs should be configured or constructed properly
		successUrl := fmt.Sprintf("http://localhost:8080/ticket/payment/success?id=%d", createdTicket.Id)
		failUrl := "http://localhost:8080/ticket/payment/fail"
		cancelUrl := "http://localhost:8080/ticket/payment/cancel"

		gatewayUrl, err := s.sslCommerz.InitPayment(req.Fare, tranID, successUrl, failUrl, cancelUrl)
		if err != nil {
			return nil, fmt.Errorf("gateway init failed: %w", err)
		}

		return &BuyTicketResponse{
			Ticket:     createdTicket,
			PaymentURL: gatewayUrl,
			Message:    "Redirect to payment gateway",
		}, nil
	}
}

func (s *service) UpdatePaymentStatus(id int64) error {
	return s.repo.UpdateStatus(id, true)
}

func (s *service) DownloadTicket(id int64) ([]byte, error) {
	// Fetch ticket
	ticket, err := s.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get ticket: %w", err)
	}

	// Generate QR Code
	qrCode, err := qrcode.Encode(ticket.QRCode, qrcode.Medium, 256)
	if err != nil {
		return nil, fmt.Errorf("failed to generate QR code: %w", err)
	}

	// Create PDF
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Swift Transit Ticket")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Ticket ID: %d", ticket.Id))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Bus Name: %s", ticket.BusName))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Route ID: %d", ticket.RouteId))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("From: %s", ticket.StartDestination))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("To: %s", ticket.EndDestination))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Fare: %.2f", ticket.Fare))
	pdf.Ln(8)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", ticket.CreatedAt))
	pdf.Ln(20)

	// Embed QR Code
	// fpdf requires an image reader or file. We can use RegisterImageOptionsReader
	imageOptions := fpdf.ImageOptions{
		ImageType: "PNG",
		ReadDpi:   true,
	}
	pdf.RegisterImageOptionsReader("qrcode.png", imageOptions, bytes.NewReader(qrCode))
	pdf.ImageOptions("qrcode.png", 10, 100, 50, 50, false, imageOptions, 0, "")

	// Output to bytes
	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return buf.Bytes(), nil
}
