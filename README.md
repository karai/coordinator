![karai_github_banner](https://user-images.githubusercontent.com/34389545/80034381-f6a14d00-84b3-11ea-857a-638322dac890.png)

[![Discord](https://img.shields.io/discord/388915017187328002?label=Join%20Discord)](http://chat.turtlecoin.lol) [![GitHub issues](https://img.shields.io/github/issues/karai/go-karai?label=Issues)](https://github.com/karai/go-karai/issues) ![GitHub stars](https://img.shields.io/github/stars/karai/go-karai?label=Github%20Stars) ![Build](https://github.com/karai/go-karai/workflows/Build/badge.svg) ![GitHub](https://img.shields.io/github/license/karai/go-karai) ![GitHub issues by-label](https://img.shields.io/github/issues/karai/go-karai/Todo) [![Go Report Card](https://goreportcard.com/badge/github.com/karai/go-karai)](https://goreportcard.com/report/github.com/karai/go-karai)

**Website:** [📝 karai.io](https://karai.io)  **Browse:** [💻 Karai Pointer Explorer](https://karaiexplorer.extrahash.org/)  **Read:** [🔗 Official Karai Blog](https://karai.io)

## Usage
>Note: Karai aims to always compile and run on **Linux** targetting the **AMD64** CPU architecture. Other operating systems and architectures may compile and run this software but should do so expecting some inconsistencies.

**Launch Karai**
```
./go-karai --coordinator --https --showip
```

**Launch Options**
```
--coordinator
```
This will launch the go-karai client in Coordinator mode. This is needed for some functions. 

```
--https
```
Use https with Coordinator API. This will use Autocert to automagically fetch the needed https certs for your channel.

```
--showip
```
Launching as Coordinator prints a lot of diagnostic information while running. This will enable showing your IP in those log messages.

>Type `menu` to view a list of functions. Functions that are darkened are disabled.

## Dependencies

-   Golang 1.13+ [[Download]](https://golang.org)
-   TurtleCoin Daemon & Wallet-API [[Download]](http://latest.turtlecoin.lol)
-   IPFS [Download](https://github.com/ipfs/go-ipfs/releases/latest)

## Operating System

Karai supports Linux on AMD64 architecture, but may compile in other settings. Differences between Linux and non-Linux installs should be expected.

## Building

```bash
git clone https://github.com/karai/go-karai

cd go-karai

go mod init github.com/karai/go-karai

go build && ./go-karai
```

**Optional:** Compile with all errors displayed, then run binary. Avoids "too many errors" from hiding error info.

`go build -gcflags="-e" && ./go-karai `

## Contributing

-   MIT License
-   `gofmt` is used on all files.
-   go modules are used to manage dependencies.

## Thanks to:

[![turtlecoin](https://user-images.githubusercontent.com/34389545/80266529-fb0b6880-8661-11ea-9a75-4cb066834775.png)](https://turtlecoin.lol)
[![IPFS](https://user-images.githubusercontent.com/34389545/80266356-0c07aa00-8661-11ea-8308-84639318213a.png)](https://ipfs.io)
[![LibP2P](https://user-images.githubusercontent.com/34389545/80266502-e4651180-8661-11ea-8367-54bf59e26470.png)](https://libp2p.io)
[![GOLANG](https://user-images.githubusercontent.com/34389545/80266422-6b65ba00-8661-11ea-836a-d1904ec15b94.png)](https://golang.org)
