package v1dto

type UserDTO struct {
	UUID  string `json:"uuid"`
	Name  string `json:"full_name"`
	Email string `json:"email_address"`
}

type GetUsersRequest struct {
	Search string `form:"search" binding:"omitempty,min=3,max=50"`
	Page   int32  `form:"page" binding:"omitempty,gte=1"`
	Limit  int32  `form:"limit" binding:"omitempty,gte=1,lte=500"`
	Order  string `form:"order_by" binding:"omitempty,oneof=user_id user_created_at"`
	Sort   string `form:"sort" binding:"omitempty,oneof=asc desc"`
}

type GetUserByUUIDRequest struct {
	UUID string `uri:"uuid" binding:"required,uuid"`
}

type AddAddressRequest struct {
	AddressLabel              string `json:"address_label" binding:"required"`
	AddressContactName        string `json:"address_contact_name" binding:"required"`
	AddressContactPhoneNumber string `json:"address_contact_phone_number" binding:"required"`
	AddressAddressLine        string `json:"address_address_line" binding:"required"`
	AddressDistrict           string `json:"address_district" binding:"required"`
	AddressProvince           string `json:"address_province" binding:"required"`
	AddressCountry            string `json:"address_country" binding:"required"`
}

type UpdateProfileRequest struct {
	FirstName   string `form:"first_name" binding:"omitempty"`
	LastName    string `form:"last_name" binding:"omitempty"`
	PhoneNumber string `form:"phone_number" binding:"omitempty"`
	BirthDay    string `form:"birth_day" binding:"omitempty"`
	Gender      string `form:"gender" binding:"omitempty"`
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
