package transport

import (
	"context"

	"example.com/car-service/car"
	"github.com/go-kit/kit/endpoint"
)

// DTO для запросов и ответов
type CreateRequest struct {
	Car car.Car `json:"car"`
}

type CreateResponse struct {
	ID  string `json:"id"`
	Err error  `json:"error,omitempty"`
}

type GetByIDRequest struct {
	ID string `json:"id"`
}

type GetByIDResponse struct {
	Car car.Car `json:"car"`
	Err error   `json:"error,omitempty"`
}

type MarkSoldRequest struct {
	ID string `json:"id"`
}

type MarkSoldResponse struct {
	Err error `json:"error,omitempty"`
}

// Endpoints — набор эндпоинтов сервиса автомобилей
type Endpoints struct {
	Create       endpoint.Endpoint
	GetByID      endpoint.Endpoint
	MarkSold     endpoint.Endpoint
}

// MakeEndpoints создаёт эндпоинты для сервиса автомобилей
func MakeEndpoints(s car.Service) Endpoints {
	return Endpoints{
		Create:       makeCreateEndpoint(s),
		GetByID:      makeGetByIDEndpoint(s),
		MarkSold:     makeMarkSoldEndpoint(s),
	}
}

func makeCreateEndpoint(s car.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateRequest)
		id, err := s.Create(ctx, req.Car)
		return CreateResponse{ID: id, Err: err}, nil
	}
}

func makeGetByIDEndpoint(s car.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetByIDRequest)
		c, err := s.GetByID(ctx, req.ID)
		return GetByIDResponse{Car: c, Err: err}, nil
	}
}

func makeMarkSoldEndpoint(s car.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(MarkSoldRequest)
		err := s.MarkSold(ctx, req.ID)
		return MarkSoldResponse{Err: err}, nil
	}
}
