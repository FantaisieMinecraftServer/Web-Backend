package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	mcstatus "main/lib"
	"main/models"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const address = "play.tensyoserver.net"

func getDBConnection() *sql.DB {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := sql.Open("mysql", os.Getenv("DB_URL"))
	if err != nil {
		panic(err.Error())
	}
	return db
}

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
	e.GET("/", getHelp)
	e.GET("/v2/", getHelp)
	e.GET("/v2/help", getHelp)
	e.GET("/v2/status", getStatus)

	// Start Server
	fmt.Printf("job ID: %v\n", nj.ID)
	e.Logger.Fatal(e.Start(":6000"))
}

func ScheduleGetStatus() {
	db := getDBConnection()
	defer db.Close()

	res_1, _ := mcstatus.GetStatusData(address, "25565")
	res_2, _ := mcstatus.GetStatusData(address, "25566")
	res_3, _ := mcstatus.GetStatusData(address, "25570")
	res_4, _ := mcstatus.GetStatusData(address, "25567")
	res_5, _ := mcstatus.GetStatusData(address, "25568")

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

// /v2/help
func getHelp(c echo.Context) error {
	data := &models.APIHelp{
		Version:       "v2.0.0",
		WhatIsThis:    "API for this server, don't ask me how to use it... :(",
		ServerAddress: "play.tensyoserver.net",
		HomePage:      "https://www.tensyoserver.net",
		Wiki:          "https://wiki.tensyoserver.net/wiki/",
	}

	return c.JSON(http.StatusOK, data)
}

// /v2/status
func getStatus(c echo.Context) error {
	db := getDBConnection()

	rows, err := db.Query("select * from ( select * from status order by date desc limit 50 ) as A order by date")

	if err != nil {
		log.Fatalf("failed to query: %v", err)
	}

	defer rows.Close()

	history := models.Status_Database{}

	var historys_proxy []models.Status_Data
	var historys_lobby []models.Status_Data
	var historys_survival []models.Status_Data
	var historys_minigame []models.Status_Data
	var historys_pve []models.Status_Data

	for rows.Next() {
		if err := rows.Scan(&history.Date, &history.Proxy, &history.Lobby, &history.Survival, &history.Minigame, &history.Pve); err != nil {
			log.Fatalf("failed to scan row: %s", err)
		}

		var proxy models.Status_Data
		var lobby models.Status_Data
		var survival models.Status_Data
		var minigame models.Status_Data
		var pve models.Status_Data

		if err := json.Unmarshal([]byte(history.Proxy), &proxy); err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal([]byte(history.Lobby), &lobby); err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal([]byte(history.Survival), &survival); err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal([]byte(history.Minigame), &minigame); err != nil {
			log.Fatal(err)
		}

		if err := json.Unmarshal([]byte(history.Pve), &pve); err != nil {
			log.Fatal(err)
		}

		historys_proxy = append(historys_proxy, proxy)
		historys_lobby = append(historys_lobby, lobby)
		historys_survival = append(historys_survival, survival)
		historys_minigame = append(historys_minigame, minigame)
		historys_pve = append(historys_pve, pve)
	}

	res_1, _ := mcstatus.GetStatusData(address, "25565")
	res_2, _ := mcstatus.GetStatusData(address, "25566")
	res_3, _ := mcstatus.GetStatusData(address, "25570")
	res_4, _ := mcstatus.GetStatusData(address, "25567")
	res_5, _ := mcstatus.GetStatusData(address, "25568")

	var data []models.Status_Server

	data = append(data, models.Status_Server{Server: "Proxy", Current: res_1, History: historys_proxy})
	data = append(data, models.Status_Server{Server: "Lobby", Current: res_2, History: historys_lobby})
	data = append(data, models.Status_Server{Server: "Survival", Current: res_3, History: historys_survival})
	data = append(data, models.Status_Server{Server: "MiniGame", Current: res_4, History: historys_minigame})
	data = append(data, models.Status_Server{Server: "PvE", Current: res_5, History: historys_pve})

	result := models.Status_Result{
		Status: "success",
		Data:   data,
	}

	return c.JSON(http.StatusOK, result)
}
