package model

type Vote struct {
	Sender  string `json:"sender" bson:"sender"`
	Target  string `json:"target" bson:"target"`
	Vote    int    `json:"vote" bson:"vote"`
	VotedAt *int64 `json:"voted_at" bson:"voted_at"`
}
