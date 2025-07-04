package model

type HealthReq struct {
	Name string `json:"name"`
}

type HealthResp struct {
	Status   int    `json:"status"`
	Response string `json:"response"`
}
