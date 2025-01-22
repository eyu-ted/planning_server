package controller

import (
	"fmt"
	"net/http"
	"plan/config"
	"plan/domain"

	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/gomail.v2"
)

type SignupController struct {
	SignupUsecase domain.SignupUsecase
	Env           *config.Env
}

// func (sc *SignupController) TokenInfo(c *gin.Contex){

// }

func (uc *SignupController) GetUserInfo(c *gin.Context) {
	var request struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	objectID, err := primitive.ObjectIDFromHex(request.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user_id format"})
		return
	}

	user, err := uc.SignupUsecase.FetchUserByID(c, objectID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
func (sc *SignupController) Signup(c *gin.Context) {
	var user domain.AuthSignup

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := sc.SignupUsecase.RegisterUser(c, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := SendApprovalEmailcontroller(user.Email, "placeholder"); err != nil {
		log.Println("failed to send approval email")
	}
	c.JSON(http.StatusOK, gin.H{"userID": userID})
}

func (sc *SignupController) Login(c *gin.Context) {
	var user domain.AuthLogin

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	token, err := sc.SignupUsecase.LoginUser(c, &user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}
func SendApprovalEmailcontroller(to string, firstName string) error {
	// SMTP server configuration
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		return err
	}

	// Get SMTP credentials from environment variables
	smtpUsername := os.Getenv("SMTPUsername")
	smtpPassword := os.Getenv("SMTPPassword")
	smtpHost := os.Getenv("SMTPHost")
	smtpPortStr := os.Getenv("SMTPPort") // Replace with your email password or app password
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	// Email content
	subject := "AASTU Planning System Account pending"
	body := fmt.Sprintf("Hello %s,\n\nYour account is wating for pending.\n\nThank you for you patience!", firstName)

	// Create email message
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", smtpUsername)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", body)

	// Create email dialer
	dialer := gomail.NewDialer(smtpHost, smtpPort, smtpUsername, smtpPassword)

	// Send the email
	return dialer.DialAndSend(mailer)
}

func (uc *SignupController) GetSubordinateUsers(c *gin.Context) {
	// Get claims from context
	claims, ok := c.MustGet("claim").(*domain.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
		return
	}

	// Extract first_name from claims
	firstName := claims.Full_Name

	// Call use case to fetch users and their count
	users, count, err := uc.SignupUsecase.GetUsersByToWhomWithCount(c.Request.Context(), firstName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"count": count,
		"users": users,
	})
}

func (uc *SignupController) GetUnverifiedUsersByToWhom(c *gin.Context) {
	claims := c.MustGet("claim").(*domain.JwtCustomClaims)
	firstName := claims.Full_Name
	fmt.Printf("First name: %s\n", claims)

	if firstName == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized or missing first name"})
		return
	}

	users, err := uc.SignupUsecase.FetchUnverifiedUsersByToWhom(c, firstName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc *SignupController) VerifyUser(c *gin.Context) {
	var request struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := uc.SignupUsecase.VerifyUser(c, request.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User verification status updated successfully"})
}

func (sc *SignupController) RejectUser(c *gin.Context) {
	var request struct {
		UserID string `json:"user_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := sc.SignupUsecase.RejectUser(c, request.UserID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		} else if err.Error() == "invalid user ID format" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User rejected and deleted successfully"})
}

// func (uc *SignupController) VerifyStatus(c *gin.Context) {
// 	// Extract user ID from JWT token (assumes middleware sets user ID in context)

// 	userID, exists := c.Get("userID")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	// Call the use case to fetch the user's verification status
// 	verify, err := uc.SignupUsecase.GetVerificationStatus(c, userID.(string))
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"verify": verify})
// }

func (sc *SignupController) GetSuperiors(c *gin.Context) {
	role := c.Query("role")
	if role == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Role is required"})
		return
	}

	superiors, err := sc.SignupUsecase.GetSuperiors(c, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"superiors": superiors})
}
