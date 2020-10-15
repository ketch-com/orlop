// Copyright (c) 2020 SwitchBit, Inc.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package orlop

//go:generate mockgen -package=orlop_test -source=cert_generation_config.go -destination=cert_generation_config_mock_test.go
//go:generate mockgen -package=orlop_test -source=client_config.go -destination=client_config_mock_test.go
//go:generate mockgen -package=orlop_test -source=credentials.go -destination=credentials_mock_test.go
//go:generate mockgen -package=orlop_test -source=enabled_config.go -destination=enabled_config_mock_test.go
//go:generate mockgen -package=orlop_test -source=file_config.go -destination=file_config_mock_test.go
//go:generate mockgen -package=orlop_test -source=key_config.go -destination=key_config_mock_test.go
//go:generate mockgen -package=orlop_test -source=server_config.go -destination=server_config_mock_test.go
//go:generate mockgen -package=orlop_test -source=server_options.go -destination=server_options_mock_test.go
//go:generate mockgen -package=orlop_test -source=tls_config.go -destination=tls_config_mock_test.go
//go:generate mockgen -package=orlop_test -source=token_config.go -destination=token_config_mock_test.go
//go:generate mockgen -package=orlop_test -source=vault_config.go -destination=vault_config_mock_test.go
