package service

import (
	"github.com/pawatOrbit/ai-mock-data-service/go/config"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/exception"
	"github.com/pawatOrbit/ai-mock-data-service/go/core/httpclient"
	"github.com/pawatOrbit/ai-mock-data-service/go/internal/repository"
	"github.com/pawatOrbit/ai-mock-data-service/go/utils"
)

type Service struct {
	Config *config.Config
	Errors *exception.MockDataServiceErrors

	TableSchemasService TableSchemasService
}

func NewService(
	repo *repository.Repository,
	config *config.Config,
	errors *exception.MockDataServiceErrors,
	utils *utils.Utils,
	lmStudioClient *httpclient.LmStudioServiceClient,
) Service {
	return Service{
		Config:              config,
		Errors:              errors,
		TableSchemasService: NewTableSchemasService(repo, errors),
	}
}
