package entity

type UserId int64

type User struct {
	Id       UserId
	Email    string
	PassHash []byte
}
