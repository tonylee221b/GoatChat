package port

type PasswordHasher interface {
	Hash(plainPW string) (string, error)
	Compare(hashedPW, plainPW string) error
}
