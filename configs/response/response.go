package response

type Errors struct {
	Code   int      `json:"code"`
	Errors []string `json:"errors"`
}

type Success struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

type Paginate struct {
	Page  int         `json:"page"`
	Limit int         `json:"limit"`
	Total int64       `json:"total"`
	Data  interface{} `json:"data"`
}

func (e Errors) Size() int {
	return len(e.Errors)
}
