package user

import (
	"context"
	"distasteful-bear/turing_machine/api/db"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// UserPageData contains all data needed to render the user page
type UserPageData struct {
	Nickname        string
	Picture         string
	GamesPlayed     int
	LeaderboardRank int
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
			LeaderboardRank: 0,
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
			LeaderboardRank: 0,
		})
		return
	}

	ctx.HTML(http.StatusOK, "user.html", UserPageData{
		Nickname:        userRecord.Nickname,
		Picture:         userRecord.Picture,
		GamesPlayed:     userRecord.GamesPlayed,
		LeaderboardRank: userRecord.LeaderboardRank,
	})
}
