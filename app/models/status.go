package models

import "time"

// mcstatus.io API
type Version struct {
	NameRaw   string `json:"name_raw"`
	NameClean string `json:"name_clean"`
	NameHtml  string `json:"name_html"`
	Protocol  int    `json:"protocol"`
}

type PlayerList struct {
	Uuid      string `json:"uuid"`
	NameRaw   string `json:"name_raw"`
	NameClean string `json:"name_clean"`
	NameHtml  string `json:"name_html"`
}

type Players struct {
	Online int           `json:"online"`
	Max    int           `json:"max"`
	List   []*PlayerList `json:"list"`
}

type Motd struct {
	Raw   string `json:"raw"`
	Clean string `json:"clean"`
	Html  string `json:"html"`
}

type Mods struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type SrvRecord struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

type Plugins struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Status_Data struct {
	AcquisitionTime string     `json:"acquisition_time"`
	Online          bool       `json:"online"`
	Host            string     `json:"host"`
	Port            int        `json:"port"`
	IpAddress       string     `json:"ip_address"`
	EulaBlocked     bool       `json:"eula_blocked"`
	RetrievedAt     int        `json:"retrieved_at"`
	ExpiresAt       int        `json:"expires_at"`
	Version         Version    `json:"version"`
	Players         Players    `json:"players"`
	Motd            Motd       `json:"motd"`
	Icon            string     `json:"icon"`
	Mods            []*Mods    `json:"mods"`
	SoftWare        string     `json:"software"`
	Plugins         []*Plugins `json:"plugins"`
	SrvRecord       SrvRecord  `json:"srv_record"`
}

// Local API
type Status_Database struct {
	Date     time.Time
	Proxy    string
	Lobby    string
	Survival string
	Minigame string
	Pve      string
}

type Status_Server struct {
	Server  string       `json:"server"`
	Current Status_Data   `json:"current"`
	History []Status_Data `json:"history"`
}

type Status_Result struct {
	Status string          `json:"status"`
	Data   []Status_Server `json:"data"`
}
