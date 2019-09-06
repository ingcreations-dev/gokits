package db

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"os"
	"strconv"
	"time"
)

var Mongodb *MongoClient

type MongoClient struct {
	database *mongo.Database
	Duration time.Duration
}

func init() {
	server := os.Getenv("MONGODB_SERVER")
	db := os.Getenv("MONGODB_DB")
	timeout,_ := strconv.Atoi(os.Getenv("MONGODB_TIMEOUT"))
	client,err := mongo.NewClient(options.Client().ApplyURI(server))
	if err != nil {
		panic(err.Error())
	}
	dur := time.Duration(timeout) * time.Second
	ctx,_ := context.WithTimeout(context.Background(),dur)
	client.Connect(ctx)

	err = client.Ping(ctx,readpref.Primary())
	if err != nil{
		panic(err.Error())
	}

	Mongodb = &MongoClient{
		database : client.Database(db),
		Duration: dur,
	}
}

func (client *MongoClient) GetCollection(tableName string){
	collections := client.database.Collection(tableName)

	println(collections.Name())
}

func (client *MongoClient) Save(tableName string,table interface{}) error{
	ctx :=client.getCtx()
	result,err := client.database.Collection(tableName).InsertOne(ctx,table)
	if err != nil{
		return err
	}
	fmt.Println("Inserted a single document: ", result.InsertedID)
	return nil
}

func (client *MongoClient) Update(tableName string,filter bson.M,setter interface{}) error {
	ctx :=client.getCtx()
	_,err := client.database.Collection(tableName).UpdateOne(ctx,filter,setter)
	return err
}

func (client *MongoClient) UpdateMany(tableName string,filter bson.M,setter interface{}) error {
	ctx :=client.getCtx()
	_,err := client.database.Collection(tableName).UpdateMany(ctx,filter,setter)
	return err
}

/**
通过条件查询一个文档
 */
func (client *MongoClient) FindOne(tableName string,filter bson.M,table interface{}) error{
	result := client.database.Collection(tableName).FindOne(client.getCtx(),filter)
	if result.Err() != nil{
		return result.Err()
	}
	err := result.Decode(table)
	if err != nil{
		return err
	}

	return nil
}

func (client *MongoClient) FindCount(tableName string,filter bson.M) (int64,error){
	return client.database.Collection(tableName).CountDocuments(client.getCtx(),filter)
}

func (client *MongoClient) Delete(tableName string,filter bson.M) error{
	_,err := client.database.Collection(tableName).DeleteOne(client.getCtx(),filter)
	return err
}

/**
通过条件查询列表
 */
func (client *MongoClient) FindAllByCondition(tableName string,filter bson.M) (*mongo.Cursor,error) {
	return client.database.Collection(tableName).Find(client.getCtx(),filter)
}

func (client *MongoClient) FindAll(tableName string)(*mongo.Cursor,error){
	return client.database.Collection(tableName).Find(client.getCtx(),nil)
}

func (client *MongoClient) getCtx() context.Context{
	ctx,_ := context.WithTimeout(context.Background(),client.Duration)
	return ctx
}
