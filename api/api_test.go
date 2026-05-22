package api

import (
	"distasteful-bear/turing_machine/api/session"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"sync"
	"testing"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func newPuzzleTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(sessions.Sessions("auth-session", cookie.NewStore([]byte("test-secret"))))
	router.Use(session.SetupSessionStoreInMem())
	setupPuzzleRoutes(router)

	return router
}

func TestSetupSessionLoadDoesNotCrash(t *testing.T) {
	requests := 100
	if os.Getenv("RUN_10M_PUZZLE_TEST") == "1" {
		requests = 10_000_000
	}
	if override := os.Getenv("PUZZLE_LOAD_REQUESTS"); override != "" {
		parsed, err := strconv.Atoi(override)
		if err != nil || parsed < 1 {
			t.Fatalf("invalid PUZZLE_LOAD_REQUESTS value %q", override)
		}
		requests = parsed
	}

	router := newPuzzleTestRouter()

	for i := range requests {
		req := httptest.NewRequest(http.MethodGet, "/setup_session", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("request %d/%d returned status %d: %s", i+1, requests, rec.Code, rec.Body.String())
		}
	}
}

func TestSetupSessionConcurrentDoesNotCrash(t *testing.T) {
	router := newPuzzleTestRouter()

	const requests = 100
	var wg sync.WaitGroup
	errs := make(chan string, requests)

	for range requests {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req := httptest.NewRequest(http.MethodGet, "/setup_session", nil)
			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			if rec.Code != http.StatusOK {
				errs <- rec.Body.String()
			}
		}()
	}

	wg.Wait()
	close(errs)

	for err := range errs {
		t.Fatalf("concurrent setup_session returned non-200 response: %s", err)
	}
}
