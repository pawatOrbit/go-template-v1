package service

import (
	"context"

	"github.com/yourorg/go-api-template/core/exception"
	"github.com/yourorg/go-api-template/internal/model"
	"github.com/yourorg/go-api-template/internal/repository"
)

// ExampleService demonstrates how to implement a service in this template
// This is just an example - replace with your actual business services
type ExampleService interface {
	GetExample(ctx context.Context, req *model.ExampleRequest) (*model.ExampleResponse, error)
	CreateExample(ctx context.Context, req *model.CreateExampleRequest) (*model.CreateExampleResponse, error)
}

type exampleService struct {
	Repo   *repository.Repository
	Errors *exception.MockDataServiceErrors
}

// NewExampleService creates a new example service
func NewExampleService(repo *repository.Repository, errors *exception.MockDataServiceErrors) ExampleService {
	return &exampleService{
		Repo:   repo,
		Errors: errors,
	}
}

// GetExample demonstrates a simple GET operation
func (s *exampleService) GetExample(ctx context.Context, req *model.ExampleRequest) (*model.ExampleResponse, error) {
	// Example business logic - replace with your actual implementation
	// In a real service, you would:
	// 1. Validate the request
	// 2. Call repository to fetch data
	// 3. Transform data if needed
	// 4. Return structured response

	data, err := s.Repo.ExampleRepository.GetExampleByID(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	return &model.ExampleResponse{
		Status: 200,
		Data: model.ExampleResponse_Data{
			ID:      data.ID,
			Message: "Example retrieved successfully",
		},
	}, nil
}

// CreateExample demonstrates a simple CREATE operation
func (s *exampleService) CreateExample(ctx context.Context, req *model.CreateExampleRequest) (*model.CreateExampleResponse, error) {
	// Example business logic - replace with your actual implementation
	// In a real service, you would:
	// 1. Validate the request
	// 2. Call repository to save data
	// 3. Handle errors appropriately
	// 4. Return structured response

	// Create data structure for repository
	data := &repository.ExampleData{
		ID:          "generated-id-123", // In real app, generate UUID
		Name:        req.Name,
		Description: req.Description,
	}

	err := s.Repo.ExampleRepository.CreateExample(ctx, data)
	if err != nil {
		return nil, err
	}

	return &model.CreateExampleResponse{
		Status: 201,
		Data: model.CreateExampleResponse_Data{
			ID:      data.ID,
			Name:    data.Name,
			Message: "Example created successfully",
		},
	}, nil
}