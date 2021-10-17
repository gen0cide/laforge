package auth

import (
	"context"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/token"
	"github.com/gin-gonic/gin"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{"user"}
var jwtKey = []byte("TReZCts6dZgXF6PJYHZ8jrdunbFquYnU9FJ6FDgoVGkdMPpvUc")

type contextKey struct {
	name string
}

// Claims Create a struct that will be encoded to a JWT.
type Claims struct {
	IssuedAt int64
	jwt.StandardClaims
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware(client *ent.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		hostname, ok := os.LookupEnv("GRAPHQL_HOSTNAME")
		if !ok {
			hostname = "localhost"
		}
		secure_cookie := false
		if env_value, exists := os.LookupEnv("HTTPS_ENABLED"); exists {
			if env_value == "true" {
				secure_cookie = true
			}
		}

		authCookie, err := ctx.Cookie("auth-cookie")
		if err != nil || authCookie == "" {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			return
		}

		// Get the JWT string from the cookie
		tknStr := authCookie

		// Initialize a new instance of `Claims`
		claims := &Claims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			if err == jwt.ErrSignatureInvalid {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		entToken, err := client.Token.Query().Where(token.TokenEQ(authCookie)).Only(ctx)
		if err != nil {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		entAuthUser, err := entToken.QueryTokenToAuthUser().Only(ctx)
		if err != nil {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		// put it in context
		c := context.WithValue(ctx.Request.Context(), userCtxKey, entAuthUser)
		ctx.Request = ctx.Request.WithContext(c)

		ctx.Next()
	}
}

// Logout decodes the share session cookie and packs the session into context
func Logout(client *ent.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hostname, ok := os.LookupEnv("GRAPHQL_HOSTNAME")
		if !ok {
			hostname = "localhost"
		}
		secure_cookie := false
		if env_value, exists := os.LookupEnv("HTTPS_ENABLED"); exists {
			if env_value == "true" {
				secure_cookie = true
			}
		}

		authCookie, err := ctx.Cookie("auth-cookie")

		// Allow unauthenticated users in
		if err != nil || authCookie == "" {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			return
		}

		// Get the JWT string from the cookie
		tknStr := authCookie

		// Initialize a new instance of `Claims`
		claims := &Claims{}

		// Parse the JWT string and store the result in `claims`.
		// Note that we are passing the key in this method as well. This method will return an error
		// if the token is invalid (if it has expired according to the expiry time we set on sign in),
		// or if the signature does not match
		tkn, err := jwt.ParseWithClaims(tknStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			if err == jwt.ErrSignatureInvalid {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		_, err = client.Token.Delete().Where(token.TokenEQ(authCookie)).Exec(ctx)
		if err != nil {
			if secure_cookie {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
			} else {
				ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			}
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

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
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, true, true)
		} else {
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
		}

		ctx.Next()
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) (*ent.AuthUser, error) {
	raw, ok := ctx.Value(userCtxKey).(*ent.AuthUser)
	if ok {
		return raw, nil
	}
	return nil, errors.New("unable to get authuser from context")

}

// ClearTokens Clears Old tokens from DB
func ClearTokens(client *ent.Client, ctx context.Context) {
	client.Token.Delete().Where(token.ExpireAtLT(time.Now().Unix())).Exec(ctx)
}
