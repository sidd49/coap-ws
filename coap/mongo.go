package coap

import (
	"coapws/coap/models"
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Mongodb struct {
	DB  *mongo.Database
	uri string
	ctx *context.Context
}

func NewMongo() *Mongodb {
	return &Mongodb{
		uri: "your-mongo-uri",
	}
}

func (m *Mongodb) Init() {
	// fetch todo context
	ctx := context.TODO()
	mongoconn := options.Client().ApplyURI(m.uri)
	// connect with the mongo server
	mongodb, err := mongo.Connect(ctx, mongoconn)
	if err != nil {
		log.Fatal("Cannot initialize mongo")
	}
	// ping the server to check if the connection is successful or not
	err = mongodb.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal("Error connecting with mongo ", err)
	}
	fmt.Println("Mongo Connection Established")
	// select the database of the application
	m.DB = mongodb.Database("coap-server")
	m.ctx = &ctx

}

func (m *Mongodb) GetThingByID(id string) (*models.Thing, error) {
	var thing models.Thing
	// fetch the thing with the deviceid from db
	err := m.DB.Collection("things").FindOne(*m.ctx, bson.D{{Key: "deviceid", Value: id}}).Decode(&thing)
	if err != nil {
		log.Println("Error in getting the thing with id : " + id + " : " + err.Error())
		return nil, err
	}
	return &thing, nil
}

func (m *Mongodb) Update(t *models.Thing) error {
	// update the thing in the db
	result := m.DB.Collection("things").FindOneAndReplace(*m.ctx, bson.D{{Key: "deviceid", Value: t.DeviceId}}, t)
	if result.Err() != nil {
		log.Println("Error in updating the event : " + t.DeviceId + " : " + result.Err().Error())
		return result.Err()
	}
	return nil
}

func (m *Mongodb) StoreThings(things []models.Thing) error {
	var insertThings []interface{}
	for _, thing := range things {
		exists, err := m.checkAndUpdateThing(&thing)
		if err != nil {
			return err
		}
		if !exists {
			insertThings = append(insertThings, thing)
		}
	}
	if len(insertThings) == 0 {
		return nil
	}
	return m.Save(insertThings)
}

func (m *Mongodb) checkAndUpdateThing(thing *models.Thing) (bool, error) {
	_, err := m.GetThingByID(thing.DeviceId)
	if err == nil {
		log.Printf("The device is already registered! Going for updation...")
		updateErr := m.Update(thing)
		if updateErr != nil {
			log.Fatalf("Error in updating the thing!")
			return false, updateErr
		}
		return true, nil
	}
	return false, nil
}

func (m *Mongodb) Save(insertThings []interface{}) error {
	_, err := m.DB.Collection("things").InsertMany(*m.ctx, insertThings)
	return err
}
