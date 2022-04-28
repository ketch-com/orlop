package cmd

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/spf13/cobra"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.ketch.com/lib/orlop/v2/config"
	"go.uber.org/fx"
)

type EmbeddedConfig struct {
	Embedded bool
}

type TestConfig struct {
	Embedded      EmbeddedConfig `config:","`
	WithDefault   string         `config:"def,default=/pki/issue"`
	Required      string         `config:"req,required"`
	SomeSlice     []string       `config:"sliced"`
	CustomParser  time.Duration  `config:"custom,default=12345s"`
	Map           map[string]string
	HexEncoded    []byte
	Base64Encoded []byte `config:",encoding=base64"`
	Ptr           *int32
	Unknown       int32
}

type TestString struct {
	String string
}

func TestRun(t *testing.T) {
	os.Setenv("TEST_CONFIG_EMBEDDED", "true")
	os.Setenv("TEST_CONFIG_REQ", "imhere")
	os.Setenv("TEST_CONFIG_SLICED", "\"a\",\"b\",\"c\"")
	os.Setenv("TEST_CONFIG_CUSTOM", "1m")
	os.Setenv("TEST_CONFIG_MAP", "[\"a=b\",\"c=d\"]")
	os.Setenv("TEST_CONFIG_HEX_ENCODED", "0102030405060708090A0B0C0D0E0F")
	os.Setenv("TEST_CONFIG_BASE_64_ENCODED", "AQIDBAUGBwgJCgsMDQ4P")
	os.Setenv("TEST_CONFIG_PTR", "123")
	os.Setenv("TEST_STRING", "string-data")

	var cfg TestConfig
	var data TestString

	var module = fx.Options(
		config.Option[TestConfig]("config"),
		config.Option[TestString](),
		fx.Invoke(func(t TestConfig, s TestString) {
			cfg = t
			data = s
			return
		}),
		fx.Invoke(
			func(lifecycle fx.Lifecycle, s fx.Shutdowner) {
				lifecycle.Append(
					fx.Hook{
						OnStart: func(_ context.Context) error {
							return s.Shutdown()
						},
					},
				)
			},
		),
	)

	Run("test", module)

	assert.True(t, cfg.Embedded.Embedded)
	assert.Equal(t, "/pki/issue", cfg.WithDefault)
	assert.Equal(t, "imhere", cfg.Required)
	assert.NotEmpty(t, cfg.SomeSlice)
	assert.Equal(t, []string{"a", "b", "c"}, cfg.SomeSlice)
	assert.Equal(t, time.Minute, cfg.CustomParser)
	assert.Equal(t, map[string]string{"a": "b", "c": "d"}, cfg.Map)
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}, cfg.HexEncoded)
	require.NotNil(t, cfg.Ptr)
	assert.Equal(t, int32(123), *cfg.Ptr)
	assert.Equal(t, "string-data", data.String)
}

func TestInit(t *testing.T) {
	var module = fx.Options(
		config.Option[TestConfig]("config"),
		config.Option[TestString](),
		fx.Invoke(
			func(lifecycle fx.Lifecycle, s fx.Shutdowner) {
				lifecycle.Append(
					fx.Hook{
						OnStart: func(_ context.Context) error {
							return s.Shutdown()
						},
					},
				)
			},
		),
	)

	var cmd = &cobra.Command{
		Use:              "test",
		TraverseChildren: true,
		SilenceUsage:     true,
	}

	NewRunner("test").SetupRoot(cmd).Setup(cmd, module)

	cmd.SetArgs([]string{"init"})
	err := cmd.Execute()
	require.NoError(t, err)
}
