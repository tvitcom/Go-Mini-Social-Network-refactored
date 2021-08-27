module my.localhost/funny/Go-Mini-Social-Network-template

go 1.15

replace my.localhost/funny/Go-Mini-Social-Network-template/config => ./config

replace my.localhost/funny/Go-Mini-Social-Network-template/routes => ./routes

require (
	github.com/badoux/checkmail v1.2.1
	github.com/gin-contrib/sessions v0.0.3
	github.com/gin-gonic/gin v1.7.4
	github.com/go-sql-driver/mysql v1.6.0
	github.com/joho/godotenv v1.3.0
	github.com/urfave/negroni v1.0.0
	golang.org/x/crypto v0.0.0-20210817164053-32db794688a5
)
