package model

type Path struct {
	ID string `json:"id" gorethink:"id,omitempty"`
	I  string `json:"i" gorethink:"i"`
	V  string `json:"v" gorethink:"v"`
}
