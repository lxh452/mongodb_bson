package builder

import (
	"github.com/iancoleman/strcase"
	"go.mongodb.org/mongo-driver/bson"
	"reflect"
	"strings"
)

const (
	_eq = "$eq"
	_ne = "$ne"

	_gt  = "$gt"
	_gte = "$gte"
	_lt  = "$lt"
	_lte = "$lte"

	_in  = "$in"
	_nin = "$nin"

	_regex = "$regex"
	_not   = "$not"

	_or = "$or"
)

// Builder represents a filter builder.
type Builder struct {
	// condMaps stores all condition maps.
	condMaps []bson.M
	// curMap represents the currently operated condition map.
	// A condition map can be a single element map also can be a multiple elements map.
	curMap bson.M
}

// New constructs a new Builder.
func New() *Builder {
	maps := []bson.M{}
	return &Builder{
		condMaps: maps,
		curMap:   bson.M{},
	}
}

// Auto constructs a new Builder with Builder.Auto method.
func Auto(val any) *Builder {
	return New().Auto(val)
}

// Flush restes the builder to initial state
func (b *Builder) Flush() *Builder {
	b.condMaps = []bson.M{}
	b.curMap = bson.M{}
	return b
}

// Auto will construct suitable eq filter as possible as it can.
//
// queryStruct should be a struct contains query fields (optionally with bson tags).
//
// If bson tag is provided on field, the tag will be used as the key of cond,
// otherwise snake case of field's name will be used as default.
//
// If it's a pointer:
//
//   - a pointer to a struct:
//     it will try to dereference it.
//
//   - a nil pointer:
//     do nothing.
//
// *Anything else will lead to a panic.
func (b *Builder) Auto(queryStruct any) *Builder {
	val := reflect.ValueOf(queryStruct)
	for val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return b
		}
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		panic("the given value is not struct")
	}

	fields := reflect.VisibleFields(val.Type())
	nonZeroFields := []reflect.StructField{}
	for _, _f := range fields {
		f := val.FieldByName(_f.Name)
		if f.IsZero() {
			continue
		}
		nonZeroFields = append(nonZeroFields, _f)
	}

	for _, _f := range nonZeroFields {

		v := val.FieldByName(_f.Name)

		bsonTag, ok := _f.Tag.Lookup("bson")
		key, _, _ := strings.Cut(bsonTag, ",")
		if key == "-" {
			continue
		}
		if !ok || key == "" {
			// use camel case as default
			b.AutoWithKey(strcase.ToSnake(_f.Name), v.Interface())
			continue
		}

		b.AutoWithKey(key, v.Interface())
	}

	return b
}

// AutoWithKey try to add eq cond with provided key and val.
//
// If val is a pointer:
//
//	Cond will be built if the val is one of following type and the val is non-zero value:
//		- array, slice
//		- bool
//		- string
//		- int, int8 and other ints.
//		- uint, uint8 and other uints.
//		- float32, float64
//
// If val is not a pointer:
//
//	Cond will be built if the pointer is not nil.
func (b *Builder) AutoWithKey(key string, val any) *Builder {
	_v := reflect.ValueOf(val)

	valFromPointer := _v.Kind() == reflect.Pointer
	for (_v.Kind() == reflect.Pointer || _v.Kind() == reflect.Interface) && !_v.IsZero() {
		_v = _v.Elem()
	}

	switch _v.Kind() {
	case
		reflect.Array, reflect.Slice,
		reflect.Bool,
		reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
	default: // stop if unsupported type
		return b
	}

	if !valFromPointer && _v.IsZero() { // non-pointer zero value
		return b
	}
	if _v.Kind() == reflect.Slice || _v.Kind() == reflect.Array {
		if _v.Len() != 0 {
			b.Any(key).In(val)
		}
	} else {
		b.Any(key).Eq(_v.Interface())
	}
	return b
}

// WantNum indicates the builder to build a condition for string type.
func (b *Builder) Str(key string) *strCond {
	return newStrCond(key, b)
}

// Num indicates the builder to build a condition for number type.
func (b *Builder) Num(key string) *numCond {
	return newNumCond(key, b)
}

func (b *Builder) Date(key string, defaultFormat ...string) *dateCond {
	return newDateCond(key, b, defaultFormat...)
}

// Oid key will use _id as default
func (b *Builder) Oid(key ...string) *oidCond {
	if len(key) != 0 {
		return newOidCond(key[0], b)
	}
	return newOidCond("_id", b)
}

