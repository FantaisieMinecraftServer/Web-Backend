package models

type APIHelp struct {
	Version    string `json:"version"`
	WhatIsThis string `json:"what_is_this"`
	Author     string `json:"author"`
	HomePage   string `json:"homepage"`
}

type Error struct {
	Reason string `json:"reason"`
}
