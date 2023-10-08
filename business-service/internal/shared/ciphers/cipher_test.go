package ciphers_test

import (
	"reflect"
	"testing"

	"github.com/lucasd-coder/fast-feet/business-service/internal/shared/ciphers"
)

func TestDecrypt(t *testing.T) {
	key := ciphers.ExtractKey([]byte("key"))
	enc, _ := ciphers.Encrypt(key, []byte("value"))
	type args struct {
		key        []byte
		ciphertext []byte
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "test decrypt",
			args: args{
				key:        key,
				ciphertext: enc,
			},
			want:    []byte("value"),
			wantErr: false,
		},
		{
			name: "test fail decrypt",
			args: args{
				key:        ciphers.ExtractKey([]byte("key2")),
				ciphertext: enc,
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ciphers.Decrypt(tt.args.key, tt.args.ciphertext)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}
