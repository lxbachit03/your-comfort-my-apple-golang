package apiresponse

import (
	"net/http"

	"github.com/gin-gonic/gin"
	errorcode "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/error_code"
)

type ApiResponse[T any] struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Data       T      `json:"data"`
}

type ValidationErrorResponse struct {
	StatusCode   int                 `json:"status_code"`
	ErrorCode    errorcode.ErrorCode `json:"error_code"`
	Message      string              `json:"message"`
	ErrorDetails any                 `json:"error_details"`
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

type ApiError struct {
	StatusCode   int
	ErrorCode    errorcode.ErrorCode
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
			"status_code": http.StatusOK,
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
		"code":  errorcode.ErrInternalServer,
	})

	// return &ApiError{
	// 	ErrorDetails: err,
	// 	// ErrorCode:    errorCode,
	// }
}

func ValidationErrorResp(ctx *gin.Context, validation ValidationErrorResponse) {
	ctx.JSON(http.StatusBadRequest, validation)
}

func httpStatusFromCode(code errorcode.ErrorCode) int {
	switch code {
	case errorcode.ErrCodeBadRequest:
		return http.StatusBadRequest
	case errorcode.ErrCodeNotFound:
		return http.StatusNotFound
	case errorcode.ErrCodeUnauthorized:
		return http.StatusUnauthorized
	case errorcode.ErrCodeTooManyRequests:
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
