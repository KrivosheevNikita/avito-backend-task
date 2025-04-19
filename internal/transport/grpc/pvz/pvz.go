package pvz

import (
	"context"

	pb "pvz-service/api/grpc/pvz/pvz_v1"

	"pvz-service/internal/db"
	"pvz-service/pkg/logger"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type PVZServer struct {
	pb.UnimplementedPVZServiceServer
}

func NewPVZServer() *PVZServer {
	return &PVZServer{}
}

var GetPvzFunc = db.GetPvz

func (s *PVZServer) GetPVZList(ctx context.Context, req *pb.GetPVZListRequest) (*pb.GetPVZListResponse, error) {
	logger.Info("Получен gRPC-запрос GetPVZList", nil)

	pvzs, err := GetPvzFunc()
	if err != nil {
		logger.Error("Ошибка при получении списка ПВЗ", err)
		return nil, err
	}

	var result []*pb.PVZ
	for _, pvz := range pvzs {
		result = append(result, &pb.PVZ{
			Id:               pvz.ID.String(),
			RegistrationDate: timestamppb.New(pvz.RegistrationDate),
			City:             pvz.City,
		})
	}

	logger.Info("Успешно возвращён список ПВЗ", map[string]interface{}{
		"count": len(result),
	})

	return &pb.GetPVZListResponse{
		Pvzs: result,
	}, nil
}
