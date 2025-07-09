package model

import "strconv"

type User struct {
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
		Username:     "lxh",
		Password:     "123456",
		Address:      address,
		RelationShip: relation,
	}
	return user
}
