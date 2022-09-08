# scripts

This folder contains build and utility scripts.

## build.sh

This script builds the binaries for the repository, if any are produced.

To build all platforms, use the following:
```shell
./scripts/build.sh
```

To build specific platform (e.g., `linux` and `darwin`), use the following:
```shell
./scripts/build.sh linux darwin
```

## bumpdep.sh

This script bumps dependencies.

If you have a clean local source, you can create a dependency branch and push to origin using the following:
```shell
./scripts/bumpdep.sh
```

To bump dependencies in your local, without creating a separate git branch, use the following:
```shell
./scripts/bumpdep.sh --local
```

## generate.sh

This script generates code and is used by `generate.go`.

To generate everything, use the following:

```shell
./scripts/generate.sh
```

To run a specific generator (e.g., `mocks`), use the following:

```shell
./scripts/generate.sh mocks
```

The following generators are available by default:
* `resources` - generates resource code from resource definitions
* `services` - generates service code from service definitions
* `types` - generates typesystems from typesystem specifications
* `gateway` - generates gateway stubs from openapi files
* `proto` - generates protobuf stubs from protobuf IDL files
* `mocks` - generates mocks using mockery

To specify multiple generators (e.g., `gateway` and `proto`), use the following:

```shell
./scripts/generate.sh gateway proto
```

## integrationtest.sh

This script runs integration tests in the repository.

To run all integration tests, use the following:

```shell
./scripts/integrationtest.sh
```

To run specific integration tests (e.g., in `./test/integration/foo/...`), use the following:

```shell
./scripts/integrationtest.sh ./test/integration/foo/...
```

## reinit.sh

This script re-runs shipbuilder init on the repository.

To reinitialize the repository, use the following:

```shell
./scripts/reinit.sh
```

Any additional arguments are added to the shipbuilder command.

## smoketest.sh

This script runs the smoketest in docker compose.

To run the smoketest, use the following:

```shell
./scripts/smoketest.sh
```

To prune the docker setup after the smoketest, use the following:

```shell
./scripts/smoketest.sh --prune
```

## unittest.sh

This script runs unit tests in the repository.

To run all integration tests, use the following:

```shell
./scripts/unittest.sh
```

To run specific unit tests (e.g., in `./pkg/db/impl/...`), use the following :

```shell
./scripts/unittest.sh ./pkg/db/impl/...
```

