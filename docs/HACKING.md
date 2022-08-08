# Hacking

## Get the code

```shell
$ git clone https://github.com/ketch-com/orlop
$ cd orlop/
```

## Getting the development dependencies

```shell
go get -u ./...
```

## Building

You can build this project using Go:

```shell
go build ./...
```

## Updating dependencies

To update the dependencies, run the following:

```shell
rm go.sum
go get -u -t ./...
```

