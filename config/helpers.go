package config

import (
	"os"
	"io"
)

// THIS FILE CONTAINS ALL THE METHODS WHICH WILL BE USED IN TEMPLATES/VIEWS

// Get function to get anything of user with ID
func Get(id interface{}, what string) string {
	db := DB()
	var RET string
	db.QueryRow("SELECT "+what+" AS RET FROM users WHERE id=?", id).Scan(&RET)
	return RET
}

// IsFollowing route
func IsFollowing(by string, to string) bool {
	db := DB()
	var followCount int
	db.QueryRow("SELECT COUNT(followID) AS followCount FROM follow WHERE followBy=? AND followTo=? LIMIT 1", by, to).Scan(&followCount)
	if followCount == 0 {
		return false
	}
	return true
}

// UsernameDecider Helper
func UsernameDecider(user int, session string) string {
	username := Get(user, "username")
	sesUsername := Get(session, "username")
	if username == sesUsername {
		return "You"
	}
	return username
}

// NoOfFollowers helper
func NoOfFollowers(user int) int {
	db := DB()
	var followersCount int
	db.QueryRow("SELECT COUNT(followID) AS followersCount FROM follow WHERE followTo=?", user).Scan(&followersCount)
	return followersCount
}

// LikedOrNot helper
func LikedOrNot(post int, user interface{}) bool {
	db := DB()
	var likeCount int
	db.QueryRow("SELECT COUNT(likeID) AS likeCount FROM likes WHERE likeBy=? AND postID=?", user, post).Scan(&likeCount)
	if likeCount == 0 {
		return false
	}
	return true
}

// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func Copyfile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
