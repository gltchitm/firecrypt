package main

import (
	"os"
	"log"
	"path"
	"strings"
	"runtime"
	"io/ioutil"
	"crypto/sha512"
	"encoding/base64"
	"firecrypt/internal/crypt"
	"firecrypt/internal/profile"
	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
)

type DecodedMessage struct {
	Name string
	Detail []string
}
func decodeMessage(message *astilectron.EventMessage) DecodedMessage {
	var data string
	message.Unmarshal(&data)

	var splitData = strings.Split(data, ",")

	var name, err = base64.StdEncoding.DecodeString(splitData[0])

    if err != nil {
		panic(err)
	}

	var detail []byte

	if detail, err = base64.StdEncoding.DecodeString(
		splitData[1],
	); err != nil {
		panic(err)
	}

	return DecodedMessage {
		Name: string(name),
		Detail: strings.Split(string(detail), ","),
	}
}

func main() {
	var logger = log.New(log.Writer(), log.Prefix(), log.Flags())

	logger.SetOutput(ioutil.Discard)

	var app, err = astilectron.New(logger, astilectron.Options {
		AppName: "Firecrypt",
		BaseDirectoryPath: "firecrypt",
	})

	if err != nil {
		panic(err)
	}

	defer app.Close()

	app.HandleSignals()

	if err = app.Start(); err != nil {
		panic(err)
	}

	var window *astilectron.Window

	if window, err = app.NewWindow("./electron/firecrypt.html", &astilectron.WindowOptions {
		Resizable: astikit.BoolPtr(false),
		Center: astikit.BoolPtr(true),
		Height: astikit.IntPtr(330),
		Width: astikit.IntPtr(330),
		AlwaysOnTop: astikit.BoolPtr(true),
		WebPreferences: &astilectron.WebPreferences {
			DevTools: astikit.BoolPtr(false),
		},
	}); err != nil {
		panic(err)
	}

	var menu = app.NewMenu([] *astilectron.MenuItemOptions {
		{
			Label: astikit.StrPtr("Testing"),
			SubMenu: [] *astilectron.MenuItemOptions {
				{ Role: astikit.StrPtr("about") },
				{ Type: astikit.StrPtr("separator") },
				{ Role: astikit.StrPtr("hide") },
				{ Role: astikit.StrPtr("hideothers") },
				{ Role: astikit.StrPtr("unhide") },
				{ Type: astikit.StrPtr("separator") },
				{ Role: astikit.StrPtr("quit") },
			},
		},
	})
	menu.Create()

	if err = window.Create(); err != nil {
		panic(err)
	}

	window.OnMessage(func (message *astilectron.EventMessage) interface {} {
		var decodedMessage = decodeMessage(message)

		if decodedMessage.Name == "is-macos" {
			return runtime.GOOS == "darwin"
		} else if decodedMessage.Name == "load-profiles" {
			return profile.GetProfiles()
		} else if decodedMessage.Name == "is-profile-open" {
			if len(os.Args) > 1	&& os.Args[1] == "--no-check-profile-open" {
				return false
			}
			return profile.IsProfileOpen(decodedMessage.Detail[0])
		} else if decodedMessage.Name == "get-hash" {
			var hashedPw []byte
            if hashedPw, err = os.ReadFile(
				path.Join(
					decodedMessage.Detail[0],
					".__firecrypt_hash__",
				),
			); err != nil {
				panic(err)
			}

			return string(hashedPw)
		} else if decodedMessage.Name == "hash-password" {
			var hashedPwBytes = sha512.Sum512([]byte(decodedMessage.Detail[0]))
			for i := 0; i < 249999; i++ {
				hashedPwBytes = sha512.Sum512(hashedPwBytes[:])
			}
			return base64.StdEncoding.EncodeToString(hashedPwBytes[:])
		} else if decodedMessage.Name == "set-password" {
			profile.SetPassword(decodedMessage.Detail[0], decodedMessage.Detail[1])
		} else if decodedMessage.Name == "lock-profile" {
			crypt.LockProfile(decodedMessage.Detail[0], decodedMessage.Detail[1])
		} else if decodedMessage.Name == "unlock-profile" {
			return crypt.UnlockProfile(decodedMessage.Detail[0], decodedMessage.Detail[1])
		} else if decodedMessage.Name == "launch-profile" {
			app.Dock().Hide()
			window.Hide()
			profile.LaunchProfile(decodedMessage.Detail[0])
			app.Dock().Show()
			window.Show()
		}

		return nil
	})

	app.Wait()
}