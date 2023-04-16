package util

type Pagination struct {
	Page int         `json:"page"`
	Data interface{} `json:"data"`
}

func CountOffset(page, limit int) int {
	return (page - 1) * limit
}
