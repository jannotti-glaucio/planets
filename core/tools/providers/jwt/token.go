package jwt

import (
	"errors"
	"time"

	"github.com/brianvoe/sjwt"
)

//TokenPayload date saved in token
type TokenPayload struct {
	UserUUID string
}

//TokenDetails details of token
type TokenDetails struct {
	AccessToken string
	AtExpires   int64
}

//Token struct for assign interface
type Token struct {
	SecretKey string
}

//New function to create a new instance
func New(secretKey string) *Token {
	return &Token{
		SecretKey: secretKey,
	}
}

//Token implements the TokenInterface
var _ IJwtProvider = &Token{}

//CreateToken create a new token
func (token *Token) CreateToken(payload TokenPayload) (*TokenDetails, error) {
	tokenDetails := TokenDetails{}
	tokenDetails.AtExpires = time.Now().Add(time.Minute * 10).Unix()

	claims := sjwt.New()
	claims.Set("user_uuid", payload.UserUUID)
	claims.Set("exp", tokenDetails.AtExpires)

	// Generate jwt
	secretKey := []byte(token.SecretKey)
	tokenDetails.AccessToken = claims.Generate(secretKey)
	return &tokenDetails, nil
}

//TokenValid check if token is valid
func (token *Token) isValid(tokenString string) bool {
	secretKey := []byte(token.SecretKey)
	// Verify that the secret signature is valid
	return sjwt.Verify(tokenString, secretKey)
}

//CheckToken responsible for check token
func (token *Token) CheckToken(tokenString string) (*TokenPayload, error) {
	if token.isValid(tokenString) {
		claims, err := sjwt.Parse(tokenString)
		if err == nil {
			err = claims.Validate()
			if err == nil {
				UserUUID, err := claims.Get("user_uuid")
				if err == nil {
					payload := TokenPayload{
						UserUUID: UserUUID.(string),
					}
					return &payload, nil
				}
			}
		}
	}

	return nil, errors.New("invalid token")
}
