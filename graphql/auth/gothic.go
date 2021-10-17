package auth

import (
	"context"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/authuser"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gitlab"
	"github.com/markbates/goth/providers/openidConnect"
	"github.com/markbates/goth/providers/slack"
	"github.com/sirupsen/logrus"
)

func InitGoth() {
	url_start := "http://"
	if env_value, exists := os.LookupEnv("HTTPS_ENABLED"); exists {
		if env_value == "true" {
			url_start = "https://"
		}
	}
	goth.UseProviders(
		github.New(os.Getenv("GITHUB_ID"), os.Getenv("GITHUB_SECRET"), url_start+os.Getenv("GRAPHQL_HOSTNAME")+"/auth/github/callback"),
		slack.New(os.Getenv("SLACK_KEY"), os.Getenv("SLACK_SECRET"), url_start+os.Getenv("GRAPHQL_HOSTNAME")+"/auth/slack/callback"),
		gitlab.New(os.Getenv("GITLAB_KEY"), os.Getenv("GITLAB_SECRET"), url_start+os.Getenv("GRAPHQL_HOSTNAME")+"/auth/gitlab/callback"),
	)

	// OpenID Connect is based on OpenID Connect Auto Discovery URL (https://openid.net/specs/openid-connect-discovery-1_0-17.html)
	// because the OpenID Connect provider initialize it self in the New(), it can return an error which should be handled or ignored
	// ignore the error for now
	openidConnect, _ := openidConnect.New(os.Getenv("OPENID_CONNECT_KEY"), os.Getenv("OPENID_CONNECT_SECRET"), url_start+os.Getenv("GRAPHQL_HOSTNAME")+"/auth/openid-connect/callback", os.Getenv("OPENID_CONNECT_DISCOVERY_URL"))
	if openidConnect != nil {
		goth.UseProviders(openidConnect)
		// FIX 3
		// The standard Goth SessionStore is not big enough to hold all the data from the OpenID connect provider
		store := sessions.NewFilesystemStore(os.TempDir(), []byte("openvpn-management"))

		// set the maxLength of the cookies stored on the disk to a larger number to prevent issues with:
		// securecookie: the value is too long
		// when using OpenID Connect , since this can contain a large amount of extra information in the id_token

		// Note, when using the FilesystemStore only the session.ID is written to a browser cookie, so this is explicit for the storage on disk
		// See: https://github.com/markbates/goth/issues/133
		store.MaxLength(math.MaxInt64)

		gothic.Store = store
	}

}

func contextWithProviderName(ctx *gin.Context, provider string) *http.Request {
	return ctx.Request.WithContext(context.WithValue(ctx.Request.Context(), "provider", provider))
}

// GothicBeginAuth decodes the share session cookie and packs the session into context
func GothicBeginAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		provider := ctx.Param("provider")

		// You have to add value context with provider name to get provider name in GetProviderName method
		ctx.Request = contextWithProviderName(ctx, provider)
		gothic.BeginAuthHandler(ctx.Writer, ctx.Request)
	}
}

// GothicCallbackHandler decodes the share session cookie and packs the session into context
func GothicCallbackHandler(client *ent.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hostname, ok := os.LookupEnv("GRAPHQL_HOSTNAME")
		if !ok {
			hostname = "localhost"
		}
		cookie_timeout := 60
		if env_value, exists := os.LookupEnv("COOKIE_TIMEOUT"); exists {
			if atio_value, err := strconv.Atoi(env_value); err == nil {
				cookie_timeout = atio_value
			}
		}
		secure_cookie := false
		if env_value, exists := os.LookupEnv("HTTPS_ENABLED"); exists {
			if env_value == "true" {
				secure_cookie = true
			}
		}
		// handle with gothic
		info, err := gothic.CompleteUserAuth(ctx.Writer, ctx.Request)
		if err != nil {
			logrus.Error(err)
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		// we're not using gothic for auth management, so end the user session immediately
		gothic.Logout(ctx.Writer, ctx.Request)

		username := info.NickName
		provider := strings.ToUpper(info.Provider)

		if provider == "OPENID-CONNECT" {
			provider = "OPENID"
		}

		entAuthUser, err := client.AuthUser.Query().Where(
			authuser.And(
				authuser.UsernameEQ(username),
				authuser.ProviderEQ(authuser.Provider(provider)),
			),
		).Only(ctx)

		if err != nil {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		expiresAt := time.Now().Add(time.Minute * time.Duration(cookie_timeout)).Unix()

		claims := &Claims{
			IssuedAt: time.Now().Unix(),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expiresAt,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error signing token"})
			return
		}

		_, err = client.Token.Create().SetTokenToAuthUser(entAuthUser).SetExpireAt(expiresAt).SetToken(tokenString).Save(ctx)
		if err != nil {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error updating token"})
			return
		}

		if secure_cookie {
			ctx.SetCookie("auth-cookie", tokenString, cookie_timeout*60, "/", hostname, true, true)
		} else {
			ctx.SetCookie("auth-cookie", tokenString, cookie_timeout*60, "/", hostname, false, false)
		}
		ctx.Redirect(http.StatusFound, "/")
		// ctx.JSON(200, entAuthUser)
	}
}
