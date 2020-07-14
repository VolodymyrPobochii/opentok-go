package pkg

type Session struct {
	ot         *OpenTok
	sessionId  string
	properties map[string]interface{}
}

func (s *Session) Id() string {
	return s.sessionId
}

// Represents an OpenTok session. Use the {@link OpenTok#CreateSession OpenTok.CreateSession()}
// method to create an OpenTok session. The <code>sessionId</code> property of the Session object
// is the session ID.
// @property {String} sessionId The session ID.
// @class Session
func NewSession(ot *OpenTok, sessionId string, properties map[string]interface{}) *Session {
	return &Session{ot, sessionId, properties}
}

func (s *Session) generateToken(options map[string]interface{}) (string, error) {
	return s.ot.GenerateToken(s.sessionId, options)
}
