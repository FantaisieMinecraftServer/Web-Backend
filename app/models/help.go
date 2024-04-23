package models

type APIHelp struct {
	Status        string `json:"status"`
	Version       string `json:"version"`
	WhatIsThis    string `json:"what_is_this"`
	ServerAddress string `json:"server_address"`
	HomePage      string `json:"home_page"`
	Wiki          string `json:"wiki"`
}
