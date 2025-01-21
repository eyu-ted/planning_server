package domain

import (
	"context"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/golang-jwt/jwt/v4"
)

type JwtCustomClaims struct {
	UserID    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Full_Name string             `json:"full_name"`
	Email     string             `json:"email"`
	Username  string             `json:"username"`
	Role      string             `json:"role"`
	To_whom   string             `json:"to_whom"`
	Status    bool               `json:"status"`
	jwt.StandardClaims
}

type AuthSignup struct {
	UserID          primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email           string             `json:"email" binding:"required"`
	Password        string             `json:"password"`
	Role            string             `json:"role"`
	To_whom         string             `json:"to_whom"`
	Verify          bool               `json:"verify"`
	Profile_Picture string             `json:"profile_picture"`
	Full_Name       string             `json:"full_name"`
}

type AuthLogin struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserRepository interface {
	CreateUser(c context.Context, user *User) error
	GetUserByUsername(ctx context.Context, username string) (*User, error)
	// GetUserByID(ctx context.Context, userID string) (*User, error)
	FindUsersByRole(c context.Context, role string) ([]User, error)
	FindUnverifiedUsersByToWhom(ctx context.Context, firstName string) ([]User, error)
	UpdateVerifyStatus(ctx context.Context, userID primitive.ObjectID, verify bool) error
	FetchByToWhom(ctx context.Context, firstName string) ([]User, error)
	DeleteUser(ctx context.Context, userID primitive.ObjectID) error
	GetUserByID(ctx context.Context, userID primitive.ObjectID) (*User, error)
}

// Role is a type for user roles
type SignupUsecase interface {
	RegisterUser(c context.Context, user *AuthSignup) (*primitive.ObjectID, error)
	LoginUser(ctx context.Context, auth *AuthLogin) (string, error)
	// GetVerificationStatus(ctx context.Context, userID string) (bool, error)
	GetSuperiors(c context.Context, role string) ([]User, error)
	FetchUnverifiedUsersByToWhom(c context.Context, firstName string) ([]User, error)
	VerifyUser(c context.Context, userID string) error
	GetUsersByToWhomWithCount(ctx context.Context, firstName string) ([]User, int, error)
	RejectUser(c context.Context, userID string) error
}

type PlanRepository interface {
	CreatePlan(c context.Context, plan *Plan) error
	GetPlan(ctx context.Context, ownerID primitive.ObjectID) (*Plan, error)

	// UpdatePlan(ctx context.Context, plan *Plan) error
	GetSubmittedPlans(ctx context.Context, supervisor_name string) ([]Plan, error)
	GetPlanByID(ctx context.Context, planID primitive.ObjectID) (*Plan, error)
	// ForwardPlanToSupervisor(ctx context.Context, plan *Plan) error
	UpdatePlanPacth(ctx context.Context, plan *Plan) error
	CreateComment(ctx context.Context, comment *Comment) error
	FetchSupervisorComments(ctx context.Context, userID primitive.ObjectID) ([]Comment, error)
	FetchCommentsByPlanID(ctx context.Context, planID primitive.ObjectID) ([]Comment, error)
	GetPlanTitlesByOwnerName(ctx context.Context, ownerName string) ([]string, error)
	GetPlansByStatusAndOwner(ctx context.Context, userID primitive.ObjectID, status string) ([]Plan, error)
	SubmitReport(ctx context.Context, report *Report) error
	GetFilteredReports(ctx context.Context, userID primitive.ObjectID, status string) ([]Report, error)
	// GetAllTitlesByUser(ctx context.Context, userID primitive.ObjectID) ([]string, error)
	CountItems(ctx context.Context, itemType string, toWhom string) (int, error)
	GetPlansBySupervisorAndStatus(ctx context.Context, supervisorName, status string) ([]Plan, error)
	GetReportsBySupervisorAndStatus(ctx context.Context, supervisorName, reportStatus string) ([]Report, error)
	UpdatePlanStatus(ctx context.Context, planID primitive.ObjectID, supervisorName, status, comment string) error
	UpdateReportStatus(ctx context.Context, reportID primitive.ObjectID, supervisorName, status, comment string) error
	UpdatePlan(ctx context.Context, planID primitive.ObjectID, updatedPlan *Plan) error
	UpdateReport(ctx context.Context, reportID primitive.ObjectID, updatedReport *Report) error
	GetAllPlansByUser(ctx context.Context, userID primitive.ObjectID) ([]Plan, error)
	CreateAnnouncement(ctx context.Context, announcement *Announcement) error
	GetAllAnnouncements(ctx context.Context) ([]Announcement, error)
	Delete(ctx context.Context, id primitive.ObjectID) error
	FindByOwnerID(ctx context.Context, ownerID primitive.ObjectID, datatype string) ([]Plan, error)
	FindByUserID(ctx context.Context, userID primitive.ObjectID, datatype string) ([]Report, error)
}
type PlanUsecase interface {
	CreatePlan(c context.Context, plan *Plan) (*primitive.ObjectID, error)
	GetPlan(ctx context.Context, ownerID primitive.ObjectID) (*Plan, error)
	// EditPlan(ctx context.Context, plan *Plan) error
	GetSubmittedPlans(ctx context.Context, supervisor_name string) ([]Plan, error)
	// ApprovePlan(ctx context.Context, planID, ownerID primitive.ObjectID) error
	AddComment(ctx context.Context, comment *Comment) error
	GetSupervisorComments(ctx context.Context, userID primitive.ObjectID) ([]Comment, error)
	GetCommentsByPlanID(ctx context.Context, planID primitive.ObjectID) ([]Comment, error)
	GetPlanTitlesByOwnerName(ctx context.Context, ownerName string) ([]string, error)
	GetPlansByStatusAndOwner(ctx context.Context, userID primitive.ObjectID, status string) ([]Plan, error)
	SubmitReport(ctx context.Context, report *Report) error
	GetFilteredReports(ctx context.Context, userID primitive.ObjectID, status string) ([]Report, error)
	// GetAllTitlesByUser(ctx context.Context, userID primitive.ObjectID) ([]string, error)
	CountItems(ctx context.Context, itemType string, toWhom string) (int, error)
	FetchPlansBySupervisorAndStatus(c context.Context, supervisorName, status string) ([]Plan, error)
	FetchReportsBySupervisorAndStatus(c context.Context, supervisorName, reportStatus string) ([]Report, error)
	UpdatePlanStatus(c context.Context, planID primitive.ObjectID, supervisorName, status, comment string) error
	UpdateReportStatus(c context.Context, reportID primitive.ObjectID, supervisorName, status, comment string) error
	UpdatePlan(c context.Context, planID string, updatedPlan *Plan) error
	UpdateReport(c context.Context, reportID string, updatedReport *Report) error
	GetAllPlansByUser(ctx context.Context, userID primitive.ObjectID) ([]Plan, error)
	PublishAnnouncement(ctx context.Context, announcement *Announcement) error
	GetAllAnnouncements(ctx context.Context) ([]Announcement, error)
	DeleteAnnouncement(ctx context.Context, id primitive.ObjectID) error
	GetPlansByOwnerID(ctx context.Context, ownerID primitive.ObjectID, datatype string) ([]Plan, error)
	GetReportsByUserID(ctx context.Context, userID primitive.ObjectID, datatype string) ([]Report, error)
}
