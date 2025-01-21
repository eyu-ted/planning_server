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

func NewPlanRouter(env *config.Env, timeout time.Duration, db database.Database, group *gin.RouterGroup) {
	ur := repository.NewPlanRepository(db, "Plan")

	sc := controller.PlanController{
		PlanUsecase: usecase.NewPlanUsecase(ur, timeout),
		Env:         env,
	}
	group.POST("/summit/plan", sc.CreatePlan)
	group.GET("/plans/title", sc.GetPlanTitlesByOwnerName)
	group.GET("/filter", sc.GetPlansByStatusAndOwner)
	group.PUT("/update/plan/:plan_id", sc.UpdatePlan)

	group.POST("/report/submit", sc.SubmitReport)
	group.GET("plan/titles", sc.GetAllPlansByUser)
	group.GET("/report/filter", sc.GetFilteredReports)
	group.PUT("/update/report/:report_id", sc.UpdateReport)

	group.GET("/plan-and-report/count", sc.CountItems)

	group.GET("/plans", sc.GetPlansByStatus)
	group.GET("/reports", sc.GetReportsByStatus)

	group.POST("/plans/update-status", sc.UpdatePlanStatus)
	group.POST("/reports/update-status", sc.UpdateReportStatus)

	group.POST("/announcements", sc.PublishAnnouncement)
	group.GET("/announcements", sc.GetAllAnnouncements)
	group.DELETE("/announcements/:id",sc.DeleteAnnouncement)

	// group.GET("/plan", sc.GetPlan)
	// group.PUT("/plan", sc.EditPlan)
	// group.PUT("/plan/:planID", sc.EditPlan)

	// group.GET("/plans", sc.GetPlans)
	// group.GET("/plans/:planID", sc.GetPlanByID)

	// Register the route in your router
	// group.GET("/plans/submissions", sc.GetSubmittedPlans)
	// group.PATCH("/plans/:planID/approve", sc.ApprovePlan)

	// Routes for comment-related operations
	// group.POST("/plans/:planID/comments", sc.AddComment)
	// group.GET("/plans/:planID/comments", sc.GetCommentsByPlanID)

}
