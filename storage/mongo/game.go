package mongo

import (
	"context"
	"github.com/pkg/errors"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"

	"github.com/peccancy/poker_dealer/model"
)

// Select game from mongo by gameID
func (r *Repo) Get(ctx context.Context, accountID, ticketID string) (*documents.Ticket, error) {
	c := r.db.Database.C(ticketCollection)

	var ticket documents.Ticket

	if err := c.Find(bson.M{"account_id": accountID, "id": ticketID}).One(&ticket); err != nil {
		if err == mgo.ErrNotFound {
			return nil, errors.New("not found")
		}
		return nil, errors.Wrap(err, "mongo get ticket")
	}

	return &ticket, nil
}

// Select tickets from mongo by accountID
func (r *Repo) GetTickets(ctx context.Context, accountID string) (*documents.Tickets, error) {
	c := r.db.Database.C(ticketCollection)

	var result documents.Tickets
	var count int
	var err error

	if count, err = c.Find(bson.M{"account_id": accountID}).Count(); err != nil {
		return nil, errors.Wrap(err, "get total count tickets")
	}

	result.TotalCount = count

	if count > 0 {
		var tickets []*documents.Ticket

		err := c.Find(bson.M{"account_id": accountID}).All(&tickets)
		if err != nil {
			return nil, errors.Wrap(err, "find tickets")
		}

		result.Tickets = tickets
	}

	return &result, nil
}

// Insert new ticket in mongo
func (r *Repo) CreateTicket(ctx context.Context, ticket *documents.Ticket) error {
	c := r.db.Database.C(ticketCollection)

	n, err := c.Find(bson.M{"account_id": ticket.AccountID, "id": ticket.ID}).Count()
	if err != nil {
		return errors.Wrap(err, "find ticket to insert")
	}

	if n > 0 {
		return errors.New("ticket already exists")
	}

	err = c.Insert(ticket)
	if err != nil {
		return errors.Wrap(err, "insert ticket")
	}

	return nil
}

// Update ticket in mongo
func (r *Repo) UpdateTicket(ctx context.Context, ticket *documents.Ticket) error {
	c := r.db.Database.C(ticketCollection)

	_, err := c.Upsert(bson.M{"account_id": ticket.AccountID, "id": ticket.ID}, ticket)

	if err != nil {
		return errors.Wrap(err, "upsert ticket")
	}

	return nil
}

// Select tickets from mongo by entity type
func (r *Repo) GetTicketsByLinkedEntities(ctx context.Context, accountID string, let *documents.EnumType) (*documents.Tickets, error) {
	c := r.db.Database.C(ticketCollection)

	var tickets []*documents.Ticket
	var count int
	var err error

	if count, err = c.Find(bson.M{"account_id": accountID, "linked_entities.entity_type": let}).Count(); err != nil {
		return nil, errors.Wrap(err, "get total count tickets by linked entities")
	}

	err = c.Find(bson.M{"account_id": accountID, "linked_entities.entity_type": let}).All(&tickets)
	if err != nil {
		return nil, errors.Wrap(err, "find by linked entities")
	}

	return &documents.Tickets{
		Tickets:    tickets,
		TotalCount: count,
	}, nil
}
