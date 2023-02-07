package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx/bsoncore"

	"golang.org/x/sys/unix"
)

var pipeFile = "/tmp/named_pipe_1"
var cursorReturned = make(chan struct{})
var cursor *mongo.Cursor

func main() {
	os.Remove(pipeFile)
	err := unix.Mkfifo(pipeFile, 0666)
	if err != nil {
		log.Fatal("Make named pipe file error:", err)
	}

	// Spin up goroutine to write to named pipe
	go writeToNamedPipe()

	// Spin up goroutine to make aggregation request to mongod
	go makeRequestToMongod()

	// Wait until the cursor is returned
	<-cursorReturned

	fmt.Println("ready to read")
	for cursor.Next(context.Background()) {
		fmt.Println(cursor.Current.String())
	}
}

func writeToNamedPipe() {
	fmt.Println("opening file")
	f, err := os.OpenFile(pipeFile, os.O_WRONLY, os.ModeNamedPipe)
	fmt.Println("file is opened")
	if err != nil {
		log.Fatalf("Open named pipe file error: %v", err)
	}

	// This is just a little test to see if the system supports deadlines on files and therefore
	// supports interrupting file I/O with a Close() call.
	log.Printf("attempting to set file deadline on OS %s with architecture %s\n", runtime.GOOS, runtime.GOARCH)
	err = f.SetDeadline(time.Now().Add(1 * time.Hour))
	if err != nil {
		fmt.Printf("error setting file deadline: %v\n", err)
	} else {
		fmt.Println("system supports file deadlines!")
	}

	_, doc := bsoncore.AppendDocumentStart(nil)
	doc = bsoncore.AppendInt32Element(doc, "a", 1)
	doc, _ = bsoncore.AppendDocumentEnd(doc, 0)

	for i := 0; i < 5; i++ {
		fmt.Println("waiting to write")
		_, err := f.Write(doc)
		if err != nil {
			log.Fatalf("error writing doc at index %d: %v", i, err)
		}
		fmt.Println("write succesful")
	}

	err = f.Close()
	if err != nil {
		log.Fatalf("error closing file: %v", err)
	}
}

func makeRequestToMongod() {
	uri := "mongodb://localhost:27017"
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}
	}()

	client.Database("test").Collection("named_pipe_1").Drop(context.Background())
	coll := client.Database("test").Collection("named_pipe_1")
	pipeline := mongo.Pipeline{bson.D{{"$match", bson.D{{"a", 1}}}}}
	options := options.Aggregate().SetCustom(bson.M{"$_externalDataSources": bson.A{bson.D{{"collName", "named_pipe_1"}, {"dataSources", bson.A{bson.D{{"url", "file://named_pipe_1"}, {"storageType", "pipe"}, {"fileType", "bson"}}}}}}})

	fmt.Println("making aggregate request")
	cursor, err = coll.Aggregate(context.Background(), pipeline, options)
	if err != nil {
		panic(err)
	}
	fmt.Println("cursor returned")
	close(cursorReturned)
}
