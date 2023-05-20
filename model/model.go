package model

type Path struct {
	ID string `json:"id" gorethink:"id,omitempty"`
	i  string `json:"i" gorethink:"i"`
	v  string `json:"v" gorethink:"v"`
}
