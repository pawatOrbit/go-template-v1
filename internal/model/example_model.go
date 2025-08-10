package model

// Example models for demonstration - replace with your actual models

// ExampleRequest represents a request to get an example
type ExampleRequest struct {
	ID string `json:"id" validate:"required"`
}

// ExampleResponse represents a response containing example data
type ExampleResponse struct {
	Status int                   `json:"status"`
	Data   ExampleResponse_Data  `json:"data"`
}

type ExampleResponse_Data struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

// CreateExampleRequest represents a request to create a new example
type CreateExampleRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=100"`
	Description string `json:"description,omitempty" validate:"max=500"`
}

// CreateExampleResponse represents a response after creating an example
type CreateExampleResponse struct {
	Status int                         `json:"status"`
	Data   CreateExampleResponse_Data  `json:"data"`
}

type CreateExampleResponse_Data struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}