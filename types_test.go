package main

import (

	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	tests := []struct{
		give	UserRequest
		want	*Error	
	}{
		{UserRequest{"test_user", "testc"}, &Error{"username and password must be longer than 8 characters." , http.StatusBadRequest}},
		{UserRequest{"usertest" , "123456"}, &Error{"username and password must be longer than 8 characters." , http.StatusBadRequest}},
	}


	for _,test := range tests {
		_,got := test.give.CreateUser()
		
		assert.Equal(t , test.want , got)
	}
}


func HashPass(pass string) string {
	hashPassword , _:= bcrypt.GenerateFromPassword([]byte(pass) , bcrypt.DefaultCost)
	return string(hashPassword)
}