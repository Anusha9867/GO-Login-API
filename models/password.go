package models

type Password struct {
	UserId   uint   `json:"user_id"`
	Password []byte `json:"-"`
}
