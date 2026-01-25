package db

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserRecord represents a user's data in Firestore
type UserRecord struct {
	UserID          string `firestore:"user_id"`
	Nickname        string `firestore:"nickname"`
	Picture         string `firestore:"picture"`
	GamesPlayed     int    `firestore:"games_played"`
	LeaderboardRank int    `firestore:"leaderboard_rank"`
}

func GetFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	return firestore.NewClientWithDatabase(ctx, "james-metz", "turing-machine")
}

func GetUserRecord(ctx context.Context, client *firestore.Client, userID string) (*firestore.DocumentSnapshot, error) {
	return client.Collection("users").Doc(userID).Get(ctx)
}

// GetOrCreateUser retrieves an existing user or creates a new one if not found
func GetOrCreateUser(ctx context.Context, client *firestore.Client, userID, nickname, picture string) (*UserRecord, error) {
	docRef := client.Collection("users").Doc(userID)
	doc, err := docRef.Get(ctx)

	if err != nil {
		// Check if document doesn't exist
		if status.Code(err) == codes.NotFound {
			// Create new user record
			newUser := &UserRecord{
				UserID:          userID,
				Nickname:        nickname,
				Picture:         picture,
				GamesPlayed:     0,
				LeaderboardRank: 0,
			}
			_, err = docRef.Set(ctx, newUser)
			if err != nil {
				return nil, err
			}
			return newUser, nil
		}
		return nil, err
	}

	// Document exists, unmarshal it
	var user UserRecord
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
