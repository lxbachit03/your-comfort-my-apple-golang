package apiresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse[T any] struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

type PaginationData[T any] struct {
	Items      []T        `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type PaginationApiResponse[T any] struct {
	StatusCode int               `json:"status_code"`
	Message    string            `json:"message"`
	Data       PaginationData[T] `json:"data"`
}

type Pagination struct {
	Page         int32 `json:"page"`
	Limit        int32 `json:"limit"`
	TotalRecords int32 `json:"total_records"`
	TotalPages   int32 `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

type ErrorCode string

const (
	ErrInternalServer      ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrCodeBadRequest      ErrorCode = "BAD_REQUEST"
	ErrCodeNotFound        ErrorCode = "NOT_FOUND"
	ErrCodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden       ErrorCode = "FORBIDDEN"
	ErrCodeTooManyRequests ErrorCode = "TOO_MANY_REQUESTS"
)

type ApiError struct {
	StatusCode   int
	ErrorCode    ErrorCode
	ErrorDetails error
}

func (ae *ApiError) Error() string {
	return ""
}

func Response(ctx *gin.Context, statusCode int, message string, data any) {
	ctx.JSON(statusCode, ApiResponse[any]{
		StatusCode: statusCode,
		Message:    message,
		Data:       data,
	})
}

func PaginationResponse[T any](ctx *gin.Context, statusCode int, message string, data PaginationApiResponse[T]) {
	data.StatusCode = statusCode
	data.Message = message
	ctx.JSON(statusCode, data)
}

func ErrorResponse(ctx *gin.Context, err error) {

	if apiError, ok := err.(*ApiError); ok {
		status := httpStatusFromCode(apiError.ErrorCode)

		response := gin.H{
			"status_code": 200,
			"error_code":  apiError.ErrorCode,
		}

		if apiError.ErrorDetails != nil {
			response["error_details"] = apiError.ErrorDetails.Error()
		}

		ctx.JSON(status, response)
		return
	}

	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
		"code":  ErrInternalServer,
	})

	// return &ApiError{
	// 	ErrorDetails: err,
	// 	// ErrorCode:    errorCode,
	// }
}

func httpStatusFromCode(code ErrorCode) int {
	switch code {
	case ErrCodeBadRequest:
		return http.StatusBadRequest
	case ErrCodeNotFound:
		return http.StatusNotFound
	case ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case ErrCodeTooManyRequests:
		return http.StatusTooManyRequests
	default:
		return http.StatusInternalServerError
	}
}

func NewPagination[T any](data []T, page int32, limit int32, totalRecords int32, totalPages int32) *PaginationApiResponse[T] {
	hasNext := page < totalPages
	hasPrev := page > 1

	return &PaginationApiResponse[T]{
		StatusCode: http.StatusOK,
		Message:    "success",
		Data: PaginationData[T]{
			Items: data,
			Pagination: Pagination{
				Page:         page,
				Limit:        limit,
				TotalRecords: totalRecords,
				TotalPages:   totalPages,
				HasNext:      hasNext,
				HasPrev:      hasPrev,
			},
		},
	}
}
