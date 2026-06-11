package v1handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	apiresponse "github.com/lxbachit03/ygz-microservices-golang/internal/pkg/api-response"
	v1dto "github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/endpoints/dtos/v1"
	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/infrastructure/utils"
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

	if err == nil && profileImage != nil {
		// 1. Validate file
		allowedExts := map[string]bool{
			".jpg":  true,
			".jpeg": true,
			".png":  true,
		}
		allowedMimeTypes := map[string]bool{
			"image/jpeg": true,
			"image/png":  true,
		}

		if err := utils.ValidateUploadedFile(profileImage, allowedExts, 5<<20, allowedMimeTypes); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 2. Save file with a random UUID name
		savedPath, err := utils.SaveUploadedFileWithRandomName(profileImage, "uploads/profile-images")
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		profileImageUrl = savedPath
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

func (arh *UserRouteHandler) UploadMultipleFiles(ctx *gin.Context) {
	form, err := ctx.MultipartForm()
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid multipart form"})
		return
	}

	images := form.File["images"]

	if len(images) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}

	// 1. Limit max files to 5 to prevent DoS
	const maxFiles = 5
	if len(images) > maxFiles {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Too many files, maximum is 5"})
		return
	}

	var successFiles []string
	var failedFiles []map[string]string

	allowedExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
	}
	allowedMimeTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
	}

	for _, image := range images {

		// 1. validate
		if err := utils.ValidateUploadedFile(image, allowedExts, 5<<20, allowedMimeTypes); err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"filename": image.Filename,
				"error":    err.Error(),
			})

			continue
		}

		// 2. save file
		savedPath, err := utils.SaveUploadedFileWithRandomName(image, "uploads/profile-images")
		if err != nil {
			failedFiles = append(failedFiles, map[string]string{
				"filename": image.Filename,
				"error":    err.Error(),
			})
			continue
		}

		successFiles = append(successFiles, savedPath)
	}

	// 2. If at least 1 file upload failed, return HTTP 400 Bad Request
	if len(failedFiles) > 0 {
		apiresponse.Response(ctx, http.StatusBadRequest, "Some or all file uploads failed", map[string]any{
			"success": successFiles,
			"failed":  failedFiles,
		})
		return
	}

	apiresponse.Response(ctx, http.StatusOK, "Files uploaded successfully", map[string]any{
		"success": successFiles,
		"failed":  failedFiles,
	})
}
