package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"main/lib"
	"main/models"
	"main/routes"

	"github.com/getsentry/sentry-go"
	sentryecho "github.com/getsentry/sentry-go/echo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	if err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://102afa6541f7ad3d029d698b1ace97cc@o4507617649426432.ingest.us.sentry.io/4507617650999296",
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	// jst, err := time.LoadLocation("Asia/Tokyo")
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// ns, err := gocron.NewScheduler(gocron.WithLocation(jst))
	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// ns.Start()

	// nj, err := ns.NewJob(
	// 	gocron.DurationJob(1*time.Minute),
	// 	gocron.NewTask(ScheduleGetStatus),
	// )

	// if err != nil {
	// 	log.Fatal(err)
	// 	return
	// }

	// Create Instance
	e := echo.New()

	// Load Database
	db, err := lib.Setup()
	if err != nil {
		e.Logger.Fatal(err)
	}

	account_handler := routes.NewAccountHandler(db)

	// Settings Middleware
	e.Use(middleware.CORS())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(sentryecho.New(sentryecho.Options{}))

	// Generate Authentication token for KeyAuth
	username := os.Getenv("AUTH_BASIC_NAME")
	password := os.Getenv("AUTH_BASIC_PASSWORD")
	identifier := os.Getenv("AUTH_BASIC_IDENTIFIER")
	secretKey := os.Getenv("AUTH_BASIC_SECRET_KEY")

	if username == "" || password == "" || identifier == "" || secretKey == "" {
		e.Logger.Fatal("AUTH_BASIC_NAME, AUTH_BASIC_PASSWORD, AUTH_BASIC_IDENTIFIER and AUTH_BASIC_SECRET_KEY environment variables must be set")
	}

	encryptedAuth, err := lib.CompressEncrypt(
		fmt.Sprintf("%s:%s:%s", username, password, secretKey),
		secretKey,
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Generated Basic Authentication key: " + encryptedAuth)

	// Set Route
	e.GET("/", routes.Get_help)

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == 404 {
				c.JSON(http.StatusNotFound, models.Error{Reason: "access to unknown endpoint"})
			}
		}
	}

	v2 := e.Group("/v2")
	v2.GET("/", routes.Get_help)
	v2.GET("/help", routes.Get_help)
	v2.POST("/contact", routes.Create_contact)

	// API endpoints related to the account
	accounts := v2.Group("/accounts")
	accounts.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup:  "header:" + echo.HeaderAuthorization,
		AuthScheme: "Basic",
		Validator: func(auth string, c echo.Context) (bool, error) {
			decryptedAuth, _ := lib.DecryptDecompress(auth, secretKey)

			return decryptedAuth == fmt.Sprintf("%s:%s:%s", username, password, secretKey), nil
		},
	}))
	accounts.POST("/", account_handler.CreateAccount)
	accounts.GET("/:accountId", account_handler.GetAccount)
	accounts.PUT("/:accountId", account_handler.UpdateAccount)
	accounts.GET("/:accountId/history", account_handler.GetHistory)
	accounts.POST("/:accountId/history", account_handler.CreateHistory)
	accounts.PUT("/:accountId/history", account_handler.UpdateHistory)
	accounts.POST("/:accountId/deposit", account_handler.Deposit)
	accounts.POST("/:accountId/withdraw", account_handler.Withdraw)

	// Start Server
	// fmt.Printf("job ID: %v\n", nj.ID())
	e.Logger.Fatal(e.Start(":8080"))
}
