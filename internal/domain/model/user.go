package model

type User struct {
	Username string `json:"username" bson:"username" validate:"required,min=4"`
	Password string `json:"password" bson:"password" validate:"required,min=8,containsany=!@#?*"`
}

func (u *User) TableName() string { return "users" }
