# rfunc

rfunc is utility functions over the network. rfunc currently provides clipboard copy and paste functions and open URLs functions.

<!-- TOC depthFrom:2 -->

- [Motivation](#motivation)
- [Features](#features)
- [Install](#install)
- [Usage](#usage)
- [Command name symlinks (busybox style)](#command-name-symlinks-busybox-style)
- [Configurations](#configurations)
- [Security](#security)
- [Development](#development)
    - [Prerequisites](#prerequisites)
    - [Development tools](#development-tools)
    - [Generate *.pb.go from *.proto files](#generate-pbgo-from-proto-files)
- [Similar projects](#similar-projects)
- [Author](#author)
- [License](#license)

<!-- /TOC -->

## Motivation

I spend most of my work time in a remote shell with tmux and vim. Sometimes, I feel like copying texts from remote shell to local desktop clipboard and opening URLs by local desktop browser. But text selecting on the terminal client is a bit hard when I split tmux panes. Therefore I want to use clipboard commands such as pbcopy/pbpaste, xclip, xsel and open commands such as xdg-open across the network.

## Features

* Cross platform (Linux, macOS and Windows)
* Clipboard copy and paste
* Open URLs
* Support TLS and client certificate

**It is strongly recommended that you use port forwarding with ssh and use server/client certificate for security reason. (See: [Security](#security))**

## Install

Use `go get` or just download [binary releases](https://github.com/yukithm/rfunc/releases).

```
go get -u github.com/yukithm/rfunc
```

## Usage

launch rfunc server on the local desktop:

```sh
rfunc server

# or daemonize
rfunc server --daemon
```

default binding address and port is 127.0.0.1:8299

login a remote server with port forwarding:

```sh
ssh -R 8299:127.0.0.1:8299 REMOTE_HOST
```

copy and paste on a remote server:

```sh
cat some.txt | rfunc copy
```

```sh
rfunc paste >clipboard.txt
```

open URL:

```sh
rfunc opne https://github.com/yukithm/rfunc
```

## Command name symlinks (busybox style)

You can create symbolic links that include sub-command names. That symbolic links are almost same as running `rfunc sub-command`.

examples:

```sh
ln -s rfunc copy

# 'copy' is equivalent to 'rfunc copy'
./copy <some.txt

# same as above (because 'pbcopy' includes 'copy' term)
ln -s rfunc pbcopy

# 'open' and 'xdg-open' are equivalent to 'rfunc open'
ln -s rfunc open
ln -s rfunc xdg-open
```

## Configurations

rfunc reads the first config file in following paths.

* ~/.config/rfunc/rfunc.toml
* ~/.rfunc.toml

Config file has following structure.
Each item is same as a command line option (See `rfunc --help`).

```toml
addr = "127.0.0.1:8299"
sock = "/path/to/socket"
logfile = "/path/to/logfile"
quiet = false
eol = "NATIVE"

[tls]
cert = "/path/to/cert.pem"
key = "/path/to/private.key"
ca = "/path/to/cacert.pem"
server-name = "override-server-name"
insecure = false

[server]
daemon = false
allow-commands = ["copy", "paste"]
```

## Security

You should use client and server certificate to protect rfunc server from other people that can access TCP ports on your PC.
If you don't use client certificate, your rfunc server accepts all access from anyone.


Run server with `--tls-***` options:

```sh
rfunc server --tls-cert=server.crt --tls-key=server.key --tls-ca=cacert-for-client.pem
```

`--tls-ca` specifies CA root for client's certificate.

Run client with `--tls-***` options:

```sh
rfunc paste --tls-cert=client.crt --tls-key=client.key --tls-ca=cacert-for-server.pem
```

`--tls-ca` specifies CA root for server's certificate.

You can set these values to configuration file.

## Development

### Prerequisites

- protoc (Protocol compiler for Protocol Buffers)  
  https://github.com/protocolbuffers/protobuf

On macOS, you can install it by Homebrew:

```sh
brew install protobuf
```

Other platforms, you can download pre-build binaries from the above URL.

### Development tools

```sh
make deps
```

### Generate *.pb.go from *.proto files

You need to regenerate `*.pb.go` and mocks when you edit *.proto files.

```sh
make proto
make mock
```

## Similar projects

* https://github.com/wincent/clipper
* https://github.com/pocke/lemonade

## Author

Yuki (@yukithm)

## License

MIT
