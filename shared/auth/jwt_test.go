package auth

import (
	"github.com/tpp/msf/model"
	"testing"
)

var tokenString = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJjb250ZXh0IjoidGhhbmgudHBwQGdtYWlsLmNvbSJ9.MXUX8NKm7dYVcvMmocsJWfHyxPSq8lAoEhwnSrC5q4U"

/*
	func TestGenerateJWTToken(t *testing.T) {
		got, _ := GenerateJWTToken(&model.User{})
		t.Error(got)
	}
*/
func TestGenerateJWTToken(t *testing.T) {
	user, err := GenerateJWTToken(&model.User{})
	if err != nil {
		t.Fail()
	}
	t.Log(user, err)
}

func TestParseJWTToken(t *testing.T) {
	user, err := ParseJWTToken(tokenString)
	if err != nil {
		t.Fail()
	}
	t.Log(user, err)
}
