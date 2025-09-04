package main

import (
	"fmt"

	common "github.com/SoeltanIT/agg-common-be"
)

func main() {
	// Create new custom error
	commonErr := common.NewError(422, 4220001, "Unprocessable entity")
	fmt.Println(commonErr.Error()) // Unprocessable entity

	// Import from defined error code
	err := common.ErrForbidden
	fmt.Println(err.HTTPStatus) // 403
	fmt.Println(err.Code)       // 4030001
	fmt.Println(err.Error())    // You do not have permission to access this resource
}
