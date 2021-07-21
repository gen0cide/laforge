package auth

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gen0cide/laforge/ent"
	"github.com/gen0cide/laforge/ent/authuser"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Username string `form:"username" json:"username" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

// LocalLogin decodes the share session cookie and packs the session into context
func LocalLogin(client *ent.Client) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hostname, ok := os.LookupEnv("GRAPHQL_HOSTNAME")
		if !ok {
			hostname = "localhost"
		}
		var loginVals login
		username := ""
		password := ""

		// TODO: Remove test login
		if err := ctx.ShouldBind(&loginVals); err != nil {
			// TODO: Change Cookie to be secure
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err})
			return
		} else {
			username = loginVals.Username
			password = loginVals.Password
		}

		entAuthUser, err := client.AuthUser.Query().Where(
			authuser.And(
				authuser.UsernameEQ(username),
				authuser.ProviderEQ(authuser.ProviderLOCAL),
			),
		).Only(ctx)

		if err != nil {
			// TODO: Change Cookie to be secure
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			return
		}
		// Compare the stored hashed password, with the hashed version of the password that was received
		if err = bcrypt.CompareHashAndPassword([]byte(entAuthUser.Password), []byte(password)); err != nil {
			// TODO: Change Cookie to be secure
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
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
			// TODO: Change Cookie to be secure
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error signing token"})
			return
		}

		_, err = client.Token.Create().SetTokenToAuthUser(entAuthUser).SetExpireAt(expiresAt).SetToken(tokenString).Save(ctx)
		if err != nil {
			// TODO: Change Cookie to be secure
			ctx.SetCookie("auth-cookie", "", 0, "/", hostname, false, false)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Error updating token"})
			return
		}

		// TODO: Change Cookie to be secure
		ctx.SetCookie("auth-cookie", tokenString, 60*60, "/", hostname, false, false)
		// Hide password so no leeks
		entAuthUser.Password = ""
		ctx.JSON(200, entAuthUser)

		ctx.Next()
	}
}
