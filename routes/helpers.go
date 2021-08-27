package routes

import (
	"fmt"
	CO "my.localhost/funny/Go-Mini-Social-Network-refactored/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func hash(password string) []byte {
	hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	CO.Err(hashErr)
	fmt.Println("DEBUG:hash:", string(hash))
	return hash
}

func renderTemplate(c *gin.Context, tmpl string, p interface{}) {
	c.HTML(http.StatusOK, tmpl+".html", p)
}

func json(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, data)
}

func ses(c *gin.Context) interface{} {
	id, username := CO.SessionsUserinfo(c)
	return map[string]interface{}{
		"id":       id,
		"username": username,
	}
}

func loggedIn(c *gin.Context, urlRedirect string) {
	var URL string
	if urlRedirect == "" {
		URL = "/login"
	} else {
		URL = urlRedirect
	}
	id, _ := CO.SessionsUserinfo(c)
	if id == nil {
		c.Redirect(http.StatusFound, URL)
	}
}

func notLoggedIn(c *gin.Context) {
	id, _ := CO.SessionsUserinfo(c)
	if id != nil {
		c.Redirect(http.StatusFound, "/")
	}
}

func invalid(c *gin.Context, what int) {
	if what == 0 {
		c.Redirect(http.StatusNotFound, "/404")
	}
}
