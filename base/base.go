package base

type Command struct {
	Cmd  string   `json:"cmd"`
	Args []string `json:"args"`
}
