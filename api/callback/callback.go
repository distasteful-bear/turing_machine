package callback

import (
	"distasteful-bear/turing_machine/api/authenticator"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Handler for our callback.
func Handler(auth *authenticator.Authenticator) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if auth == nil {
			ctx.Redirect(http.StatusTemporaryRedirect, "/login")
			return
		}

		session := sessions.Default(ctx)
		if oauthErr := ctx.Query("error"); oauthErr != "" {
			log.Printf("auth callback returned error: %s: %s", oauthErr, ctx.Query("error_description"))
			ctx.String(http.StatusUnauthorized, "Authorization failed.")
			return
		}

		sessionState, ok := session.Get("state").(string)
		if !ok || ctx.Query("state") != sessionState {
			ctx.String(http.StatusBadRequest, "Invalid state parameter.")
			return
		}

		code := ctx.Query("code")
		if code == "" {
			ctx.String(http.StatusBadRequest, "Missing authorization code.")
			return
		}

		// Exchange an authorization code for a token.
		token, err := auth.Exchange(ctx.Request.Context(), code)
		if err != nil {
			log.Printf("failed to exchange authorization code for token: %v", err)
			ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
			return
		}

		idToken, err := auth.VerifyIDToken(ctx.Request.Context(), token)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
			return
		}

		var profile map[string]interface{}
		if err := idToken.Claims(&profile); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		session.Set("access_token", token.AccessToken)
		session.Set("profile", profile)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Redirect to logged in page.
		ctx.Redirect(http.StatusTemporaryRedirect, "/user")
	}
}
