package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	//
	// Prepare DataBase
	//
	ctx := context.Background()
	Client = NewMongoDBClient(ctx)
	DB = Client.Database("testdatabase")
	Tasks = DB.Collection("tasks")
	defer Client.Disconnect(ctx)

	//
	// Delete All Tasks
	//
	if err := DeleteAllTasks(ctx); err != nil {
		log.Fatal("DeleteAllTasks", err)
	}

	t := NewTaskDocumentGateway()
	t.Title = "Task0"
	t.Username = "garry.tasker"

	//
	// Insert
	//
	taskID, err := t.Insert(ctx)
	if err != nil {
		log.Fatal("Insert", err)
	}
	fmt.Println("taskID", taskID)

	//
	// FindByID
	//
	f := NewTaskFinder()
	foundTask, err := f.FindByID(ctx, taskID)
	if err != nil {
		log.Fatal("FindByID", err)
	}

	taskJSON, _ := json.MarshalIndent(foundTask, "  ", "  ")
	fmt.Println("taskJSON", string(taskJSON))

	//
	// Update task
	//
	t.Progress = 25
	if err := t.Update(ctx); err != nil {
		log.Fatal("Update", err)
	}

	countTasks, err := f.Count(ctx)
	if err != nil {
		log.Fatal("Count", err)
	}
	fmt.Println("countTasks", countTasks)

	//
	// Find By Username
	//
	updatedTasks, err := f.FindByUsername(ctx, "garry.tasker")
	if err != nil {
		log.Fatal("FindByID", err)
	}
	if len(updatedTasks) == 0 {
		log.Fatal("len(updatedTasks) = ", updatedTasks)
	}
	updatedTask := updatedTasks[0]

	updatedTaskJSON, _ := json.MarshalIndent(updatedTask, "  ", "  ")
	fmt.Println("updatedTaskJSON", string(updatedTaskJSON))

	//
	// Delete
	//
	if err := t.Delete(ctx); err != nil {
		log.Fatal("Delete", err)
	}

	countAfterDeleteTasks, err := f.Count(ctx)
	if err != nil {
		log.Fatal("Count", err)
	}
	fmt.Println("countAfterDeleteTasks", countAfterDeleteTasks)
}

const (
	uri = "mongodb://127.0.0.1:27017"
)

var (
	Client *mongo.Client
	DB     *mongo.Database
	Tasks  *mongo.Collection
)

// NewMongoDBClient - ...
func NewMongoDBClient(ctx context.Context) *mongo.Client {
	opts := options.Client().ApplyURI(uri)

	c, err := mongo.NewClient(opts)
	if err != nil {
		log.Fatal("NewClient", err)
	}

	if err := c.Connect(ctx); err != nil {
		log.Fatal("Connect", err)
	}

	if err := c.Ping(ctx, nil); err != nil {
		log.Fatal("Ping", err)
	}

	return c
}

// NewTaskDocumentGateway - ...
func NewTaskDocumentGateway() *TaskDocumentGateway {
	t := new(TaskDocumentGateway)
	t.CreatedDate = time.Now()
	return t
}

// TaskDocumentGateway - ...
type TaskDocumentGateway struct {
	ObjectID    primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title,omitempty"`
	CreatedDate time.Time          `json:"createddate" bson:"createddate,omitempty"`
	Username    string             `json:"username" bson:"username,omitempty"`
	Progress    int                `json:"progress" bson:"progress,omitempty"`
}

// Insert - ..
func (g *TaskDocumentGateway) Insert(ctx context.Context) (interface{}, error) {
	r, err := Tasks.InsertOne(ctx, g)
	if err != nil {
		return 0, err
	}
	g.ObjectID = r.InsertedID.(primitive.ObjectID)
	return g.ObjectID, nil
}

// Update - ...
func (g *TaskDocumentGateway) Update(ctx context.Context) error {
	filter := bson.M{"_id": g.ObjectID}
	_, err := Tasks.ReplaceOne(ctx, filter, g)
	return err
}

// Delete - ...
func (g *TaskDocumentGateway) Delete(ctx context.Context) error {
	_, err := Tasks.DeleteOne(ctx, bson.M{"_id": g.ObjectID})
	return err
}

// NewTaskFinder - ...
func NewTaskFinder() *TaskFinder {
	f := new(TaskFinder)
	return f
}

// TaskFinder - ...
type TaskFinder struct{}

// FindByID - ...
func (f *TaskFinder) FindByID(ctx context.Context, id interface{}) (*TaskDocumentGateway, error) {
	r := Tasks.FindOne(ctx, bson.M{"_id": id})
	if err := r.Err(); err != nil {
		return nil, err
	}
	var t *TaskDocumentGateway
	if err := r.Decode(&t); err != nil {
		return nil, err
	}
	return t, nil
}

// FindByUsername - ..
func (f *TaskFinder) FindByUsername(ctx context.Context, username string) ([]*TaskDocumentGateway, error) {
	var ts []*TaskDocumentGateway
	r, err := Tasks.Find(ctx, bson.M{"username": username})
	if err != nil {
		return ts, err
	}
	if err := r.All(ctx, &ts); err != nil {
		return ts, err
	}
	return ts, nil
}

// Count - ...
func (f *TaskFinder) Count(ctx context.Context) (int64, error) {
	return Tasks.CountDocuments(ctx, bson.M{})
}

// DeleteAllTasks - ...
func DeleteAllTasks(ctx context.Context) error {
	_, err := Tasks.DeleteMany(ctx, bson.M{})
	return err
}
