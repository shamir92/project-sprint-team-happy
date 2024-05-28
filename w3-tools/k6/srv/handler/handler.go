package handler

import (
	"context"
	"fmt"

	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/entity/pb"
	"github.com/nandanugg/HaloSusterTestCasesPSW3B2/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	srv *service.NipService
	pb.UnimplementedNIPServiceServer
}

func NewRequestHandler(srv *service.NipService) *Handler {
	return &Handler{srv: srv}
}
func (h *Handler) ResetAll(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	err := h.srv.Reset(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Println("Reset all!")
	return nil, nil
}
