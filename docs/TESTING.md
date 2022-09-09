# Testing

## Test Fixtures

All test fixtures are stored in the `test/fixtures` folder and any subfolders required for organization. These files are
made available using the following import and code:

```go
import "go.ketch.com/lib/orlop/v2/test"

func TestSomething(t *testing.T) {
	// This opens the test/fixtures/foo/bar.json file from assets
    f, err := test.Fixtures.Open("foo/bar.json")
}

```

## Unit testing

All unit tests are in files sitting in the same package folder as the units under test.

To unit test this repository, run the dependencies as described in [RUNNING](RUNNING.md).

All unit tests should be created with the following build tag:

```go
//go:build unit && !integration && !smoke

```

Then you can run unit tests using Go Test:

```shell
go test -v --tags unit ./...
```

You can also set the `unit` build tag in your IDE.
