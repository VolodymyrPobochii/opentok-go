package pkg

import (
	"reflect"
	"testing"
	"time"
)

func TestOpenTok_decodeSessionId(t *testing.T) {
	type args struct {
		sessionId string
	}
	time1 := time.Unix(0, 1584806881261*int64(time.Millisecond))
	time2 := time.Unix(0, 1584807675066*int64(time.Millisecond))
	tests := []struct {
		name string
		args args
		want *SessionInfo
	}{
		{
			"decode_session1",
			args{sessionId: "2_MX40NjUxMzYwMn5-MTU4NDgwNjg4MTI2MX55NG5zMzBaN1loUi9YVHVmV1pkRkNkRTV-UH4"},
			&SessionInfo{"46513602", "", &time1},
		},
		{
			"decode_session2",
			args{sessionId: "2_MX40NjUxMzYwMn5-MTU4NDgwNzY3NTA2Nn51WVJVSUtXQkZ2U3o2ZmhyWkp2QW1qTW1-fg"},
			&SessionInfo{"46513602", "", &time2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ot := &OpenTok{}
			if got, _ := ot.decodeSessionId(tt.args.sessionId); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeSessionId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewOpenTok(t *testing.T) {
	type args struct {
		apiKey    string
		apiSecret string
		env       interface{}
	}
	tests := []struct {
		name string
		args args
		want *OpenTok
	}{
		{
			name: "NewOpenTok no Env",
			args: args{
				apiKey:    "46513602",
				apiSecret: "d06eaf53e214c105f02f2615170a04e08bf39aa6",
				env:       nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ot := NewOpenTok(tt.args.apiKey, tt.args.apiSecret, tt.args.env)
			session, err := ot.CreateSession(make(map[string]interface{}))
			if err != nil {
				t.Errorf("CreateSession() = %v", err)
			}
			token, err := ot.GenerateToken(session.sessionId, make(map[string]interface{}, 0))
			if err != nil {
				t.Errorf("GenerateToken() = %v", err)
			}
			t.Log("token:", token)
			if !reflect.DeepEqual(ot, tt.want) {
				t.Errorf("NewOpenTok() = %v, want %v", ot, tt.want)
			}
		})
	}
}

func TestOpenTok_createSession(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
		env       interface{}
		client    *Client
	}
	type args struct {
		options map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Session
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ot := &OpenTok{
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				env:       tt.fields.env,
				client:    tt.fields.client,
			}
			got, err := ot.CreateSession(tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateSession() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateSession() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenTok_generateJwt(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
		env       interface{}
		client    *Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ot := &OpenTok{
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				env:       tt.fields.env,
				client:    tt.fields.client,
			}
			got, err := ot.generateJwt()
			if (err != nil) != tt.wantErr {
				t.Errorf("generateJwt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateJwt() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOpenTok_generateToken(t *testing.T) {
	type fields struct {
		apiKey    string
		apiSecret string
		env       interface{}
		client    *Client
	}
	type args struct {
		sessionId string
		options   map[string]interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ot := &OpenTok{
				apiKey:    tt.fields.apiKey,
				apiSecret: tt.fields.apiSecret,
				env:       tt.fields.env,
				client:    tt.fields.client,
			}
			got, err := ot.GenerateToken(tt.args.sessionId, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}
