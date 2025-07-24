package builder

import "go.mongodb.org/mongo-driver/bson/primitive"

// oidCond represents a ObjectId-type condition builder.
type oidCond struct {
	*cond
}

func newOidCond(key string, builderRef *Builder) *oidCond {
	return &oidCond{
		cond: newCond(key, builderRef),
	}
}

func (c *oidCond) Eq(oid string) *Builder {
	id, err := primitive.ObjectIDFromHex(oid)
	if err != nil {
		panic(err)
	}
	return c.cond.Eq(id)
}
