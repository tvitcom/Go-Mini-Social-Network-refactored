package config

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)

// AllSessions function to return all the sessions
func SessionsUserinfo(c *gin.Context) (interface{}, interface{}) {
	session := sessions.Default(c)
	id := session.Get("id")
	username := session.Get("username")
	return id, username
}
