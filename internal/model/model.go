package model

type User struct {
	ID       int
	Email    string
	PassHash []byte
}
