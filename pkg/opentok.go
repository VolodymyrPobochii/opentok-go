package pkg

import (
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"time"
)

var (
	ErrorNoSessionId    = errors.New("token cannot be generated without a sessionId parameter")
	ErrorNoApiKey       = errors.New("token cannot be generated unless the session belongs to the API Key")
	ErrorWrongMediaMode = errors.New("a session with always archive mode must also have the routed media mode")
	ErrorInvalidIPv4    = errors.New("invalid arguments when calling CreateSession, location must be an IPv4 address")
)

type SessionInfo struct {
	apiKey     string
	location   string
	createTime *time.Time
}

type OpenTok struct {
	apiKey    string
	apiSecret string
	env       interface{}
	client    *Client
}

func (ot *OpenTok) ApiKey() string {
	return ot.apiKey
}

func NewOpenTok(apiKey, apiSecret string, env interface{}) *OpenTok {
	client := NewClient(apiKey, apiSecret)

	//apiConfig := map[string]interface{}{
	//	"apiEndpoint": "https://api.opentok.com",
	//	"apiKey":      apiKey,
	//	"apiSecret":   apiSecret,
	//	"auth":        map[string]interface{}{"expire": int64(300)},
	//}
	apiConfig := &ApiConfig{
		ApiEndpoint: "https://api.opentok.com",
		ApiKey:      apiKey,
		ApiSecret:   apiSecret,
		Auth:        &Auth{Expire: 300},
	}
	// env can be either an object with a bunch of DI options, or a simple string for the apiUrl
	//clientConfig := map[string]interface{}{
	//	"request": make(map[string]interface{}),
	//}
	clientConfig := &ClientConfig{}

	if env, ok := env.(string); ok {
		clientConfig.ApiUrl = env
		apiConfig.ApiEndpoint = env
	}

	if env, ok := env.(map[string]interface{}); ok {
		if apiUrl, ok := env["apiUrl"].(string); ok {
			clientConfig.ApiUrl = apiUrl
			apiConfig.ApiEndpoint = apiUrl
		}

		if proxy, ok := env["proxy"].(string); ok {
			clientConfig.Request.Proxy = proxy
			apiConfig.Proxy = proxy
		}

		if uaAddendum, ok := env["uaAddendum"].(string); ok {
			clientConfig.UaAddendum = uaAddendum
			apiConfig.UaAddendum = uaAddendum
		}
	}

	client.configure(clientConfig)

	return &OpenTok{apiKey, apiSecret, env, client}
}

// decodes a sessionId into the metadata that it contains
// @param     {string}         sessionId
// @returns   {?SessionInfo}    sessionInfo
func (ot *OpenTok) decodeSessionId(sessionId string) (*SessionInfo, error) {
	// remove sentinal (e.g. '1_', '2_')
	sessionId = sessionId[2:]

	// decode base64 with padding
	sessionIdBytes := []byte(sessionId)
	shift := 4 - len(sessionIdBytes)%4
	for i := 0; i < shift; i++ {
		sessionIdBytes = append(sessionIdBytes, byte('='))
	}
	sessionId = string(sessionIdBytes)
	bytes, err := base64.URLEncoding.DecodeString(sessionId)
	if err != nil {
		return nil, err
	}

	// separate fields
	fields := strings.Split(string(bytes), "~")
	timestamp, err := strconv.ParseInt(fields[3], 10, 64)
	if err != nil {
		return nil, err
	}
	ttime := time.Unix(0, timestamp*int64(time.Millisecond))
	return &SessionInfo{
		apiKey:     fields[1],
		location:   fields[2],
		createTime: &ttime,
	}, nil
}

func (ot *OpenTok) CreateSession(options map[string]interface{}) (*Session, error) {
	// whitelist the keys allowed
	src := map[string]interface{}{"mediaMode": "relayed", "archiveMode": "manual"}
	keys := []string{"mediaMode", "archiveMode", "location"}
	if options == nil {
		options = make(map[string]interface{}, len(src))
	}
	options = Pick(Defaults(options, src), keys)

	if options["mediaMode"] != "routed" && options["mediaMode"] != "relayed" {
		options["mediaMode"] = "relayed"
	}

	if options["archiveMode"] != "manual" && options["archiveMode"] != "always" {
		options["archiveMode"] = "manual"
	}

	if options["archiveMode"] == "always" && options["mediaMode"] != "routed" {
		return nil, ErrorWrongMediaMode
	}

	if location, ok := options["location"].(string); ok && net.ParseIP(location) == nil {
		return nil, ErrorInvalidIPv4
	}

	// rename mediaMode -> p2p.preference
	// store backup for use in constructing Session
	backupOpts := Clone(options)
	// avoid mutating passed in options
	options = Clone(options)
	mediaModeToParam := map[string]string{"routed": "disabled", "relayed": "enabled"}
	options["p2p.preference"] = mediaModeToParam[options["mediaMode"].(string)]
	delete(options, "mediaMode")

	sessionId, err := ot.client.createSession(options)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to CreateSession. %v", err))
	}
	log.Println("created session:", sessionId)
	return NewSession(ot, sessionId, backupOpts), nil
}

