package models

type Path struct {
	ID string `json:"id" gorethink:"id,omitempty"`
	I  string `json:"i" gorethink:"i"`
	V  string `json:"v" gorethink:"v"`
}

type PathResponse struct {
	Path   []Point `json:"path"`
	UserId string  `json:"userId"`
}

type Point struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
