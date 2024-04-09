package model

type PaginationParams struct {
	Page  int
	Limit int
	Skip  int
}

type PaginationResponse struct {
	CurrentPage int `json:"page"`
	Data        any `json:"data"`
	Limit       int `json:"limit"`
	NextPage    int `json:"nextPage"`
}
