package config

import (
	"context"
	"os"
	"sort"
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

func Test_providerImpl_List(t *testing.T) {
	var config TestConfig

	p := New(Params{
		Environ: env.NewEnviron("test"),
		Prefix:  "test",
		Defs: []Definition{
			{
				Name:   "config",
				Config: &config,
			},
		},
	})

	vars, err := p.List(context.Background())
	require.NoError(t, err)

	sort.Strings(vars)

	assert.Equal(t, []string{
		"TEST_CONFIG_BASE64_ENCODED=# [v1, v2, v3]",
		"TEST_CONFIG_CUSTOM=12345s",
		"TEST_CONFIG_DEF=/pki/issue",
		"TEST_CONFIG_EMBEDDED=false # bool",
		"TEST_CONFIG_HEX_ENCODED=# [v1, v2, v3]",
		"TEST_CONFIG_MAP=# [k=v, k=v, k=v]",
		"TEST_CONFIG_PTR=",
		"TEST_CONFIG_REQ=# string",
		"TEST_CONFIG_SLICED=# [v1, v2, v3]",
		"TEST_CONFIG_UNKNOWN=0 # int",
	}, vars)
}

func Test_providerImpl_Get(t *testing.T) {
	var config TestConfig
	ctx := context.Background()

	p := New(Params{
		Environ: env.NewEnviron("test"),
		Prefix:  "test",
		Defs: []Definition{
			{
				Name:   "config",
				Config: &config,
			},
		},
	})

	os.Setenv("TEST_CONFIG_EMBEDDED", "true")
	os.Setenv("TEST_CONFIG_REQ", "imhere")
	os.Setenv("TEST_CONFIG_SLICED", "\"a\",\"b\",\"c\"")
	os.Setenv("TEST_CONFIG_CUSTOM", "1m")
	os.Setenv("TEST_CONFIG_MAP", "[\"a=b\",\"c=d\"]")
	os.Setenv("TEST_CONFIG_HEX_ENCODED", "0102030405060708090A0B0C0D0E0F")
	os.Setenv("TEST_CONFIG_BASE_64_ENCODED", "AQIDBAUGBwgJCgsMDQ4P")
	os.Setenv("TEST_CONFIG_PTR", "123")

	val, err := p.Get(ctx, "config")
	require.NoError(t, err)

	if c, ok := val.(*TestConfig); ok {
		assert.True(t, c.Embedded.Embedded)
		assert.Equal(t, "/pki/issue", c.WithDefault)
		assert.Equal(t, "imhere", c.Required)
		assert.NotEmpty(t, c.SomeSlice)
		assert.Equal(t, []string{"a", "b", "c"}, c.SomeSlice)
		assert.Equal(t, time.Minute, config.CustomParser)
		assert.Equal(t, map[string]string{"a": "b", "c": "d"}, c.Map)
		assert.Equal(t, []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F}, c.HexEncoded)
		require.NotNil(t, c.Ptr)
		assert.Equal(t, int32(123), *c.Ptr)
	} else {
		t.Fatal("invalid type")
	}

}

func Test_providerImpl_load(t *testing.T) {
	os.Setenv("TEST_CONFIG_EMBEDDED", "true")
	os.Setenv("TEST_CONFIG_REQ", "imhere")
	os.Setenv("TEST_CONFIG_SLICED", "\"a\",\"b\",\"c\"")
	os.Setenv("TEST_CONFIG_CUSTOM", "1m")
	os.Setenv("TEST_CONFIG_MAP", "[\"a=b\",\"c=d\"]")
	os.Setenv("TEST_CONFIG_HEX_ENCODED", "0102030405060708090A0B0C0D0E0F")
	os.Setenv("TEST_CONFIG_BASE_64_ENCODED", "AQIDBAUGBwgJCgsMDQ4P")
	os.Setenv("TEST_CONFIG_PTR", "123")

	var config TestConfig
	p := &providerImpl{
		configs: map[string]value{
			"config": {
				isPopulated: false,
				value:       &config,
			},
		},
		environ: env.NewEnviron("test"),
		prefix:  "test",
	}

	err := p.load(context.Background(), "config", &config)
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
