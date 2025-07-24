package bsonStruct

import "go.mongodb.org/mongo-driver/bson"

// 罗列该制造器所需的key
const (
	project = "$project"
	group   = "$group"
	push    = "$push"
	sum     = "$sum"
	lookup  = "$lookup"
	unwind  = "$unwind"
)

// CondBuild 用于制作一个条件制造器
type CondBuild struct {
	condsMap []bson.M
	curMap   bson.M
}

// New 创造一个新的制造器
func (condBuild *CondBuild) New() *CondBuild {
	maps := []bson.M{}
	return &CondBuild{
		condsMap: maps,
		curMap:   bson.M{},
	}
}
