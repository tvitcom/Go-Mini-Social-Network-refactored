package main

import (
	R "my.localhost/funny/Go-Mini-Social-Network-refactored/routes"
	"os"

	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-contrib/sessions"

	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/urfave/negroni"
)

func init() {
	godotenv.Load()
}

func main() {
	router := gin.Default()

	// Session store init
	sessionStore := cookie.NewStore([]byte("secretu"), []byte("conf.KEY32123412341234123412dfrt"))
	sess := sessions.Sessions("bin", sessionStore)
	router.Use(sess) //Название ключа в куках

	router.LoadHTMLGlob("views/*.html")

	auth := router.Group("/user")
	{
		auth.POST("/signup", R.UserSignup)
		auth.POST("/login", R.UserLogin)
	}

	router.GET("/", R.Index)
	router.GET("/welcome", R.Welcome)
	router.GET("/explore", R.Explore)
	router.GET("/404", R.NotFound)
	router.GET("/signup", R.Signup)
	router.GET("/login", R.Login)
	router.GET("/logout", R.Logout)
	router.GET("/deactivate", R.Deactivate)
	router.GET("/edit_profile", R.EditProfile)
	router.GET("/create_post", R.CreatePost)

	router.GET("/profile/:id", R.Profile)
	router.GET("/profile", R.NotFound)

	router.GET("/view_post/:id", R.ViewPost)
	router.GET("/view_post", R.NotFound)

	router.GET("/edit_post/:id", R.EditPost)
	router.GET("/edit_post", R.NotFound)

	router.GET("/followers/:id", R.Followers)
	router.GET("/followers", R.NotFound)

	router.GET("/followings/:id", R.Followings)
	router.GET("/followings", R.NotFound)

	router.GET("/likes/:id", R.Likes)
	router.GET("/likes", R.NotFound)

	api := router.Group("/api")
	{
		api.POST("/create_new_post", R.CreateNewPost)
		api.POST("/delete_post", R.DeletePost)
		api.POST("/update_post", R.UpdatePost)
		api.POST("/update_profile", R.UpdateProfile)
		api.POST("/change_avatar", R.ChangeAvatar)
		api.POST("/follow", R.Follow)
		api.POST("/unfollow", R.Unfollow)
		api.POST("/like", R.Like)
		api.POST("/unlike", R.Unlike)
		api.POST("/deactivate-account", R.DeactivateAcc)
	}

	server := negroni.Classic()
	server.UseHandler(router)

	server.Run(os.Getenv("PORT"))

}
