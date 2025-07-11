package main

import (
	"context"
	"fmt"
	"git.jsyiot.com/jsytech/lego-lib/xbson/agop"
	builder "git.jsyiot.com/jsytech/mongo-filter-builder"
	"go.mongodb.org/mongo-driver/bson"
	"refeject_bson/data"
	"refeject_bson/model"
	"refeject_bson/utils"
)

func main() {
	var ctx context.Context
	userColl := data.UserColl
	var user model.User
	var customer model.Customer
	user = user.NewUser()

	customer = customer.NewCustomer()
	bsonM := utils.NewStructIntoBsonM()
	//构建分组条件
	userBsonM := bsonM.Change(user, "$user.")
	customerBsonM := bsonM.Change(customer, "$customer.")

	// 构建显示条件
	project := bsonM.MergerIntoBsonM(userBsonM.Project, customerBsonM.Project)
	fmt.Println(project)
	fmt.Println(userBsonM.Group)

	filter := builder.New()
	filter.Str("address.city").Eq("广州1")
	pipe := bson.A{
		agop.Match(filter.Build()),
		bson.M{
			"$group": bson.M{
				"_id":      "$user_id",
				"customer": bson.M{"$push": customerBsonM.Group},
				"user":     bson.M{"$push": userBsonM.Group},
			},
		},
		bson.M{
			"$project": project,
		},
	}
	testFilter, err := userColl.Aggregate(ctx, pipe)
	if err != nil {
		fmt.Println(err)
		return
	}
	var a []interface{}
	err = testFilter.All(ctx, &a)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(a)
}
