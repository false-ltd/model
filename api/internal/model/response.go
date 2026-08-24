package model

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PagedResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    interface{} `json:"meta,omitempty"`
}

type PageMeta struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

type CountMeta struct {
	Count int `json:"count"`
}

func ErrorResponse(code int, msg string) Response {
	return Response{Code: code, Message: msg, Data: nil}
}

func SuccessResponse(data interface{}) Response {
	return Response{Code: 0, Message: "success", Data: data}
}

func PagedSuccessResponse(data interface{}, total int64, page, pageSize int) PagedResponse {
	totalPages := int(total) / pageSize
	if int(total)%pageSize > 0 {
		totalPages++
	}
	return PagedResponse{
		Code:    0,
		Message: "success",
		Data:    data,
		Meta: &PageMeta{
			Total:      total,
			Page:       page,
			PageSize:   pageSize,
			TotalPages: totalPages,
		},
	}
}

func CountSuccessResponse(data interface{}, count int) PagedResponse {
	return PagedResponse{
		Code:    0,
		Message: "success",
		Data:    data,
		Meta:    CountMeta{Count: count},
	}
}
