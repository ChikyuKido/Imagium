package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"imagu/db/model"
	"imagu/db/repo"
	"imagu/util"
	"net/http"
)

// AuthMiddleware is a middleware function that handles authentication.
// It checks for a JWT token in the "Authorization" header. If the token is not present
// or is "guest", it retrieves and sets a "guest" user. Otherwise, it validates the token
// and sets the authenticated user in the context.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, _ := c.Cookie("jwt")
		// guest login
		if tokenString == "" || tokenString == "guest" {
			guest, err := repo.GetUserByName("guest")
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "No guest user found."})
				c.Abort()
				return
			}
			c.Set("user", guest)
			c.Next()
			return
		}

		token, err := util.GetToken(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		userName := claims["username"].(string)
		user, err := repo.GetUserByName(userName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

// AuthPermission is a middleware function that checks if the authenticated user
// has a specific permission. If the user does not have the required permission, it
// either redirects them to the login page or responds with a 403 Forbidden status.

func AuthPermission(permission string, redirect bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
			c.Abort()
			return
		}
		u, ok := user.(*model.User)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user type"})
			c.Abort()
			return
		}
		if repo.HasRole(u, permission) {
			c.Next()
		} else {
			if redirect {
				c.Redirect(http.StatusPermanentRedirect, "/login")
			} else {
				c.JSON(http.StatusForbidden, gin.H{"message": "You do not have permission to access this site"})
			}
			c.Abort()
		}
	}
}
