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
