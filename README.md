# RGB LED matrix

This repository contains the current Go applications and will host their gradual Rust replacements.

## Layout

- `go/servers/rpi`: production Raspberry Pi server
- `go/servers/emulator`: macOS emulator server
- `go/servers/terminal`: true-color terminal emulator server
- `go/clients/cli`: command-line client
- `go/internal`: packages shared by the Go applications
- `rust`: incremental Rust replacements and new clients
- `assets`: tracked assets such as fonts
- `third_party/rpi-rgb-led-matrix`: shared C++ hardware driver
- `deployments`: service definitions
- `scripts`: development and hardware helper scripts

The root `go.work` file makes it possible to run Go commands from the repository root.

## Development

```sh
go run ./go/servers/emulator
go run ./go/servers/terminal
go run ./go/clients/cli --help
go test ./go/...
```

The emulator expects `config.toml` and image assets under `images`. Both are intentionally ignored by Git.

The Raspberry Pi server requires Linux, cgo, and the `with_cgo` build tag.
