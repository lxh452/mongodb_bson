package model

import (
	"context"
	"fmt"
	builder "refeject_bson/jsy_builder"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *mongo.Database

const dbCollName = "_filter_builder_test"

func setup() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	d := client.Database(dbCollName)
	db = d
	err = d.Drop(ctx)
	if err != nil {
		panic(err)
	}

	coll := d.Collection(dbCollName)

	_, err = coll.InsertMany(ctx, []interface{}{
		bson.M{"name": "a", "capName": "A", "age": 8, "birthdate": time.Now()},
		bson.M{"name": "aa", "capName": "AA", "age": 58, "birthdate": time.Now()},
		bson.M{"name": "b", "capName": "B", "age": 18, "birthdate": time.Now().Add(24 * time.Hour)},
		bson.M{"name": "bb", "capName": "BB", "age": 18, "birthdate": time.Now().Add(24 * time.Hour)},
		bson.M{"name": "c", "capName": "C", "age": 38, "birthdate": time.Now().Add(3 * 24 * time.Hour)},
		bson.M{"name": "d", "capName": "D", "age": 28, "birthdate": time.Now().Add(6 * 24 * time.Hour)},
		bson.M{"name": "ab", "capName": "AB", "age": 5, "birthdate": time.Now()},
	})
	if err != nil {
		panic(err)
	}
}

func tearDown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := db.Drop(ctx)
	if err != nil {
		panic(err)
	}
}

func fetchData(filter interface{}) []interface{} {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	c, err := db.Collection(dbCollName).Find(context.Background(), filter)
	if err != nil {
		panic(err)
	}
	defer c.Close(ctx)
	var res []interface{}
	for c.Next(ctx) {
		var result bson.D
		err := c.Decode(&result)
		if err != nil {
			panic(err)
		}
		res = append(res, result)
	}
	if err := c.Err(); err != nil {
		panic(err)
	}
	return res
}

func mustEqual(t *testing.T, val1, val2 interface{}) {
	res1, res2 := fetchData(val1), fetchData(val2)
	// t.Logf("\n%v\n%v\n", res1, res2)
	assert.Equal(t, reflect.DeepEqual(res1, res2), true)
}

func TestBuilder_AutoWithKey(t *testing.T) {
	b := builder.New().AutoWithKey("test_key", 1).Build()
	c := builder.New().Num("test_key").Eq(1).Build()
	assert.Equal(t, c, b)

	b = builder.New().AutoWithKey("testKey", 0).Build()
	c = builder.New().Build()
	assert.Equal(t, c, b)

	b = builder.New().AutoWithKey("test_key", "123").Build()
	c = builder.New().Str("test_key").Eq("123").Build()
	assert.Equal(t, c, b)

	b = builder.New().AutoWithKey("testKey", "").Build()
	c = builder.New().Build()
	assert.Equal(t, c, b)

	b = builder.New().AutoWithKey("test_key", "123").Build()
	c = builder.New().Str("test_key").Eq("123").Build()
	assert.Equal(t, c, b)

	b = builder.New().AutoWithKey("testKey", "").Build()
	c = builder.New().Build()
	assert.Equal(t, c, b)

	b = builder.New().AutoWithKey("test_key", []int{1, 2, 3}).Build()
	c = builder.New().Num("test_key").In([]int{1, 2, 3}).Build()
	assert.Equal(t, c, b)

	b = builder.New().AutoWithKey("test_key", []int{}).Build()
	c = builder.New().Build()
	assert.Equal(t, c, b)

	// pointer val
	var val = 2
	var ptr = &val

	b = builder.New().AutoWithKey("test_key", ptr).Build()
	c = builder.New().Num("test_key").Eq(2).Build()
	assert.Equal(t, c, b)

	ptr = nil
	b = builder.New().AutoWithKey("test_key", ptr).Build()
	c = builder.New().Build()
	assert.Equal(t, c, b)

}

func TestBuilder_Auto(t *testing.T) {
	type Query struct {
		Age     int
		Name    string
		FooBar  []int
		Ignored any       `bson:"-"`
		Bar     *struct { // this kind of struct will be ignored if used like Builder.Auto(Query{Bar{...}})
			NameBar string
		}
		HasKey string `bson:"do_has_key"`
	}
	b := builder.New().Auto(Query{Age: 2}).Build()
	c := builder.New().Num("age").Eq(2).Build()
	assert.Equal(t, c, b)

	b = builder.New().Auto(Query{Age: 2, Name: "tester"}).Build()
	c = builder.New().Num("age").Eq(2).Str("name").Eq("tester").Build()
	assert.Equal(t, c, b)

	b = builder.New().Auto(Query{FooBar: []int{1, 2, 3}, Bar: &struct{ NameBar string }{"21312412"}}).Build()
	c = builder.New().Num("foo_bar").In([]int{1, 2, 3}).Build()
	assert.Equal(t, c, b)

	q := Query{Bar: &struct{ NameBar string }{"123"}}
	b = builder.New().Auto(q.Bar).Build()
	c = builder.New().Str("name_bar").Eq("123").Build()
	assert.Equal(t, c, b)

	q = Query{
		Name:    "queryName",
		Ignored: 123,
	}
	b = builder.New().Auto(q).Build()
	c = builder.New().Str("name").Eq("queryName").Build()
	assert.Equal(t, c, b)

	q = Query{
		HasKey: "123",
	}
	b = builder.New().Auto(q).Build()
	c = builder.New().Str("do_has_key").Eq("123").Build()
	assert.Equal(t, c, b)
}

func TestBuilder_RemoveCond(t *testing.T) {
	b := builder.New().
		Str("test_key").Eq("123").
		Str("test_key_2").Eq("321").
		RemoveCond("test_key"). // remove test_key
		Build()
	c := builder.New().
		Str("test_key_2").Eq("321").
		Build()
	assert.Equal(t, c, b)

	b = builder.New().
		Str("test_key").Eq("1").
		Str("test_key_2").Eq("1").
		Or().
		Str("test_key").Eq("2").
		Str("test_key_2").Eq("2").
		RemoveCond("test_key", true). // remove test_key across Or conds
		Build()
	c = builder.New().
		Str("test_key_2").Eq("1").
		Or().
		Str("test_key_2").Eq("2").
		Build()
	assert.Equal(t, c, b)

	b = builder.New().
		Str("test_key").Eq("1").
		Str("test_key_2").Eq("1").
		Or().
		Str("test_key").Eq("2").
		Str("test_key_2").Eq("2").
		RemoveCond("test_key"). // remove test_key with Or conds' preserved
		Build()
	c = builder.New().
		Str("test_key").Eq("1").
		Str("test_key_2").Eq("1").
		Or().
		Str("test_key_2").Eq("2").
		Build()
	assert.Equal(t, c, b)

}

func TestBuilder_Build(t *testing.T) {
	b := builder.New()
	b.IsAdd("test", "eq", "name")
	fmt.Println(b)
}
