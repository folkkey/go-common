package gocommon

type PagingResult struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

type PagingQuery struct {
	Page    int     `query:"page" form:"page"`
	Size    int     `query:"size" form:"size"`
	OrderBy *string `query:"order_by" form:"order_by" default:"desc"`
	SortBy  *string `query:"sort_by" form:"sort_by" default:"created_at"`
}

func NewPagingQuery() PagingQuery {
	return PagingQuery{
		Page: 0,
		Size: 25,
	}
}

type PagedResultDto[Dto any] struct {
	Total int64 `json:"total"`
	Items []Dto `json:"items"`
}