// Creates a token for connecting to an OpenTok session. In order to authenticate a user
// connecting to an OpenTok session, the client passes a token when connecting to the session.
// <p>
// For testing, you can also generate a token by logging into your
// <a href="https://tokbox.com/account">TokBox account</a>.
//
// @param sessionId The session ID corresponding to the session to which the user will connect.
//
// @param options An object that defines options for the token (each of which is optional):
//
// <ul>
//    <li><code>role</code> (String) &mdash; The role for the token. Each role defines a set of
//      permissions granted to the token:
//
//        <ul>
//           <li> <code>'subscriber'</code> &mdash; A subscriber can only subscribe to streams.</li>
//
//           <li> <code>'publisher'</code> &mdash; A publisher can publish streams, subscribe to
//              streams, and signal. (This is the default value if you do not specify a role.)</li>
//
//           <li> <code>'moderator'</code> &mdash; In addition to the privileges granted to a
//             publisher, in clients using the OpenTok.js library, a moderator can call the
//             <code>forceUnpublish()</code> and <code>forceDisconnect()</code> method of the
//             Session object.</li>
//        </ul>
//
//    </li>
//
//    <li><code>expireTime</code> (Number) &mdash; The expiration time for the token, in seconds
//      since the UNIX epoch. The maximum expiration time is 30 days after the creation time. If
//      a fractional number is specified, then it is rounded down to the nearest whole number.
//      The default expiration time of 24 hours after the token creation time.
//    </li>
//
//    <li><code>data</code> (String) &mdash; A string containing connection metadata describing the
//      end-user.For example, you can pass the user ID, name, or other data describing the end-user.
//      The length of the string is limited to 1000 characters. This data cannot be updated once it
//      is set.
//    </li>
//
//    <li><code>initialLayoutClassList</code> (Array) &mdash; An array of class names (strings)
//      to be used as the initial layout classes for streams published by the client. Layout
//      classes are used in customizing the layout of videos in
//      <a href="https://tokbox.com/developer/guides/broadcast/live-streaming/">live streaming
//      broadcasts</a> and
//      <a href="https://tokbox.com/developer/guides/archiving/layout-control.html">composed
//      archives</a>.
//    </li>
//
// </ul>
//
// @return The token string.
func (ot *OpenTok) GenerateToken(sessionId string, options map[string]interface{}) (string, error) {
	now := time.Now().UnixNano() / int64(time.Second)
	if options == nil {
		options = make(map[string]interface{})
	}
	// avoid mutating passed in options
	// todo: copy props ?
	options = Clone(options)

	if len(sessionId) == 0 {
		return "", ErrorNoSessionId
	}

	// validate the sessionId belongs to the apiKey of this OpenTok instance
	decoded, err := ot.decodeSessionId(sessionId)
	if err != nil {
		return "", err
	}
	if decoded.apiKey != ot.apiKey {
		return "", ErrorNoApiKey
	}

	// combine defaults, opts, and whitelisted property names to create tokenData
	if expireTime, ok := options["expire_time"].(int64); ok {
		// Automatic rounding to help out people who pass in a fractional expireTime
		options["expire_time"] = expireTime
	}
	if expireTime, ok := options["expire_time"].(string); ok {
		// Automatic rounding to help out people who pass in a fractional expireTime
		exp, err := strconv.ParseInt(expireTime, 10, 64)
		if err != nil {
			return "", err
		}
		options["expire_time"] = exp
	}
	if data, ok := options["data"].(string); ok {
		if len(data) > 1024 {
			return "", errors.New("invalid data for token generation, must be a string with maximum length 1024")
		}
		options["connection_data"] = data
	}
	if initialLayoutClassList, ok := options["initialLayoutClassList"].([]string); ok {
		joinedClassList := strings.Join(initialLayoutClassList, " ")
		if len(joinedClassList) > 1024 {
			return "", errors.New("invalid initial layout class list for token generation, must have concatenated length of less than 1024'")
		}
		options["initial_layout_class_list"] = joinedClassList
	}
	if initialLayoutClassList, ok := options["initialLayoutClassList"].(string); ok {
		if len(initialLayoutClassList) > 1024 {
			return "", errors.New("invalid initial layout class list for token generation, must have concatenated length of less than 1024'")
		}
		options["initial_layout_class_list"] = initialLayoutClassList
	}

	srcData := map[string]interface{}{
		"session_id":                sessionId,
		"create_time":               now,
		"expire_time":               now + int64(24*time.Hour), // 1 day
		"nonce":                     rand.Int63(),
		"role":                      "publisher",
		"initial_layout_class_list": "",
	}
	tokenData := Pick(
		Defaults(options, srcData),
		[]string{"session_id", "create_time",
			"nonce", "role", "expire_time",
			"connection_data", "initial_layout_class_list"})

	// validate tokenData
	if !Includes([]interface{}{"publisher", "subscriber", "moderator"}, tokenData["role"]) {
		return "", errors.New(fmt.Sprintf("invalid role for token generation: %s", tokenData["role"]))
	}

	if _, ok := tokenData["expire_time"].(int64); !ok {
		return "", errors.New(fmt.Sprintf("invalid expireTime for token generation: %s", tokenData["expire_time"]))
	}

	expireTime := tokenData["expire_time"].(int64)
	if expireTime < now {
		return "", errors.New(fmt.Sprintf("invalid expireTime for token generation, time cannot be in the past: %v < %v", expireTime, now))
	}

	if connectionData, ok := tokenData["connection_data"]; ok && connectionData != nil {
		if connectionData, ok := connectionData.(string); !ok || len(connectionData) > 1024 {
			return "", errors.New("invalid data for token generation, must be a string with maximum length 1024")
		}
	}

	if layoutClassList, ok := tokenData["initial_layout_class_list"].(string); ok && len(layoutClassList) > 1024 {
		return "", errors.New("invalid initial layout class list for token generation, must have concatenated length of less than 1024")
	}

	return EncodeToken(tokenData, ot.apiKey, ot.apiSecret)
}

// decodes a sessionId into the metadata that it contains
// @param     none
// @returns   {string}    JWT
func (ot *OpenTok) GenerateJwt() (string, error) {
	return GenerateJwt(ot.client.config)
}
