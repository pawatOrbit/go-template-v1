package model

type GetDatabaseSchemaTableNamesRequest struct{}

type GetDatabaseSchemaTableNamesResponse_Data struct {
	TableNames []string `json:"table_names"`
}

type GetDatabaseSchemaTableNamesResponse struct {
	Status int                                      `json:"status"`
	Data   GetDatabaseSchemaTableNamesResponse_Data `json:"data"`
}
