package utils

import "golang.org/x/crypto/bcrypt"

func HashAndSaltPassword(in string) string {
	bs, _ := bcrypt.GenerateFromPassword([]byte(in), bcrypt.DefaultCost)
	return string(bs)
}

func CompareHashAndPassword(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}
