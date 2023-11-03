package model

type User struct {
	Username string `json:"username" bson:"username" validate:"required,min=4"`
	Password string `json:"password" bson:"password" validate:"required,min=8,containsany=!@#?*%"`
	Role     string `json:"role" query:"role" form:"role" bson:"role" validate:"required"`
}
