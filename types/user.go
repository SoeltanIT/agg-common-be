package types

import (
	"database/sql/driver"
	"fmt"
)

type UserGender string

const (
	Male   UserGender = "male"
	Female UserGender = "female"
	Other  UserGender = "other"
)

type UserStatus string

const (
	Active   UserStatus = "active"
	Inactive UserStatus = "inactive"
)

type UserType string

const (
	UserTypeGuest UserType = "guest"
	UserTypeUser  UserType = "user"
)

func (t *UserType) Scan(value interface{}) error {
	b, ok := value.(string)
	if !ok {
		return fmt.Errorf("invalid type for FavouriteObjectType: %T", value)
	}
	*t = UserType(b)
	return nil
}

func (t UserType) Value() (driver.Value, error) {
	return string(t), nil
}
