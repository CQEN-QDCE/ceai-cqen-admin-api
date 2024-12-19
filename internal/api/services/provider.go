package services

type Provider interface {
	GetUsers()

	GetUser(string)

	CreateUser()

	UpdateUser()

	DeleteUser(string)
}
