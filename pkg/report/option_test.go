package report

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestConfig_Build(t *testing.T) {
	config := DefaultConfig()
	want := &Report{
		config:   &config,
		Reporter: NewHTTPReport(&config),
	}
	t.Run("", func(t *testing.T) {
		if got := config.Build(); !reflect.DeepEqual(got.config, want.config) {
			t.Errorf("Build() = %#v, want %#v", got, want)
		}
		assert.Equal(t, "dev", config.Env)
	})

}
