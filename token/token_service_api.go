package token

type PayLoad any
type CustomToken any

// Token.Manager interface is an API for creating and verifying tokens
type Manager interface {
	CreateToken(key PayLoad) (token CustomToken, err error)
	VerifyToken(token CustomToken, key PayLoad) (err error)
}
