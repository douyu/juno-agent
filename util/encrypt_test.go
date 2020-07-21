package util

import "testing"

func TestAesEncrypt(t *testing.T) {
	type args struct {
		orig string
		key  string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "1",
			args: args{
				orig: "123", key: "12341234123412341234123412341234",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := AESCBCEncrypt(tt.args.orig, tt.args.key)
			dec, _ := AESCBCDecrypt(got, tt.args.key)
			if tt.args.orig != dec {
				t.Errorf("AESCBCEncrypt() = %v, AESCBCDecrypt() = %v", got, dec)
			}
		})
	}
}
