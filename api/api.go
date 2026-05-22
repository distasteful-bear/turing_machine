package api

import (
	"distasteful-bear/turing_machine/api/authenticator"
	"distasteful-bear/turing_machine/api/callback"
	"distasteful-bear/turing_machine/api/db"
	"distasteful-bear/turing_machine/api/leaderboard"
	"distasteful-bear/turing_machine/api/login"
	"distasteful-bear/turing_machine/api/logout"
	"distasteful-bear/turing_machine/api/middleware"
	"distasteful-bear/turing_machine/api/session"
	"distasteful-bear/turing_machine/api/user"
	"distasteful-bear/turing_machine/utils"
	"distasteful-bear/turing_machine/verifiers"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func SetupRouter() *gin.Engine {
	// Initialize Gin router
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// Setup server-side session store for puzzles
	router.Use(session.SetupSessionStoreInMem())

	// Load HTML templates from the src directory
	router.LoadHTMLGlob("src/*.html")

	// Serve static files (CSS)
	router.Static("/static", "./src")

	// HTML
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/welcome", func(c *gin.Context) {
		c.HTML(200, "welcome.html", nil)
	})
	router.GET("/guess", func(c *gin.Context) {
		c.HTML(200, "guess.html", nil)
	})
	router.GET("/final", func(c *gin.Context) {
		c.HTML(200, "final.html", nil)
	})
	router.GET("/failure", func(c *gin.Context) {
		c.HTML(200, "failure.html", nil)
	})
	router.GET("/success", func(c *gin.Context) {
		c.HTML(200, "success.html", nil)
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "welcome.html", nil)
	})

	// auth routes
	auth, err := authenticator.New()
	if err != nil {
		log.Printf("Auth0 is not available; falling back to local dev login: %v", err)
		auth = nil
	}
	if auth == nil {
		log.Println("Auth0 env vars are not fully configured; /login will use a local dev user.")
	}
	router.GET("/login", login.Handler(auth))
	router.GET("/callback", callback.Handler(auth))
	router.GET("/logout", logout.Handler)
	router.GET("/user", middleware.IsAuthenticated, user.Handler)
	router.GET("/leaderboard", leaderboard.Handler)

	// Check authentication status
	router.GET("/auth_status", func(c *gin.Context) {
		session := sessions.Default(c)
		profile := session.Get("profile")

		if profile == nil {
			c.JSON(200, gin.H{"logged_in": false})
			return
		}

		profileMap, ok := profile.(map[string]interface{})
		if !ok {
			c.JSON(200, gin.H{"logged_in": false})
			return
		}

		c.JSON(200, gin.H{
			"logged_in": true,
			"nickname":  profileMap["nickname"],
			"picture":   profileMap["picture"],
		})
	})

	setupPuzzleRoutes(router)

	return router
}

