package main

import (
	"fmt"
	"log"
	"os"

	"main/lib"
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
		Dsn:              os.Getenv("SENTRY_URL"),
		TracesSampleRate: 1.0,
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	// Create Instance
	e := echo.New()

	// Load Database
	db, mongo, err := lib.Setup()
	if err != nil {
		e.Logger.Fatal(err)
	}

	items_handler := routes.NewItemsHandler(mongo)
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

	accounts.POST("", account_handler.CreateAccount)
	accounts.GET("/:accountId", account_handler.GetAccount)
	accounts.PUT("/:accountId", account_handler.UpdateAccount)
	accounts.GET("/:accountId/history", account_handler.GetHistory)
	accounts.POST("/:accountId/history", account_handler.CreateHistory)
	accounts.PUT("/:accountId/history", account_handler.UpdateHistory)
	accounts.POST("/:accountId/deposit", account_handler.Deposit)
	accounts.POST("/:accountId/withdraw", account_handler.Withdraw)

	// Items API
	items := v2.Group("/items")

	items.POST("", items_handler.CreateItem)
	items.GET("", items_handler.GetItems)
	items.GET("/:id", items_handler.GetItem)
	items.PUT("/:id", items_handler.UpdateItem)
	items.DELETE("/:id", items_handler.DeleteItem)

	// Start Server
	e.Logger.Fatal(e.Start(":8080"))
}
