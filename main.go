package main

import (
	"tele-sticker-finder/config"
	"tele-sticker-finder/db"
	"tele-sticker-finder/migration"
	"tele-sticker-finder/pkg/bot"
)

func main() {
	//e := echo.New()
	//e.GET("/", func(c echo.Context) error {
	//	return c.String(http.StatusOK, "Hello, World!")
	//})
	//
	//e.Logger.Fatal(e.Start(":1323"))

	//pkg.ReadFile(context.Background())

	// Setup Config
	cfg, err := config.NewLoadConfig()
	if err != nil {
		panic(err)
	}

	// Setup DB connection
	dbConn, err := db.NewDatabaseConnection(cfg)
	if err != nil {
		panic(err)
	}

	redisConn := db.NewRedisConnection(cfg)

	migration.AutoMigration(dbConn)

	bot.StartBot(cfg, dbConn, redisConn)
}
