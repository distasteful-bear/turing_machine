package user

import (
	"context"
	"distasteful-bear/turing_machine/api/db"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// calculateStats computes win rate and average guesses per game
func calculateStats(gamesPlayed, gamesWon, totalGuesses int) (winRate string, avgGuesses string) {
	if gamesPlayed == 0 {
		return "--", "--"
	}
	winRateVal := float64(gamesWon) / float64(gamesPlayed) * 100
	avgGuessesVal := float64(totalGuesses) / float64(gamesPlayed)
	return fmt.Sprintf("%.0f%%", winRateVal), fmt.Sprintf("%.1f", avgGuessesVal)
}

// UserPageData contains all data needed to render the user page
type UserPageData struct {
	Nickname        string
	Picture         string
	GamesPlayed     int
	GamesWon        int
	GamesLost       int
	TotalGuesses    int
	LeaderboardRank int
	WinRate         string
	AvgGuesses      string
}

// Handler for our logged-in user page.
func Handler(ctx *gin.Context) {
	session := sessions.Default(ctx)
	profile := session.Get("profile")

	// Extract profile data
	profileMap, ok := profile.(map[string]interface{})
	if !ok {
		ctx.HTML(http.StatusInternalServerError, "user.html", gin.H{"error": "Invalid profile data"})
		return
	}

	// Get user identifiers from profile
	userID, _ := profileMap["sub"].(string)
	nickname, _ := profileMap["nickname"].(string)
	picture, _ := profileMap["picture"].(string)

	// Get or create user in Firestore
	firestoreCtx := context.Background()
	client, err := db.GetFirestoreClient(firestoreCtx)
	if err != nil {
		log.Printf("Error creating Firestore client: %v", err)
		ctx.HTML(http.StatusOK, "user.html", UserPageData{
			Nickname:        nickname,
			Picture:         picture,
			GamesPlayed:     0,
			GamesWon:        0,
			GamesLost:       0,
			TotalGuesses:    0,
			LeaderboardRank: 0,
			WinRate:         "--",
			AvgGuesses:      "--",
		})
		return
	}
	defer client.Close()

	userRecord, err := db.GetOrCreateUser(firestoreCtx, client, userID, nickname, picture)
	if err != nil {
		log.Printf("Error getting/creating user record: %v", err)
		ctx.HTML(http.StatusOK, "user.html", UserPageData{
			Nickname:        nickname,
			Picture:         picture,
			GamesPlayed:     0,
			GamesWon:        0,
			GamesLost:       0,
			TotalGuesses:    0,
			LeaderboardRank: 0,
			WinRate:         "--",
			AvgGuesses:      "--",
		})
		return
	}

	winRate, avgGuesses := calculateStats(userRecord.GamesPlayed, userRecord.GamesWon, userRecord.TotalGuesses)

	ctx.HTML(http.StatusOK, "user.html", UserPageData{
		Nickname:        userRecord.Nickname,
		Picture:         userRecord.Picture,
		GamesPlayed:     userRecord.GamesPlayed,
		GamesWon:        userRecord.GamesWon,
		GamesLost:       userRecord.GamesLost,
		TotalGuesses:    userRecord.TotalGuesses,
		LeaderboardRank: userRecord.LeaderboardRank,
		WinRate:         winRate,
		AvgGuesses:      avgGuesses,
	})
}

func IsUserLoggedIn(c *gin.Context) (string, error) {
	// return uid and optional err
	session := sessions.Default(c)
	profile := session.Get("profile")
	if profile == nil {
		return "", errors.New("no profile")
	}

	profileMap, ok := profile.(map[string]interface{})
	if !ok {
		return "", errors.New("invalid profile")
	}

	userId, ok := profileMap["sub"].(string)

	if ok {
		return userId, nil
	} else {
		return "", errors.New("could not locate uid")
	}
}
