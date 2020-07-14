package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Endpoints struct {
	CreateSession       string
	GetStream           string
	ListStreams         string
	SetArchiveLayout    string
	SetStreamClassLists string
	Dial                string
	StartBroadcast      string
	StopBroadcast       string
	GetBroadcast        string
	SetBroadcastLayout  string
	ListBroadcasts      string
}

type Request struct {
	Timeout int64
	Proxy   string
}
type Auth struct {
	Expire int64
}

type Config struct {
	ApiKey    string
	ApiSecret string
	ApiUrl    string
	Endpoints *Endpoints
	Request   *Request
	Auth      *Auth
	ClientConfig
}

type ClientConfig struct {
	ApiUrl     string
	Request    *Request
	UaAddendum string
}

type ApiConfig struct {
	ApiEndpoint string
	ApiKey      string
	ApiSecret   string
	Auth        *Auth
	Proxy       string
	UaAddendum  string
}

func defaultConfig() *Config {
	return &Config{
		ApiKey:    "",
		ApiSecret: "",
		ApiUrl:    "https://api.opentok.com",
		Endpoints: &Endpoints{
			CreateSession:       "/session/create",
			GetStream:           "/v2/project/%s/session/%s/stream/%s", //<%apiKey%>,<%sessionId%>,<%streamId%>
			ListStreams:         "/v2/project/%s/session/%s/stream",    //<%apiKey%>,<%sessionId%>
			SetArchiveLayout:    "/v2/project/%s/archive/%s/layout",    //<%apiKey%>,<%archiveId%>
			SetStreamClassLists: "/v2/project/%s/session/%s/stream",    //<%apiKey%>,<%sessionId%>
			Dial:                "/v2/project/%s/dial",                 //<%apiKey%>
			StartBroadcast:      "/v2/project/%s/broadcast",            //<%apiKey%>
			StopBroadcast:       "/v2/project/%s/broadcast/%s/stop",    //<%apiKey%>,<%broadcastId%>
			GetBroadcast:        "/v2/project/%s/broadcast/%s",         //<%apiKey%>,<%broadcastId%>
			SetBroadcastLayout:  "/v2/project/%s/broadcast/%s/layout",  //<%apiKey%>,<%broadcastId%>
			ListBroadcasts:      "/v2/project/%s/broadcast",            //<%apiKey%>
		},
		Request: &Request{Timeout: 20000}, // 20 seconds
		Auth:    &Auth{Expire: 300},
	}
}

type Client struct {
	config     *Config
	httpClient *http.Client
	//config    map[string]interface{}
}

type CreateSessionResponse struct {
	SessionId           string      `json:"session_id"`
	SessionSegmentId    string      `json:"session_segment_id"`
	SessionStatus       string      `json:"session_status"`
	StatusInvalid       interface{} `json:"status_invalid"`
	ProjectId           string      `json:"project_id"`
	PartnerId           string      `json:"partner_id"`
	CreateDT            string      `json:"create_dt"`
	MediaServerUrl      string      `json:"media_server_url"`
	MediaServerHostname string      `json:"media_server_hostname"`
	MessagingServerUrl  string      `json:"messaging_server_url"`
	MessagingUrl        string      `json:"messaging_url"`
	IceCredentialExp    int64       `json:"ice_credential_expiration"`
	IceServer           string      `json:"ice_server"`
	SymphonyAddress     string      `json:"symphony_address"`
	IceServers          interface{} `json:"ice_servers"`
	Properties          interface{} `json:"properties"`
}

type CreateSessionError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (c *Client) createSession(options map[string]interface{}) (string, error) {
	body, err := json.Marshal(options)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("%v%v", c.config.ApiUrl, c.config.Endpoints.CreateSession)
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	err = c.generateHeaders(&request.Header)
	if err != nil {
		return "", err
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", errors.New(fmt.Sprintf("the request failed: %v", err))
	}

	defer response.Body.Close()

	decoder := json.NewDecoder(response.Body)

	var responseError *CreateSessionError
	// handle client errors
	if response.StatusCode == 403 {
		_ = decoder.Decode(&responseError)
		return "", errors.New(fmt.Sprintf("an authentication error occurred: (%d) %v", response.StatusCode, responseError))
	}

	// handle server errors
	if response.StatusCode >= 500 && response.StatusCode <= 599 {
		_ = decoder.Decode(&responseError)
		return "", errors.New(fmt.Sprintf("a server error occurred: (%d) %v", response.StatusCode, responseError))
	}

	var sessionResponse []*CreateSessionResponse
	err = decoder.Decode(&sessionResponse)
	if err != nil {
		return "", err
	}

	return sessionResponse[0].SessionId, nil
}

func (c *Client) generateHeaders(header *http.Header) error {
	jwt, err := GenerateJwt(c.config)
	if err != nil {
		return errors.New(fmt.Sprintf("failed to generate JWT: %v", err))
	}
	header.Set("X-OPENTOK-AUTH", jwt)
	header.Set("Accept", "application/json")
	header.Set("User-Agent", "OpenTok-GO-SDK/v0.0.1")
	return nil
}

func (c *Client) configure(clientConfig *ClientConfig) *Config {
	// merge configs
	c.config.ClientConfig = *clientConfig
	if len(c.config.Endpoints.Dial) != 0 && len(c.config.ApiKey) != 0 {
		c.config.Endpoints.Dial = fmt.Sprintf(c.config.Endpoints.Dial, c.config.ApiKey)
	}
	if c.config.Request != nil {
		c.httpClient.Timeout = time.Duration(c.config.Request.Timeout)
	}
	return c.config
}

func NewClient(apiKey, apiSecret string) *Client {
	//config := make(map[string]interface{})
	config := defaultConfig()
	config.ApiKey = apiKey
	config.ApiSecret = apiSecret
	return &Client{config: config, httpClient: &http.Client{}}
}
