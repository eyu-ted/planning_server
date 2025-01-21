package controller

import (
	"plan/config"
	"plan/domain"

	"net/http"

	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

type PlanController struct {
	PlanUsecase domain.PlanUsecase
	Env         *config.Env
}

func (ac *PlanController) DeleteAnnouncement(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	// Convert ID to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// Call the usecase
	err = ac.PlanUsecase.DeleteAnnouncement(c, objectID)
	if err != nil {
		if err.Error() == "announcement not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Announcement not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Announcement deleted successfully"})
}

func (rc *PlanController) GetAllAnnouncements(c *gin.Context) {
	announcements, err := rc.PlanUsecase.GetAllAnnouncements(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "Announcements retrieved successfully",
		"announcements": announcements,
	})
}

func (rc *PlanController) PublishAnnouncement(c *gin.Context) {
	var announcement domain.Announcement

	if err := c.ShouldBindJSON(&announcement); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	announcement.ID = primitive.NewObjectID()
	announcement.CreatedTime = time.Now()
	announcement.Type = "announcement"

	err := rc.PlanUsecase.PublishAnnouncement(c, &announcement)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Announcement published successfully",
		"data":    announcement,
	})
}
func (rc *PlanController) UpdateReport(c *gin.Context) {
	reportID := c.Param("report_id")
	if reportID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "report_id is required"})
		return
	}

	var updatedReport domain.Report
	if err := c.ShouldBindJSON(&updatedReport); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the status is always set to "Pending"
	updatedReport.Status = "Pending"

	// Call the usecase to update the report
	err := rc.PlanUsecase.UpdateReport(c, reportID, &updatedReport)
	if err != nil {
		if err.Error() == "report not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report updated successfully"})
}

func (pc *PlanController) UpdatePlan(c *gin.Context) {
	planID := c.Param("plan_id")
	if planID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "plan_id is required"})
		return
	}

	var updatedPlan domain.Plan
	if err := c.ShouldBindJSON(&updatedPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the status is always set to "Pending"
	updatedPlan.Status = "Pending"

	// Call the usecase to update the plan
	err := pc.PlanUsecase.UpdatePlan(c, planID, &updatedPlan)
	if err != nil {
		if err.Error() == "plan not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan updated successfully"})
}

func (rc *PlanController) UpdateReportStatus(c *gin.Context) {
	var request struct {
		ReportID string `json:"report_id" binding:"required"`
		Status   string `json:"status" binding:"required"` // Should be "Approved" or "Rejected"
		Comment  string `json:"comment" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims) // Get user info from token
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	supervisorName := user.Full_Name // Extract supervisor's name from token

	reportID, err := primitive.ObjectIDFromHex(request.ReportID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Report ID"})
		return
	}

	err = rc.PlanUsecase.UpdateReportStatus(c, reportID, supervisorName, request.Status, request.Comment)
	if err != nil {
		if err.Error() == "report not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		} else if err.Error() == "unauthorized access" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report status updated successfully"})
}

func (pc *PlanController) UpdatePlanStatus(c *gin.Context) {
	var request struct {
		PlanID  string `json:"plan_id" binding:"required"`
		Status  string `json:"status" binding:"required"` // Should be "Approved" or "Rejected"
		Comment string `json:"comment" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims) // Get user info from token
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	supervisorName := user.Full_Name // Extract supervisor's name from token

	planID, err := primitive.ObjectIDFromHex(request.PlanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Plan ID"})
		return
	}

	err = pc.PlanUsecase.UpdatePlanStatus(c, planID, supervisorName, request.Status, request.Comment)
	if err != nil {
		if err.Error() == "plan not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
		} else if err.Error() == "unauthorized access" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Unauthorized access"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Plan status updated successfully"})
}

func (rc *PlanController) GetReportsByStatus(c *gin.Context) {
	reportStatus := c.Query("report_status") // Get the status from query parameters

	if reportStatus == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Report status is required"})
		return
	}

	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims) // Extract user from the token
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	supervisorName := user.Full_Name // Extract supervisor name from the token

	// Call usecase with both report_status and supervisor name
	reports, err := rc.PlanUsecase.FetchReportsBySupervisorAndStatus(c, supervisorName, reportStatus)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}