// Any constructs a condition without type restricted.
func (b *Builder) Any(key string) *cond {
	return newCond(key, b)
}

// Or appends b.curMap to b.condMaps, and b.curMap will be assigned to a new empty map.
// Thus if finally b.condMaps's len is bigger than 1, then the final filter will wraps all maps into a $or condition.
func (b *Builder) Or() *Builder {
	if len(b.curMap) == 0 {
		return b
	}
	b.condMaps = append(b.condMaps, b.curMap)
	b.curMap = bson.M{}
	return b
}

// AnyMap will set the given map to current condition.
func (b *Builder) AnyMap(key string, m bson.M) *Builder {
	b.curMap[key] = m
	return b
}

// RemoveCond removes given key that has been added to the builder.
// Elimination will across all conditions if acrossOrCond is given true.
func (b *Builder) RemoveCond(key string, acrossOrCond ...bool) *Builder {
	for k := range b.curMap {
		if k == key {
			delete(b.curMap, k)
		}
	}

	if len(acrossOrCond) != 0 && acrossOrCond[0] {
		for _, condMap := range b.condMaps {
			for k := range condMap {
				if k == key {
					delete(condMap, k)
				}
			}
		}
	}

	return b
}

// Build builds final filter and returns it as bson.M.
func (b *Builder) Build() bson.M {
	if len(b.curMap) != 0 {
		b.condMaps = append(b.condMaps, b.curMap)
	}
	if len(b.condMaps) > 1 {

		return bson.M{_or: b.condMaps}
	}
	res := b.curMap
	return res
}

func (b *Builder) IsIntAdd() {

}

func (b *Builder) IsFloatAdd() {

}

func (b *Builder) IsStringAdd() {

}

func (b *Builder) IsArrayDateAdd(cond, key string, value []string) {
	if len(value) == 2 {
		b.Date(key).BetweenStr(value[0], value[1])
	}
	switch cond {
	case "nin":
		b.Str(key).Nin(value...)
	case "in":
		b.Str(key).In(value...)
	}

}

func (b *Builder) IsAdd(value any, cond, key string) *Builder {
	val := reflect.ValueOf(value)
	switch val.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Float32, reflect.Float64:
		{
			if !val.IsZero() {
				return b.isAddNum(value.(int), cond, key)
			}
		}
	case reflect.String:
		{
			if !val.IsZero() {
				return b.isAddStr(value.(string), cond, key)
			}
		}

	case reflect.Slice, reflect.Array:
		{
			if !val.IsNil() {
				return b.isAddSlice(value, cond, key)
			}
		}
	case reflect.Pointer:
		if !val.IsNil() {

		}
		b.IsAdd(&value, cond, key)

	default:
	}
	return b
}

type Request struct {
	Num *int
}

func (b *Builder) isAddNum(value any, conds, key string) *Builder {
	switch conds {
	case "eq":
		return b.Num(key).Eq(value)
	case "ne":
		return b.Num(key).Ne(value)
	case "lt":
		return b.Num(key).Lt(value)
	case "lte":
		return b.Num(key).Lte(value)
	case "gt":
		return b.Num(key).Gt(value)
	case "gte":
		return b.Num(key).Gte(value)
	case "in":
		return b.Num(key).In(value)
	default:
	}
	return b
}

func (b *Builder) isAddStr(value any, conds, key string) *Builder {
	val := value.(string)
	switch conds {
	case "eq":
		return b.Str(key).Eq(val)
	case "ne":
		return b.Str(key).Ne(val)
	case "lt":
		return b.Str(key).Lt(val)
	case "lte":
		return b.Str(key).Lte(val)
	case "gte":
		return b.Str(key).Gte(val)
	case "in":
		return b.Str(key).In(val)
	case "like":
		return b.Str(key).Like(val)
	case "not":
		return b.Str(key).Not(val)
	case "notLike":
		return b.Str(key).NotLike(val)
	default:
	}
	return b
}

func (b *Builder) isAddSlice(value interface{}, conds, key string) *Builder {
	switch value.(type) {
	case []string:
		{
			val := value.([]string)
			switch conds {
			case "nin":
				return b.Str(key).Nin(val...)
			case "in":
				return b.Str(key).In(val...)
			}
		}
	case []int:
		{
			val := value.([]int)
			if conds == "between" && len(val) == 2 && val[0] < val[1] {
				return b.Num(key).Between(val[0], val[1])
			} else if conds == "in" {
				return b.Num(key).In(val)
			}
		}
	}

	return b
}
