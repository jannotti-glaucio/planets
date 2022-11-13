package hash

//IHashProvider interface of hash
type IHashProvider interface {
	Create(payload string) string
	Compare(hashed string, payload string) bool
}
