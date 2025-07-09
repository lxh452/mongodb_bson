package main

import (
	"context"
	"fmt"
	"git.jsyiot.com/jsytech/lego-lib/xbson/agop"
	builder "git.jsyiot.com/jsytech/mongo-filter-builder"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"refeject_bson/model"
	"reflect"
)

func main() {
	var ctx context.Context
	//连接数据库
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil {
		panic(err)
	}
	userColl := client.Database("admin").Collection("user")
	var user model.User
	var customer model.Customer
	user = user.NewUser()
	customer = customer.NewCustomer()
	userGroup, userProject := changeIntoBsonM(user, "$user.")
	customerGroup, _ := changeIntoBsonM(customer, "$customer.")
	filter := builder.New()
	filter.Str("address.city").Eq("广州1")
	pipe := bson.A{
		agop.Match(filter.Build()),
		bson.M{
			"$group": bson.M{
				"_id":      "$user_id",
				"customer": bson.M{"$push": customerGroup},
				"user":     bson.M{"$push": userGroup},
			},
		},
		bson.M{
			"$project": userProject,
		},
	}
	testfilter, err := userColl.Aggregate(ctx, pipe)
	if err != nil {
		fmt.Println(err)
		return
	}
	var a []interface{}
	err = testfilter.All(ctx, &a)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(a)
}

func changeIntoBsonM(need interface{}, str string) (bson.M, bson.M) {

	//获取反射条件
	v := reflect.ValueOf(need)
	t := v.Type()
	filter := bson.M{}
	project := bson.M{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			changeIntoBsonM(v.Field(i).Interface(), str)
		}
		project["_id"] = 1
		filter[field.Tag.Get("bson")] = "$" + field.Tag.Get("bson")
		project[field.Tag.Get("bson")] = str + field.Tag.Get("bson")
	}
	return filter, project
}
