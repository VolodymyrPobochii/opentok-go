package pkg

import (
	"encoding/base64"
	"errors"
	"strings"
)

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

func (se *SessionCallback) ParseData() (uid string, chatId string, videoCallId string, err error) {
	var dataString string
	if connection := se.Connection; connection != nil {
		dataString = connection.Data
	}

	if stream := se.Stream; stream != nil {
		dataString = stream.Connection.Data
	}

	if len(dataString) == 0 {
		err = errors.New("connection/stream is not present")
		return
	}

	if rawData, err := base64.StdEncoding.DecodeString(dataString); err == nil {
		dataSplit := strings.Split(string(rawData), "&")
		uid = dataSplit[0]
		chatId = dataSplit[1]
		videoCallId = dataSplit[2]
	}
	return
}
