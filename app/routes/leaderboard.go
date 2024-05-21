package routes

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"main/lib"
	"main/models"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func GetUUIDtoName(uuid string) (models.Profile, error) {
	apiURL := fmt.Sprintf("https://sessionserver.mojang.com/session/minecraft/profile/%s", uuid)

	var data models.Profile

	res, err := http.Get(apiURL)
	if err != nil {
		return data, fmt.Errorf("error making request to mojang.com API: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return data, fmt.Errorf("error reading response body: %v", err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return data, fmt.Errorf("error decoding JSON: %v", err)
	}

	return data, nil
}

func GetLeaderBoard(c echo.Context) error {
	data_type := c.QueryParam("type")

	types := []string{"blockBreak", "blockPlace", "balance", "playTime"}
	db_table := slices.Contains(types, data_type)

	if !db_table {
		return c.JSON(http.StatusNotFound, models.LeaderBoard_Error{Status: "failed", Message: "Type does not match"})
	}

	var result []models.LeaderBoardData

	db := lib.GetLBDBConnection()

	rows, err := db.Query(fmt.Sprintf("SELECT * FROM Player ORDER BY LPAD(%s,64,0) DESC", data_type))
	if err != nil {
		log.Fatalf("failed to query: %v", err)
	}

	defer rows.Close()

	history := models.Player_Database{}

	for rows.Next() {
		if err := rows.Scan(&history.UUID, &history.LastName, &history.Balance, &history.BlockBreak, &history.BlockPlace, &history.PlayTime); err != nil {
			log.Fatalf("failed to scan row: %s", err)
		}

		var iconURL string

		if history.LastName == nil {
			continue
		} else {
			if strings.HasPrefix(*history.LastName, ".") {
				iconURL = "https://crafthead.net/helm/MHF_Steve/40"
			} else {
				iconURL = fmt.Sprintf("https://crafthead.net/helm/%s/40", *history.LastName)
			}
		}

		switch data_type {
		case "blockBreak":
			result = append(result, models.LeaderBoardData{Name: *history.LastName, Icon: iconURL, Value: strconv.Itoa(int(history.BlockBreak))})
		case "blockPlace":
			result = append(result, models.LeaderBoardData{Name: *history.LastName, Icon: iconURL, Value: strconv.Itoa(int(history.BlockPlace))})
		case "balance":
			result = append(result, models.LeaderBoardData{Name: *history.LastName, Icon: iconURL, Value: strconv.Itoa(int(history.Balance))})
		case "playTime":
			result = append(result, models.LeaderBoardData{Name: *history.LastName, Icon: iconURL, Value: strconv.Itoa(int(history.PlayTime))})
		}
	}

	data := &models.LeaderBoard{
		Status: "success",
		Data:   result,
	}

	return c.JSON(http.StatusOK, data)
}
