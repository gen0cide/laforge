package auth

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/authuser"
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

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// Claims Create a struct that will be encoded to a JWT.
type Claims struct {
	IssuedAt int64
	jwt.StandardClaims
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware(client *ent.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		hostname, ok := os.LookupEnv("SERVER_HOSTNAME")
		if !ok {
			hostname = "localhost"
		}

		authCookie, err := ctx.Cookie("auth-cookie")
		if err != nil || authCookie == "" {
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
			// TODO: Change Cookie to be secure
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			if err == jwt.ErrSignatureInvalid {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			// TODO: Change Cookie to be secure
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		entToken, err := client.Token.Query().Where(token.TokenEQ(authCookie)).Only(ctx)
		if err != nil {
			// TODO: Change Cookie to be secure
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		entAuthUser, err := entToken.QueryTokenToAuthUser().Only(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}
		// put it in context
		c := context.WithValue(ctx.Request.Context(), userCtxKey, entAuthUser)
		ctx.Request = ctx.Request.WithContext(c)

		ctx.Next()
	}
}

// Login decodes the share session cookie and packs the session into context
func Login(client *ent.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hostname, ok := os.LookupEnv("SERVER_HOSTNAME")
		if !ok {
			hostname = "localhost"
		}
		var loginVals login
		username := ""
		password := ""
		if err := ctx.ShouldBind(&loginVals); err != nil {
			username = "black"
			password = "black"
		} else {
			username = loginVals.Username
			password = loginVals.Password
		}

		entAuthUser, err := client.AuthUser.Query().Where(
			authuser.And(
				authuser.UsernameEQ(username),
				authuser.PasswordEQ(password),
			),
		).Only(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}

		expiresAt := time.Now().Add(time.Hour * time.Duration(1)).Unix()

		claims := &Claims{
			IssuedAt: time.Now().Unix(),
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: expiresAt,
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error signing token"})
			return
		}

		_, err = client.Token.Create().SetTokenToAuthUser(entAuthUser).SetCreatedAt(int(expiresAt)).SetToken(tokenString).Save(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error updating token"})
			return
		}

		// TODO: Change Cookie to be secure
		ctx.SetCookie("auth-cookie", tokenString, 60*60, "/", hostname, false, false)

		ctx.Next()
	}
}

// Logout decodes the share session cookie and packs the session into context
func Logout(client *ent.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hostname, ok := os.LookupEnv("SERVER_HOSTNAME")
		if !ok {
			hostname = "localhost"
		}

		authCookie, err := ctx.Cookie("auth-cookie")

		// Allow unauthenticated users in
		if err != nil || authCookie == "" {
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
			// TODO: Change Cookie to be secure
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			if err == jwt.ErrSignatureInvalid {
				ctx.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			ctx.AbortWithStatus(http.StatusBadRequest)
			return
		}
		if !tkn.Valid {
			// TODO: Change Cookie to be secure
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		_, err = client.Token.Delete().Where(token.TokenEQ(authCookie)).Exec(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		}

		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error updating token"})
			return
		}

		ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)

		ctx.Next()
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *ent.AuthUser {
	raw, _ := ctx.Value(userCtxKey).(*ent.AuthUser)
	return raw
}
