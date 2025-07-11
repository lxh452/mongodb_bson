package data

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	Client   *mongo.Client
	UserColl *mongo.Collection
)

func init() {
	//连接数据库
	var err error
	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	if err = Client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}

	UserColl = Client.Database("admin").Collection("user")
}
