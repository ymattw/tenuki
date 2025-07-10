# Tenuki - Go elsewhere, to the terminal

**Tenuki** is a terminal-based OGS ([Online-Go Server](https://online-go.com))
client built with Go and the [tview](https://github.com/rivo/tview) UI library.
It lets you escape the browser and enjoy your correspondence games right from
your favorite terminal.

## Features

- List your active games
- Play and chat
- Watch top live games

## Limitations

Tenuki is designed primarily for **correspondence games**. The following
limitations apply (for now):

- No support for creating or accepting challenges
- No auto match functionality

## Usage

### Requirements

- A terminal that supports emoji rendering and a font with good Unicode
  coverage
- An [OGS OAuth2 Application](https://online-go.com/oauth2/applications/), with
  `Authorization grant type` set to **Resource owner password-based**

  <kbd>
    <img alt="Register a new application" src="https://github.com/ymattw/tenuki/blob/main/screenshots/register.png?raw=true"  />
  </kbd>

### Run the app

#### Run from binary

Download the binary from [Release](https://github.com/ymattw/tenuki/releases)
that are built from automated GitHub Actions, decompress, and run from your
terminal.

```bash
./tenuki -help
./tenuki
```

#### Run from source

Requires Go **1.18+** ([Download and install](https://go.dev/doc/install)).
Clone the code, enter the directory and then:

```bash
go run . -help
go run .
```

## Screenshots

Screenshots taken on macOS using iTerm2 with the Monaco font (size 14).

*Play mode with default night board theme:*

<img alt="Play mode with night board theme" src="https://github.com/ymattw/tenuki/blob/main/screenshots/play-night-theme.png?raw=true" width="500" />

*Play mode with oak board theme:*

<img alt="Play mode with oak board theme" src="https://github.com/ymattw/tenuki/blob/main/screenshots/play-oak-theme.png?raw=true" width="500" />

*Home page showing your active games:*

<img alt="Home page showing your active games" src="https://github.com/ymattw/tenuki/blob/main/screenshots/home.png?raw=true" width="500" />

*Watch page showing top live games:*

<img alt="Watch page showing top live games" src="https://github.com/ymattw/tenuki/blob/main/screenshots/watch.png?raw=true" width="500" />

## Acknowledgments

This project is inspired by and a spiritual successor to
[termsuji](https://github.com/lvank/termsuji), with OGS library decoupled to
[googs](https://github.com/ymattw/googs).
