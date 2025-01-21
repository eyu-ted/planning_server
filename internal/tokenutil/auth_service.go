package tokenutil
import (
    // "context"
    "errors"
    // "time"

    "github.com/golang-jwt/jwt/v4"
    "plan/domain"
    // "plan/repository"
    "fmt"
)

// type AuthService struct {
//     tokenRepo repository.MongoTokenRepository
//     jwtSecret string
// }

// func NewAuthService(tokenRepo repository.MongoTokenRepository, jwtSecret string) *AuthService {
//     return &AuthService{
//         tokenRepo: tokenRepo,
//         jwtSecret: jwtSecret,
//     }
// }

// func (s *AuthService) ValidateAccessToken(tokenString string) (*domain.Token, error) {
//     token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//         return []byte(s.jwtSecret), nil
//     })

//     if err != nil || !token.Valid {
//         return nil, errors.New("invalid access token")
//     }

//     claims, ok := token.Claims.(jwt.MapClaims)
//     if !ok || claims.Valid() != nil {
//         return nil, errors.New("invalid token claims")
//     }

//     return s.tokenRepo.FindTokenByAccessToken(context.Background(), tokenString)
// }

// func (s *AuthService) RefreshAccessToken(refreshToken string) (string, error) {
//     storedToken, err := s.tokenRepo.FindTokenByRefreshToken(context.Background(), refreshToken)
//     if err != nil || storedToken == nil {
//         return "", errors.New("invalid refresh token")
//     }

//     if time.Now().After(storedToken.ExpiresAt) {
//         _ = s.tokenRepo.DeleteToken(context.Background(), storedToken.ID)
//         return "", errors.New("refresh token expired")
//     }

//     newAccessToken, err := s.generateAccessToken(storedToken.UserID.Hex())
//     if err != nil {
//         return "", err
//     }

//     storedToken.CreatedAt = time.Now()
//     err = s.tokenRepo.SaveToken(context.Background(), storedToken)
//     if err != nil {
//         return "", err
//     }

//     return newAccessToken, nil
// }

// func (s *AuthService) generateAccessToken(userID string) (string, error) {
//     claims := jwt.MapClaims{
//         "id":        userID,
//         "exp":       time.Now().Add(time.Hour * 1).Unix(), // 1 hour expiry
//     }

//     token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//     return token.SignedString([]byte(s.jwtSecret))
// }
func VerifyToken(tokenString string) (*domain.JwtCustomClaims, error) {
	claims := &domain.JwtCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte("ts"), nil
	})
	fmt.Println("Parsed Token:", token)
	fmt.Println("Token Claims:", claims)
	if err != nil {
		fmt.Println("Token Verification Error:", err)
		return nil, err
	}
	if !token.Valid {
		fmt.Println("Token is invalid")
		return nil, errors.New("token is invalid")
	}

	return claims, nil
}
