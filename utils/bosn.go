package utils

import (
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
)

type StructIntoBsonM struct {
	Group   bson.M
	Project bson.M
}

func NewStructIntoBsonM() *StructIntoBsonM {
	return &StructIntoBsonM{}
}

func (s *StructIntoBsonM) Change(need interface{}, str string) StructIntoBsonM {
	return changeIntoBsonM(need, str)
}

func (s *StructIntoBsonM) MergerIntoBsonM(bsonM ...bson.M) bson.M {
	return mergerIntoBsonM(bsonM)
}

func changeIntoBsonM(need interface{}, str string) StructIntoBsonM {

	//获取反射条件
	v := reflect.ValueOf(need)
	t := v.Type()
	group := bson.M{}
	project := bson.M{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			changeIntoBsonM(v.Field(i).Interface(), str)
		}
		project["_id"] = 1
		group[field.Tag.Get("bson")] = "$" + field.Tag.Get("bson")
		project[field.Tag.Get("bson")] = str + field.Tag.Get("bson")
	}
	return StructIntoBsonM{Group: group, Project: project}

}

func mergerIntoBsonM(bsonM []bson.M) bson.M {
	result := make(bson.M)
	for _, v := range bsonM {
		if v == nil {
			continue
		}
		for k1, v1 := range v {
			result[k1] = v1
		}
	}
	return result
}
