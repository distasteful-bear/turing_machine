package api

import (
	"distasteful-bear/turing_machine/utils"
	"distasteful-bear/turing_machine/verifiers"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

type SessionToken struct {
	CreationTime   time.Time
	ExpirationTime time.Time
	TokenStr       string
}

type SessionStore struct {
	ActiveTokens  []SessionToken
	ActivePuzzles map[string]verifiers.Puzzle
}

func SetupSessionStoreInMem() gin.HandlerFunc {
	store := &SessionStore{
		ActiveTokens:  []SessionToken{},
		ActivePuzzles: map[string]verifiers.Puzzle{},
	}

	return func(c *gin.Context) {
		c.Set("session_store", store)
		c.Next()
	}
}

func SetupRouter() *gin.Engine {
	// Initialize Gin router
	router := gin.Default()
	store := cookie.NewStore([]byte("03g3iq2n4fp2wo23n1pnic9f0422fjuP"))
	router.Use(sessions.Sessions("globalsession", store))

	// Load HTML templates from the html directory
	router.LoadHTMLGlob("html/*.html")

	// HTML
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "welcome.html", nil)
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

	// puzzle data management
	router.GET("/setup_session", func(c *gin.Context) {
		session := sessions.Default(c)
		cache := session.Get("puzzle")
		if cache != nil {
			session.Delete("puzzle")
		}
		// new session
		puzzle := verifiers.GenerateRandomPuzzle()
		session.Set("puzzle", puzzle)
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
		cache := session.Get("puzzle")

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
		if cache == nil {
			c.JSON(400, gin.H{"error": "no puzzle cached"})
			return
		}
		puzzle, ok := cache.(verifiers.Puzzle)
		if !ok {
			c.JSON(400, gin.H{"error": "invalid cached puzzle"})
			return
		}
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
		c.JSON(200, gin.H{"status": "success", "puzzle_summary": puzzleSummary})
	})
	router.GET("/check_final", func(c *gin.Context) {

		session := sessions.Default(c)
		cache := session.Get("puzzle")

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

		if cache == nil {
			c.JSON(400, gin.H{"error": "no puzzle cached"})
			return
		}

		puzzle, ok := cache.(verifiers.Puzzle)
		if !ok {
			c.JSON(400, gin.H{"error": "invalid cached puzzle"})
			return
		}

		if proposedSolution.Display == puzzle.Sol.Display {
			c.JSON(200, gin.H{"status": "success"})
			return
		} else {
			c.JSON(200, gin.H{"status": "failure"})
			return
		}
	})

	return router
}
