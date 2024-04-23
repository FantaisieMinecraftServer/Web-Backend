package mcstatus

import (
	"backend/app/models"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func GetStatusData(address string, port string) (models.Status_Data, error) {
	apiURL := fmt.Sprintf("https://api.mcstatus.io/v2/status/java/%s:%s", address, port)

	var data models.Status_Data

	res, err := http.Get(apiURL)
	if err != nil {
		return data, fmt.Errorf("error making request to mcstatus.io API: %v", err)
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

	data.AcquisitionTime = time.Now().Format("2006-01-02 15:04:05")

	return data, nil
}
