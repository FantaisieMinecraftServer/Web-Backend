package models

type Propertie struct {
	Name      string `json:"name"`
	Value     string `json:"value"`
	Signature string `json:"signature"`
}

type Profile struct {
	Id             string      `json:"id"`
	Name           string      `json:"name"`
	Properties     []Propertie `json:"properties"`
	ProfileActions []string    `json:"profileActions"`
}

type LeaderBoardData struct {
	Name  string `json:"name"`
	Icon  string `json:"icon"`
	Value string `json:"value"`
}

type LeaderBoard struct {
	Status string            `json:"status"`
	Data   []LeaderBoardData `json:"data"`
}

type LeaderBoard_Error struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Player_Database struct {
	UUID       string
	LastName   *string
	Balance    float64
	BlockBreak int64
	BlockPlace int64
	PlayTime   int64
}
