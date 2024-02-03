package models

type WebsocketRequest struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

type WebsocketResponse struct {
	Answer string `json:"answer"`
	Data   any    `json:"data"`
}
