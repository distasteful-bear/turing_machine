package session

import (
	"distasteful-bear/turing_machine/verifiers"
	"errors"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
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
	mu            sync.Mutex
	lastPrune     time.Time
}

const (
	maxActivePuzzles = 1_000
	puzzleTTL        = time.Hour
	pruneInterval    = time.Minute
)

func SetupSessionStoreInMem() gin.HandlerFunc {
	store := &SessionStore{
		ActiveTokens:  []SessionToken{},
		ActivePuzzles: map[string]PuzzleWithExpiration{},
		lastPrune:     time.Now(),
	}

	return func(c *gin.Context) {
		c.Set("session_store", store)
		c.Next()
	}
}

func (s *SessionStore) AddPuzzle(id string, puzzle verifiers.Puzzle, now time.Time) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if now.Sub(s.lastPrune) >= pruneInterval {
		s.pruneExpiredLocked(now)
		s.lastPrune = now
	}

	if len(s.ActivePuzzles) >= maxActivePuzzles {
		return errors.New("too many active games")
	}

	s.ActivePuzzles[id] = PuzzleWithExpiration{
		Puzzle:     puzzle,
		Expiration: now.Add(puzzleTTL),
	}
	return nil
}

func (s *SessionStore) GetPuzzle(id string, now time.Time) (verifiers.Puzzle, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	puzzleWithExp, ok := s.ActivePuzzles[id]
	if !ok {
		return verifiers.Puzzle{}, false
	}
	if puzzleWithExp.Expiration.Before(now) {
		delete(s.ActivePuzzles, id)
		return verifiers.Puzzle{}, false
	}
	return puzzleWithExp.Puzzle, true
}

func (s *SessionStore) DeletePuzzle(id string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.ActivePuzzles, id)
}

func (s *SessionStore) pruneExpiredLocked(now time.Time) {
	for key, p := range s.ActivePuzzles {
		if p.Expiration.Before(now) {
			delete(s.ActivePuzzles, key)
		}
	}
}
