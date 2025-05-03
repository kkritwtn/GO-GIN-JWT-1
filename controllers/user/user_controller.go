package user

import (
	"go-jwt-crud/config"
	"go-jwt-crud/dto"
	"go-jwt-crud/models"
	"go-jwt-crud/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"fmt"
)

func LoginUser(c *gin.Context) {
	var req dto.LoginUserRequest
	var user models.User
	//
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	fmt.Println("Request Body:", req.Password)

	// Fetch the user from the database
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error2": "Invalid credentials"})
		return
	}

	// Compare password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	fmt.Println("check ====>", user)
	// Generate JWT token after successful authentication
	token, _ := utils.GenerateJWT(user.Email, "user", user.ID)

	response := dto.ResponseLogin{
		Token: token,
	}
	c.JSON(http.StatusOK, response)

}

func GetUsers(c *gin.Context) {
	var users []models.User
	config.DB.Find(&users)
	c.JSON(http.StatusOK, users)
}

func RegisterUser(c *gin.Context) {
	var req dto.RegisterUserRequest
	fmt.Println("Request Body:", req)
	// Bind and validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)

	// Create user
	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
	})
}
func GetProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	var user struct {
		Name string `json:"name"`
		//Password string `json:"password"`
	}

	if err := config.DB.
		Model(&models.User{}).
		Select("name", "password").
		Where("id = ?", userID).
		First(&user).Error; err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
		return
	}

	var user models.User
	if err := config.DB.First(&user, userID).Error; err == gorm.ErrRecordNotFound {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Input struct including email
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update name
	if input.Name != "" {
		user.Name = input.Name
	}

	// Check for email uniqueness if trying to update email
	if input.Email != "" && input.Email != user.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ?", input.Email).First(&existingUser).Error; err == nil {
			// Found another user with the same email
			c.JSON(http.StatusBadRequest, gin.H{"error": "Email already in use"})
			return
		}
		user.Email = input.Email
	}

	// Update password if provided
	if input.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
		user.Password = string(hashedPassword)
	}

	if err := config.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	// Response
	type UserResponse struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user": UserResponse{
			Name:     user.Name,
			Email:    user.Email,
			Password: user.Password,
		},
	})
}

func DeleteAccount(c *gin.Context) {
	// userID, exists := c.Get("user_id")
	// fmt.Println("user", userID)

	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found in token"})
	// 	return
	// }
	// config.DB.Delete(&models.User{}, userID)
	// c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}
