package hash

type HashService interface {
	CompareHashAndPassword(passwordHash, password string) error
}

type hashService struct{}

func NewHashService() HashService {
	return &hashService{}
}

func (hs *hashService) CompareHashAndPassword(passwordHash, password string) error {
	return nil
}