func setupPuzzleRoutes(router gin.IRoutes) {
	// puzzle data management
	router.GET("/setup_session", func(c *gin.Context) {
		curSession := sessions.Default(c)

		puzzle := verifiers.GenerateRandomPuzzle()

		storeInterface, ok := c.Get("session_store")
		if !ok {
			c.JSON(500, gin.H{"error": "session store not initialized"})
			return
		}
		store := storeInterface.(*session.SessionStore)

		// Generate a unique puzzle ID for this session
		sessionID := uuid.New().String()

		if err := store.AddPuzzle(sessionID, puzzle, time.Now()); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"error": err.Error()})
			return
		}

		// Store only the session ID in the session cookie and reset guess count
		curSession.Set("puzzle_id", sessionID)
		curSession.Set("guess_count", 0)
		curSession.Save()

		type verSummary struct {
			Id          int    `json:"id"`
			Description string `json:"description"`
			Result      string `json:"result"`
		}
		puzzleSummary := []verSummary{}
		for i, v := range puzzle.Vers {
			puzzleSummary = append(puzzleSummary, verSummary{
				Id:          i,
				Description: v.Desc,
				Result:      "",
			})
		}

		c.JSON(200, gin.H{"status": "success", "puzzle_summary": puzzleSummary})
	})

	router.GET("/check_guess", func(c *gin.Context) {

		curSession := sessions.Default(c)

		// Retrieve puzzle ID from session
		puzzleId := curSession.Get("puzzle_id")
		if puzzleId == nil {
			fmt.Println("No puzzle session found")
			c.JSON(400, gin.H{"error": "no puzzle session found"})
			return
		}

		// Get session store
		storeInterface, exists := c.Get("session_store")
		if !exists {
			fmt.Println("Session store not initialized")
			c.JSON(500, gin.H{"error": "session store not initialized"})
			return
		}
		store := storeInterface.(*session.SessionStore)

		// Retrieve puzzle from server-side storage
		puzzle, ok := store.GetPuzzle(puzzleId.(string), time.Now())
		if !ok {
			c.JSON(400, gin.H{"error": "puzzle not found in store or was expired"})
			return
		}

		guess := c.Query("guess")
		if guess == "" {
			c.JSON(400, gin.H{"error": "no guess"})
			return
		}
		proposedSolution, err := utils.SanitizeGuess(guess)
		if err != nil {
			c.JSON(400, gin.H{"error": "error parsing guess"})
			return
		}

		// Increment guess count in session
		guessCount := curSession.Get("guess_count")
		if guessCount == nil {
			guessCount = 0
		}
		newGuessCount := guessCount.(int) + 1
		curSession.Set("guess_count", newGuessCount)
		curSession.Save()

		type verSummary struct {
			Id          int    `json:"id"`
			Description string `json:"description"`
			Result      string `json:"result"`
		}
		puzzleSummary := []verSummary{}
		for i, v := range puzzle.Vers {
			passTest := v.VerifierFunc(proposedSolution)
			if passTest {
				puzzleSummary = append(puzzleSummary, verSummary{
					Id:          i,
					Description: v.Desc,
					Result:      "true",
				})
			} else {
				puzzleSummary = append(puzzleSummary, verSummary{
					Id:          i,
					Description: v.Desc,
					Result:      "false",
				})
			}
		}
		c.JSON(200, gin.H{"status": "success", "puzzle_summary": puzzleSummary, "guess_count": newGuessCount})
	})

	router.GET("/check_final", func(c *gin.Context) {
		curSession := sessions.Default(c)
		storeInterface, exists := c.Get("session_store")
		if !exists {
			c.JSON(500, gin.H{"error": "session store not initialized"})
			return
		}
		store := storeInterface.(*session.SessionStore)

		// query params
		puzzleId := curSession.Get("puzzle_id").(string)
		if puzzleId == "" {
			c.JSON(400, gin.H{"error": "no puzzle session found"})
			return
		}
		guessCount, ok := curSession.Get("guess_count").(int)
		if !ok {
			guessCount = 0
		}
		guess := c.Query("guess")
		if guess == "" {
			c.JSON(400, gin.H{"error": "no guess"})
			return
		}
		proposedSolution, err := utils.SanitizeGuess(guess)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		// retrieve puzzle from server-side storage
		puzzle, ok := store.GetPuzzle(puzzleId, time.Now())
		defer store.DeletePuzzle(puzzleId)
		if !ok {
			c.JSON(400, gin.H{"error": "puzzle not found in store"})
			return
		}

		success := proposedSolution.Display == puzzle.Sol.Display

		// log results if user is logged in
		userId, err := user.IsUserLoggedIn(c)
		recordedGame := false
		if err == nil {
			err := db.RecordGameCompletion(c.Request.Context(), userId, success, guessCount)
			if err != nil {
				log.Printf("Error recording game completion: %v", err)
			} else {
				recordedGame = true
			}
		}

		// compute leaderboard rankings
		if recordedGame {
			err = db.ComputeLeaderboardRankings(c.Request.Context())
			if err != nil {
				log.Printf("Error computing leaderboard rankings: %v", err)
			}
		}

		// Convert rune array to string for JSON response
		solutionStr := string(puzzle.Sol.Display[:])

		if success {
			c.JSON(200, gin.H{"status": "success", "guess_count": guessCount, "solution": solutionStr})
			return
		} else {
			c.JSON(200, gin.H{"status": "failure", "guess_count": guessCount, "solution": solutionStr})
			return
		}
	})

}
