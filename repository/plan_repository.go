package repository

import (
	"context"
	"fmt"
	"plan/database"
	"plan/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"errors"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type planRepository struct {
	database   database.Database
	collection string
}

func NewPlanRepository(db database.Database, collection string) domain.PlanRepository {
	return &planRepository{
		database:   db,
		collection: collection,
	}
}

func (repo *planRepository) FindByOwnerID(ctx context.Context, ownerID primitive.ObjectID, datatype string) ([]domain.Plan, error) {
	var plans []domain.Plan
	cursor, err := repo.database.Collection(repo.collection).Find(ctx, bson.M{"owner_id": ownerID}, options.Find())
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &plans); err != nil {
		return nil, err
	}
	return plans, nil
}

func (repo *planRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID, datatype string) ([]domain.Report, error) {
	var reports []domain.Report
	cursor, err := repo.database.Collection(repo.collection).Find(ctx, bson.M{"report_user_id": userID}, options.Find())
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, err
	}
	return reports, nil
}
func (repo *planRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	collection := repo.database.Collection(repo.collection)

	// Delete the document with the given ID
	result, err := collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	if result == 0 {
		return errors.New("announcement not found")
	}

	return nil
}
func (rr *planRepository) GetAllAnnouncements(ctx context.Context) ([]domain.Announcement, error) {
	collection := rr.database.Collection(rr.collection)

	var announcements []domain.Announcement
	cursor, err := collection.Find(ctx, bson.M{"type": "announcement"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var announcement domain.Announcement
		if err := cursor.Decode(&announcement); err != nil {
			return nil, err
		}
		announcements = append(announcements, announcement)
	}

	return announcements, nil
}

func (rr *planRepository) CreateAnnouncement(ctx context.Context, announcement *domain.Announcement) error {
	collection := rr.database.Collection(rr.collection)
	_, err := collection.InsertOne(ctx, announcement)
	return err
}
func (rr *planRepository) UpdateReport(ctx context.Context, reportID primitive.ObjectID, updatedReport *domain.Report) error {
	collection := rr.database.Collection(rr.collection)

	filter := bson.M{"_id": reportID}
	update := bson.M{
		"$set": bson.M{
			"report_title":      updatedReport.ReportTitle,
			"acomplished_value": updatedReport.AccomplishedValue,
			"report_details":    updatedReport.ReportDetails,
			"type":              updatedReport.Type,
			"supervisor_name":   updatedReport.SupervisorName,
			"report_status":     updatedReport.Status, // Always "Pending"
			"updated_at":        time.Now(),
			"comment":           "",
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("report not found")
	}

	return nil
}

func (pr *planRepository) UpdatePlan(ctx context.Context, planID primitive.ObjectID, updatedPlan *domain.Plan) error {
	collection := pr.database.Collection(pr.collection)

	filter := bson.M{"_id": planID}
	update := bson.M{
		"$set": bson.M{
			"title":           updatedPlan.Title,
			"description":     updatedPlan.Description,
			"priority":        updatedPlan.Priority,
			"which_quarter":   updatedPlan.WhichQuarter,
			"quantify":        updatedPlan.Quantify,
			"aligned_pillary": updatedPlan.AlignedPillary,
			"start_date":      updatedPlan.StartDate,
			"end_date":        updatedPlan.EndDate,
			"type":            updatedPlan.Type,
			"status":          updatedPlan.Status, // Always "Pending"
			"updated_at":      time.Now(),
			"comment":         "",
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("plan not found")
	}

	return nil
}

func (rr *planRepository) UpdateReportStatus(ctx context.Context, reportID primitive.ObjectID, supervisorName, status, comment string) error {
	collection := rr.database.Collection(rr.collection)

	// Ensure the supervisor is authorized
	filter := bson.M{
		"_id":             reportID,
		"supervisor_name": supervisorName,
	}

	update := bson.M{
		"$set": bson.M{
			"report_status": status,
			"comment":       comment,
			"updated_at":    time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("report not found")
	}

	return nil
}

func (pr *planRepository) UpdatePlanStatus(ctx context.Context, planID primitive.ObjectID, supervisorName, status, comment string) error {
	collection := pr.database.Collection(pr.collection)

	// Ensure the supervisor is authorized
	filter := bson.M{
		"_id":             planID,
		"supervisor_name": supervisorName,
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"comment":    comment,
			"updated_at": time.Now(),
		},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("plan not found")
	}

	return nil
}

func (rr *planRepository) GetReportsBySupervisorAndStatus(ctx context.Context, supervisorName, reportStatus string) ([]domain.Report, error) {
	collection := rr.database.Collection(rr.collection)

	// Filter by report_status and supervisor_name
	filter := bson.M{
		"report_status":   reportStatus,
		"supervisor_name": supervisorName,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var reports []domain.Report
	if err = cursor.All(ctx, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}

func (pr *planRepository) GetPlansBySupervisorAndStatus(ctx context.Context, supervisorName, status string) ([]domain.Plan, error) {
	collection := pr.database.Collection(pr.collection)

	// Filter by status and supervisor name
	filter := bson.M{
		"status":          status,
		"supervisor_name": supervisorName,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plans []domain.Plan
	if err = cursor.All(ctx, &plans); err != nil {
		return nil, err
	}

	return plans, nil
}

func (pr *planRepository) CreatePlan(c context.Context, plan *domain.Plan) error {
	fmt.Printf("plan: %v\n", plan)
	plan.Type = "plan"
	collection := pr.database.Collection(pr.collection)
	plan.ID = primitive.NewObjectID() // Create a new ID for the plan
	_, err := collection.InsertOne(c, plan)
	return err
}
func (pr *planRepository) GetPlanTitlesByOwnerName(ctx context.Context, ownerName string) ([]string, error) {
	collection := pr.database.Collection(pr.collection)

	filter := bson.M{"owner_name": ownerName}
	projection := bson.M{"title": 1, "_id": 0} // Only fetch the "title" field

	findOptions := options.Find().SetProjection(projection)
	cursor, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plans []struct {
		Title string `bson:"title"`
	}
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, err
	}

	titles := make([]string, len(plans))
	for i, plan := range plans {
		titles[i] = plan.Title
	}

	return titles, nil
}
func (pr *planRepository) GetPlansByStatusAndOwner(ctx context.Context, userID primitive.ObjectID, status string) ([]domain.Plan, error) {
	collection := pr.database.Collection(pr.collection)

	// Create the filter
	filter := bson.M{
		"owner_id": userID,
		"status":   status,
		"type":     "plan",
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plans []domain.Plan
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, err
	}

	return plans, nil
}
func (rr *planRepository) SubmitReport(ctx context.Context, report *domain.Report) error {
	report.Status = "Pending"
	collection := rr.database.Collection(rr.collection)
	report.ID = primitive.NewObjectID()

	_, err := collection.InsertOne(ctx, report)
	return err

}

func (rr *planRepository) GetFilteredReports(ctx context.Context, userID primitive.ObjectID, status string) ([]domain.Report, error) {
	filter := bson.M{
		"report_user_id": userID,
		"status":         status,
		"type":           "report",
	}

	var reports []domain.Report
	cursor, err := rr.database.Collection(rr.collection).Find(ctx, filter, options.Find())
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &reports); err != nil {
		return nil, err
	}

	return reports, nil
}
func (pr *planRepository) GetAllPlansByUser(ctx context.Context, userID primitive.ObjectID) ([]domain.Plan, error) {
	filter := bson.M{
		"owner_id": userID,
		"type":     "plan",
		"status":   "Approved", // Assuming status "Approved" is required
	}

	cursor, err := pr.database.Collection(pr.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var plans []domain.Plan
	if err := cursor.All(ctx, &plans); err != nil {
		return nil, err
	}

	return plans, nil
}

func (r *planRepository) CountItems(ctx context.Context, itemType string, toWhom string) (int, error) {
	collection := r.database.Collection(r.collection) // Adjust `database` and `collection` initialization

	// Create the filter
	filter := bson.M{
		"type":            itemType,
		"supervisor_name": toWhom,
		"status":          "Pending",
	}

	// Count the documents that match the filter
	count, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// FetchByToWhom retrieves users whose To_whom matches the specified firstName.

func (pr *planRepository) GetPlan(ctx context.Context, ownerID primitive.ObjectID) (*domain.Plan, error) {
	filter := bson.M{
		"owner_id": ownerID,
	}

	var plan domain.Plan
	err := pr.database.Collection(pr.collection).FindOne(ctx, filter).Decode(&plan)

	if err != nil {
		return nil, err
	}

	return &plan, nil
}

// func (pr *planRepository) UpdatePlan(ctx context.Context, plan *domain.Plan) error {
// 	filter := bson.M{
// 		"owner_id": plan.OwnerID,
// 	}

// 	update := bson.M{
// 		"$set": bson.M{
// 			"title":         plan.Title,
// 			"description":   plan.Description,
// 			"priority":      plan.Priority,
// 			"which_quarter": plan.WhichQuarter,
// 			"quantify":      plan.Quantify,
// 			"updated_at":    plan.UpdatedAt,
// 			"status":        plan.Status,
// 		},
// 	}

// 	result, err := pr.database.Collection(pr.collection).UpdateOne(ctx, filter, update)

// 	if err != nil {
// 		return err
// 	}

// 	if result.MatchedCount == 0 {
// 		return domain.ErrPlanNotFound
// 	}

// 	return nil
// }

func (pr *planRepository) GetSubmittedPlans(ctx context.Context, supervisor_name string) ([]domain.Plan, error) {
	var plans []domain.Plan

	// Query for plans where `SupervisorPlanID` equals the user's ID
	filter := bson.M{
		"$or": []bson.M{
			{"supervisor_name": supervisor_name},
		},
	}

	cursor, err := pr.database.Collection(pr.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &plans); err != nil {
		return nil, err
	}

	return plans, nil
}

func (pr *planRepository) GetPlanByID(ctx context.Context, planID primitive.ObjectID) (*domain.Plan, error) {
	var plan domain.Plan

	// Find the plan by its ID
	err := pr.database.Collection(pr.collection).FindOne(ctx, bson.M{"_id": planID}).Decode(&plan)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return nil, errors.New("plan not found")
		}

		return nil, err
	}

	return &plan, nil
}

func (pr *planRepository) UpdatePlanPacth(ctx context.Context, plan *domain.Plan) error {
	filter := bson.M{"_id": plan.ID}
	update := bson.M{"$set": plan}

	// cursor, err := pr.database.Collection(pr.collection).Find(ctx, filter)

	_, err := pr.database.Collection(pr.collection).UpdateOne(ctx, filter, update)
	return err
}

// func (pr *planRepository) ForwardPlanToSupervisor(ctx context.Context, plan *domain.Plan) error {
// 	if plan.SupervisorPlanID == nil {
// 		return errors.New("no supervisor to forward to")
// 	}

// 	err := pr.UpdatePlan(ctx, plan)
// 	if err != nil {
// 		panic(err)
// 	}

// 	// Duplicate the plan for the supervisor and update its owner
// 	// newPlan := *plan
// 	// newPlan.ID = primitive.NewObjectID() // Create a new ID for the forwarded plan
// 	// newPlan.OwnerID = *plan.SupervisorPlanID
// 	// newPlan.Status = "Pending"
// 	// newPlan.CreatedAt = time.Now()
// 	// newPlan.UpdatedAt = time.Now()

// 	// Insert the duplicated plan into the database
// 	// cursor, err := pr.database.Collection(pr.collection).Find(ctx, filter)
// 	// _, err = pr.database.Collection(pr.collection).InsertOne(ctx, newPlan)
// 	return err
// }

func (cr *planRepository) CreateComment(ctx context.Context, comment *domain.Comment) error {
	_, err := cr.database.Collection(cr.collection).InsertOne(ctx, comment)
	return err
}

func (cr *planRepository) FetchSupervisorComments(ctx context.Context, userID primitive.ObjectID) ([]domain.Comment, error) {
	filter := bson.M{
		"commenter": "supervisor", // Assuming "supervisor" is stored as the commenter's username
		"plan_id": bson.M{
			"$in": bson.A{
				bson.M{"owner_id": userID}, // Plans owned by the user
			},
		},
	}

	cursor, err := cr.database.Collection(cr.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var comments []domain.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}

func (cr *planRepository) FetchCommentsByPlanID(ctx context.Context, planID primitive.ObjectID) ([]domain.Comment, error) {
	// Filter to match the specified PlanID
	filter := bson.M{"plan_id": planID}

	// Find comments matching the PlanID
	cursor, err := cr.database.Collection(cr.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(ctx)

	var comments []domain.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return nil, err
	}

	return comments, nil
}
