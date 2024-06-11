package handler

import (
	"github.com/nandanugg/BeliMangTestCasesPB2W4/caddy/entity/pb"
	"github.com/nandanugg/BeliMangTestCasesPB2W4/caddy/service"
)

type Handler struct {
	srv *service.MerchantService
	pb.UnimplementedMerchantServiceServer
}

func NewHandler(srv *service.MerchantService) *Handler {
	return &Handler{srv: srv}
}
