package leaderboard

import (
	"context"
	"distasteful-bear/turing_machine/api/db"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// LeaderboardEntry represents a user on the leaderboard
type LeaderboardEntry struct {
	Rank          int
	Nickname      string
	Picture       string
	GamesWon      int
	GamesLost     int
	WinRate       string
	AvgGuesses    string
	IsCurrentUser bool
}

// LeaderboardPageData contains all data needed to render the leaderboard page
type LeaderboardPageData struct {
	IsLoggedIn  bool
	CurrentUser *LeaderboardEntry
	TopUsers    []LeaderboardEntry
}

// calculateStats computes win rate and average guesses per game
func calculateStats(gamesPlayed, gamesWon, totalGuesses int) (winRate string, avgGuesses string) {
	if gamesPlayed == 0 {
		return "--", "--"
	}
	winRateVal := float64(gamesWon) / float64(gamesPlayed) * 100
	avgGuessesVal := float64(totalGuesses) / float64(gamesPlayed)
	return fmt.Sprintf("%.0f%%", winRateVal), fmt.Sprintf("%.1f", avgGuessesVal)
}

// Handler for the leaderboard page
func Handler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	pageData := LeaderboardPageData{
		IsLoggedIn:  false,
		CurrentUser: nil,
		TopUsers:    []LeaderboardEntry{},
	}

	var currentUserID string

	// Check if user is logged in
	if profile != nil {
		if profileMap, ok := profile.(map[string]interface{}); ok {
			pageData.IsLoggedIn = true
			currentUserID, _ = profileMap["sub"].(string)
		}
	}

	// Get Firestore client
	firestoreCtx := context.Background()
	client, err := db.GetFirestoreClient(firestoreCtx)
	if err != nil {
		log.Printf("Error creating Firestore client: %v", err)
		ctx.HTML(http.StatusOK, "leaderboard.html", pageData)
		return
	}
	defer client.Close()

	// Get current user's stats if logged in
	if pageData.IsLoggedIn && currentUserID != "" {
		userRecord, err := db.GetOrCreateUser(firestoreCtx, client, currentUserID, "", "")
		if err == nil {
			winRate, avgGuesses := calculateStats(userRecord.GamesPlayed, userRecord.GamesWon, userRecord.TotalGuesses)
			pageData.CurrentUser = &LeaderboardEntry{
				Rank:          0, // Will be determined by position in leaderboard
				Nickname:      userRecord.Nickname,
				Picture:       userRecord.Picture,
				GamesWon:      userRecord.GamesWon,
				GamesLost:     userRecord.GamesLost,
				WinRate:       winRate,
				AvgGuesses:    avgGuesses,
				IsCurrentUser: true,
			}
		}
	}

	// Get top 10 users
	topUsers, err := db.GetTopUsers(firestoreCtx, client, 10)
	if err != nil {
		log.Printf("Error getting top users: %v", err)
		ctx.HTML(http.StatusOK, "leaderboard.html", pageData)
		return
	}

	// Convert to LeaderboardEntry with ranks
	for i, user := range topUsers {
		winRate, avgGuesses := calculateStats(user.GamesPlayed, user.GamesWon, user.TotalGuesses)
		entry := LeaderboardEntry{
			Rank:          i + 1,
			Nickname:      user.Nickname,
			Picture:       user.Picture,
			GamesWon:      user.GamesWon,
			GamesLost:     user.GamesLost,
			WinRate:       winRate,
			AvgGuesses:    avgGuesses,
			IsCurrentUser: user.UserID == currentUserID,
		}
		pageData.TopUsers = append(pageData.TopUsers, entry)

		// Update current user's rank if found in top 10
		if entry.IsCurrentUser && pageData.CurrentUser != nil {
			pageData.CurrentUser.Rank = i + 1
		}
	}

	ctx.HTML(http.StatusOK, "leaderboard.html", pageData)
}
