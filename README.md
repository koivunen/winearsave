# Windows Ear Saver

Gradually lower windows volume to potentially save your ears when you're too focused on things to do it yourself when loud music starts playing. A hello world for myself to test go and to make an experiment.

## Prerequisites

- Go 1.8 or later
- `go-wca` (github.com/moutend/go-wca)

## Build

```shell
go build -ldflags "-H windowsgui" 
```

## Usage

```shell
winearsave.exe
```

## TODO

- Handle sound card vanishing (bluetooth)
- Handle multiple sound cards
- Disable in games
- More techniques for lowering volume
- Log volume statistics

## Credits
- moutend for Windows Core Audio API bindings and example
- Wikipedia for logo