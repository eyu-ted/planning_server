package usecase

import (
	"fmt"
	"plan/domain"

	// "plan/internal/tokenutil"
	"context"
	"errors"

	// "plan/internal/userutil"

	// "net/smtp"
	"time"
	// "github.com/dgrijalva/jwt-go"
	// jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type planUsecaseStruct struct {
	planRepository domain.PlanRepository
	contextTimeout time.Duration
}

func NewPlanUsecase(planRepositoryPAR domain.PlanRepository, timeout time.Duration) domain.PlanUsecase {
	return &planUsecaseStruct{
		planRepository: planRepositoryPAR,
		contextTimeout: timeout,
	}
}
func (uc *planUsecaseStruct) GetPlansByOwnerID(ctx context.Context, ownerID primitive.ObjectID,datatype string) ([]domain.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	return uc.planRepository.FindByOwnerID(ctx, ownerID,datatype)
}

func (uc *planUsecaseStruct) GetReportsByUserID(ctx context.Context, userID primitive.ObjectID,datatype string) ([]domain.Report, error) {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()
	return uc.planRepository.FindByUserID(ctx, userID,datatype)
}
func (uc *planUsecaseStruct) DeleteAnnouncement(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, uc.contextTimeout)
	defer cancel()

	// Call the repository to delete
	err := uc.planRepository.Delete(ctx, id)
	if err != nil {
		if err.Error() == "announcement not found" {
			return errors.New("announcement not found")
		}
		return err
	}

	return nil
}
func (ru *planUsecaseStruct) GetAllAnnouncements(ctx context.Context) ([]domain.Announcement, error) {
	ctx, cancel := context.WithTimeout(ctx, ru.contextTimeout)
	defer cancel()

	return ru.planRepository.GetAllAnnouncements(ctx)
}

func (ru *planUsecaseStruct) PublishAnnouncement(ctx context.Context, announcement *domain.Announcement) error {
	ctx, cancel := context.WithTimeout(ctx, ru.contextTimeout)
	defer cancel()

	if announcement.Title == "" || announcement.Description == "" {
		return errors.New("title and description are required")
	}

	announcement.CreatedTime = time.Now()
	return ru.planRepository.CreateAnnouncement(ctx, announcement)
}
func (ru *planUsecaseStruct) UpdateReport(c context.Context, reportID string, updatedReport *domain.Report) error {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(reportID)
	if err != nil {
		return errors.New("invalid report ID format")
	}

	return ru.planRepository.UpdateReport(ctx, objectID, updatedReport)
}

func (pu *planUsecaseStruct) UpdatePlan(c context.Context, planID string, updatedPlan *domain.Plan) error {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	objectID, err := primitive.ObjectIDFromHex(planID)
	if err != nil {
		return errors.New("invalid plan ID format")
	}

	return pu.planRepository.UpdatePlan(ctx, objectID, updatedPlan)
}

func (ru *planUsecaseStruct) UpdateReportStatus(c context.Context, reportID primitive.ObjectID, supervisorName, status, comment string) error {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()

	// Validate status
	if status != "Approved" && status != "Rejected" {
		return fmt.Errorf("invalid status")
	}

	// Update the report in the repository
	return ru.planRepository.UpdateReportStatus(ctx, reportID, supervisorName, status, comment)
}

func (pu *planUsecaseStruct) UpdatePlanStatus(c context.Context, planID primitive.ObjectID, supervisorName, status, comment string) error {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	// Validate status
	if status != "Approved" && status != "Rejected" {
		return fmt.Errorf("invalid status")
	}

	// Update the plan in the repository
	return pu.planRepository.UpdatePlanStatus(ctx, planID, supervisorName, status, comment)
}

func (ru *planUsecaseStruct) FetchReportsBySupervisorAndStatus(c context.Context, supervisorName, reportStatus string) ([]domain.Report, error) {
	ctx, cancel := context.WithTimeout(c, ru.contextTimeout)
	defer cancel()

	return ru.planRepository.GetReportsBySupervisorAndStatus(ctx, supervisorName, reportStatus)
}

func (pu *planUsecaseStruct) FetchPlansBySupervisorAndStatus(c context.Context, supervisorName, status string) ([]domain.Plan, error) {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	return pu.planRepository.GetPlansBySupervisorAndStatus(ctx, supervisorName, status)
}

func (pu *planUsecaseStruct) CreatePlan(c context.Context, plan *domain.Plan) (*primitive.ObjectID, error) {
	ctx, cancel := context.WithTimeout(c, pu.contextTimeout)
	defer cancel()

	plan.CreatedAt = time.Now()
	plan.UpdatedAt = time.Now()
	plan.Status = "Pending" // Status starts as Pending

	// Insert the plan into the repository
	err := pu.planRepository.CreatePlan(ctx, plan)
	if err != nil {
		return nil, err
	}

	return &plan.ID, nil
}
func (pu *planUsecaseStruct) GetPlansByStatusAndOwner(ctx context.Context, userID primitive.ObjectID, status string) ([]domain.Plan, error) {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	return pu.planRepository.GetPlansByStatusAndOwner(c, userID, status)
}

