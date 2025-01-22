package repository

import (
	"context"
	"fmt"
	"plan/database"
	"plan/domain"

	"go.mongodb.org/mongo-driver/bson"

	"errors"

	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type userRepository struct {
	database   database.Database
	collection string
}

func NewUserRepository(db database.Database, collection string) domain.UserRepository {
	return &userRepository{
		database:   db,
		collection: collection,
	}
}


func (u *userRepository) CreateUser(c context.Context, user *domain.User) error {
	user.Created_At = primitive.NewDateTimeFromTime(time.Now())
	collection := u.database.Collection(u.collection)
	_, err := collection.InsertOne(c, user)
	fmt.Println(user)
	return err
}
func (ur *userRepository) FindUnverifiedUsersByToWhom(ctx context.Context, firstName string) ([]domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	filter := bson.M{
		"verify":  false,
		"to_whom": firstName,
	}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var users []domain.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}
func (ur *userRepository) FetchByToWhom(ctx context.Context, firstName string) ([]domain.User, error) {
	var users []domain.User

	// Query to match To_whom with firstName
	filter := bson.M{
		"to_whom": firstName,
		"verify":  true,
	}
	cursor, err := ur.database.Collection(ur.collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Decode the results
	if err := cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func (ur *userRepository) GetUserByID(ctx context.Context, userID primitive.ObjectID) (*domain.User, error) {
	collection := ur.database.Collection(ur.collection)

	filter := bson.M{"_id": userID}
	var user domain.User
	err := collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (ur *userRepository) UpdateVerifyStatus(ctx context.Context, userID primitive.ObjectID, verify bool) error {
	collection := ur.database.Collection(ur.collection)

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": bson.M{"verify": verify}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, userID primitive.ObjectID) error {
	collection := ur.database.Collection(ur.collection)

	filter := bson.M{"_id": userID}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}
	if result == 0 {
		return errors.New("user not found")
	}

	return nil
}

func (ur *userRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var user domain.User
	filter := bson.M{"email": username}

	err := ur.database.Collection(ur.collection).FindOne(ctx, filter).Decode(&user)
	if err != nil {
		fmt.Println("user not found")
		return nil, errors.New("user not found")
	}

	return &user, nil
}

// func (ur *userRepository) GetUserByID(ctx context.Context, userID string) (*domain.User, error) {
// 	var user domain.User
// 	objectID, err := primitive.ObjectIDFromHex(userID)
// 	if err != nil {
// 		return nil, errors.New("invalid user ID")
// 	}

// 	filter := bson.M{"_id": objectID}

// 	err = ur.database.Collection(ur.collection).FindOne(ctx, filter).Decode(&user)
// 	if err != nil {
// 		return nil, errors.New("user not found")
// 	}

// 	return &user, nil
// }

func (u *userRepository) FindUsersByRole(c context.Context, role string) ([]domain.User, error) {
	collection := u.database.Collection(u.collection)

	filter := bson.M{
		"role":   role,
		"verify": true,
	}
	cursor, err := collection.Find(c, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(c)

	var users []domain.User
	for cursor.Next(c) {
		var user domain.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err != nil {
		return nil, err
	}

	return users, nil
}
