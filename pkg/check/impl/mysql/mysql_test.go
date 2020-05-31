package mysql

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	DsnConfig     = `root:saybye.1314LRJ@tcp(127.0.0.1:3306)/mark?charset=utf8&parseTime=True&loc=Local&timeout=2s`
	DsnJSONConfig = `{"dsn":"root:saybye.1314LRJ@tcp(127.0.0.1:3306)/mark?charset=utf8&parseTime=True&loc=Local&timeout=2s"}`
)

func ExampleMysqlHealthCheck() {
	// get instance
	mysqlHealthCheck := NewMysqlHealthCheck()
	// LoadExtConfig
	err := mysqlHealthCheck.LoadExtConfig(DsnConfig)
	if err != nil {
		fmt.Println(err)
	}
	// DoHealthCheck
	resHealthCheck, err := mysqlHealthCheck.DoHealthCheck()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(resHealthCheck)
}

var healthCheck *MysqlHealthCheck

func TestMain(m *testing.M) {
	healthCheck = NewMysqlHealthCheck()
	m.Run()
}
func TestParseDSN(t *testing.T) {
	cfg, err := ParseDSN(DsnConfig)
	assert.Nil(t, err)
	t.Logf("dsn:%v", cfg)
}
func Test_mysqlHealthCheck_LoadExtConfig(t *testing.T) {
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
			args:    args{extConfig: DsnJSONConfig},
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
func TestMysqlHealthCheck_DoHealthCheck(t *testing.T) {
	err := healthCheck.LoadExtConfig(DsnJSONConfig)
	assert.Nil(t, err)
	resHealthCheck, err := healthCheck.DoHealthCheck()
	assert.Nil(t, err)
	assert.Equal(t, true, resHealthCheck.CheckResult.IsSuccess)
	t.Logf("resHealthCheck:%+v,checkResult:%+v", resHealthCheck, resHealthCheck.CheckResult)

}
