package ticket

import "swift_transit/ticket"

type Service interface {
	BuyTicket(req ticket.BuyTicketRequest) (*ticket.BuyTicketResponse, error)
	UpdatePaymentStatus(id int64) error
	DownloadTicket(id int64) ([]byte, error)
}
