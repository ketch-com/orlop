package cmd

import (
	"context"
	"os"
	"os/exec"
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

	os.Setenv("TEST_SECOND_EMBEDDED", "false")
	os.Setenv("TEST_SECOND_REQ", "imherealso")
	os.Setenv("TEST_SECOND_SLICED", "\"e\",\"f\",\"g\"")
	os.Setenv("TEST_SECOND_CUSTOM", "5m")
	os.Setenv("TEST_SECOND_MAP", "[\"e=f\",\"g=h\"]")
	os.Setenv("TEST_SECOND_HEX_ENCODED", "0102030405060708090A0B0C0D0E0F")
	os.Setenv("TEST_SECOND_BASE_64_ENCODED", "AQIDBAUGBwgJCgsMDQ4P")
	os.Setenv("TEST_SECOND_PTR", "456")

	os.Setenv("TEST_STRING", "string-data")

	var cfg1, cfg2 TestConfig
	var data TestString

	type testParams struct {
		fx.In

		T1 TestConfig
		T2 TestConfig `name:"2nd"`
		S  TestString
	}

	var module = fx.Options(
		config.Option[TestConfig]("config"),
		config.Option[TestConfig]("second", "2nd"),
		config.Option[TestString](),

		fx.Invoke(
			func(p testParams) {
				cfg1, cfg2, data = p.T1, p.T2, p.S
				return
			},
		),
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

	assert.True(t, cfg1.Embedded.Embedded)
	assert.Equal(t, "/pki/issue", cfg1.WithDefault)
	assert.Equal(t, "imhere", cfg1.Required)
	assert.NotEmpty(t, cfg1.SomeSlice)
	assert.Equal(t, []string{"a", "b", "c"}, cfg1.SomeSlice)
	assert.Equal(t, time.Minute, cfg1.CustomParser)
	assert.Equal(t, map[string]string{"a": "b", "c": "d"}, cfg1.Map)
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}, cfg1.HexEncoded)
	require.NotNil(t, cfg1.Ptr)
	assert.Equal(t, int32(123), *cfg1.Ptr)

	assert.False(t, cfg2.Embedded.Embedded)
	assert.Equal(t, "/pki/issue", cfg2.WithDefault)
	assert.Equal(t, "imherealso", cfg2.Required)
	assert.NotEmpty(t, cfg2.SomeSlice)
	assert.Equal(t, []string{"e", "f", "g"}, cfg2.SomeSlice)
	assert.Equal(t, 5*time.Minute, cfg2.CustomParser)
	assert.Equal(t, map[string]string{"e": "f", "g": "h"}, cfg2.Map)
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}, cfg2.HexEncoded)
	require.NotNil(t, cfg2.Ptr)
	assert.Equal(t, int32(456), *cfg2.Ptr)

	assert.Equal(t, "string-data", data.String)
}

func TestInit(t *testing.T) {
	var module = fx.Options(
		config.Option[TestConfig]("config"),
		config.Option[TestConfig]("second", "2nd"),
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

func TestOptions(t *testing.T) {
	module := fx.Options(
		fx.Invoke(
			func(lifecycle fx.Lifecycle, s fx.Shutdowner) {
				lifecycle.Append(
					fx.Hook{
						OnStart: func(_ context.Context) error {
							return s.Shutdown()
						},
						OnStop: func(ctx context.Context) error {
							time.Sleep(fx.DefaultTimeout + time.Second)
							return nil
						},
					},
				)
			},
		),
	)

	if os.Getenv("CRASH") == "0" {
		testOptions(module, fx.StopTimeout(2*fx.DefaultTimeout))
		return
	} else if os.Getenv("CRASH") == "1" {
		testOptions(module)
		return
	}

	crashCmd := exec.Command(os.Args[0], "-test.run=TestOptions")
	crashCmd.Env = append(os.Environ(), "CRASH=1")
	err := crashCmd.Run()
	assert.Error(t, err, "process did not return an error", err)
	assert.IsType(t, &exec.ExitError{}, err, "process did not return an ExitError", err)

	successCmd := exec.Command(os.Args[0], "-test.run=TestOptions")
	successCmd.Env = append(os.Environ(), "CRASH=0")
	err = successCmd.Run()
	assert.NoError(t, err, "process returned error %v, want exit status 0", err)
}

func testOptions(module fx.Option, options ...fx.Option) error {
	options = append(options, module)

	var cmd = &cobra.Command{
		Use:              "test",
		TraverseChildren: true,
	}

	NewRunner("test").SetupRoot(cmd).Setup(cmd, options...)

	return cmd.Execute()
}
