package usecase

import (
	"fmt"
	"log"
	"os"
	"plan/domain"
	"strconv"

	// "github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"

	// "plan/internal/tokenutil"
	"context"
	"errors"
	"plan/internal/userutil"

	"gopkg.in/gomail.v2"

	// "net/smtp"
	"time"

	// "github.com/dgrijalva/jwt-go"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type signupUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

func NewSignupUsecase(userRepository domain.UserRepository, timeout time.Duration) domain.SignupUsecase {
	return &signupUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}

func (su *signupUsecase) RegisterUser(c context.Context, user *domain.AuthSignup) (*primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	hashedPassword, err := userutil.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}

	adduser := &domain.User{
		ID:         primitive.NewObjectID(),
		First_Name: user.First_Name,
		Last_Name:  user.Last_Name,
		Username:   user.Username,
		Email:      user.Email,
		Password:   hashedPassword,
		Role:       user.Role,
		To_whom:    user.To_whom,
		Verify:     false,
	}
	err = su.userRepository.CreateUser(ctx, adduser)
	return &adduser.ID, err
}

func (su *signupUsecase) LoginUser(ctx context.Context, auth *domain.AuthLogin) (string, error) {
	// Fetch user from the repository
	user, err := su.userRepository.GetUserByUsername(ctx, auth.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Verify the password
	// if !userutil.ComparePassword(auth.Password, user.Password) {
	//     return "", errors.New("invalid credentials")
	// }
	err = userutil.ComparePassword(user.Password, auth.Password)
	if err != nil {
		fmt.Println("password not match")
		return "", errors.New("invalid credentials")
	}

	// Check if the user is verified
	if !user.Verify {
		return "", errors.New("your account is pending verification")
	}

	// Generate JWT token
	token, err := GenerateJWTToken(user)
	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return token, nil
}

func (uc *signupUsecase) FetchUnverifiedUsersByToWhom(c context.Context, firstName string) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	return uc.userRepository.FindUnverifiedUsersByToWhom(ctx, firstName)
}

// func (uu *signupUsecase) GetVerificationStatus(ctx context.Context, userID string) (bool, error) {

// 	user, err := uu.userRepository.GetUserByID(ctx, userID)
// 	if err != nil {
// 		return false, errors.New("user not found")
// 	}

//		return user.Verify, nil
//	}
func (uc *signupUsecase) GetUsersByToWhomWithCount(ctx context.Context, firstName string) ([]domain.User, int, error) {
	// Delegate to repository
	users, err := uc.userRepository.FetchByToWhom(ctx, firstName)
	if err != nil {
		return nil, 0, err
	}

	// Return users and their count
	return users, len(users), nil
}

func (uc *signupUsecase) RejectUser(c context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	// Retrieve user details for logging or notifications (optional)
	user, err := uc.userRepository.GetUserByID(ctx, objectID)
	if err != nil {
		return errors.New("user not found")
	}

	// Delete user
	err = uc.userRepository.DeleteUser(ctx, objectID)
	if err != nil {
		return err
	}

	// Optional: Send rejection email or log the action
	if err := SendRejectionEmail(user.Email, user.First_Name); err != nil {
		return errors.New("failed to send rejection email")
	}

	return nil
}

func GenerateJWTToken(user *domain.User) (string, error) {
	claims := &domain.JwtCustomClaims{
		First_Name: user.First_Name,
		UserID:     user.ID,
		Email:      user.Email,
		Username:   user.Username,
		Role:       user.Role,
		To_whom:    user.To_whom,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(72)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte("ts"))
}

func (su *signupUsecase) GetSuperiors(c context.Context, role string) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()

	// Define the hierarchy map
	roleHierarchy := map[string]string{
		"vice_president":   "strategy_and_planning_manager",
		"regular_employee": "team_lead",
		"team_lead":        "director",
		"director":         "vice_president",
	}

	superiorRole, exists := roleHierarchy[role]
	if !exists {
		return nil, fmt.Errorf("invalid role: %s", role)
	}

	return su.userRepository.FindUsersByRole(ctx, superiorRole)
}
func (uc *signupUsecase) VerifyUser(c context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(c, uc.contextTimeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return errors.New("invalid user ID format")
	}

	// Retrieve user details for email
	user, err := uc.userRepository.GetUserByID(ctx, objectID)
	if err != nil {
		return errors.New("failed to retrieve user details")
	}

	// Update verification status
	if err := uc.userRepository.UpdateVerifyStatus(ctx, objectID, true); err != nil {
		return err
	}

	// Send approval email
	if err := SendApprovalEmail(user.Email, user.First_Name); err != nil {
		return errors.New("failed to send approval email")
	}

	return nil
}

// SendApprovalEmail sends an email to a user notifying them of their approval.
func SendApprovalEmail(to string, firstName string) error {
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
	subject := "AASTU Planning System Account Approved"
	body := fmt.Sprintf("Hello %s,\n\nYour account has been approved. You can now log in and start using our services.\n\nThank you!", firstName)

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

func SendRejectionEmail(to string, firstName string) error {
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
	smtpPortStr := os.Getenv("SMTPPort")
	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	// Email content
	subject := "AASTU Planning System Account Rejected"
	body := fmt.Sprintf("Hello %s,\n\nWe regret to inform you that your account has been rejected. For further details, please contact support.\n\nThank you!", firstName)

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
