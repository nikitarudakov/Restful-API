package model

type User struct {
	Username       string `json:"username" form:"username" bson:"username"`
	HashedPassword string `json:"password" form:"password" bson:"password"`
}

func (User) TableName() string { return "users" }
