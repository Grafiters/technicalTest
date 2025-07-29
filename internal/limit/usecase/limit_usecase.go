package usecase

import (
	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type limitUsecase struct {
	limitRepo domain.LimitRepository
	logger    *configs.LoggerFormat
}

func NewLimitUsecase(
	limitRepo domain.LimitRepository,
	logger *configs.LoggerFormat,
) domain.LimitUsecase {
	return &limitUsecase{limitRepo, logger}
}

// Create implements domain.LimitUsecase.
func (l *limitUsecase) BulkCreateLimit(ctx *fiber.Ctx, data *domain.BulkLimitInput) error {
	_, err := l.limitRepo.BulkCreateLimit(data)
	if err != nil {
		l.logger.Error("failed to bulk create limit, err: %+v", err)
		return err
	}
	return nil
}

// GetByCustommerID implements domain.LimitUsecase.
func (l *limitUsecase) GetByCustommerID(ctx *fiber.Ctx, CustomerID int64) ([]*domain.Limit, error) {
	limit, err := l.limitRepo.GetByCustommerID(CustomerID)
	if err != nil {
		l.logger.Error("failed to get limit by customer, err: %+v", err)
		return nil, err
	}

	return limit, err
}

// Update implements domain.LimitUsecase.
func (l *limitUsecase) Update(ctx *fiber.Ctx, newData *domain.BulkLimitInput) ([]*domain.Limit, error) {
	_, err := l.limitRepo.BulkCreateLimit(newData)
	if err != nil {
		l.logger.Error("failed to bulk create limit, err: %+v", err)
		return nil, err
	}

	limit, err := l.limitRepo.GetByCustommerID(newData.CustomerID)
	if err != nil {
		l.logger.Error("failed to get limit by customer, err: %+v", err)
		return nil, err
	}

	return limit, nil
}
