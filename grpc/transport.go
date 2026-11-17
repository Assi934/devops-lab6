package grpc

import (
	"context"

	"example.com/car-service/car"
	"example.com/car-service/transport"

	"github.com/go-kit/kit/log"
	kitgrpc "github.com/go-kit/kit/transport/grpc"

	pb "example.com/car-service/pb"
)

// NewGRPCServer связывает эндпоинты Go kit с gRPC-транспортом
func NewGRPCServer(endpoints transport.Endpoints, logger log.Logger) pb.CarServiceServer {
	options := []kitgrpc.ServerOption{
		kitgrpc.ServerErrorLogger(logger),
	}
	return &grpcServer{
		create:   kitgrpc.NewServer(endpoints.Create, decodeCreateRequest, encodeCreateResponse, options...),
		getByID:  kitgrpc.NewServer(endpoints.GetByID, decodeGetByIDRequest, encodeGetByIDResponse, options...),
		markSold: kitgrpc.NewServer(endpoints.MarkSold, decodeMarkSoldRequest, encodeMarkSoldResponse, options...),
	}
}

type grpcServer struct {
	pb.UnimplementedCarServiceServer
	create   kitgrpc.Handler
	getByID  kitgrpc.Handler
	markSold kitgrpc.Handler
}

func (s *grpcServer) Create(ctx context.Context, req *pb.CreateRequest) (*pb.CreateResponse, error) {
	_, rep, err := s.create.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.CreateResponse), nil
}

func (s *grpcServer) GetByID(ctx context.Context, req *pb.GetByIDRequest) (*pb.GetByIDResponse, error) {
	_, rep, err := s.getByID.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.GetByIDResponse), nil
}

func (s *grpcServer) MarkSold(ctx context.Context, req *pb.MarkSoldRequest) (*pb.MarkSoldResponse, error) {
	_, rep, err := s.markSold.ServeGRPC(ctx, req)
	if err != nil {
		return nil, err
	}
	return rep.(*pb.MarkSoldResponse), nil
}

func decodeCreateRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.CreateRequest)
	return transport.CreateRequest{Car: carFromProto(req.Car)}, nil
}

func encodeCreateResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.CreateResponse)
	return &pb.CreateResponse{Id: resp.ID}, nil
}

func decodeGetByIDRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.GetByIDRequest)
	return transport.GetByIDRequest{ID: req.Id}, nil
}

func encodeGetByIDResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(transport.GetByIDResponse)
	return &pb.GetByIDResponse{Car: carToProto(resp.Car)}, nil
}

func decodeMarkSoldRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb.MarkSoldRequest)
	return transport.MarkSoldRequest{ID: req.Id}, nil
}

func encodeMarkSoldResponse(_ context.Context, response interface{}) (interface{}, error) {
	return &pb.MarkSoldResponse{}, nil
}

func carFromProto(c *pb.Car) car.Car {
	return car.Car{
		ID:    c.Id,
		Brand: c.Brand,
		Model: c.Model,
		Year:  int(c.Year),
		Price: c.Price,
		Sold:  c.Sold,
	}
}

func carToProto(c car.Car) *pb.Car {
	return &pb.Car{
		Id:    c.ID,
		Brand: c.Brand,
		Model: c.Model,
		Year:  int32(c.Year),
		Price: c.Price,
		Sold:  c.Sold,
	}
}
