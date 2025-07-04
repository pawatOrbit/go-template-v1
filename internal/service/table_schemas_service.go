package service

import (
	"context"

	"github.com/pawatOrbit/ai-mock-data-service/go/core/exception"
	"github.com/pawatOrbit/ai-mock-data-service/go/internal/model"
	"github.com/pawatOrbit/ai-mock-data-service/go/internal/repository"
)

type TableSchemasService interface {
	GetDatabaseSchemaTableNames(ctx context.Context, req *model.GetDatabaseSchemaTableNamesRequest) (*model.GetDatabaseSchemaTableNamesResponse, error)
}

type tableSchemasService struct {
	Repo   *repository.Repository
	Errors *exception.MockDataServiceErrors
}

func NewTableSchemasService(repo *repository.Repository, errors *exception.MockDataServiceErrors) TableSchemasService {
	return &tableSchemasService{
		Repo:   repo,
		Errors: errors,
	}
}

func (s *tableSchemasService) GetDatabaseSchemaTableNames(ctx context.Context, req *model.GetDatabaseSchemaTableNamesRequest) (*model.GetDatabaseSchemaTableNamesResponse, error) {
	// Call the repository method to get the table names
	tableNames, err := s.Repo.TableSchemasRepository.GetDatabaseSchemaTableNames(ctx)
	if err != nil {
		return nil, err
	}

	// Create the response object

	return &model.GetDatabaseSchemaTableNamesResponse{
		Status: 200,
		Data: model.GetDatabaseSchemaTableNamesResponse_Data{
			TableNames: tableNames,
		},
	}, nil
}
