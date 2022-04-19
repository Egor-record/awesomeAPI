package storage

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type EventStorage struct{}

type Event struct {
	Type       string    `bson:"type"`
	State      bool      `bson:"state"`
	StartedAt  time.Time `bson:"started_at"`
	FinishedAt time.Time `bson:"finished_at"`
}

func (s *EventStorage) StartEvent(eventName string) error {
	collection, client, ctx, cancel, err := s.connectAndCheckDbConnection("events")
	if err != nil {
		return err
	}
	defer s.closeDB(client, ctx, cancel)

	_, err = s.findEvent(eventName, false, *collection)
	if err != nil {
		switch err {
		case mongo.ErrNoDocuments:
			err := s.insertEvent(eventName, *collection)
			if err != nil {
				return errors.New("error inserting new type")
			}
			return nil
		default:
			return errors.New("error searching new document")
		}
	}
	return nil

}

func (s *EventStorage) EndEvent(eventName string) error {
	collection, client, ctx, cancel, err := s.connectAndCheckDbConnection("events")
	if err != nil {
		return err
	}
	defer s.closeDB(client, ctx, cancel)

	event, err := s.findEvent(eventName, false, *collection)

	if err != nil {
		return err
	}

	err = s.changeEventStatus(event.Type, true, *collection)
	if err != nil {
		return err
	}
	return nil

}

func (s *EventStorage) changeEventStatus(eventName string, status bool, collection mongo.Collection) error {
	filter := bson.M{
		"state": bson.M{
			"$eq": !status, //
		},
		"type": bson.M{
			"$eq": eventName,
		},
	}

	update := bson.M{
		"$set": bson.M{
			"state":       status,
			"finished_at": time.Now(),
		},
	}

	_, err := collection.UpdateOne(
		context.Background(),
		filter,
		update,
	)

	if err != nil {
		fmt.Println("UpdateOne() result ERROR:", err)
		return err
	}

	return nil
}

func (s *EventStorage) insertEvent(eventName string, collection mongo.Collection) error {

	event := Event{eventName, false, time.Now(), time.Now()}
	insertResult, err := collection.InsertOne(context.TODO(), event)

	if err != nil {
		return err
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return nil
}

func (s *EventStorage) findEvent(eventName string, status bool, collection mongo.Collection) (*Event, error) {
	var result Event
	err := collection.FindOne(context.Background(), bson.M{"type": eventName, "state": status}).Decode(&result)
	if err != nil {
		return nil, err
	}
	fmt.Println("FindOne() result:", result)
	return &result, nil
}

func (s *EventStorage) connectAndCheckDbConnection(collectionName string) (*mongo.Collection, *mongo.Client, context.Context, context.CancelFunc, error) {
	client, ctx, cancel, err := s.connectDB("mongodb+srv://" + USERNAME + ":" + PASS + "@cluster0.1wcc1.mongodb.net/" + DB + "?retryWrites=true&w=majority")
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection := client.Database("eventsDB").Collection(collectionName)
	return collection, client, ctx, cancel, nil
}

func (s *EventStorage) connectDB(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(uri).SetServerAPIOptions(serverAPIOptions)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	return client, ctx, cancel, err
}

func (s *EventStorage) closeDB(client *mongo.Client, ctx context.Context,
	cancel context.CancelFunc) {

	defer cancel()

	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}
