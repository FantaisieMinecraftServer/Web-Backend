package routes

import (
	"encoding/json"
	"log"
	"main/lib"
	"main/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetStatus(c echo.Context) error {
	address := "play.tensyoserver.net"

	db := lib.GetDBConnection()

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

	res_1, _ := lib.GetStatusData(address, "25565")
	res_2, _ := lib.GetStatusData(address, "25566")
	res_3, _ := lib.GetStatusData(address, "25570")
	res_4, _ := lib.GetStatusData(address, "25567")
	res_5, _ := lib.GetStatusData(address, "25568")

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
