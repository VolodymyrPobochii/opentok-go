package pkg

const (
	EventConnectionCreated    = "connectionCreated"
	EventConnectionDestroyed  = "connectionDestroyed"
	EventStreamCreated        = "streamCreated"
	EventStreamDestroyed      = "streamDestroyed"
	ReasonClientDisconnected  = "clientDisconnected"
	ReasonForceDisconnected   = "forceDisconnected"
	ReasonForceUnpublished    = "forceUnpublished"
	ReasonMediaStopped        = "mediaStopped"
	ReasonNetworkDisconnected = "networkDisconnected"
)

type Callback struct {
	SessionID string `json:"sessionId"`
	ProjectID string `json:"projectId"`
	Event     string `json:"event"`
	Reason    string `json:"reason"`
	Timestamp int64  `json:"timestamp"`
}

type Connection struct {
	ID        string `json:"id"`
	CreatedAt int64  `json:"createdAt"`
	Data      string `json:"data"`
}

type ConnectionCallback struct {
	Callback
	Connection *Connection `json:"connection"`
}

type Stream struct {
	ID         string      `json:"id"`
	Connection *Connection `json:"connection"`
	CreatedAt  int64       `json:"createdAt"`
	Name       string      `json:"name"`
	VideoType  string      `json:"videoType"`
}

type StreamCallback struct {
	Callback
	Stream *Stream `json:"stream"`
}

type SessionCallback struct {
	Callback
	Connection *Connection `json:"connection,omitempty"`
	Stream     *Stream     `json:"stream,omitempty"`
}
