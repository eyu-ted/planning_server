package domain

import (
	"time"

	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Quantify represents the quantifiable metrics for the plan's success
type Quantify struct {
	Target   string    `bson:"target" json:"target"`     // What is the target (e.g., sales target, project completion percentage)
	Actual   string    `bson:"actual" json:"actual"`     // The actual result (can be updated over time)
	Unit     string    `bson:"unit" json:"unit"`         // Unit of measurement (e.g., percentage, hours, count)
	Deadline time.Time `bson:"deadline" json:"deadline"` // Target deadline for the quantifiable metric
}

// Plan represents the structure of a plan in the system
type Plan struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty" json:"plan_id"`                           // MongoDB Object ID
	Title            string              `bson:"title" json:"title"`                                     // Title of the plan
	Description      string              `bson:"description" json:"description"`                         // Description of the plan
	Priority         string              `bson:"priority" json:"priority"`                               // Priority of the plan (e.g., High, Medium, Low)
	OwnerName        string              `bson:"owner_name" json:"owner_name"`                           // Owner of the plan (name of the person responsible)
	OwnerRole        string              `bson:"owner_role" json:"owner_role"`                           // Owner of the plan (name of the person responsible)
	SupervisorName   string              `bson:"supervisor_name" json:"supervisor_name"`                 // Supervisor's name (1 level higher in hierarchy)
	WhichQuarter     string              `bson:"which_quarter" json:"which_quarter"`                     // The quarter (e.g., Q1, Q2, etc.)
	Quantify         Quantify            `bson:"quantify" json:"quantify"`                               // Metrics and quantifiable targets
	CreatedBy        string              `bson:"created_by" json:"created_by"`                           // ID of the user who created the plan
	SupervisorPlanID *primitive.ObjectID `bson:"supervisor_plan_id,omitempty" json:"supervisor_plan_id"` // Parent plan (supervisor's plan)
	Status           string              `bson:"status" json:"status"`                                   // Status of the plan (e.g., Pending, Approved, Completed)
	CreatedAt        time.Time           `bson:"created_at" json:"created_at"`                           // Time when the plan was created
	UpdatedAt        time.Time           `bson:"updated_at" json:"updated_at"`
	OwnerID          primitive.ObjectID  `bson:"owner_id" json:"owner_id"`               // ID of the user who created the plan
	SuperiorPlan     string              `bson:"superior_plan" json:"superior_plan"`     // ID of the user who created the plan
	AlignedPillary   string              `bson:"aligned_pillary" json:"aligned_pillary"` // ID of the user who created the plan
	Quarter          int                 `bson:"quarter" json:"quarter"`                 // ID of the user who created the plan
	StartDate        time.Time           `bson:"start_date" json:"start_date"`           // ID of the user who created the plan
	EndDate          time.Time           `bson:"end_date" json:"end_date"`               // ID of the user who created the plan
	Type             string              `bson:"type" json:"type"`

	Comment string  `bson:"comment" json:"comment"`
	Value   float64 `bson:"value" json:"value"`
	// ID of the user who created the plan
}

type Report struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"report_id"`
	ReportUserID      primitive.ObjectID `bson:"report_user_id" json:"report_user_id"`
	ReportTitle       string             `bson:"report_title" json:"report_title"`     // Title of the plan
	AccomplishedValue string             `bson:"acomplished_value" json:"acomplished"` // Description of the plan
	ReportDetails     string             `bson:"report_details" json:"report_details"` // Priority of the plan (e.g., High, Medium, Low)
	Status            string             `bson:"status" json:"status"`                 // Owner of the plan (name of the person responsible)
	Type              string             `bson:"type" json:"type"`                     // Owner of the plan (name of the person responsible)

	SupervisorName  string `bson:"supervisor_name" json:"supervisor_name"` // Supervisor's name (1 level higher in hierarchy)
	Comment         string `bson:"comment" json:"comment"`      
	         // ID of the user who created the plan
// Supervisor's name (1 level higher in hierarchy)
	Value float64 `bson:"value" json:"value"`
	PlanID 		 primitive.ObjectID `bson:"plan_id" json:"plan_id"`
}

// Comment represents a comment on a plan.
type Comment struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"comment_id"` // MongoDB Object ID for the comment
	PlanID    primitive.ObjectID `bson:"plan_id" json:"plan_id"`          // The plan that this comment is related to
	Commenter string             `bson:"commenter" json:"commenter"`      // Name of the user commenting
	Content   string             `bson:"content" json:"content"`          // The actual comment text
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`    // Time when the comment was created
}

type Announcement struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title       string             `bson:"title" json:"title" binding:"required"`
	Description string             `bson:"description" json:"description" binding:"required"`
	CreatedTime time.Time          `bson:"created_time" json:"created_time"`
	Type        string             `bson:"type" json:"type"`
}

// PlanResponse represents the response returned when fetching a plan.

var (
	ErrPlanNotFound = errors.New("plan not found")
)
