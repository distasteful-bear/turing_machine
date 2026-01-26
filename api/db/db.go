package db

import (
	"context"
	"sort"

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
	GamesWon        int    `firestore:"games_won"`
	GamesLost       int    `firestore:"games_lost"`
	TotalGuesses    int    `firestore:"total_guesses"`
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

// RecordGameCompletion updates a user's stats after completing a game
func RecordGameCompletion(ctx context.Context, client *firestore.Client, userID string, won bool, guessCount int) error {
	docRef := client.Collection("users").Doc(userID)
	doc, err := docRef.Get(ctx)
	if err != nil {
		return err
	}

	var user UserRecord
	if err := doc.DataTo(&user); err != nil {
		return err
	}

	user.GamesPlayed++
	user.TotalGuesses += guessCount
	if won {
		user.GamesWon++
	} else {
		user.GamesLost++
	}

	_, err = docRef.Set(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

// GetTopUsers retrieves the top N users sorted by leaderboard rank
func GetTopUsers(ctx context.Context, client *firestore.Client, limit int) ([]UserRecord, error) {
	docs, err := client.Collection("users").
		Where("leaderboard_rank", ">", 0).
		OrderBy("leaderboard_rank", firestore.Asc).
		Limit(limit).
		Documents(ctx).
		GetAll()

	if err != nil {
		return nil, err
	}

	users := make([]UserRecord, 0, len(docs))
	for _, doc := range docs {
		var user UserRecord
		if err := doc.DataTo(&user); err != nil {
			continue
		}
		users = append(users, user)
	}

	return users, nil
}

// rankedUser holds user data with computed score for ranking
type rankedUser struct {
	UserID string
	Score  float64
}

// ComputeLeaderboardRankings calculates and updates leaderboard rankings
// Ranking formula: games_played / avg_guesses_per_game
// This rewards playing many games while encouraging fewer guesses per game
func ComputeLeaderboardRankings(ctx context.Context, client *firestore.Client) error {
	// Step 1: Query top 100 users by games played
	docs, err := client.Collection("users").
		OrderBy("games_played", firestore.Desc).
		Limit(10).
		Documents(ctx).
		GetAll()

	if err != nil {
		return err
	}

	// Step 2: Calculate ranking score for each user
	var rankedUsers []rankedUser
	for _, doc := range docs {
		var user UserRecord
		if err := doc.DataTo(&user); err != nil {
			continue
		}

		// Calculate average guesses per game
		avgGuesses := float64(user.TotalGuesses) / float64(user.GamesPlayed)

		// Avoid division by zero
		if avgGuesses == 0 {
			avgGuesses = 10
		}

		// Score = games_played / avg_guesses
		// Higher is better (more games, fewer guesses)
		score := float64(user.GamesPlayed) / avgGuesses

		rankedUsers = append(rankedUsers, rankedUser{
			UserID: user.UserID,
			Score:  score,
		})
	}

	// Sort by score descending (higher score = better rank)
	sort.Slice(rankedUsers, func(i, j int) bool {
		return rankedUsers[i].Score > rankedUsers[j].Score
	})

	// Step 3: Write ranks to user documents
	// First, clear ranks for users who might have dropped out of top 100
	// We'll do this by setting rank to 0 for all users first, then updating ranked ones
	bulkWriter := client.BulkWriter(ctx)

	// Get all users who currently have a rank
	rankedDocs, err := client.Collection("users").
		Where("leaderboard_rank", ">", 0).
		Documents(ctx).
		GetAll()

	if err != nil {
		return err
	}

	// Clear existing ranks
	for _, doc := range rankedDocs {
		bulkWriter.Update(doc.Ref, []firestore.Update{
			{Path: "leaderboard_rank", Value: 0},
		})
	}
	bulkWriter.Flush()

	// separate bulk writer for updating ranks
	bulkWriter = client.BulkWriter(ctx)

	// Set new ranks for top users
	for rank, ru := range rankedUsers {
		docRef := client.Collection("users").Doc(ru.UserID)
		bulkWriter.Update(docRef, []firestore.Update{
			{Path: "leaderboard_rank", Value: rank + 1}, // 1-indexed rank
		})
	}

	// Commit all updates
	bulkWriter.Flush()
	return nil
}
