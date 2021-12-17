package main

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/programzheng/language-repository/admin"
	"github.com/programzheng/language-repository/dictionary"
	"github.com/programzheng/language-repository/orm"
	"github.com/programzheng/language-repository/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/gofiber/template/html"
	_ "github.com/joho/godotenv/autoload"
)

var (
	_, b, _, _ = runtime.Caller(0)
	basepath   = filepath.Dir(b)
)

func getLoggerFile() *os.File {
	//log directory
	path := filepath.Join(basepath, "./logs")
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		//make nested directories
		err = os.MkdirAll(path, 0700)
		if err != nil {
			log.Fatal("create log directory error:", err)
		}
	}

	t := time.Now().Local()
	date := t.Format("2006-01-02")
	file, err := os.OpenFile("./logs/"+date+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}

	return file
}

func setFileLogger(file *os.File) fiber.Handler {

	return logger.New(logger.Config{
		Format: "[${time}] - ${latency} - ${ip} - ${status} - [${method}] ${path}\n${body}\n",
		Output: file,
	})
}

func getCors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins: os.Getenv("CORS_ALLOW_ORIGINS"),
		AllowHeaders: "Origin, Content-Type, Accept",
	})
}

func getAdminJwtWare() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_ADMIN_SECRET")),
	})
}

func getUserJwtWare() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: []byte(os.Getenv("JWT_USER_SECRET")),
	})
}

func main() {
	orm.InitDatabase()

	engine := html.New("./dist", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	logFile := getLoggerFile()
	//logger
	app.Use(setFileLogger(logFile))
	defer logFile.Close()

	//cors
	app.Use(getCors())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	apiGroup := app.Group("/api")
	v1Group := apiGroup.Group("/v1")

	adminGroup := v1Group.Group("/admin")
	adminGroup.Post("login", admin.Login)
	adminGroup.Post("", admin.NewAdmin)

	userGroup := v1Group.Group("/user")
	userGroup.Post("login", user.Login)
	userGroup.Use(getAdminJwtWare()).Post("", user.NewUser)

	dictionaryGroup := v1Group.Group("/dictionary")
	dictionaryGroup.Use(getUserJwtWare())
	dictionaryGroup.Get("", dictionary.GetDictionaries)
	dictionaryGroup.Get(":id", dictionary.GetDictionary)
	dictionaryGroup.Post("", dictionary.NewDictionary)
	dictionaryGroup.Put(":id", dictionary.UpdateDictionary)
	dictionaryGroup.Delete(":id", dictionary.DeleteDictionary)

	port := os.Getenv("PORT")
	app.Listen(":" + port)
}
