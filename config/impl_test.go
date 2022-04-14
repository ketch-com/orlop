package config

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.ketch.com/lib/orlop/v2/env"
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

func Test_providerImpl_Load(t *testing.T) {
	os.Setenv("TEST_CONFIG_EMBEDDED", "true")
	os.Setenv("TEST_CONFIG_REQ", "imhere")
	os.Setenv("TEST_CONFIG_SLICED", "\"a\",\"b\",\"c\"")
	os.Setenv("TEST_CONFIG_CUSTOM", "1m")
	os.Setenv("TEST_CONFIG_MAP", "[\"a=b\",\"c=d\"]")
	os.Setenv("TEST_CONFIG_HEX_ENCODED", "0102030405060708090A0B0C0D0E0F")
	os.Setenv("TEST_CONFIG_BASE_64_ENCODED", "AQIDBAUGBwgJCgsMDQ4P")
	os.Setenv("TEST_CONFIG_PTR", "123")

	var config TestConfig

	err := New(env.NewEnviron("test")).Load(context.Background(), "config", &config)
	require.NoError(t, err)

	assert.True(t, config.Embedded.Embedded)
	assert.Equal(t, "/pki/issue", config.WithDefault)
	assert.Equal(t, "imhere", config.Required)
	assert.NotEmpty(t, config.SomeSlice)
	assert.Equal(t, []string{"a", "b", "c"}, config.SomeSlice)
	assert.Equal(t, time.Minute, config.CustomParser)
	assert.Equal(t, map[string]string{"a": "b", "c": "d"}, config.Map)
	assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}, config.HexEncoded)
	require.NotNil(t, config.Ptr)
	assert.Equal(t, int32(123), *config.Ptr)
}
