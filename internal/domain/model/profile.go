package model

type Profile struct {
	Nickname  string `json:"nickname" bson:"nickname" validate:"required,sha256"`
	FirstName string `json:"first_name" query:"first_name" form:"first_name" bson:"first_name" validate:"required"`
	LastName  string `json:"last_name" query:"last_name" form:"last_name" bson:"last_name" validate:"required"`
	CreatedAt *int64 `json:"created_at" bson:"created_at"`
	UpdatedAt *int64 `json:"updated_at" bson:"updated_at"`
	DeletedAt *int64 `json:"deleted_at" bson:"deleted_at"`
}

func (p *Profile) TableName() string { return "profiles" }
