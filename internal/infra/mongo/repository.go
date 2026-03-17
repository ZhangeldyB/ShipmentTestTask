package mongo

import (
	"context"
	"errors"

	"github.com/ZhangeldyB/ShipmentTestTask/internal/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	shipmentsCollection      = "shipments"
	shipmentEventsCollection = "shipment_events"
)

type Repository struct {
	db *mongo.Database
}

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{db: db}
}

// EnsureIndexes creates required indexes on startup.
func EnsureIndexes(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection(shipmentsCollection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "reference_number", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return err
	}

	_, err = db.Collection(shipmentEventsCollection).Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: "shipment_id", Value: 1}},
	})
	return err
}

func (r *Repository) Save(ctx context.Context, s *domain.Shipment) error {
	doc := toShipmentDoc(s)
	_, err := r.db.Collection(shipmentsCollection).ReplaceOne(
		ctx,
		bson.M{"_id": doc.ID},
		doc,
		options.Replace().SetUpsert(true),
	)
	return err
}

func (r *Repository) FindByID(ctx context.Context, id string) (*domain.Shipment, error) {
	var doc shipmentDoc
	err := r.db.Collection(shipmentsCollection).FindOne(ctx, bson.M{"_id": id}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrShipmentNotFound
		}
		return nil, err
	}
	return toDomainShipment(&doc), nil
}

func (r *Repository) FindByReferenceNumber(ctx context.Context, ref string) (*domain.Shipment, error) {
	var doc shipmentDoc
	err := r.db.Collection(shipmentsCollection).FindOne(ctx, bson.M{"reference_number": ref}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, domain.ErrShipmentNotFound
		}
		return nil, err
	}
	return toDomainShipment(&doc), nil
}

func (r *Repository) SaveEvent(ctx context.Context, e *domain.ShipmentEvent) error {
	doc := toEventDoc(e)
	_, err := r.db.Collection(shipmentEventsCollection).InsertOne(ctx, doc)
	return err
}

func (r *Repository) FindEventsByShipmentID(ctx context.Context, shipmentID string) ([]*domain.ShipmentEvent, error) {
	opts := options.Find().SetSort(bson.D{{Key: "occurred_at", Value: 1}})
	cursor, err := r.db.Collection(shipmentEventsCollection).Find(ctx, bson.M{"shipment_id": shipmentID}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var docs []shipmentEventDoc
	if err := cursor.All(ctx, &docs); err != nil {
		return nil, err
	}

	events := make([]*domain.ShipmentEvent, len(docs))
	for i := range docs {
		events[i] = toDomainEvent(&docs[i])
	}
	return events, nil
}
