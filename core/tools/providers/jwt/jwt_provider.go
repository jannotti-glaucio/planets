package jwt

//IJwtProvider for sign token struct
type IJwtProvider interface {
	CreateToken(payload TokenPayload) (*TokenDetails, error)
	CheckToken(tokenString string) (*TokenPayload, error)
}
