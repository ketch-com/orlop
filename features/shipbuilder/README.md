# shipbuilder

The shipbuilder allows configuring shipbuilder operations without modifying the scripts (usually).

## Configuring

The following variables are provided to control Shipbuilder:

| variable                      | description                                          |
|-------------------------------|------------------------------------------------------|
| `shipbuilder_overwrite`       | space separated list of paths to overwrite on reinit |
| `shipbuilder_exclude`         | space separated list of paths to exclude on reinit   |
| `shipbuilder_generate`        | targets to generate in `./scripts/generate.sh`       |
| `shipbuilder_go_version`      | pinned Go version                                    |
| `shipbuilder_go_os`           | operating systems to build go binaries for           |
| `shipbuilder_go_linux_arch`   | architectures to build Linux binaries for            |
| `shipbuilder_go_windows_arch` | architectures to build Windows binaries for          |
| `shipbuilder_go_darwin_arch`  | architectures to build Mac binaries for              |
| `shipbuilder_go_cgo_enabled`  | 1 to enable cgo and 0 to disable cgo                 |
| `shipbuilder_build`           | which build scripts to run                           |
| `shipbuilder_type`            | type of repository (library, service)                |
| `shipbuilder_module`          | primary module of the repository                     |
