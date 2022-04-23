package scf

type Reply struct {
	Msg       string      `json:"msg"`
	Data      interface{} `json:"data"`
	RequestId string      `json:"request_id"`
	UnixTime  int64       `json:"unix_time"`
}
