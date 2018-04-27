# Windows Ear Saver

Gradually lower windows volume to potentially save your ears when you're too focused on things to do it yourself when loud music starts playing. A hello world for myself to test go and to make an experiment.

## Warning

Unresearched, may actually contribute to hearing damage instead of preventing it.

## Prerequisites

- Go

## Build

```shell
go build -ldflags "-H windowsgui" 
```

## Usage

```shell
winearsave.exe
```

## TODO

- When no audio, do not lower
- Handle sound card vanishing (bluetooth)
- Handle multiple sound cards
- Disable in games
- More techniques for lowering volume
- Log volume statistics

## Credits
- moutend for Windows Core Audio API bindings and example
