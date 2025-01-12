package entity

type User struct {
	ID       int
	Email    string
	Password string
}

type Credentials struct {
	Email    string
	Password string
}