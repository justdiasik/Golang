package models

import (
	"regexp"
	"unicode"
)

type Laptop struct {
	ID int `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	
	Year int `json:"year" db:"year"`
	Display string `json:"display" db:"display"`

	Memory string `json:"memory" db:"memory"`
	Storage string `json:"storage" db:"storage"`

	Condition string `json:"condition" db:"condition"`
	Description string `json:"description" db:"description"`
	
	Phonenumber  string `json:"phonenumber" db:"phonenumber"`
	
	Price string `json:"price" db:"price"`
}

type LaptopsFilter struct {
	Query *string `json:"query"`
}

/////////////////////////

type Snowboard struct {
	ID int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`

	Size string `json:"size" db:"size"`
	Condition string `json:"condition" db:"condition"`

	Description string `json:"description" db:"description"`

	Phonenumber  string `json:"phonenumber" db:"phonenumber"`
	Price string `json:"price" db:"price"`
}

type SnowboardsFilter struct {
	Query *string `json:"query"`
}

type Shirt struct {
	ID int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`

	Color string `json:"color" db:"color"`
	Size string `json:"size" db:"size"`
	Condition string `json:"condition" db:"condition"`

	Description string `json:"description" db:"description"`
	Phonenumber  string `json:"phonenumber" db:"phonenumber"`
	Price string `json:"price" db:"price"`
}

type ShirtsFilter struct {
	Query *string `json:"query"`
}

type Toy struct {
	ID int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`

	Condition string `json:"condition" db:"condition"`
	Description string `json:"description" db:"description"`

	Phonenumber  string `json:"phonenumber" db:"phonenumber"`
	Price string `json:"price" db:"price"`
}

type ToysFilter struct {
	Query *string `json:"query"`
}

type User struct {
	Username string `json:"username" db:"username"`
	Email string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type UsersFilter struct {
	Query *string `json:"query"`
}

func (u *User) IsEmailValid() bool {
   emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
   return emailRegex.MatchString(u.Email)
}

func (u *User) IsPasswordValid() bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(u.Password) >= 7 {
		hasMinLen = true
	}
	for _, char := range u.Password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}





