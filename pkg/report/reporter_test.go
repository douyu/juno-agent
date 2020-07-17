package report

import (
	"github.com/douyu/juno-agent/pkg/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

var report *Report
var config Config

func instantiationReport() {
	buildConfig()
	report = &Report{
		config:   &config,
		Reporter: NewHTTPReport(&config),
	}
}
func buildConfig() {
	config = DefaultConfig()
	config.Enable = true
	config.Addr = "http://127.0.0.1:8080/report"
	config.Internal = 60
	config.RegionCode = "WUHAN"
	config.ZoneName = "Test"
}
func TestMain(m *testing.M) {
	instantiationReport()
	m.Run()
}
func TestReport_ReportAgentStatus(t *testing.T) {
	req := model.AgentReportRequest{
		Hostname:     report.config.HostName,
		IP:           appIP,
		AgentVersion: "0.0.0.1",
		RegionCode:   report.config.RegionCode,
		RegionName:   report.config.RegionName,
		ZoneCode:     report.config.ZoneCode,
		ZoneName:     report.config.ZoneName,
		Env:          report.config.Env,
	}
	reporterResp := report.Reporter.Report(req)
	assert.Equal(t, 0, reporterResp.Err)
}
