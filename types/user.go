package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost         = 12
	minimumFistNameLen = 2
	minimumLastNameLen = 2
	minimumPasswordLen = 7
)

type CreateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (p UpdateUserParams) ToBSON() bson.M {
	m := bson.M{}
	if len(p.FirstName) > 0 {
		m["firstName"] = p.FirstName
	}
	if len(p.LastName) > 0 {
		m["lastName"] = p.LastName
	}
	return m
}

func (p CreateUserParams) Validate() map[string]string {

	errors := make(map[string]string)

	if len(p.FirstName) < minimumFistNameLen {
		errors["firstName"] = fmt.Sprintf("must be at least %d characters long", minimumFistNameLen)
	}
	if len(p.LastName) < minimumLastNameLen {
		errors["lastName"] = fmt.Sprintf("must be at least %d characters long", minimumLastNameLen)
	}
	if len(p.Password) < minimumPasswordLen {
		errors["password"] = fmt.Sprintf("must be at least %d characters long", minimumPasswordLen)
	}
	if !isEmailValid(p.Email) {
		errors["email"] = "must be a valid email address"
	}
	return errors
}

func IsValidPassword(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

func isEmailValid(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"     json:"id,omitempty"`
	FirstName         string             `bson:"firstName"         json:"firstName"`
	LastName          string             `bson:"lastName"          json:"lastName"`
	Email             string             `bson:"email"             json:"email"`
	EncryptedPassword string             `bson:"encryptedPassword" json:"-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
