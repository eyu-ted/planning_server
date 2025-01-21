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

func NewSignupRouter(env *config.Env, timeout time.Duration, db database.Database, group *gin.RouterGroup) {
	ur := repository.NewUserRepository(db, "Staff")

	sc := controller.SignupController{
		SignupUsecase: usecase.NewSignupUsecase(ur, timeout),
		Env:           env,
	}
	group.POST("/signup", sc.Signup)
	group.POST("/login", sc.Login)

	// group.GET("/UNVERIFIED_USERS", sc.UNVERIFIED_USERS)
	// group.PATCH("/verify/:userID", sc.VerifyUser)

	group.GET("/hierarchy", sc.GetSuperiors)
	

}
