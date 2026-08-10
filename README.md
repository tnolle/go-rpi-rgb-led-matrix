# RGB LED matrix

This repository contains the current Go applications and will host their gradual Rust replacements.

## Layout

- `go/servers/rpi`: production Raspberry Pi server
- `go/servers/terminal`: true-color terminal emulator server
- `go/clients/led`: command-line client
- `go/internal`: packages shared by the Go applications
- `rust`: incremental Rust replacements and new clients
- `assets`: tracked assets such as fonts
- `third_party/rpi-rgb-led-matrix`: shared C++ hardware driver
- `deployments`: service definitions
- `scripts`: development and hardware helper scripts

## Development

```sh
go -C go build -o ../.bin/server-terminal ./servers/terminal
./.bin/server-terminal
go -C go run ./clients/led --help
go -C go test ./...
```

The CLI reads available content from the selected server:

```sh
go -C go run ./clients/led get
go -C go run ./clients/led get image
go -C go run ./clients/led get gif
go -C go run ./clients/led get dashboard
go -C go run ./clients/led get animation
```

Set the content displayed by the selected server:

```sh
go -C go run ./clients/led set image autodarts
go -C go run ./clients/led set gif celebration
go -C go run ./clients/led set gif --once success
go -C go run ./clients/led set dashboard clock
go -C go run ./clients/led set animation plasma
```

The emulator expects `config.toml` and image assets under `images`. Both are intentionally ignored by Git.

The Raspberry Pi server requires Linux, cgo, and the `with_cgo` build tag.

The server's HTTP interface is documented in [docs/api.md](docs/api.md).
