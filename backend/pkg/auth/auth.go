//go:build web

package auth

import (
	"encoding/base64"

	"github.com/mjiee/world-news/backend/pkg/errorx"
	"github.com/mjiee/world-news/backend/pkg/httpx"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
)

// BasicAuth returns a Basic HTTP Authorization middleware. It takes as argument a map[string]string where
func BasicAuth(accounts gin.Accounts) gin.HandlerFunc {
	authPairs := make([]string, len(accounts), len(accounts))

	for user, password := range accounts {
		authPairs = append(authPairs, authData(user, password))
	}

	return func(c *gin.Context) {
		var (
			auth   = c.GetHeader("Authorization")
			unAuth = true
		)

		if auth != "" {
			for _, pair := range authPairs {
				if pair == auth {
					unAuth = false
					break
				}
			}
		}

		if unAuth {
			httpx.WebResp(c, nil, errorx.Unauthorized)

			c.Abort()

			return
		}

		c.Next()
	}
}

func authData(user, password string) string {
	base := user + ":" + password
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(base))
}
