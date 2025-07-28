package delivery

import (
	"fmt"
	"strings"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"gorm.io/gorm"
)

type mysqlLimitRepository struct {
	db     *gorm.DB
	logger *configs.LoggerFormat
}

func NewLimitRepository(
	db *gorm.DB,
	logger *configs.LoggerFormat,
) domain.LimitRepository {
	return &mysqlLimitRepository{db, logger}
}

// Get implements domain.LimitRepository.
func (m *mysqlLimitRepository) Get(filter *domain.LimitFilter) ([]*domain.Limit, error) {
	var (
		limit            []*domain.Limit
		whereQuery, args = getWhereClause(filter)
	)
	query := m.db.Where(whereQuery, args...)

	if err := query.Find(&limit).Error; err != nil {
		m.logger.Error("failed to query limit data, err: %+v", err)
		return nil, err
	}

	return limit, nil
}

func (m *mysqlLimitRepository) BulkCreateLimit(data *domain.BulkLimitInput) ([]*int64, error) {
	err := m.bulkDeleteLimit(data)
	if err != nil {
		m.logger.Error("failed to deleted data limit customer: err %+v", err)
		return nil, err
	}

	var limits []domain.Limit
	for _, input := range data.LimitTenor {
		limits = append(limits, domain.Limit{
			CustomerID: data.CustomerID,
			Tenor:      input.Tenor,
			Amount:     input.Amount,
		})
	}

	if err := m.db.Create(&limits).Error; err != nil {
		m.logger.Error("failed to bulk insert limits: %+v", err)
		return nil, err
	}

	var ids []*int64
	for _, limit := range limits {
		id := limit.ID
		ids = append(ids, &id)
	}

	return ids, nil
}

func (m *mysqlLimitRepository) BulkUpdateLimit(data *domain.BulkLimitInput) ([]*int64, error) {
	err := m.bulkDeleteLimit(data)
	if err != nil {
		m.logger.Error("failed to deleted data limit customer: err %+v", err)
		return nil, err
	}

	var limits []domain.Limit
	for _, input := range data.LimitTenor {
		limits = append(limits, domain.Limit{
			CustomerID: data.CustomerID,
			Tenor:      input.Tenor,
			Amount:     input.Amount,
		})
	}

	if err := m.db.Create(&limits).Error; err != nil {
		m.logger.Error("failed to bulk insert limits: %+v", err)
		return nil, err
	}

	var ids []*int64
	for _, limit := range limits {
		id := limit.ID
		ids = append(ids, &id)
	}

	return ids, nil
}

func (m *mysqlLimitRepository) GetByCustommerID(CustomerID int64) ([]*domain.Limit, error) {
	var (
		limit []*domain.Limit
	)
	query := m.db.Where("customer_id = ?", CustomerID)

	if err := query.Find(&limit).Error; err != nil {
		m.logger.Error("failed to query limit data, err: %+v", err)
		return nil, err
	}

	return limit, nil
}

func (m *mysqlLimitRepository) bulkDeleteLimit(data *domain.BulkLimitInput) error {
	err := m.db.Where("customer_id = ?", data.CustomerID).Delete(&domain.Limit{})
	return err.Error
}

func getWhereClause(filter *domain.LimitFilter) (string, []interface{}) {
	var (
		whereClause  string
		whereClauses []string
		args         []interface{}
	)

	if filter.Tenor != 0 {
		query := fmt.Sprintf("tenor = ?")
		whereClauses = append(whereClauses, query)
		args = append(args, filter.Tenor)
	}

	if filter.CustomerID != 0 {
		query := fmt.Sprintf("customer_id = ?")
		whereClauses = append(whereClauses, query)
		args = append(args, filter.CustomerID)
	}

	if filter.Amount != 0 {
		query := fmt.Sprintf("amount >= ?")
		whereClauses = append(whereClauses, query)
		args = append(args, filter.Amount)
	}

	// Build where clause based on provided filter.
	if len(whereClauses) > 0 {
		whereClause = fmt.Sprintf("WHERE %s\n", strings.Join(whereClauses, ` AND `))
	}

	return whereClause, args
}
