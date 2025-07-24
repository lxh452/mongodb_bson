package bsonStruct

import (
	"go.mongodb.org/mongo-driver/bson"
)

// cond 条件创造器
type cond struct {
	key string

	m bson.M

	CondBuild *CondBuild
}

func (this *cond) AddToCondBuild() {
	var v interface{}
	var ok bool

	if v, ok = this.CondBuild.curMap[this.key]; !ok {
		this.CondBuild.curMap[this.key] = this.m
		return
	}
	conditionsMap := v.(bson.M)

	for k, v := range this.m {
		conditionsMap[k] = v
	}
}

func newCond(key string, condBuild *CondBuild) *cond {
	return &cond{
		key:       key,
		m:         bson.M{},
		CondBuild: condBuild,
	}
}

func (c *cond) Push(key bson.M) *CondBuild {
	c.m[push] = key
	c.AddToCondBuild()
	return c.CondBuild
}

func (c *cond) Sum(key bson.M) *CondBuild {
	c.m[sum] = key
	c.AddToCondBuild()
	return c.CondBuild
}

func (c *cond) Group(key bson.M) *CondBuild {
	c.m[group] = key
	c.AddToCondBuild()
	return c.CondBuild
}

func (c *cond) Project(key bson.M) *CondBuild {
	c.m[project] = key
	c.AddToCondBuild()
	return c.CondBuild
}
