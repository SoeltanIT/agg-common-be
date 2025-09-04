package main

import (
	"fmt"

	common "github.com/SoeltanIT/agg-common-be"
)

type createUserRequest struct {
	Username string `validate:"required"`
	Password string `validate:"required"`
}

func main() {
	newUser := createUserRequest{
		Username: "john",
		Password: "123456",
	}

	nilerr := common.Validator().Struct(newUser)
	fmt.Println(nilerr) // <nil>

	newUserInvalid := createUserRequest{
		Username: "john",
	}
	err := common.Validator().Struct(newUserInvalid)
	fmt.Println(err) // Key: 'createUserRequest.Password' Error:Field validation for 'Password' failed on the 'required' tag
}
