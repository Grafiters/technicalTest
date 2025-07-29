package usecase

import (
	"fmt"
	"strings"

	"github.com/Grafiters/archive/configs"
	"github.com/Grafiters/archive/internal/domain"
	"github.com/Grafiters/archive/utils"
	"github.com/gofiber/fiber/v2"
)

type customerUsecase struct {
	customerRepo domain.CustomerRepository
	limitRepo    domain.LimitRepository
	minio        *configs.MinioConfig
	logger       *configs.LoggerFormat
}

func NewCustomerUsecase(
	customerRepo domain.CustomerRepository,
	limitRepo domain.LimitRepository,
	minio *configs.MinioConfig,
	logger *configs.LoggerFormat,
) domain.CustomerUsecase {
	return &customerUsecase{
		customerRepo: customerRepo,
		limitRepo:    limitRepo,
		minio:        minio,
		logger:       logger,
	}
}

func (c *customerUsecase) Create(ctx *fiber.Ctx, data *domain.CustomerInput) (*domain.CustomerResponse, error) {
	customer, err := c.customerRepo.Create(data)
	if err != nil {
		c.logger.Error("failed to create customer, err: %+v", err)
		return &domain.CustomerResponse{}, fmt.Errorf(utils.ProsessError)
	}

	if data.KTPImage != "" {
		err := c.validateImage(data.KTPImage)
		if err != nil {
			c.logger.Error("failed to validate ktp image, err: ", err)
			return nil, err
		}

		sizeTrue := c.validateImageSize(data.KTPImage)
		if !sizeTrue {
			c.logger.Error("file ktp image to large, max: ", configs.MaxSizeMB)
			return nil, fmt.Errorf("image ktp to large")
		}
	}

	if data.SelfieImage != "" {
		err := c.validateImage(data.SelfieImage)
		if err != nil {
			c.logger.Error("failed to validate ktp image, err: ", err)
			return nil, err
		}

		sizeTrue := c.validateImageSize(data.SelfieImage)
		if !sizeTrue {
			c.logger.Error("file ktp image to large, max: ", configs.MaxSizeMB)
			return nil, fmt.Errorf("image ktp to large")
		}
	}

	customerResponse := customer.ToCustomerResponse()
	return customerResponse, nil
}

func (c *customerUsecase) Get(ctx *fiber.Ctx, ID int64) (*domain.CustomerResponse, error) {
	customer, err := c.customerRepo.GetByID(ID)
	if err != nil {
		c.logger.Error("failed to get customer, err: %+v", err)
		return &domain.CustomerResponse{}, err
	}

	customerResponse := customer.ToCustomerResponse()

	limit, err := c.limitRepo.GetByCustommerID(customer.ID)
	if err != nil {
		c.logger.Error("failed to limit customer, err: %+v", err)
		return &domain.CustomerResponse{}, err
	}

	limitResponse := []*domain.LimitResponse{}
	for _, val := range limit {
		limitResponse = append(limitResponse, val.ToLimitResponse())
	}

	customerResponse.Limit = limitResponse

	ktpImage := c.getImagePublisherData(customer.KTPImageUrl)
	selfieImage := c.getImagePublisherData(customer.SelfieImageUrl)

	customerResponse.KTPImageUrl = ktpImage
	customerResponse.SelfieImageUrl = selfieImage

	return customerResponse, nil
}

