# orlop

[![CircleCI](https://circleci.com/gh/switch-bit/orlop.svg?style=svg&circle-token=549d5ece4247811389e5c7c0689ffce808778d5f)](https://circleci.com/gh/switch-bit/orlop)

Orlop is the base deck in a ship where the cables are stowed.

It is SwitchBit's standard (opinionated) library that all of our projects include.
* Configuration
* Logging (Logrus)
* Metrics (Prometheus)
* Server setup
* TLS - leveraging Vault
* Vault - secrets and certificates

It supports gRPC clients and servers, Swagger and automatic certificate generation
for mTLS.
