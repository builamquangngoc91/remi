package up

type OffsetPaging struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
