package repo

import (
	"swift_transit/domain"
	"swift_transit/ticket"
	"swift_transit/utils"

	"github.com/jmoiron/sqlx"
)

type TicketRepo interface {
	ticket.TicketRepo
}

type ticketRepo struct {
	dbCon       *sqlx.DB
	utilHandler *utils.Handler
}

func NewTicketRepo(dbcon *sqlx.DB, utilHandler *utils.Handler) TicketRepo {
	return &ticketRepo{
		dbCon:       dbcon,
		utilHandler: utilHandler,
	}
}

func (r *ticketRepo) Create(ticket domain.Ticket) (*domain.Ticket, error) {
	query := `
		INSERT INTO tickets (user_id, route_id, bus_name, start_destination, end_destination, fare, paid_status, qr_code, created_at)
		VALUES (:user_id, :route_id, :bus_name, :start_destination, :end_destination, :fare, :paid_status, :qr_code, :created_at)
		RETURNING id
	`
	rows, err := r.dbCon.NamedQuery(query, ticket)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&ticket.Id)
		if err != nil {
			return nil, err
		}
	}

	return &ticket, nil
}

func (r *ticketRepo) UpdateStatus(id int64, status bool) error {
	query := `UPDATE tickets SET paid_status = $1 WHERE id = $2`
	_, err := r.dbCon.Exec(query, status, id)
	return err
}

func (r *ticketRepo) Get(id int64) (*domain.Ticket, error) {
	var ticket domain.Ticket
	query := `SELECT * FROM tickets WHERE id = $1`
	err := r.dbCon.Get(&ticket, query, id)
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}
