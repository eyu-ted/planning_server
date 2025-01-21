package route

import (
	"plan/config"
	"plan/database"
	"plan/delivery/controller"

	// "plan/domain"
	"plan/repository"
	"plan/usecase"
	"time"

	"github.com/gin-gonic/gin"
)

// Setup sets up the routes for the application

func NewProtectedRouter(env *config.Env, timeout time.Duration, db database.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, "Staff")

	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, timeout),
		Env:           env,
	}

	// group.GET("/verify-status", sc.VerifyStatus)
	group.GET("/users/unverified", sc.GetUnverifiedUsersByToWhom)
	group.POST("/verify", sc.VerifyUser)
	group.DELETE("/reject", sc.RejectUser)
	group.GET("/users/subordinates", sc.GetSubordinateUsers)



}
