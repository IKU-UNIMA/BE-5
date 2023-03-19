package util

type Pagination struct {
	Page int         `json:"page"`
	Data interface{} `json:"data"`
}

func CountOffset(page int) int {
	return (page - 1) * 10
}