func (c *customerUsecase) Update(ctx *fiber.Ctx, ID int64, data *domain.CustomerUpdate) (*domain.CustomerResponse, error) {
	if data.KTPImage != "" {
		err := c.validateImage(data.KTPImage)
		if err != nil {
			c.logger.Error("failed to validate ktp image, err: ", err)
			return nil, err
		}

		sizeTrue := c.validateImageSize(data.KTPImage)
		if !sizeTrue {
			c.logger.Error("file ktp image to large, max: ", configs.MaxSizeMB)
			return nil, fmt.Errorf("image ktp to large")
		}
	}

	if data.SelfieImage != "" {
		err := c.validateImage(data.SelfieImage)
		if err != nil {
			c.logger.Error("failed to validate ktp image, err: ", err)
			return nil, err
		}

		sizeTrue := c.validateImageSize(data.SelfieImage)
		if !sizeTrue {
			c.logger.Error("file ktp image to large, max: ", configs.MaxSizeMB)
			return nil, fmt.Errorf("image ktp to large")
		}
	}

	ktp, selfie, err := c.handleUploadImage(ID, data.KTPImage, data.SelfieImage)
	if err != nil {
		c.logger.Error("failed uplaod image to minio, err: ", err)
		return nil, err
	}

	data.KTPImage = ktp
	data.SelfieImage = selfie

	customer, err := c.customerRepo.Update(ID, data)
	if err != nil {
		c.logger.Error("failed to update customer, err: %+v", err)
		return &domain.CustomerResponse{}, fmt.Errorf(utils.ProsessError)
	}

	customerResponse := customer.ToCustomerResponse()
	ktpImage := c.getImagePublisherData(customer.KTPImageUrl)
	selfieImage := c.getImagePublisherData(customer.SelfieImageUrl)

	customerResponse.KTPImageUrl = ktpImage
	customerResponse.SelfieImageUrl = selfieImage
	return customerResponse, nil
}

func (c *customerUsecase) UpdateSalary(ctx *fiber.Ctx, ID int64, salary *domain.CustomerUpdateSalary) (*domain.CustomerResponse, error) {
	customer, err := c.customerRepo.UpdateSalary(ID, *salary)
	if err != nil {
		c.logger.Error("failed to update customer, err: %+v", err)
		return &domain.CustomerResponse{}, fmt.Errorf(utils.ProsessError)
	}

	err = c.handleCerateLimit(customer)
	if err != nil {
		c.logger.Error("failed to generate new data limit, err: ", err)
	}
	customerResponse := customer.ToCustomerResponse()
	return customerResponse, nil
}

func (c *customerUsecase) handleUploadImage(customerID int64, ktpImage, selfieImage string) (string, string, error) {
	ktpUpload, err := configs.Minio.UplaodFile(customerID, ktpImage)
	if err != nil {
		c.logger.Error("errored upload ktp image", err)
		return "", "", err
	}

	selfieUpload, err := configs.Minio.UplaodFile(customerID, selfieImage)
	if err != nil {
		c.logger.Error("errored upload selfie image", err)
		return "", "", err
	}

	return ktpUpload, selfieUpload, nil
}

func (c *customerUsecase) validateImage(image string) error {
	if !strings.Contains(image, "data:image/") {
		return fmt.Errorf("image invalid string must be contrain data:image/")
	}

	include := false
	for _, i := range configs.ImageAllowList {
		if strings.Contains(image, i) {
			include = true
		}
	}

	if !include {
		return fmt.Errorf("failed to validate format type image, must be %+v", configs.ImageAllowList)
	}

	return nil
}

func (c *customerUsecase) validateImageSize(image string) bool {
	data := configs.RemoveDataURLPrefix(image)

	size := len(data)
	sizeInMB := float64(size) / (1024 * 1024)

	return sizeInMB <= configs.MaxSizeMB
}

func (c *customerUsecase) handleCerateLimit(data *domain.Customer) error {
	tenorLimit := domain.BuildTenorFactor(data.Salary)
	limitInput := &domain.BulkLimitInput{
		CustomerID: data.ID,
		LimitTenor: tenorLimit,
	}

	_, err := c.limitRepo.BulkCreateLimit(limitInput)
	if err != nil {
		c.logger.Error("failed to bulk create limit, err: %+v", err)
		return err
	}

	return nil

}

func (c *customerUsecase) getImagePublisherData(image string) string {
	var err error
	if strings.Contains(image, "http") || strings.Contains(image, "https") || image == "" {
		image = image
	} else {
		image, err = configs.Minio.GetUrlPublishData(image)
		if err != nil {
			configs.Logger.Error("Error to add banner url -> ", err.Error())
		}
	}

	return image
}
