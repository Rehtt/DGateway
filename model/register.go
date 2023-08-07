package model

type Register struct {
	Uid    string   `json:"uid"`
	Port   int      `json:"port"`
	Scheme string   `json:"scheme"`
	Routes []*Route `json:"routes"`
	Remote string   `json:"remote,omitempty"`
}
type Route struct {
	Method string `json:"method"`
	Uri    string `json:"uri"`
}
