package searchservice

import "time"

type SortOrder int

const (
	Ascending SortOrder = iota
	Descending
)

type SearchRequest struct {
	Responder  string
	Start      time.Time
	End        time.Time
	LimitCount int
	Order      SortOrder
}

type SearchResponse struct {
	CurrentPage int   `json:"current_page"`
	TotalCount  int   `json:"total_count"`
	Result      []CFS `json:"cfs_list"`
}

type ISearchService interface {
	Search(*SearchRequest) *SearchResponse
}

type SearchService struct{}

func (*SearchService) Search(request *SearchRequest) *SearchResponse {
	return nil
}
