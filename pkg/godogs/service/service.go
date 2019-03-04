package service

import "fmt"

type LoginService interface {
	Login(username, password string) string
}

type MyLoginService struct {
	name string
}

func (s *MyLoginService) Login(username, password string) string {
	fmt.Println(fmt.Sprintf("This is %s.", s.name))
	if username == "username" && password == "password" {
		return "OK"
	} else {
		return "NOT-OK"
	}
}

func NewService() LoginService {
	return &MyLoginService{name: "cool login service"}
}
