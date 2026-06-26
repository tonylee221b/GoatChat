package in

type IdentityHandler struct{}

func NewIdentityHandler() (*IdentityHandler, error) {
	return &IdentityHandler{}, nil
}
