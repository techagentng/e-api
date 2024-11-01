package server

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/techagentng/ecommerce-api/models"
	"github.com/techagentng/ecommerce-api/server/response"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func createS3Client() (*s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(os.Getenv("AWS_REGION")))
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}
	return s3.NewFromConfig(cfg), nil
}

// A map to hold content types based on file extensions
var contentTypes = map[string]string{
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".png":  "image/png",
	".mp4":  "video/mp4",
	".avi":  "video/x-msvideo",
}

func uploadFileToS3(client *s3.Client, file multipart.File, bucketName, key string) (string, error) {
	defer file.Close()

	// Read the file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read file content: %v", err)
	}

	// Determine file extension and content type
	extension := filepath.Ext(key)
	contentType, exists := contentTypes[extension]
	if !exists {
		contentType = "application/octet-stream" // Default content type
	}

	// Remove whitespace from the file content
	trimmedContent := strings.TrimSpace(string(fileContent))
	fileContent = []byte(trimmedContent)
	key = strings.ReplaceAll(key, " ", "_")
	// Upload the file to S3
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        bytes.NewReader(fileContent),
		ContentType: aws.String(contentType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %v", err)
	}

	// Return the S3 URL of the uploaded file
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, os.Getenv("AWS_REGION"), key)
	return fileURL, nil
}

func (s *Server) handleSignup() gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := c.Request.ParseMultipartForm(10 << 20); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}

		var filePath string // This will hold the S3 URL

		// Get the profile image from the form
		file, handler, err := c.Request.FormFile("profile_image")
		if err == nil {
			defer file.Close()

			// Create S3 client
			s3Client, err := createS3Client()
			if err != nil {
				response.JSON(c, "", http.StatusInternalServerError, nil, err)
				return
			}

			// Generate unique filename
			userID := c.PostForm("user_id")
			filename := fmt.Sprintf("%s_%s", userID, handler.Filename)

			// Upload file to S3
			filePath, err = uploadFileToS3(s3Client, file, os.Getenv("AWS_BUCKET"), filename)
			if err != nil {
				response.JSON(c, "", http.StatusInternalServerError, nil, err)
				return
			}
		} else if err == http.ErrMissingFile {
			filePath = "uploads/default-profile.png" // Adjust this to a default S3 URL if necessary
		} else {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}

		var user models.User
		user.Fullname = c.PostForm("fullname")
		user.Username = c.PostForm("username")
		user.Telephone = c.PostForm("telephone")
		user.Email = c.PostForm("email")
		user.Password = c.PostForm("password")
		user.ThumbNailURL = filePath

		// Fetch the UUID for the role
		role, err := s.AuthService.GetRoleByName("User")
		if err != nil {
			response.JSON(c, "", http.StatusInternalServerError, nil, err)
			return
		}
		log.Printf("Fetched role ID for 'User': %s", role.ID.String())

		// Assign the role UUID directly to RoleID
		user.RoleID = role.ID

		// Validate the user data using the validator package
		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			response.JSON(c, "", http.StatusBadRequest, nil, err)
			return
		}

		userResponse, err := s.AuthService.SignupUser(&user)
		if err != nil {
			response.HandleErrors(c, err)
			return
		}

		response.JSON(c, "Signup successful, check your email for verification", http.StatusCreated, userResponse, nil)
	}
}
