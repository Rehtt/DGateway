package model

type Register struct {
	Uid        string   `json:"uid"`
	Port       int      `json:"port"`
	Scheme     string   `json:"scheme"`
	Routes     []*Route `json:"routes"`
	RemoteBase string   `json:"remote_base,omitempty"`
}
type Route struct {
	Method string `json:"method"`
	Uri    string `json:"uri"`
	Match  bool   `json:"match"` // 正则匹配
}
