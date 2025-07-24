package builder

// strCond represents a string-type condition builder.
type strCond struct {
	*cond
}

// newStrCond constructs a new strCond.
func newStrCond(key string, builderRef *Builder) *strCond {
	return &strCond{
		cond: newCond(key, builderRef),
	}
}

// Eq adds `$eq: val` to the strc.m
func (strc *strCond) Eq(val string) *Builder {
	return strc.cond.Eq(val)
}

// Ne adds `$ne: val` to the strc.m
func (strc *strCond) Ne(val string) *Builder {
	return strc.cond.Ne(val)
}

// Regex adds `$regex: exp, $options: ""` to the strc.m
func (strc *strCond) Regex(exp string) *Builder {
	return strc.cond.Regex(exp)
}

// RegexWithOpt adds `$regex: exp, $options: opt` to the strc.m
func (strc *strCond) RegexWithOpt(exp string, opt string) *Builder {
	return strc.cond.RegexWithOpt(exp, opt)
}

// Like calls strc.Regex(val) under the wood
func (strc *strCond) Like(val string) *Builder {
	return strc.Regex(val)
}

// NotLike calls strc.Not(val) under the wood
func (strc *strCond) NotLike(val string) *Builder {
	return strc.Not(val)
}

// Not adds `$not: exp, $options: ""` to the strc.m
func (strc *strCond) Not(exp string) *Builder {
	return strc.cond.Not(exp)
}

// NotWithOpt adds `$not: exp, $options: opt` to the strc.m
func (strc *strCond) NotWithOpt(exp string, opt string) *Builder {
	return strc.cond.NotWithOpt(exp, opt)
}

// In adds `$in: vals` to the strc.m
func (strc *strCond) In(vals ...string) *Builder {
	return strc.cond.In(vals)
}

// In adds `$nin: vals` to the strc.m
func (strc *strCond) Nin(vals ...string) *Builder {
	return strc.cond.Nin(vals)
}
