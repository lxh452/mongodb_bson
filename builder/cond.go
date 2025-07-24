package builder

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type cond struct {
	// key is the root name for the current condtion, eq: {key: {$eq: ...}}.
	key string
	// TODO consider remove the m, populate the builder.curMap directly.
	// m is the current map stores the value of condition.
	m bson.M
	// builder refers to the current builder.
	builder *Builder
}

func newCond(key string, builder *Builder) *cond {
	return &cond{
		key:     key,
		m:       bson.M{},
		builder: builder,
	}
}

// addMapToBuilder adds baseCond.m to the referenced builder's map with key as baseCond.key.
// If the same key is set again, it will try to merge two map.
func (baseCond *cond) addMapToBuilder() {
	var v interface{}
	var ok bool

	if v, ok = baseCond.builder.curMap[baseCond.key]; !ok {
		baseCond.builder.curMap[baseCond.key] = baseCond.m
		return
	}
	preMap := v.(bson.M)

	for k, v := range baseCond.m {
		preMap[k] = v
	}
}

// Eq adds `$Eq: val` to the baseCond.m
func (baseCond *cond) Eq(val interface{}) *Builder {
	baseCond.m[_eq] = val
	baseCond.addMapToBuilder()
	return baseCond.builder
}

// Ne adds `$Ne: val` to the baseCond.m
func (baseCond *cond) Ne(val interface{}) *Builder {
	baseCond.m[_ne] = val
	baseCond.addMapToBuilder()
	return baseCond.builder
}

// Lt adds `$Lt: val` to the baseCond.m
func (baseCond *cond) Lt(val interface{}) *Builder {
	baseCond.m[_lt] = val
	baseCond.addMapToBuilder()
	return baseCond.builder
}

// Lte adds `$Lte: val` to the baseCond.m
func (baseCond *cond) Lte(val interface{}) *Builder {
	baseCond.m[_lte] = val
	baseCond.addMapToBuilder()
	return baseCond.builder
}

// gt adds `$gt: val` to the baseCond.m
func (baseCond *cond) gt(val interface{}) *Builder {
	baseCond.m[_gt] = val
	baseCond.addMapToBuilder()
	return baseCond.builder
}

// Gte adds `$Gte: val` to the baseCond.m
func (baseCond *cond) Gte(val interface{}) *Builder {
	baseCond.m[_gte] = val
	baseCond.addMapToBuilder()
	return baseCond.builder
}

// Regex adds `$Regex: exp, $options: ""` to the baseCond.m
func (baseCond *cond) Regex(exp string) *Builder {
	baseCond.RegexWithOpt(exp, "")
	return baseCond.builder
}

// RegexWithOpt adds `$regex: exp, $options: opt` to the baseCond.m
func (baseCond *cond) RegexWithOpt(exp string, opt string) *Builder {
	baseCond.m[_regex] = primitive.Regex{Pattern: exp, Options: opt}
	baseCond.addMapToBuilder()
	return baseCond.builder
}

// Not adds `$not: exp, $options: ""` to the baseCond.m
func (baseCond *cond) Not(exp string) *Builder {
	baseCond.NotWithOpt(exp, "")
	return baseCond.builder
}

// NotWithOpt adds `$not: exp, $options: opt` to the baseCond.m
func (baseCond *cond) NotWithOpt(exp string, opt string) *Builder {
	baseCond.m[_not] = primitive.Regex{Pattern: exp, Options: opt}
	baseCond.addMapToBuilder()
	return baseCond.builder
}

// In adds `$In: vals` to the baseCond.m
func (baseCond *cond) In(vals interface{}) *Builder {
	baseCond.m[_in] = vals
	baseCond.addMapToBuilder()
	return baseCond.builder
}

// Nin adds `$Nin: vals` to the baseCond.m
func (baseCond *cond) Nin(vals interface{}) *Builder {
	baseCond.m[_nin] = vals
	baseCond.addMapToBuilder()
	return baseCond.builder
}