func (pu *planUsecaseStruct) GetPlanTitlesByOwnerName(ctx context.Context, ownerName string) ([]string, error) {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	return pu.planRepository.GetPlanTitlesByOwnerName(c, ownerName)
}

func (ru *planUsecaseStruct) SubmitReport(ctx context.Context, report *domain.Report) error {
	c, cancel := context.WithTimeout(ctx, ru.contextTimeout)
	defer cancel()

	return ru.planRepository.SubmitReport(c, report)
}
func (ru *planUsecaseStruct) GetFilteredReports(ctx context.Context, userID primitive.ObjectID, status string) ([]domain.Report, error) {
	c, cancel := context.WithTimeout(ctx, ru.contextTimeout)
	defer cancel()

	return ru.planRepository.GetFilteredReports(c, userID, status)
}
func (pu *planUsecaseStruct) GetAllPlansByUser(ctx context.Context, userID primitive.ObjectID) ([]domain.Plan, error) {
	c, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	return pu.planRepository.GetAllPlansByUser(c, userID)
}

func (uc *planUsecaseStruct) CountItems(ctx context.Context, itemType string, toWhom string) (int, error) {
	if itemType != "plan" && itemType != "report" {
		return 0, fmt.Errorf("invalid item type")
	}

	count, err := uc.planRepository.CountItems(ctx, itemType, toWhom)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// GetUsersByToWhomWithCount fetches users whose To_whom matches the given firstName and returns their count.

func (pu *planUsecaseStruct) GetPlan(ctx context.Context, ownerID primitive.ObjectID) (*domain.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	plan, err := pu.planRepository.GetPlan(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	return plan, nil
}

// func (pu *planUsecaseStruct) EditPlan(ctx context.Context, plan *domain.Plan) error {
// 	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
// 	defer cancel()

// 	existingPlan, err := pu.planRepository.GetPlan(ctx, plan.OwnerID)
// 	if err != nil {
// 		if err == domain.ErrPlanNotFound {
// 			return domain.ErrPlanNotFound
// 		}
// 		return err
// 	}

// 	// Ensure only specific fields are updated
// 	existingPlan.Title = plan.Title
// 	existingPlan.Description = plan.Description
// 	existingPlan.Priority = plan.Priority
// 	existingPlan.WhichQuarter = plan.WhichQuarter
// 	existingPlan.Quantify = plan.Quantify
// 	existingPlan.UpdatedAt = time.Now()

// 	return pu.planRepository.UpdatePlan(ctx, existingPlan)
// }

func (pu *planUsecaseStruct) GetSubmittedPlans(ctx context.Context, supervisor_name string) ([]domain.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
	defer cancel()

	// Fetch plans submitted to the user or one level below
	plans, err := pu.planRepository.GetSubmittedPlans(ctx, supervisor_name)
	if err != nil {
		return nil, err
	}

	return plans, nil
}

// func (pu *planUsecaseStruct) ApprovePlan(ctx context.Context, planID, ownerID primitive.ObjectID) error {
// 	ctx, cancel := context.WithTimeout(ctx, pu.contextTimeout)
// 	defer cancel()

// 	// Fetch the plan
// 	plan, err := pu.planRepository.GetPlanByID(ctx, planID)
// 	if err != nil {
// 		return err
// 	}

// 	// Check if the user is authorized to approve this plan
// 	// if plan.OwnerID != ownerID {
// 	// 	return errors.New("unauthorized: only the plan owner can approve it")
// 	// }

// 	// Update the status to "Approved"
// 	plan.Status = "Approved"
// 	plan.UpdatedAt = time.Now()

// 	// If the plan has a supervisor, forward it to the supervisor
// 	if plan.SupervisorName != "" {
// 		err = pu.planRepository.ForwardPlanToSupervisor(ctx, plan)
// 		if err != nil {
// 			return err
// 		}
// 	} else {
// 		// Just save the updated status if there's no supervisor
// 		err = pu.planRepository.UpdatePlanPacth(ctx, plan)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

func (cu *planUsecaseStruct) AddComment(ctx context.Context, comment *domain.Comment) error {
	ctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	// Add the comment to the database
	return cu.planRepository.CreateComment(ctx, comment)
}

func (cu *planUsecaseStruct) GetSupervisorComments(ctx context.Context, userID primitive.ObjectID) ([]domain.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	// Fetch comments made by the supervisor on plans related to the user
	return cu.planRepository.FetchSupervisorComments(ctx, userID)
}

func (cu *planUsecaseStruct) GetCommentsByPlanID(ctx context.Context, planID primitive.ObjectID) ([]domain.Comment, error) {
	ctx, cancel := context.WithTimeout(ctx, cu.contextTimeout)
	defer cancel()

	// Fetch comments from the repository
	return cu.planRepository.FetchCommentsByPlanID(ctx, planID)
}
