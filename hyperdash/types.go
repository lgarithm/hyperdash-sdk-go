// https://github.com/hyperdashio/hyperdash-sdk-py/blob/master/hyperdash/sdk_message.py
package hyperdash

type (
	MsgType     string
	FinalStatus string
	LogLevel    string
)

const (
	kTypeRunStarted MsgType = `run_started`
	kTypeRunEnded   MsgType = `run_ended`
	kTypeLog        MsgType = `log`
	kTypeHeartbeat  MsgType = `heartbeat`
	kTypeMetric     MsgType = `metric`
	kTypeParam      MsgType = `param`

	kOutcomeSuccess      FinalStatus = `success`
	kOutcomeFailure      FinalStatus = `failure`
	kOutcomeUserCanceled FinalStatus = `user_canceled`

	kLevelInfo LogLevel = `INFO`
)

type SDKMessage struct {
	Type       MsgType     `json:"type"`
	Timestamp  int64       `json:"timestamp"`
	SDKRunUUID string      `json:"sdk_run_uuid"`
	Payload    interface{} `json:"payload"`
}

type runStarted struct {
	JobName string `json:"job_name"`
}

type runEnded struct {
	FinalStatus FinalStatus `json:"final_status"`
}

type userLog struct {
	UUID  string   `json:"uuid"`
	Level LogLevel `json:"level"`
	Body  string   `json:"body"`
}

type metric struct {
	Name       string      `json:"name"`
	Timestamp  int64       `json:"timestamp"`
	Value      interface{} `json:"value"`
	IsInternal bool        `json:"is_internal"`
}

type param struct {
	Params     map[string]interface{} `json:"params"`
	IsInternal bool                   `json:"is_internal"`
}
