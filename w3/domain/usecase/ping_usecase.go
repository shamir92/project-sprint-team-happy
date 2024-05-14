package usecase

// import "halosuster/domain/repository"

type pingUsecase struct {
}

type IPingUsecase interface {
	GetPing() (bool, error)
}

func NewPingUsecase() *pingUsecase {
	return &pingUsecase{}
}

func (puc *pingUsecase) GetPing() (bool, error) {
	return true, nil
}
