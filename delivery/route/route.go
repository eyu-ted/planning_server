package route

import (
	"plan/config"
	"plan/database"

	"plan/delivery/middleware"
	"time"

	"github.com/gin-gonic/gin"
	// "github.com/google/generative-ai-go/genai"
)

func Setup(env *config.Env, timeout time.Duration, db database.Database, gin *gin.Engine) {
	publicRouter := gin.Group("")
	NewSignupRouter(env, timeout, db, publicRouter)

	protectedRouter := gin.Group("")
	protectedRouter.Use(middleware.AuthMidd)

	
	NewProtectedRouter(env, timeout, db, protectedRouter)

	NewPlanRouter(env, timeout, db, protectedRouter)

}
