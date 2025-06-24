# Tenuki - Go elsewhere, to the terminal

**Tenuki** is a terminal-based OGS ([Online-Go Server](https://online-go.com))
client built with Go and the [tview](https://github.com/rivo/tview) UI library.
It lets you escape the browser and enjoy your correspondence games right from
your favorite terminal.

## ✨ Features

- List your active games
- Play and chat
- Watch top live games

## ⚠️ Limitations

Tenuki is designed primarily for **correspondence games**. The following
limitations apply (for now):

- No support for creating or accepting challenges
- No auto match functionality
- Only the **byo-yomi** clock system is fully supported

## 🚀 Usage

### Requirements

- Go **1.18+**
- An [OGS application](https://online-go.com/oauth2/applications/), with
  `Authorization grant type` set to *Resource owner password-based*
- A terminal that supports emoji rendering and a font with good Unicode
  coverage

### Run the app

```bash
go run .
```

## 📸 Screenshots

Screenshots taken on macOS using iTerm2 with the Monaco font (size 14).

![Home page showing game list](https://github.com/ymattw/tenuki/blob/main/screenshots/home.png?raw=true "Home page showing game list")
![Watch page showing top live games](https://github.com/ymattw/tenuki/blob/main/screenshots/watch.png?raw=true "Watch page showing top live games")
![Play mode with night board theme](https://github.com/ymattw/tenuki/blob/main/screenshots/play-night-theme.png?raw=true "Play mode with night board theme")
![Play mode with oak board theme](https://github.com/ymattw/tenuki/blob/main/screenshots/play-oak-theme.png?raw=true "Play mode with oak board theme")

## 🙏 Acknowledgments

This project is inspired by and a spiritual successor to
[termsuji](https://github.com/lvank/termsuji), with OGS library decoupled to
[googs](https://github.com/ymattw/googs).
