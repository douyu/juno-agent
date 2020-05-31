package tcp

import (
	"testing"
)

const (
	TcpJsonConfig = `{"addr":"127.0.0.1:8000","network":"tcp4"}`
)

func TestTCPHealthCheck_DoHealthCheck(t *testing.T) {
	// get instance
	tcpHealthCheck := NewTCPHealthCheck()
	type args struct {
		extConfig string
	}
	test := struct {
		name    string
		args    args
		wantErr bool
	}{
		name:    "tcp check",
		args:    args{extConfig: TcpJsonConfig},
		wantErr: false,
	}
	t.Run(test.name, func(t *testing.T) {
		if err := tcpHealthCheck.LoadExtConfig(test.args.extConfig); (err != nil) != test.wantErr {
			t.Errorf("LoadExtConfig() error = %v, wantErr %v", err, test.wantErr)
		}
	})
	t.Run(test.name, func(t *testing.T) {
		if _, err := tcpHealthCheck.DoHealthCheck(); err != nil {
			t.Errorf("DoHealthCheck() error = %v ", err.Error())
		}
	})
	t.Log("name: ", test.name, "config: ", tcpHealthCheck.Addr)
}
