package service

import (
	"context"
	"sync"

	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/entity"
	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/entity/pb"
)

type UsedAccountRepository interface {
	PostUsedITAccount(ctx context.Context, usr *entity.UsedUser) error
	PostUsedNurseAccount(ctx context.Context, usr *entity.UsedUser) error
	GetUsedITAccount(ctx context.Context) (*entity.UsedUser, error)
	GetUsedNurseAccount(ctx context.Context) (*entity.UsedUser, error)
	Reset(ctx context.Context) error
}

type NipService struct {
	pb.UnsafeNIPServiceServer
	itNIPs            []uint64
	nurseNIPs         []uint64
	itUsedIndexNIP    int
	nurseUsedIndexNIP int
	itIndexMutex      *sync.Mutex
	nurseIndexMutex   *sync.Mutex
	repo              UsedAccountRepository
}

// Create a new NipServiceServer
func NewNipService(
	itNIPs, nurseNIPs []uint64,
	itIndexMutex *sync.Mutex,
	nurseIndexMutex *sync.Mutex,
	repo UsedAccountRepository,
) *NipService {
	return &NipService{
		itNIPs:          itNIPs,
		nurseNIPs:       nurseNIPs,
		itIndexMutex:    itIndexMutex,
		nurseIndexMutex: nurseIndexMutex,
		repo:            repo,
	}
}
