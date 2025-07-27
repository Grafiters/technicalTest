package domain

type ListResponse struct {
	Code    int           // number
	Data    []interface{} // list of data
	Message string        // string
}

type SingleResponse struct {
	Code    int         // number
	Data    interface{} // single data
	Message string      // string
}

type PaginationResponse struct {
	Code    int         // number
	Data    interface{} // list of data
	Message string      // string
}

type PaginationData struct {
	Content    []interface{} `json:"content"`    // List of objects (dynamic content)
	First      bool          `json:"first"`      // Boolean flag indicating if it's the first page
	Last       bool          `json:"last"`       // Boolean flag indicating if it's the last page
	Page       int           `json:"page"`       // Page number, starting from 0
	PageSize   int           `json:"pageSize"`   // Number of items per page (10, 20, 50, 100)
	Sort       []SortObject  `json:"sort"`       // List of sort objects
	TotalSize  int           `json:"totalSize"`  // Total number of items
	TotalPages int           `json:"totalPages"` // Total number of pages
}

type SortObject struct {
	SortBy    string `json:"sortBy"`    // Field name to sort by
	SortOrder string `json:"sortOrder"` // Sort order, either ASC or DESC
}
