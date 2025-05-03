package admin

// import (
// 	"go-jwt-crud/config"
// 	"go-jwt-crud/models"
// 	"go-jwt-crud/utils"
// 	"net/http"
// 	"fmt"
// 	"github.com/gin-gonic/gin"
// )

import (
	"fmt"
	"net/http"

	"go-jwt-crud/config"
	"go-jwt-crud/dto"
	"go-jwt-crud/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"

	"go-jwt-crud/utils"
)

func RegisterAdmin(c *gin.Context) {

	var req dto.RegisterAdminRequest

	// Bind and validate input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), 14)

	// Create user
	user := models.Admin{
		Username: req.Username,
		Password: string(hashedPassword),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return success
	c.JSON(http.StatusCreated, gin.H{
		"message": "Admin registered successfully",
	})
}

func LoginAdmin(c *gin.Context) {

	var req dto.LoginAdminRequest
	var user models.Admin
	//
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
		return
	}

	fmt.Println("Request Body:", req.Password)

	// Fetch the user from the database
	if err := config.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error2": "Invalid credentials"})
		return
	}

	// Compare the provided password with the hashed password in the database
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	// Generate JWT token after successful authentication
	token, _ := utils.GenerateJWT(user.Username, "admin", user.ID)

	response := dto.ResponseLogin{
		Token: token,
	}
	c.JSON(http.StatusOK, response)

}

func AdminDashboard(c *gin.Context) {

}
func GetAllUsers(c *gin.Context) {
	var users []models.User

	if err := config.DB.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	// Optional: Hide sensitive fields (like password) before returning
	var response []struct {
		ID    uint   `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}

	for _, user := range users {
		response = append(response, struct {
			ID    uint   `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		})
	}

	c.JSON(http.StatusOK, response)
}
