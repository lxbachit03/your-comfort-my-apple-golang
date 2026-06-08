package v1dto

type GetUsersParams struct {
	Search string `form:"search" binding:"omitempty,min=3,max=50"`
	Page   int32  `form:"page" binding:"omitempty,gte=1"`
	Limit  int32  `form:"limit" binding:"omitempty,gte=1,lte=500"`
	Order  string `form:"order_by" binding:"omitempty,oneof=user_id user_created_at"`
	Sort   string `form:"sort" binding:"omitempty,oneof=asc desc"`
}

type UserDTO struct {
	UUID  string `json:"uuid"`
	Name  string `json:"full_name"`
	Email string `json:"email_address"`
}

type Pagination[T any] struct {
	Page         int32 `json:"page"`
	Limit        int32 `json:"limit"`
	TotalRecords int32 `json:"total_records"`
	TotalPages   int32 `json:"total_pages"`
	HasNext      bool  `json:"has_next"`
	HasPrev      bool  `json:"has_prev"`
}

func MapUserToDTO(user any) *UserDTO {
	return &UserDTO{
		UUID:  "test",
		Name:  "test",
		Email: "test",
	}
}

func MapUsersToDTO(users []UserDTO) []UserDTO {
	dtos := make([]UserDTO, 0, len(users))

	for _, users := range users {
		dtos = append(dtos, *MapUserToDTO(users))
	}

	return dtos
}
