package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"refeject_bson/model"
	"reflect"
)

func main() {
	//连接数据库
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	usercoll := client.Database("admin").Collection("user")

	var user model.User
	user = user.NewUser()
	filter := changeIntoBsonM(user)

	one, err := usercoll.InsertOne(context.TODO(), filter)
	if err != nil {
		return
	}
	fmt.Printf("Inserted a single document: %+v\n", one)
}

func changeIntoBsonM(need interface{}) bson.M {

	//获取反射条件
	v := reflect.ValueOf(need)
	t := v.Type()
	filter := bson.M{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			changeIntoBsonM(v.Field(i).Interface())
		}
		filter[field.Tag.Get("bson")] = v.Field(i).Interface()
	}
	return filter
}
