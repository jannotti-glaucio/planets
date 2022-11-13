package hash

import (
	"golang.org/x/crypto/bcrypt"
)

//PasswordHash struct for assign IHasg
type PasswordHash struct{}

//New create a new instance
func New() *PasswordHash {
	return &PasswordHash{}
}

//Assign interface
var _ IHashProvider = &PasswordHash{}

//Create hashed password from payload
func (hash *PasswordHash) Create(payload string) string {
	crypted, _ := bcrypt.GenerateFromPassword([]byte(payload), bcrypt.DefaultCost)
	return string(crypted)
}

//Compare hashed and playload
func (hash *PasswordHash) Compare(hashed string, payload string) bool {
	check := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(payload))
	return check == nil
}
