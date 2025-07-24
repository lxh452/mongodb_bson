package builder

import (
	"fmt"
	"time"
)

// dateCond represents a date-type condition builder.
type dateCond struct {
	*cond
	defaultFormat string
}

// newDateCond retunrs a new dateCond.
// RFC3339Nano format is used as default time format.
func newDateCond(key string, builder *Builder, format ...string) *dateCond {
	defaultFormat := time.RFC3339Nano
	if len(format) != 0 {
		defaultFormat = format[0]
	}
	return &dateCond{
		cond:          newCond(key, builder),
		defaultFormat: defaultFormat,
	}
}

func (c *dateCond) Eq(val time.Time) *Builder {
	return c.cond.Eq(val)
}

func (c *dateCond) EqStr(val string, format ...string) *Builder {
	t := c.mustParse(val, format...)
	return c.Eq(t)
}

func (c *dateCond) Ne(val time.Time) *Builder {
	return c.cond.Ne(val)
}

func (c *dateCond) NeStr(val string, format ...string) *Builder {
	t := c.mustParse(val, format...)
	return c.Ne(t)
}

func (c *dateCond) Lt(val time.Time) *Builder {
	return c.cond.Lt(val)
}

func (c *dateCond) LtStr(val string, format ...string) *Builder {
	t := c.mustParse(val, format...)
	return c.Lt(t)
}

func (c *dateCond) Lte(val time.Time) *Builder {
	return c.cond.Lte(val)
}

func (c *dateCond) LteStr(val string, format ...string) *Builder {
	t := c.mustParse(val, format...)
	return c.Lte(t)
}

func (c *dateCond) Gt(val time.Time) *Builder {
	return c.cond.gt(val)
}

func (c *dateCond) GtStr(val string, format ...string) *Builder {
	t := c.mustParse(val, format...)
	return c.Gt(t)
}

func (c *dateCond) Gte(val time.Time) *Builder {
	return c.cond.Gte(val)
}

func (c *dateCond) GteStr(val string, format ...string) *Builder {
	t := c.mustParse(val, format...)
	return c.Gte(t)
}

// Between side-inclusive time range query.
//
// * Since this method used to build cond like x_time: {$gte: ..., $lte: ...},
// it will remove the existing key to do overwirte any existing cond.
func (c *dateCond) Between(min, max time.Time) *Builder {
	c.builder.RemoveCond(c.key, false) //
	c.cond.Gte(min)
	return c.cond.Lte(max)
}

func (c *dateCond) BetweenStr(min, max string, format ...string) *Builder {
	minT := c.mustParse(min, format...)
	maxT := c.mustParse(max, format...)
	return c.Between(minT, maxT)
}

// RangeStr try to use rg[0], rg[1] as the input of BetweenStr.
func (c *dateCond) RangeStr(rg []string, format ...string) *Builder {
	if len(rg) < 2 || rg[0] == "" || rg[1] == "" {
		return c.builder
	}
	return c.BetweenStr(rg[0], rg[1], format...)
}

// mustParse parses time from string with format.
func (c *dateCond) mustParse(timeStr string, format ...string) time.Time {
	f := c.defaultFormat
	if len(format) != 0 {
		f = format[0]
	}

	t, err := time.Parse(f, timeStr)
	if err != nil {
		panic(fmt.Errorf("filterBuilder: failed to parse time from string: %s with format: %s, err: %v", timeStr, f, err))
	}
	return t
}
