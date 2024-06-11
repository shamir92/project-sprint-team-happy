package handler

import (
	"context"

	"github.com/nandanugg/BeliMangTestCasesPB2W4/caddy/entity"
	"github.com/nandanugg/BeliMangTestCasesPB2W4/caddy/entity/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) InitMerchantNearestLocations(ctx context.Context, req *pb.InitPregenerated) (*emptypb.Empty, error) {
	h.srv.InitMerchantNearestLocations(int(req.GenerateCount))
	return &emptypb.Empty{}, nil
}

func (h *Handler) InitPegeneratedTSPMerchants(ctx context.Context, req *pb.InitPregenerated) (*emptypb.Empty, error) {
	h.srv.InitPegeneratedTSPMerchants(int(req.GenerateCount))
	return &emptypb.Empty{}, nil
}

func (h *Handler) InitZonesWithPregeneratedMerchants(ctx context.Context, req *pb.InitZonesRequest) (*emptypb.Empty, error) {
	h.srv.InitZonesWithPregeneratedMerchants(entity.MerchantZoneOpts{
		Area:                     float64(req.Area),
		Gap:                      float64(req.Gap),
		NumberOfZones:            int(req.NumberOfZones),
		NumberOfMerchantsPerZone: int(req.NumberOfMerchantsPerZone),
	})
	return &emptypb.Empty{}, nil
}

func (h *Handler) ResetAll(ctx context.Context, _ *emptypb.Empty) (*emptypb.Empty, error) {
	h.srv.ResetAll()
	return &emptypb.Empty{}, nil
}
