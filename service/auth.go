package service

import "github.com/dgrijalva/jwt-go"

// Create the JWT key used to create the signature
var JwtKey = []byte("my_secret_key")


// Create a struct that will be encoded to a JWT.
// We add jwt.StandardClaims as an embedded type, to provide fields like expiry time
type Claims struct {
	ID string `json:"id"`
	FullName string `json:"full_name"`
	PhoneNumber string `json:"phone_number"`
	IsAdmin bool `json:"is_admin"`
	jwt.StandardClaims
}

