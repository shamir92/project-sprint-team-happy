package handler

import (
	"context"

	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/entity"
	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/entity/pb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *Handler) GetUsedIt(ctx context.Context, req *emptypb.Empty) (*pb.PostUsedAcc, error) {
	usr, err := s.srv.GetItUsedAccount(ctx)

	if err != nil {
		return nil, err
	}
	if usr == nil {
		return nil, nil
	}

	return &pb.PostUsedAcc{
		Nip:      usr.Nip,
		Password: usr.Password,
	}, nil
}

func (s *Handler) PostUsedIT(ctx context.Context, req *pb.PostUsedAcc) (*emptypb.Empty, error) {
	err := s.srv.AddItUsedAccount(ctx, &entity.UsedUser{
		Nip:      req.Nip,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (s *Handler) GetItNip(ctx context.Context, req *emptypb.Empty) (*pb.GetNipResponse, error) {
	return &pb.GetNipResponse{
		Nip: s.srv.GetItNip(),
	}, nil
}
