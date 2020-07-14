package pkg

import "testing"

func TestEncodeToken(t *testing.T) {
	type args struct {
		tokenData map[string]interface{}
		apiKey    string
		apiSecret string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "encode_token",
			args: args{
				tokenData: map[string]interface{}{"iss": "iss"},
				apiKey:    "1234567",
				apiSecret: "secret",
			},
			want:    "T1==cGFydG5lcl9pZD0xMjM0NTY3JnNpZz0xNzNiYjVkMjU5NWM5MjhiYjg0MmJjNDk4ZWE1ODkzMjAzMDZhZDY4Omlzcz1pc3MmY3JlYXRlX3RpbWU9MTU4NTQ4NzMzNyZleHBpcmVfdGltZT0xNTg1NTczNzM3Jm5vbmNlPTE1ODU0ODczMzY4MzQwMCZyb2xlPXB1Ymxpc2hlcg==",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeToken(tt.args.tokenData, tt.args.apiKey, tt.args.apiSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("EncodeToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_signString(t *testing.T) {
	type args struct {
		unsigned string
		secret   string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "sign_string",
			args: args{
				unsigned: "unsigned_data",
				secret:   "secretkey",
			},
			want:    "19c185a9c0d81cc6bc554d3e926710dc500ffb22",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := signString(tt.args.unsigned, tt.args.secret)
			if (err != nil) != tt.wantErr {
				t.Errorf("signString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("signString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
