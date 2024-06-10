package services

import (
	"context"
	"fmt"
	"strconv"

	"github.com/LilLebowski/loyalty-system/internal/storage"
)

type Order struct {
	db storage.Storage
}

func OrderInit(storage storage.Storage) *Order {
	return &Order{db: storage}
}

func (os *Order) GetOrdersNotProcessed() ([]string, error) {
	orders, err := os.db.GetOrdersInProgress(context.Background())
	if err != nil {
		return nil, err
	}
	numbers := make([]string, 0)
	for _, order := range orders {
		num, err := strconv.ParseInt(order.Number, 10, 64)
		if err != nil {
			return nil, err
		}
		numbers = append(numbers, strconv.FormatInt(num, 10))
	}
	if len(numbers) == 0 {
		return nil, fmt.Errorf("empty list")
	}
	return numbers, nil
}

func (os *Order) UpdateOrder(accrual *storage.Accrual) error {
	order := storage.Accrual{Order: accrual.Order, Accrual: accrual.Accrual, Status: accrual.Status}
	err := os.db.UpdateOrder(context.Background(), order)
	if err != nil {
		return err
	}
	return nil
}
