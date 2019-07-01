classzz
====
[![Build Status](https://travis-ci.org/bourbaki-czz/classzz.png?branch=master)](https://travis-ci.org/bourbaki-czz/classzz)
[![Go Report Card](https://goreportcard.com/badge/github.com/bourbaki-czz/classzz)](https://goreportcard.com/report/github.com/bourbaki-czz/classzz)
[![ISC License](http://img.shields.io/badge/license-ISC-blue.svg)](http://copyfree.org)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](http://godoc.org/github.com/bourbaki-czz/classzz)

classzz is an alternative full node bitcoin cash implementation written in Go (golang).

This project is a port of the [btcd](https://github.com/btcsuite/btcd) codebase to Bitcoin Cash. It provides a high powered
and reliable blockchain server which makes it a suitable backend to serve blockchain data to lite clients and block explorers
or to power your local wallet.

classzz does not include any wallet functionality by design as it makes the codebase more modular and easy to maintain. 
The [czzwallet](https://github.com/bourbaki-czz/czzwallet) is a separate application that provides a secure Bitcoin Cash wallet 
that communicates with your running classzz instance via the API.

## Table of Contents

- [Requirements](#requirements)
- [Install](#install)
  - [Install prebuilt packages](#install-pre-built-packages)
  - [Build from Source](#build-from-source)
- [Getting Started](#getting-started)
- [Documentation](#documentation)
- [Contributing](#contributing)
- [License](#license)

## Requirements

[Go](http://golang.org) 1.9 or newer.

## Install

### Install Pre-built Packages

The easiest way to run the server is to download a pre-built binary. You can find binaries of our latest release for each operating system at the [releases page](https://github.com/bourbaki-czz/classzz/releases).

### Build from Source

If you prefer to install from source do the following:

- Install Go according to the installation instructions here:
  http://golang.org/doc/install

- Run the following commands to obtain btcd, all dependencies, and install it:

```bash
go get github.com/bourbaki-czz/classzz
```

This will download and compile `classzz` and put it in your path.

If you are a classzz contributor and would like to change the default config file (`classzz.conf`), make any changes to `sample-classzz.conf` and then run the following commands:

```bash
go-bindata sample-classzz.conf  # requires github.com/go-bindata/go-bindata/
gofmt -s -w bindata.go
```

## Getting Started

To start classzz with default options just run:

```bash
./classzz
```

You'll find a large number of runtime options with the help flag. All of them can also be set in a config file.
See the [sample config file](https://github.com/bourbaki-czz/classzz/blob/master/sample-classzz.conf) for an example of how to use it.

```bash
./classzz --help
```

You can use the common json RPC interface through the `czzctl` command:

```bash
./czzctl --help

./czzctl --listcommands
```

Classzz separates the node and the wallet. Commands for the wallet will work when you are also running
[czzwallet](https://github.com/bourbaki-czz/czzwallet):

```bash
./czzctl -u username -P password --wallet getnewaddress
```

## Docker

Building and running `classzz` in docker is quite painless. To build the image:

```
docker build . -t classzz
```

To run the image:

```
docker run classzz
```

To run `czzctl` and connect to your `classzz` instance:

```
# Find the running classzz container.
docker ps

# Exec czzctl.
docker exec <container> czzctl <command>
```

## Documentation

The documentation is a work-in-progress.  It is located in the [docs](https://github.com/bourbaki-czz/classzz/tree/master/docs) folder.

## Contributing

Contributions are definitely welcome! Please read the contributing [guidelines](https://github.com/bourbaki-czz/classzz/blob/master/docs/code_contribution_guidelines.md) before starting.

## Security Disclosures

To report security issues please contact:

Chris Pacia (ctpacia@gmail.com) - GPG Fingerprint: 0150 2502 DD3A 928D CE52 8CB9 B895 6DBF EE7C 105C

or

Josh Ellithorpe (quest@mac.com) - GPG Fingerprint: B6DE 3514 E07E 30BB 5F40  8D74 E49B 7E00 0022 8DDD 

## License

classzz is licensed under the [copyfree](http://copyfree.org) ISC License.
