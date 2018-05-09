# rfunc

rfunc is utility functions over the network. rfunc currently provides clipboard copy and paste functions and open URLs functions.

## Motivation

I spend most of my work time in a remote shell with tmux and vim. Sometimes, I feel like copying texts from remote shell to local desktop clipboard and opening URLs by local desktop browser. But text selecting on the terminal client is a bit hard when I split tmux panes. Therefore I want to use clipboard commands such as pbcopy/pbpaste, xclip, xsel and open commands such as xdg-open across the network.

## Features

* Cross platform (currently only support Linux and macOS)
* Clipboard copy and paste
* Open URLs

**rfunc does NOT support encryption. It is strongly recommended that you use port forwarding with ssh.**

## Install

Use `go get` or just download [binary releases](https://github.com/yukithm/rfunc/releases).

```
go get github.com/yukithm/rfunc
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
```

## Similar projects

* https://github.com/wincent/clipper
* https://github.com/pocke/lemonade

## Author

Yuki (@yukithm)

## License

MIT