func (pc *PlanController) GetPlansByStatus(c *gin.Context) {
	status := c.Query("status") // Get the status from query parameters

	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status is required"})
		return
	}

	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims) // Extract user from the token
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	supervisorName := user.Full_Name // Supervisor name from the token

	// Call usecase with both status and supervisor name
	plans, err := pc.PlanUsecase.FetchPlansBySupervisorAndStatus(c, supervisorName, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func (pc *PlanController) CreatePlan(c *gin.Context) {
	var plan domain.Plan
	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims) // Get the user from the JWT token

	if !ok {
		// Handle error if the type assertion fails
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := c.ShouldBindJSON(&plan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Fill the plan information based on the user role and the JWT token
	plan.OwnerRole = user.Role
	plan.OwnerID = user.UserID
	plan.OwnerName = user.Full_Name
	plan.SupervisorName = user.To_whom
	plan.CreatedBy = user.Username
	// plan.SupervisorPlanID = &user.UserID // Assuming To_whom is the supervisor's ID

	planID, err := pc.PlanUsecase.CreatePlan(c, &plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plan_id": planID})
}
func (pc *PlanController) GetPlansByStatusAndOwner(c *gin.Context) {
	// Get the status from query params
	status := c.Query("status")
	if status == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status query parameter is required"})
		return
	}

	// Get the user from JWT claims
	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Call the usecase
	plans, err := pc.PlanUsecase.GetPlansByStatusAndOwner(c, user.UserID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func (pc *PlanController) GetPlanTitlesByOwnerName(c *gin.Context) {
	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims) // Get the user from the JWT token
	if !ok {
		// Handle error if the type assertion fails
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ownerName := user.To_whom // Get the "To_whom" field from the JWT token

	titles, err := pc.PlanUsecase.GetPlanTitlesByOwnerName(c, ownerName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"titles": titles})
}

func (rc *PlanController) SubmitReport(c *gin.Context) {
	var report domain.Report

	// Bind the JSON request body to the report struct
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Extract user ID from JWT claims
	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	report.SupervisorName = user.To_whom
	// Attach the user ID to the report
	report.ReportUserID = user.UserID
	report.Type = "report"



	// Call the usecase to submit the report
	err := rc.PlanUsecase.SubmitReport(c, &report)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Report submitted successfully"})
}
func (rc *PlanController) GetFilteredReports(c *gin.Context) {
	status := c.Query("status")

	// Validate status parameter
	if status != "Pending" && status != "Approved" && status != "Rejected" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid status query parameter"})
		return
	}

	// Extract user ID from JWT claims
	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch reports from usecase
	reports, err := rc.PlanUsecase.GetFilteredReports(c, user.UserID, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"reports": reports})
}
func (pc *PlanController) GetAllPlansByUser(c *gin.Context) {
	// Extract user ID from JWT claims
	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Fetch all plans submitted by the user from the usecase
	plans, err := pc.PlanUsecase.GetAllPlansByUser(c, user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plans": plans})
}

func (c *PlanController) CountItems(ctx *gin.Context) {
	itemType := ctx.Query("type") // Retrieve the query parameter

	if itemType != "plan" && itemType != "report" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid type. Use 'plan' or 'report'."})
		return
	}

	// Extract "to_whom" from claims
	claims, ok := ctx.MustGet("claim").(*domain.JwtCustomClaims)
	if !ok {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid token claims"})
		return
	}
	toWhom := claims.Full_Name // Assuming `FirstName` is part of the claims
	
	// Call the use case
	count, err := c.PlanUsecase.CountItems(ctx.Request.Context(), itemType, toWhom)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count items"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"type":  itemType,
		"count": count,
	})
}

// GetSubordinateUsers fetches users whose To_whom matches the current user's first name.

func (pc *PlanController) GetPlan(c *gin.Context) {
	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims) // Extract user claims from the JWT
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ownerID := user.UserID // Extract OwnerID from JWT claims

	// Fetch the plan using the OwnerID
	plan, err := pc.PlanUsecase.GetPlan(c, ownerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"plan": plan})
}

// func (pc *PlanController) EditPlan(c *gin.Context) {
// 	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims) // Extract user claims from the JWT
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	ownerID := user.UserID // Extract OwnerID from JWT claims

// 	var updatedPlan domain.Plan
// 	if err := c.ShouldBindJSON(&updatedPlan); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 		return
// 	}

// 	updatedPlan.OwnerID = ownerID // Ensure the OwnerID matches the authenticated user

// 	err := pc.PlanUsecase.EditPlan(c, &updatedPlan)
// 	if err != nil {
// 		if err == domain.ErrPlanNotFound {
// 			c.JSON(http.StatusNotFound, gin.H{"error": "Plan not found"})
// 		} else {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		}
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Plan updated successfully"})
// }

func (pc *PlanController) GetSubmittedPlans(c *gin.Context) {
	// Extract user claims from the context
	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	supervisor_name := user.To_whom
	plans, err := pc.PlanUsecase.GetSubmittedPlans(c, supervisor_name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, plans)
}

// func (pc *PlanController) ApprovePlan(c *gin.Context) {
// 	// Extract user claims
// 	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims)
// 	if !ok {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	// Get planID from route parameters
// 	planIDHex := c.Param("planID")
// 	planID, err := primitive.ObjectIDFromHex(planIDHex)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
// 		return
// 	}

// 	// Call the usecase to approve the plan
// 	err = pc.PlanUsecase.ApprovePlan(c, planID, user.UserID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Plan approved and forwarded to supervisor if applicable"})
// }

func (cc *PlanController) AddComment(c *gin.Context) {
	// Extract plan ID
	planIDHex := c.Param("planID")
	planID, err := primitive.ObjectIDFromHex(planIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	// Parse the request body for comment content
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Extract user claims
	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Call the usecase to add the comment
	comment := domain.Comment{
		PlanID:    planID,
		Commenter: user.Username,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}
	err = cc.PlanUsecase.AddComment(c, &comment)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Comment added successfully"})
}

func (cc *PlanController) GetSupervisorComments(c *gin.Context) {
	// Extract user claims
	user, ok := c.MustGet("claim").(*domain.JwtCustomClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	// Call the usecase to fetch supervisor comments
	comments, err := cc.PlanUsecase.GetSupervisorComments(c, user.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func (cc *PlanController) GetCommentsByPlanID(c *gin.Context) {
	// Extract the plan ID from the request
	planIDHex := c.Param("planID")
	planID, err := primitive.ObjectIDFromHex(planIDHex)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plan ID"})
		return
	}

	// Call the usecase to fetch comments by PlanID
	comments, err := cc.PlanUsecase.GetCommentsByPlanID(c, planID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}
