package entities

import (
	"time"

	"github.com/google/uuid"
)

//User ...
type User struct {
	UUID      string     `json:"uuid" bson:"_id"`
	Name      string     `json:"name" bson:"name"`
	Email     string     `json:"email" bson:"email"`
	Password  string     `json:"password" bson:"password"`
	CreatedAt time.Time  `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty" bson:"updated_at"`
}

//Users ...
type Users []User

//UserPublic ...
type UserPublic struct {
	UUID  string `json:"uuid" bson:"_id"`
	Name  string `json:"name" bson:"name"`
	Email string `json:"email" bson:"email"`
}

//New return a object
func (User *User) New() *User {
	User.UUID = uuid.New().String()
	User.CreatedAt = time.Now()
	return User
}

//PublicUser return a User without non important attributes
func (User *User) PublicUser() interface{} {
	return UserPublic{
		UUID:  User.UUID,
		Name:  User.Name,
		Email: User.Email,
	}
}

//PublicUsers return a list of users without important attributes
func (users *Users) PublicUsers() []interface{} {
	result := make([]interface{}, len(*users))
	for index, user := range *users {
		result[index] = user.PublicUser()
	}
	return result
}
