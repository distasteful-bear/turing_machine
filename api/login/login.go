package login

import (
	"crypto/rand"
	"distasteful-bear/turing_machine/api/authenticator"

	"encoding/base64"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler for our login.
func Handler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if auth == nil {
			session := sessions.Default(ctx)
			session.Set("profile", map[string]interface{}{
				"sub":      "local-dev-user",
				"nickname": "Local Dev",
				"picture":  "",
			})
			if err := session.Save(); err != nil {
				ctx.String(http.StatusInternalServerError, err.Error())
				return
			}
			ctx.Redirect(http.StatusTemporaryRedirect, "/user")
			return
		}

		state, err := generateRandomState()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Save the state inside the session.
		session := sessions.Default(ctx)
		session.Set("state", state)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, auth.AuthCodeURL(state))
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}
