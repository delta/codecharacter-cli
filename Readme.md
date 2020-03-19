# Code Character CLI

CLI for Code Character to compile, execute and render the game offline

-------------------------------------------

## Requirements
- Docker

## Usage instructions

To compile, execute and render :
```bash
./codecharacter-cli [ -p PORT ] [ -m MAP_FILE ] PLAYER_1_CODE PLAYER_2_CODE
```
To render :

```bash
./codecharacter-cli -r [ -p PORT ] LOG_DIRECTORY 
```
Given directory must contain the following files : game.log, player_1.dlog, player_2.dlog

## Dev Instructions
- Copy the built `index.html` and `libpack.js` from renderer.
- Run `go-bindata index.html libpack.js`.
- Change compiler and runner images in `constants.go` if needed.

## Build Instructions
- Run `go build` to build for current platform.
- Run `build.sh` to build for common platforms.
