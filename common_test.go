package opentok

import (
	"reflect"
	"testing"
)

func TestClone(t *testing.T) {
	type args struct {
		src map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Clone(tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Clone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDefaults(t *testing.T) {
	type args struct {
		dst map[string]interface{}
		src map[string]interface{}
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Defaults(tt.args.dst, tt.args.src); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Defaults() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIncludes(t *testing.T) {
	type args struct {
		src    []interface{}
		target interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Includes(tt.args.src, tt.args.target); got != tt.want {
				t.Errorf("Includes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPick(t *testing.T) {
	type args struct {
		src  map[string]interface{}
		keys []string
	}
	tests := []struct {
		name string
		args args
		want map[string]interface{}
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Pick(tt.args.src, tt.args.keys); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Pick() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNonce(t *testing.T) {
	type args struct {
		length int
	}
	tests := []struct {
		name string
		args args
		want int64
		wantErr bool
	}{
		{"nonce", args{0}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Nonce(tt.args.length)()
			if (err != nil) != tt.wantErr {
				t.Errorf("Nonce() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Nonce() got = %v, want %v", got, tt.want)
			}
		})
	}
}