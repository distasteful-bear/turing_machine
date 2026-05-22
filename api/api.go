package api

import (
	"distasteful-bear/turing_machine/api/authenticator"
	"distasteful-bear/turing_machine/api/callback"
	"distasteful-bear/turing_machine/api/db"
	"distasteful-bear/turing_machine/api/leaderboard"
	"distasteful-bear/turing_machine/api/login"
	"distasteful-bear/turing_machine/api/logout"
	"distasteful-bear/turing_machine/api/middleware"
	"distasteful-bear/turing_machine/api/user"
	"distasteful-bear/turing_machine/utils"
	"distasteful-bear/turing_machine/verifiers"
	"encoding/gob"
	"fmt"
	"log"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type SessionToken struct {
	CreationTime   time.Time
	ExpirationTime time.Time
	TokenStr       string
}
type PuzzleWithExpiration struct {
	Puzzle     verifiers.Puzzle
	Expiration time.Time
}
type SessionStore struct {
	ActiveTokens  []SessionToken
	ActivePuzzles map[string]PuzzleWithExpiration
}

func SetupSessionStoreInMem() gin.HandlerFunc {
	store := &SessionStore{
		ActiveTokens:  []SessionToken{},
		ActivePuzzles: map[string]PuzzleWithExpiration{},
	}

	return func(c *gin.Context) {
		c.Set("session_store", store)
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	// Initialize Gin router
	router := gin.Default()

	// To store custom types in our cookies,
	// we must first register them using gob.Register
	gob.Register(map[string]interface{}{})

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("auth-session", store))

	// Setup server-side session store for puzzles
	router.Use(SetupSessionStoreInMem())

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
		fmt.Print("Error creating authenticator session.")
		panic(err)
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

	// puzzle data management
	router.GET("/setup_session", func(c *gin.Context) {
		session := sessions.Default(c)

		// Generate a new puzzle
		puzzle := verifiers.GenerateRandomPuzzle()

		// Get or create session store
		storeInterface, exists := c.Get("session_store")
		if !exists {
			c.JSON(500, gin.H{"error": "session store not initialized"})
			return
		}
		store := storeInterface.(*SessionStore)

		// Generate a unique puzzle ID for this session
		sessionID := uuid.New().String()

		puzzleWithExp := PuzzleWithExpiration{
			Puzzle:     puzzle,
			Expiration: time.Now().Add(time.Hour),
		}
		store.ActivePuzzles[sessionID] = puzzleWithExp

		// Store only the session ID in the session cookie and reset guess count
		session.Set("puzzle_id", sessionID)
		fmt.Println("Puzzle stored with ID:", sessionID) // Testing
		session.Set("guess_count", 0)
		session.Save()

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

		session := sessions.Default(c)

		// Retrieve puzzle ID from session
		puzzleId := session.Get("puzzle_id")
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
		store := storeInterface.(*SessionStore)

		// Retrieve puzzle from server-side storage
		puzzleWithExp, ok := store.ActivePuzzles[puzzleId.(string)]
		if !ok || puzzleWithExp.Expiration.Before(time.Now()) {
			delete(store.ActivePuzzles, puzzleId.(string))
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
		guessCount := session.Get("guess_count")
		if guessCount == nil {
			guessCount = 0
		}
		newGuessCount := guessCount.(int) + 1
		session.Set("guess_count", newGuessCount)
		session.Save()

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
		// session
		session := sessions.Default(c)
		storeInterface, exists := c.Get("session_store")
		if !exists {
			c.JSON(500, gin.H{"error": "session store not initialized"})
			return
		}
		store := storeInterface.(*SessionStore)

		// query params
		puzzleId := session.Get("puzzle_id").(string)
		if puzzleId == "" {
			c.JSON(400, gin.H{"error": "no puzzle session found"})
			return
		}
		guessCount, ok := session.Get("guess_count").(int)
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
		puzzleWithExp, ok := store.ActivePuzzles[puzzleId]
		defer delete(store.ActivePuzzles, puzzleId)
		if !ok || puzzleWithExp.Expiration.Before(time.Now()) {
			delete(store.ActivePuzzles, puzzleId)
			c.JSON(400, gin.H{"error": "puzzle not found in store"})
			return
		}

		puzzle := puzzleWithExp.Puzzle
		success := proposedSolution.Display == puzzle.Sol.Display

		// log results if user is logged in
		userId, err := user.IsUserLoggedIn(c)
		if err != nil {
			err := db.RecordGameCompletion(c.Request.Context(), userId, success, guessCount)
			if err != nil {
				log.Printf("Error recording game completion: %v", err)
				c.JSON(500, gin.H{"error": "failed to record game"})
				return
			}
		}

		// compute leaderboard rankings
		err = db.ComputeLeaderboardRankings(c.Request.Context())
		if err != nil {
			log.Printf("Error computing leaderboard rankings: %v", err)
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

	return router
}
