package service

import (
	"github.com/yourorg/go-api-template/config"
	"github.com/yourorg/go-api-template/core/exception"
	"github.com/yourorg/go-api-template/core/httpclient"
	"github.com/yourorg/go-api-template/internal/repository"
	"github.com/yourorg/go-api-template/utils"
)

type Service struct {
	Config *config.Config
	Errors *exception.MockDataServiceErrors

	// Core services
	HealthService  HealthServiceInterface
	
	// Example services - replace with your actual services
	ExampleService ExampleService
}

func NewService(
	repo *repository.Repository,
	config *config.Config,
	errors *exception.MockDataServiceErrors,
	utils *utils.Utils,
	lmStudioClient *httpclient.LmStudioServiceClient,
) Service {
	return Service{
		Config: config,
		Errors: errors,

		// Core services
		HealthService: NewHealthService(repo),

		// Example services - replace with your actual services
		ExampleService: NewExampleService(repo, errors),
	}
}
