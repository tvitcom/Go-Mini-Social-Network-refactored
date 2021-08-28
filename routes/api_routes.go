package routes

import (
	"my.localhost/funny/Go-Mini-Social-Network-refactored/config"
	"os"
	"strings"
	"time"
	"github.com/badoux/checkmail"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
)

// CreateNewPost route
func CreateNewPost(c *gin.Context) {

	title := strings.TrimSpace(c.PostForm("title"))
	content := strings.TrimSpace(c.PostForm("content"))
	id, _ := config.SessionsUserinfo(c)

	db := config.DB()

	stmt, _ := db.Prepare("INSERT INTO posts(title, content, createdBy, createdAt) VALUES (?, ?, ?, ?)")
	rs, iErr := stmt.Exec(title, content, id, time.Now())
	config.Err(iErr)

	insertID, _ := rs.LastInsertId()

	resp := map[string]interface{}{
		"postID": insertID,
		"mssg":   "Post Created!!",
	}
	json(c, resp)
}

// DeletePost route
func DeletePost(c *gin.Context) {
	post := c.PostForm("post")
	db := config.DB()

	_, dErr := db.Exec("DELETE FROM posts WHERE postID=?", post)
	config.Err(dErr)

	json(c, map[string]interface{}{
		"mssg": "Post Deleted!!",
	})
}

// UpdatePost route
func UpdatePost(c *gin.Context) {
	postID := c.PostForm("postID")
	title := c.PostForm("title")
	content := c.PostForm("content")

	db := config.DB()
	db.Exec("UPDATE posts SET title=?, content=? WHERE postID=?", title, content, postID)

	json(c, map[string]interface{}{
		"mssg": "Post Updated!!",
	})
}

// UpdateProfile route
func UpdateProfile(c *gin.Context) {
	resp := make(map[string]interface{})

	id, _ := config.SessionsUserinfo(c)
	username := strings.TrimSpace(c.PostForm("username"))
	email := strings.TrimSpace(c.PostForm("email"))
	bio := strings.TrimSpace(c.PostForm("bio"))

	mailErr := checkmail.ValidateFormat(email)
	db := config.DB()

	if username == "" || email == "" {
		resp["mssg"] = "Some values are missing!!"
	} else if mailErr != nil {
		resp["mssg"] = "Invalid email format!!"
	} else {
		_, iErr := db.Exec("UPDATE users SET username=?, email=?, bio=? WHERE id=?", username, email, bio, id)
		config.Err(iErr)

		session := sessions.Default(c)
		session.Set("username", username)
		session.Save()

		resp["mssg"] = "Profile updated!!"
		resp["success"] = true
	}

	json(c, resp)
}

// ChangeAvatar route
func ChangeAvatar(c *gin.Context) {
	resp := make(map[string]interface{})
	id, _ := config.SessionsUserinfo(c)

	dir, _ := os.Getwd()
	users_dir := dir + "/public/users/" + id.(string)
	dest := users_dir + "/avatar.png"

  	// Make user dir if not exist
  	if _, errPath := os.Stat(users_dir); errPath != nil {
		config.Err(os.Mkdir(users_dir, 0777))
  	}
	if _, err := os.Stat(dest); err == nil {
		dErr := os.Remove(dest)
		config.Err(dErr)
  	}

	file, _ := c.FormFile("avatar")
	upErr := c.SaveUploadedFile(file, dest)

	if upErr != nil {
		resp["mssg"] = "An error occured!!"
	} else {
		resp["mssg"] = "Avatar changed!!"
		resp["success"] = true
	}

	json(c, resp)
}

// Follow route
func Follow(c *gin.Context) {
	id, _ := config.SessionsUserinfo(c)
	user := c.PostForm("user")
	username := config.Get(user, "username")

	db := config.DB()
	stmt, _ := db.Prepare("INSERT INTO follow(followBy, followTo, followTime) VALUES(?, ?, ?)")
	_, exErr := stmt.Exec(id, user, time.Now())
	config.Err(exErr)

	json(c, gin.H{
		"mssg": "Followed " + username + "!!",
	})
}

// Unfollow route
func Unfollow(c *gin.Context) {
	id, _ := config.SessionsUserinfo(c)
	user := c.PostForm("user")
	username := config.Get(user, "username")

	db := config.DB()
	stmt, _ := db.Prepare("DELETE FROM follow WHERE followBy=? AND followTo=?")
	_, dErr := stmt.Exec(id, user)
	config.Err(dErr)

	json(c, gin.H{
		"mssg": "Unfollowed " + username + "!!",
	})
}

// Like post route
func Like(c *gin.Context) {
	post := c.PostForm("post")
	db := config.DB()
	id, _ := config.SessionsUserinfo(c)

	stmt, _ := db.Prepare("INSERT INTO likes(postID, likeBy, likeTime) VALUES (?, ?, ?)")
	_, err := stmt.Exec(post, id, time.Now())
	config.Err(err)

	json(c, gin.H{
		"mssg": "Post Liked!!",
	})
}

// Unlike post route
func Unlike(c *gin.Context) {
	post := c.PostForm("post")
	id, _ := config.SessionsUserinfo(c)
	db := config.DB()

	stmt, _ := db.Prepare("DELETE FROM likes WHERE postID=? AND likeBy=?")
	_, err := stmt.Exec(post, id)
	config.Err(err)

	json(c, gin.H{
		"mssg": "Post Unliked!!",
	})
}

// DeactivateAcc route post method
func DeactivateAcc(c *gin.Context) {
	session := sessions.Default(c)
	id, _ := config.SessionsUserinfo(c)
	db := config.DB()
	var postID int

	db.Exec("DELETE FROM profile_views WHERE viewBy=?", id)
	db.Exec("DELETE FROM profile_views WHERE viewTo=?", id)
	db.Exec("DELETE FROM follow WHERE followBy=?", id)
	db.Exec("DELETE FROM follow WHERE followTo=?", id)
	db.Exec("DELETE FROM likes WHERE likeBy=?", id)

	rows, _ := db.Query("SELECT postID FROM posts WHERE createdBy=?", id)
	for rows.Next() {
		rows.Scan(&postID)
		db.Exec("DELETE FROM likes WHERE postID=?", postID)
	}

	db.Exec("DELETE FROM posts WHERE createdBy=?", id)
	db.Exec("DELETE FROM users WHERE id=?", id)

	dir, _ := os.Getwd()
	userPath := dir + "/public/users/" + id.(string)

	rmErr := os.RemoveAll(userPath)
	config.Err(rmErr)

	session.Delete("id")
	session.Delete("username")
	session.Save()

	json(c, gin.H{
		"mssg": "Deactivated your account!!",
	})
}
