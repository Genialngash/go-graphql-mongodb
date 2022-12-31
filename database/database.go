package database

import (
	"context"
	"log"
	"time"

	"github.com/Genialngash/go-mongo-graphql/graph/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var connectionString string = "mongodb://localhost:27017"

type DB struct {
	client *mongo.Client
}

func Connect() *DB {
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
		return nil
	}

	c, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err = client.Connect(c); err != nil {
		log.Fatal(err)
	}

	if err = client.Ping(c, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	return &DB{
		client: client,
	}

}

func (db *DB) GetJob(id string) (*model.JobListing, error) {
	jobCollect := db.client.Database("job-graph").Collection("job")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	_id, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": _id}

	var joblisting model.JobListing

	if err := jobCollect.FindOne(ctx, filter).Decode(&joblisting); err != nil {
		log.Fatal(err)
		return nil, err
	}

	return &joblisting, nil
}

func (db *DB) GetJobs() ([]*model.JobListing, error) {
	jobCollect := db.client.Database("job-graph").Collection("job")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cursor, err := jobCollect.Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	var joblisting []*model.JobListing

	if err := cursor.All(context.TODO(), &joblisting); err != nil {
		return nil, err
	}

	return joblisting, nil
}
func (db *DB) CreateJobListing(jobInfo model.CreateJobListingInput) *model.JobListing {

	jobCollect := db.client.Database("job-graph").Collection("job")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	inserted, err := jobCollect.InsertOne(ctx, bson.M{

		"title":       jobInfo.Title,
		"description": jobInfo.Description,
		"url":         jobInfo.URL,
		"company":     jobInfo.Company,
	})

	if err != nil {
		log.Fatal(err)
	}


	insertedId := inserted.InsertedID.(primitive.ObjectID).Hex()

	returnListing := model.JobListing{ID: insertedId, Title: jobInfo.Title, Company: jobInfo.Company, Description: jobInfo.Description, URL: jobInfo.URL}

	return &returnListing

}

func (db *DB) UpdateJobListing(jobId string, jobinfo model.UpdateJobListingInput) *model.JobListing {
	jobCollect := db.client.Database("job-graph").Collection("job")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	updateJobInfo := bson.M{}
	if jobinfo.Title != nil {
		updateJobInfo["title"] = jobinfo.Title
		
	}
	if jobinfo.Description != nil {
		updateJobInfo["description"] = jobinfo.Description
	}
	_id ,_ := primitive.ObjectIDFromHex(jobId)
	filter := bson.M{"_id":_id}
	update := bson.M{"$set":updateJobInfo}
	results := jobCollect.FindOneAndUpdate(ctx,filter,update,options.FindOneAndUpdate().SetReturnDocument(1))

	var joblisting model.JobListing
	if err := results.Decode(&joblisting); err != nil {
		log.Fatal(err)
		
	}
	return &joblisting
}
func (db *DB) DeleteJobListing(jobId string) *model.DeleteJobRespone {
	return &model.DeleteJobRespone{DeleteJobID: jobId}
}
