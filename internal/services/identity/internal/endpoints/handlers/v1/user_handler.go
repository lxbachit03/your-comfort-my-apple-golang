package v1handler

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apiresponse "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/api-response"
	v1dto "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/dtos/v1"
	identity_validator "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils/validation"
	usecase "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases"
	command "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/commands"
	query "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/usecases/users/queries"
)

type UserRouteHandler struct {
	uc *usecase.Usecase
}

func NewUserRouteHandler(uc *usecase.Usecase) *UserRouteHandler {
	return &UserRouteHandler{
		uc: uc,
	}
}

func (arh *UserRouteHandler) GetUsers(ctx *gin.Context) {
	var request v1dto.GetUsersParams
	if err := ctx.ShouldBindQuery(&request); err != nil {
		apiresponse.ValidationErrorResp(ctx, identity_validator.HandleValidationError(err))
		return
	}

	query := query.GetUsersQuery{
		Search: request.Search,
		Page:   request.Page,
		Limit:  request.Limit,
		Order:  request.Order,
		Sort:   request.Sort,
	}

	result, err := arh.uc.Queries.GetUsersHandler.Handle(ctx.Request.Context(), query)
	if err != nil {
		apiresponse.ErrorResponse(ctx, err)
		return
	}

	apiresponse.PaginationResponse(ctx, http.StatusOK, "Test GetUsers", result)
}

func (arh *UserRouteHandler) GetUserByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "uuid is required"})
		return
	}

	query := query.GetUserByUUIDQuery{
		UUID: uuid,
	}

	result, err := arh.uc.Queries.GetUserByUUIDHandler.Handle(ctx.Request.Context(), query)
	if err != nil {
		apiresponse.ErrorResponse(ctx, err)
		return
	}

	apiresponse.Response(ctx, http.StatusOK, "User retrieved successfully", result)
}

func (arh *UserRouteHandler) AddAddress(ctx *gin.Context) {
	var request v1dto.AddAddressRequest

	// Test by using x-www-form-urlencoded
	if err := ctx.ShouldBind(&request); err != nil {
		apiresponse.ValidationErrorResp(ctx, identity_validator.HandleValidationError(err))
		return
	}

	cmd := command.AddAddressCommand{
		Label:              request.AddressLabel,
		ContactName:        request.AddressContactName,
		ContactPhoneNumber: request.AddressContactPhoneNumber,
		AddressLine:        request.AddressAddressLine,
		District:           request.AddressDistrict,
		Province:           request.AddressProvince,
		Country:            request.AddressCountry,
	}

	result, err := arh.uc.Commands.AddAddressHandler.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		apiresponse.ErrorResponse(ctx, err)
		return
	}

	apiresponse.Response(ctx, http.StatusOK, "Address added successfully", result)
}

func (arh *UserRouteHandler) UpdateProfile(ctx *gin.Context) {
	var request v1dto.UpdateProfileRequest
	var profileImageUrl string
	if err := ctx.ShouldBind(&request); err != nil {
		apiresponse.ValidationErrorResp(ctx, identity_validator.HandleValidationError(err))
		return
	}

	profileImage, err := ctx.FormFile("profile_image")

	if profileImage != nil {
		// validate file name
		var allowExts = map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}

		log.Printf("profileImage.Filename %s", profileImage.Filename)

		ext := strings.ToLower(filepath.Ext(profileImage.Filename))
		if !allowExts[ext] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file extension"})
			return
		}

		// check file size (5MB)
		// 1 << 20 = 1 * 2^20 = 1 * 1048576 = 1MB
		// 5 << 20 = 5 * 2^20 = 5 * 1048576 = 5MB
		if profileImage.Size > 5<<20 {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "File too large (5 MB)"})
			return
		}

		// check file type
		var allowMimeTypes = map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
		}
		file, err := profileImage.Open()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot open file"})
			return
		}
		defer file.Close()

		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "cannot read file"})
			return
		}

		mimeType := http.DetectContentType(buffer)
		if !allowMimeTypes[mimeType] {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid MIME type"})
			return
		}

		// Change filename (UUID).jpg
		filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)

		// save file
		saveDir := "uploads/profile-images/"
		// dest := path.Join(utils.GetWorkingDir(), saveDir, filename)
		// if err := os.MkdirAll(saveDir, 0755); err != nil {
		// 	ctx.JSON(http.StatusBadRequest, gin.H{"error": "failed to create directory"})
		// 	return
		// }

		profileImageUrl = saveDir + filename
	}

	cmd := command.UpdateProfileCommand{
		ProfileId:       ctx.Param("profile_id"),
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		PhoneNumber:     request.PhoneNumber,
		BirthDay:        request.BirthDay,
		Gender:          request.Gender,
		ProfileImageUrl: profileImageUrl,
	}

	result, err := arh.uc.Commands.UpdateProfileHandler.Handle(ctx.Request.Context(), cmd)
	if err != nil {
		apiresponse.ErrorResponse(ctx, err)
		return
	}

	var messsage string

	if result {
		messsage = "Profile updated successfully"
	} else {
		messsage = "Profile update failed"
	}

	apiresponse.Response(ctx, http.StatusOK, messsage, result)
}
