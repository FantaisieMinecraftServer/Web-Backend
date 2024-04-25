package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"main/lib"
	"main/routes"

	"github.com/go-co-op/gocron/v2"
	_ "github.com/go-sql-driver/mysql"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		log.Fatal(err)
		return
	}

	ns, err := gocron.NewScheduler(gocron.WithLocation(jst))
	if err != nil {
		log.Fatal(err)
		return
	}

	ns.Start()

	nj, err := ns.NewJob(
		gocron.DurationJob(1*time.Minute),
		gocron.NewTask(ScheduleGetStatus),
	)

	if err != nil {
		log.Fatal(err)
		return
	}

	// Create Instance
	e := echo.New()

	// Settings Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Set Route
	e.GET("/", routes.GetHelp)
	e.GET("/v2/", routes.GetHelp)
	e.GET("/v2/help", routes.GetHelp)
	e.GET("/v2/status", routes.GetStatus)

	// Start Server
	fmt.Printf("job ID: %v\n", nj.ID())
	e.Logger.Fatal(e.Start(":6000"))
}

func ScheduleGetStatus() {
	address := "play.tensyoserver.net"

	db := lib.GetDBConnection()
	defer db.Close()

	res_1, _ := lib.GetStatusData(address, "25565")
	res_2, _ := lib.GetStatusData(address, "25566")
	res_3, _ := lib.GetStatusData(address, "25570")
	res_4, _ := lib.GetStatusData(address, "25567")
	res_5, _ := lib.GetStatusData(address, "25568")

	ins, err := db.Prepare("INSERT INTO status(date,proxy,lobby,survival,minigame,pve) VALUES(?,?,?,?,?,?)")

	if err != nil {
		log.Fatal(err)
	}

	data_1, _ := json.Marshal(res_1)
	data_2, _ := json.Marshal(res_2)
	data_3, _ := json.Marshal(res_3)
	data_4, _ := json.Marshal(res_4)
	data_5, _ := json.Marshal(res_5)

	ins.Exec(time.Now(), string(data_1), string(data_2), string(data_3), string(data_4), string(data_5))
}
