package model

import (
	"strconv"
	"time"
)

type User struct {
	UserId       string    `bson:"user_id"`
	Username     string    `json:"username" bson:"username"`
	Password     string    `json:"password" bson:"password"`
	Address      []Address `json:"address" bson:"address"`
	RelationShip []User    `json:"relation" bson:"relation"`
}
type Address struct {
	City  string `json:"city" bson:"city"`
	State string `json:"state" bson:"state"`
}

// 创建用户对象
func (u *User) NewUser() User {
	var address []Address
	var relation []User
	for i := 0; i < 5; i++ {
		relation = append(relation, User{
			Username:     "用户" + strconv.Itoa(i),
			Password:     "123456" + strconv.Itoa(i),
			Address:      nil,
			RelationShip: nil,
		})
		address = append(address, Address{
			City:  "广州" + strconv.Itoa(i),
			State: "番禺" + strconv.Itoa(i),
		})
	}
	user := User{
		UserId:       "1",
		Username:     "lxh",
		Password:     "123456",
		Address:      address,
		RelationShip: relation,
	}
	return user
}

type Customer struct {
	CustomerId string   `json:"customer_id" bson:"customer_id"`
	UserId     string   `json:"user_id" bson:"user_id"`
	Wallet     []Wallet `json:"wallet" bson:"wallet"`
}
type Wallet struct {
	PayType     string    `json:"pay_type" bson:"pay_type"`
	Count       float64   `json:"count" bson:"count"`
	LastPayTime time.Time `json:"last_pay_time" bson:"last_pay_time"`
}

// 创建用户对象
func (c *Customer) NewCustomer() Customer {
	var wallet []Wallet
	for i := 0; i < 5; i++ {
		wallet = append(wallet, Wallet{
			PayType:     "银行卡" + strconv.Itoa(i),
			Count:       float64(i*i - i + i*i),
			LastPayTime: time.Now(),
		})
	}
	customer := Customer{
		CustomerId: "1",
		UserId:     "1",
		Wallet:     wallet,
	}
	return customer
}
