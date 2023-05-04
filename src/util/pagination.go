package util

type Pagination struct {
	Limit       int         `json:"limit"`
	Page        int         `json:"page"`
	TotalPage   int         `json:"total_page"`
	TotalResult int         `json:"total_result"`
	Data        interface{} `json:"data"`
}

func CountOffset(page, limit int) int {
	return (page - 1) * limit
}

func CountTotalPage(totalResult, limit int) int {
	total := totalResult / limit
	if totalResult%limit > 0 {
		total++
	}
	return total
}
