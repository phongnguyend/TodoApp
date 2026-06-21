module github.com/todo/backend/go

go 1.23

require (
	github.com/gin-gonic/gin v1.10.0
	github.com/joho/godotenv v1.5.1
	github.com/swaggo/files v1.0.1
	github.com/swaggo/gin-swagger v1.6.0
	github.com/swaggo/swag v1.16.4
	gorm.io/driver/sqlite v1.5.7
	gorm.io/gorm v1.25.12
)

// To add PostgreSQL support, also run:
//   go get gorm.io/driver/postgres
// To add MySQL support:
//   go get gorm.io/driver/mysql
