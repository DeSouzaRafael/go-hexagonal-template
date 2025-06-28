package port

type TokenGenerator interface {
	GenerateToken(userID string) (string, error)
}
