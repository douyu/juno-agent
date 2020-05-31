package redis

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var healthCheck *RedisHealthCheck

const (
	JSONRedisDialOptions = `
	{"addr":"127.0.0.1:6379"}
	`
)

func TestMain(m *testing.M) {
	healthCheck = NewRedisHealthCheck()
	m.Run()
}

func Test_redixHealthCheck_LoadExtConfig(t *testing.T) {
	type args struct {
		extConfig string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "testing",
			args:    args{extConfig: JSONRedisDialOptions},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := healthCheck.LoadExtConfig(tt.args.extConfig); (err != nil) != tt.wantErr {
				t.Errorf("LoadExtConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redixHealthCheck_DoHealthCheck(t *testing.T) {
	err := healthCheck.LoadExtConfig(JSONRedisDialOptions)
	assert.Nil(t, err)
	_, err = healthCheck.DoHealthCheck()
	assert.Nil(t, err)
}
