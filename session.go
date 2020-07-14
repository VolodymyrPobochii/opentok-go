package opentok

type Session struct {
	ot         *OpenTok
	sessionId  string
	properties map[string]interface{}
}

// Represents an OpenTok session. Use the {@link OpenTok#createSession OpenTok.createSession()}
// method to create an OpenTok session. The <code>sessionId</code> property of the Session object
// is the session ID.
// @property {String} sessionId The session ID.
// @class Session
func NewSession(ot *OpenTok, sessionId string, properties map[string]interface{}) *Session {
	return &Session{ot, sessionId, properties}
}

func (s *Session) generateToken(options map[string]interface{}) (string, error) {
	return s.ot.generateToken(s.sessionId, options)
}