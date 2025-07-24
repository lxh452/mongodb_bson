package builder

// numCond represents a numeric-type condition builder.
// For convince, the val passed in MUST to be a numeric type,
// APIs WILL NOT check its type.
type numCond struct {
	*cond
}

func newNumCond(key string, builder *Builder) *numCond {
	return &numCond{
		cond: newCond(key, builder),
	}
}

// Eq adds `$eq: val` to the c.m
func (c *numCond) Eq(val interface{}) *Builder {
	return c.cond.Eq(val)
}

// Ne adds `$ne: val` to the c.m
func (c *numCond) Ne(val interface{}) *Builder {
	return c.cond.Ne(val)
}

// Lt adds `$lt: val` to the c.m
func (c *numCond) Lt(val interface{}) *Builder {
	return c.cond.Lt(val)
}

// Lte adds `$lte: val` to the c.m
func (c *numCond) Lte(val interface{}) *Builder {
	return c.cond.Lte(val)
}

// Gt adds `$gt: val` to the c.m
func (c *numCond) Gt(val interface{}) *Builder {
	return c.cond.gt(val)
}

// Gte adds `$gte: val` to the c.m
func (c *numCond) Gte(val interface{}) *Builder {
	return c.cond.Gte(val)
}

// Between => [min, max]
func (c *numCond) Between(min interface{}, max interface{}) *Builder {
	c.cond.Gte(min)
	return c.cond.Lte(max)
}

// In adds `$nin: vals` to the c.m
func (c *numCond) In(nums interface{}) *Builder {
	return c.cond.In(nums)
}
