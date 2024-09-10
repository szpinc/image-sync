package types

type CmdRequest struct {
	Host              string `json:"host"`
	Port              int    `json:"port"`
	App               string `json:"app"`
	Repository        string `json:"repository"`
	Tag               string `json:"tag"`
	DockerComposeFile string `json:"dockerComposeFile"`
	Deploy            bool   `json:"deploy"`
}

type ErrorNotAuthed error
