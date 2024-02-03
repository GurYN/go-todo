package models

type ActionRequest struct {
	Action string `json:"action"`
	Data   any    `json:"data"`
}

type TodoList struct {
	Answer string `json:"answer"`
	Data   any    `json:"data"`
}

type TodoId struct {
	Id string `json:"id"`
}
