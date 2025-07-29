package repository_test

import (
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/auth/repository"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	cleanup := func() {
		db.Close()
	}

	return gormDB, mock, cleanup
}

func TestLogin_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	logger := &configs.LoggerFormat{}
	repo := repository.NewAuthRepository(db, logger)

	email := "test@example.com"
	expected := &domain.Customer{
		ID:    1,
		Email: email,
	}

	rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(expected.ID, expected.Email)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE email = $1 ORDER BY "customers"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)

	result, err := repo.Login(&domain.AuthRequest{Email: email})
	assert.NoError(t, err)
	assert.Equal(t, expected.Email, result.Email)
}

func TestLogin_NotFound(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	logger := &configs.LoggerFormat{}
	repo := repository.NewAuthRepository(db, logger)

	email := "notfound@example.com"

	mock.ExpectQuery(`SELECT * FROM "customers" WHERE email = $1 ORDER BY "customers"."id" LIMIT $2`).
		WithArgs("notfound@example.com", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := repo.Login(&domain.AuthRequest{Email: email})
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestRegister_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	logger := &configs.LoggerFormat{}
	defer cleanup()

	repo := repository.NewAuthRepository(db, logger)

	mockCustomer := &domain.Customer{
		Email:    "test@example.com",
		Password: "securepassword",
	}

	// Expect begin transaction
	mock.ExpectBegin()

	mock.ExpectQuery(regexp.QuoteMeta(`
	INSERT INTO "customers" ("email","password","nik","full_name","legal_name","birth_place","birth_date","salary","ktp_image_url","selfie_image_url","created_at","updated_at") 
	VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12) RETURNING "id"`)).
		WithArgs(
			mockCustomer.Email,
			mockCustomer.Password,
			mockCustomer.NIK,
			mockCustomer.FullName,
			mockCustomer.LegalName,
			mockCustomer.BirthPlace,
			mockCustomer.BirthDate,
			mockCustomer.Salary,
			mockCustomer.KTPImageUrl,
			mockCustomer.SelfieImageUrl,
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// Expect commit
	mock.ExpectCommit()

	result, err := repo.Register(&domain.RegisterRequest{
		Email:    mockCustomer.Email,
		Password: mockCustomer.Password,
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(1), result.ID)
	assert.Equal(t, mockCustomer.Email, result.Email)
}
func TestGetByEmail_Success(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	logger := &configs.LoggerFormat{}
	repo := repository.NewAuthRepository(db, logger)

	email := "test@example.com"
	expected := &domain.Customer{
		ID:    1,
		Email: email,
	}

	rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(expected.ID, expected.Email)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE email = $1 ORDER BY "customers"."id" LIMIT $2`)).
		WithArgs(email, 1).
		WillReturnRows(rows)

	result, err := repo.Login(&domain.AuthRequest{Email: email})
	assert.NoError(t, err)
	assert.Equal(t, expected.Email, result.Email)
}

func TestGetByEmail_Failed(t *testing.T) {
	db, mock, cleanup := setupMockDB(t)
	defer cleanup()

	logger := &configs.LoggerFormat{}
	repo := repository.NewAuthRepository(db, logger)

	email := "error@example.com"

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "customers" WHERE email = $1 ORDER BY "customers"."id" LIMIT 1`)).
		WithArgs(email).
		WillReturnError(errors.New("query error"))

	result, err := repo.GetByEmail(email)
	assert.Error(t, err)
	assert.Nil(t, result)
}
