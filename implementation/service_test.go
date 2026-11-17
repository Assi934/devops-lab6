package implementation

import (
	"context"
	"testing"

	"example.com/car-service/car"
	"example.com/car-service/repository"

	"github.com/go-kit/kit/log"
)

func newTestService(t *testing.T) car.Service {
	t.Helper()
	return NewService(repository.NewMemoryRepository(), log.NewNopLogger())
}

// TestCreateAndGetByID проверяет создание и чтение автомобиля
func TestCreateAndGetByID(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()

	id, err := svc.Create(ctx, car.Car{Brand: "Toyota", Model: "Camry", Year: 2023, Price: 2500000})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if id == "" {
		t.Fatal("Create() вернул пустой id")
	}

	c, err := svc.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if c.Brand != "Toyota" || c.Model != "Camry" || c.Year != 2023 {
		t.Errorf("неожиданный автомобиль: %+v", c)
	}
	if c.Sold {
		t.Error("новый автомобиль не должен быть продан")
	}
}

// TestGetByIDNotFound проверяет ошибку при чтении несуществующего автомобиля
func TestGetByIDNotFound(t *testing.T) {
	svc := newTestService(t)

	if _, err := svc.GetByID(context.Background(), "missing"); err != car.ErrCarNotFound {
		t.Errorf("GetByID(missing) error = %v, want %v", err, car.ErrCarNotFound)
	}
}

// TestMarkSold проверяет перевод автомобиля в состояние «продан»
func TestMarkSold(t *testing.T) {
	svc := newTestService(t)
	ctx := context.Background()

	id, err := svc.Create(ctx, car.Car{Brand: "BMW", Model: "X5", Year: 2024, Price: 8500000})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if err := svc.MarkSold(ctx, id); err != nil {
		t.Fatalf("MarkSold() error = %v", err)
	}

	c, err := svc.GetByID(ctx, id)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if !c.Sold {
		t.Error("после MarkSold автомобиль должен быть продан")
	}
}

// TestMarkSoldNotFound проверяет ошибку для несуществующего автомобиля
func TestMarkSoldNotFound(t *testing.T) {
	svc := newTestService(t)

	if err := svc.MarkSold(context.Background(), "missing"); err != car.ErrCmdRepository {
		t.Errorf("MarkSold(missing) error = %v, want %v", err, car.ErrCmdRepository)
	}
}
