package handler

import (
	"context"

	"github.com/nandanugg/BeliMangTestCasesPB2W4/caddy/entity/pb"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (h *Handler) GetAllMerchantNearestLocations(ctx context.Context, _ *emptypb.Empty) (*pb.AllMerchantNearestRecord, error) {
	res, err := h.srv.GetAllMerchantNearestLocations()
	if err != nil {
		return nil, err
	}

	var zones []*pb.MerchantNearestRecordZone
	for _, zone := range res {
		var recordsRes []*pb.MerchantNearestRecord
		merchantsOrder := make(map[int32]*pb.Merchant)
		for _, record := range zone.Records {
			for order, merchant := range record.MerchantOrder {
				merchantsOrder[int32(order)] = &pb.Merchant{
					MerchantId:     merchant.MerchantId,
					PregeneratedId: merchant.PregeneratedId,
					Location: &pb.LocationPoint{
						Lat:  float32(merchant.Location.Lat),
						Long: float32(merchant.Location.Long),
					},
				}
			}

			recordsRes = append(recordsRes, &pb.MerchantNearestRecord{
				StartingPoint: &pb.LocationPoint{
					Lat:  float32(record.StartingPoint.Lat),
					Long: float32(record.StartingPoint.Long),
				},
				Merchants: merchantsOrder,
			})
		}
		zones = append(zones, &pb.MerchantNearestRecordZone{
			Records: recordsRes,
		})
	}
	return &pb.AllMerchantNearestRecord{
		Zones: zones,
	}, nil
}

func (h *Handler) GetAllMerchantRoutes(ctx context.Context, _ *emptypb.Empty) (*pb.AllGeneratedRoutes, error) {
	srvRes, err := h.srv.GetAllMerchantRoutes()
	if err != nil {
		return nil, err
	}

	zones := []*pb.RouteZone{}
	for _, zone := range srvRes {
		var allGeneratedRoutes []*pb.GeneratedRoutes
		generatedRoute := make(map[int32]*pb.Merchant)
		for _, zoneGeneratedRoutes := range zone.GeneratedTSPRoutes {
			for order, r := range zoneGeneratedRoutes.GeneratedRoutes {
				generatedRoute[int32(order)] = &pb.Merchant{
					MerchantId:     r.MerchantId,
					PregeneratedId: r.PregeneratedId,
					Location: &pb.LocationPoint{
						Lat:  float32(r.Location.Lat),
						Long: float32(r.Location.Long),
					},
				}
			}
			allGeneratedRoutes = append(allGeneratedRoutes, &pb.GeneratedRoutes{
				GeneratedRoutes: generatedRoute,
				StartingPoint: &pb.LocationPoint{
					Lat:  float32(zoneGeneratedRoutes.StartingPoint.Lat),
					Long: float32(zoneGeneratedRoutes.StartingPoint.Long),
				},
				TotalDuration: int32(zoneGeneratedRoutes.TotalDurationInMinute),
				TotalDistance: float32(zoneGeneratedRoutes.TotalDistance),
				StartingIndex: int32(zoneGeneratedRoutes.StartingIndex),
			})
		}
		zones = append(zones, &pb.RouteZone{
			Routes: allGeneratedRoutes,
		})
	}

	return &pb.AllGeneratedRoutes{
		Zone: zones,
	}, nil
}

func (h *Handler) GetTwoZoneMerchantRoutes(ctx context.Context, _ *emptypb.Empty) (*pb.AllGeneratedRoutes, error) {
	srvRes, err := h.srv.GetTwoZoneMerchantRoutes()
	if err != nil {
		return nil, err
	}

	zones := []*pb.RouteZone{}
	for _, zone := range srvRes {
		var allGeneratedRoutes []*pb.GeneratedRoutes
		generatedRoute := make(map[int32]*pb.Merchant)
		for _, zoneGeneratedRoutes := range zone.GeneratedTSPRoutes {
			for order, r := range zoneGeneratedRoutes.GeneratedRoutes {
				generatedRoute[int32(order)] = &pb.Merchant{
					MerchantId:     r.MerchantId,
					PregeneratedId: r.PregeneratedId,
					Location: &pb.LocationPoint{
						Lat:  float32(r.Location.Lat),
						Long: float32(r.Location.Long),
					},
				}
			}
			allGeneratedRoutes = append(allGeneratedRoutes, &pb.GeneratedRoutes{
				GeneratedRoutes: generatedRoute,
				StartingPoint: &pb.LocationPoint{
					Lat:  float32(zoneGeneratedRoutes.StartingPoint.Lat),
					Long: float32(zoneGeneratedRoutes.StartingPoint.Long),
				},
				TotalDuration: int32(zoneGeneratedRoutes.TotalDurationInMinute),
				TotalDistance: float32(zoneGeneratedRoutes.TotalDistance),
				StartingIndex: int32(zoneGeneratedRoutes.StartingIndex),
			})
		}
		zones = append(zones, &pb.RouteZone{
			Routes: allGeneratedRoutes,
		})
	}

	return &pb.AllGeneratedRoutes{
		Zone: zones,
	}, nil
}

func (h *Handler) GetMerchantRoutes(ctx context.Context, _ *emptypb.Empty) (*pb.RouteZone, error) {
	srvRes, err := h.srv.GetMerchantRoutes()
	if err != nil {
		return nil, err
	}
	res := &pb.RouteZone{
		Routes: []*pb.GeneratedRoutes{},
	}

	for _, tspRoute := range srvRes.GeneratedTSPRoutes {
		resRoute := &pb.GeneratedRoutes{
			GeneratedRoutes: make(map[int32]*pb.Merchant),
			StartingPoint: &pb.LocationPoint{
				Lat:  float32(tspRoute.StartingPoint.Lat),
				Long: float32(tspRoute.StartingPoint.Long),
			},
			TotalDistance: float32(tspRoute.TotalDistance),
			TotalDuration: int32(tspRoute.TotalDurationInMinute),
			StartingIndex: int32(tspRoute.StartingIndex),
		}
		for order, merchant := range tspRoute.GeneratedRoutes {
			resRoute.GeneratedRoutes[int32(order)] = &pb.Merchant{
				MerchantId:     merchant.MerchantId,
				PregeneratedId: merchant.PregeneratedId,
				Location: &pb.LocationPoint{
					Lat:  float32(merchant.Location.Lat),
					Long: float32(merchant.Location.Long),
				},
			}
		}
		res.Routes = append(res.Routes, resRoute)
	}

	return res, nil
}

func (h *Handler) GetAllPregeneratedMerchants(ctx context.Context, _ *emptypb.Empty) (*pb.PregeneratedMerchant, error) {
	srvRes := h.srv.GetAllPregeneratedMerchants()

	z := []*pb.Merchant{}
	for _, v := range srvRes {
		z = append(z, &pb.Merchant{
			MerchantId:     v.MerchantId,
			PregeneratedId: v.PregeneratedId,
			Location: &pb.LocationPoint{
				Lat:  float32(v.Location.Lat),
				Long: float32(v.Location.Long),
			},
		})
	}

	return &pb.PregeneratedMerchant{
		Merchant: z,
	}, nil
}

func (h *Handler) GetPregeneratedMerchant(ctx context.Context, _ *emptypb.Empty) (*pb.Merchant, error) {
	srvRes := h.srv.GetPregeneratedMerchant()
	if srvRes == nil {
		return nil, nil
	}

	return &pb.Merchant{
		MerchantId:     srvRes.MerchantId,
		PregeneratedId: srvRes.PregeneratedId,
		Location: &pb.LocationPoint{
			Lat:  float32(srvRes.Location.Lat),
			Long: float32(srvRes.Location.Long),
		},
	}, nil
}

func (h *Handler) GetMerchantNearestLocations(ctx context.Context, _ *emptypb.Empty) (*pb.MerchantNearestRecord, error) {
	zone, err := h.srv.GetMerchantNearestLocations()
	if err != nil {
		return nil, err
	}
	res := &pb.MerchantNearestRecord{
		Merchants: make(map[int32]*pb.Merchant),
	}
	for _, record := range zone.Records {
		res.StartingPoint = &pb.LocationPoint{
			Lat:  float32(record.StartingPoint.Lat),
			Long: float32(record.StartingPoint.Long),
		}
		for i, merchant := range record.MerchantOrder {
			res.Merchants[int32(i)] = &pb.Merchant{
				MerchantId:     merchant.MerchantId,
				PregeneratedId: merchant.PregeneratedId,
				Location: &pb.LocationPoint{
					Lat:  float32(merchant.Location.Lat),
					Long: float32(merchant.Location.Long),
				},
			}
		}
	}
	return res, nil
}

func (h *Handler) AssignPregeneratedMerchant(ctx context.Context, req *pb.AssignMerchant) (*emptypb.Empty, error) {
	err := h.srv.AssignPregeneratedMerchant(req.PregeneratedId, req.MerchantId)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
