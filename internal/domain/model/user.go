package model

type User struct {
	HashID string `json:"hash_id" bson:"hash_id"`
}

func (u *User) TableName() string { return "users" }
