package Database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Task struct {

	ID      primitive.ObjectID	`bson:"_id,omitempty"`
	Name	string			`bson:"name"`
	Done	bool			`bson:"done"`
}

func GetTasksTable(user, pass, uri string) (*mongo.Collection) {
	opts := options.Client().ApplyURI(uri).SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-1", // or "SCRAM-SHA-1"
		AuthSource:    "admin",         // the database where the user is defined
		Username:      user,
		Password:      pass,
	})

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	prod := client.Database("htmx_todo_crud")
	tasksCollection := prod.Collection("tasks")
	return tasksCollection
}

func DropAll(user, pass, uri string) {
	opts := options.Client().ApplyURI(uri).SetAuth(options.Credential{
		AuthMechanism: "SCRAM-SHA-1", // or "SCRAM-SHA-1"
		AuthSource:    "admin",         // the database where the user is defined
		Username:      user,
		Password:      pass,
	})

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	prod := client.Database("htmx_todo_crud")


	prod.Collection("tasks").Drop(context.TODO())
	prod.CreateCollection(context.TODO(), "tasks")
}

func AddTask(collection *mongo.Collection, name string) error {
	newTask := Task {
		ID: primitive.NewObjectID(),
		Name: name,
		Done: false,
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	_, err := collection.InsertOne(ctx, newTask)

	return err
}

func GetAllTasks(collection *mongo.Collection) ([]Task) {
	cur, err := collection.Find(context.TODO(), bson.D{{}}, options.Find())

	if err != nil {
		log.Fatalf("Failed to get all tasks from the collection: %v", err)
	}

	var results []Task
	for cur.Next(context.TODO()) {
		var element Task
		err := cur.Decode(&element)
		if err != nil {
			log.Fatalf("Failed to decode element: %v", err)
		}
		results = append(results, element)
	}
	return results
}

