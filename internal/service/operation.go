package service

import (
	"context"
	"errors"
	"usersegments/internal/entity"
	"usersegments/internal/repository/repoerrors"

	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=operation.go -destination=mocks/operation_mock.go
type OperationRepo interface {
	GetOperationHistory(ctx context.Context, user_id int, month, year string) ([]entity.Operation, error)
}

type OperationService struct {
	operationRepo OperationRepo
}

func NewOperationService(operationRepo OperationRepo) *OperationService {
	return &OperationService{
		operationRepo: operationRepo,
	}
}

type GetOperationHistoryInput struct {
	Month string
	Year  string
}

func (o *OperationService) GetOperationHistory(ctx context.Context, user_id int, input GetOperationHistoryInput) ([]entity.Operation, error) {
	operations, err := o.operationRepo.GetOperationHistory(ctx, user_id, input.Month, input.Year)
	if err != nil {
		if errors.Is(err, repoerrors.ErrOperationHistoryNotFound) {
			return nil, ErrOperationHistoryNotFound
		}
		logrus.Errorf("service - o.operationRepo.GetOperationHistory. %v", err)
		return nil, ErrCannotGetOperationHistory
	}
	return operations, nil
}
