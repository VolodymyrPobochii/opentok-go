package pkg

import "testing"

func TestGenerateJwt(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "generate JWT",
			args: args{config: &Config{
				ApiKey:    "1234567",
				ApiSecret: "bfguyaergfubfgabeabev",
				Auth:      &Auth{Expire: 1},
			}},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODU0ODk3MDMsImlhdCI6MTU4NTQ4OTcwMiwiaXNzIjoiMTIzNDU2NyIsImlzdCI6InByb2plY3QifQ.YKP8Xy_7UD3nLA5W5pvlLexkrX9V3lZ5wY78NnhIwAE",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateJwt(tt.args.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateJwt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateJwt() got = %v, want %v", got, tt.want)
			}
		})
	}
}
