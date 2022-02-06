package model

import "github.com/dgrijalva/jwt-go"

type User struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
	UserType string `json:"userType"`
}

type UserDTO struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserCredentialsDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

type CustomClaims struct {
	UserType string `json:"userType"`
	jwt.StandardClaims
}
