package port

type CSRFService interface {
	GenerateToken() (string, error)
	VerifyToken(token, expectedToken string) bool
}
