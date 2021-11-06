package main

import (
	"encoding/base64"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/asticode/go-astikit"
	"github.com/asticode/go-astilectron"
	"github.com/gltchitm/firecrypt/internal/crypt"
	"github.com/gltchitm/firecrypt/internal/profile"
)

type DecodedMessage struct {
	Name   string
	Detail []string
}

func decodeMessage(message *astilectron.EventMessage) DecodedMessage {
	var data string

	err := message.Unmarshal(&data)
	if err != nil {
		panic(err)
	}

	splitData := strings.Split(data, ",")

	name, err := base64.StdEncoding.DecodeString(splitData[0])
	if err != nil {
		panic(err)
	}

	detail, err := base64.StdEncoding.DecodeString(splitData[1])
	if err != nil {
		panic(err)
	}

	return DecodedMessage{
		Name:   string(name),
		Detail: strings.Split(string(detail), ","),
	}
}
func main() {
	logger := log.New(log.Writer(), log.Prefix(), log.Flags())

	logger.SetOutput(ioutil.Discard)

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	app, err := astilectron.New(logger, astilectron.Options{
		AppName:           "Firecrypt",
		BaseDirectoryPath: "firecrypt",
		VersionElectron:   "15.3.0",
		AppIconDarwinPath: path.Join(wd, "./electron/resources/icon/icon.icns"),
	})
	if err != nil {
		panic(err)
	}

	defer app.Close()

	app.HandleSignals()

	err = app.Start()
	if err != nil {
		panic(err)
	}

	window, err := app.NewWindow("./electron/firecrypt.html", &astilectron.WindowOptions{
		Resizable:   astikit.BoolPtr(false),
		Center:      astikit.BoolPtr(true),
		Height:      astikit.IntPtr(330),
		Width:       astikit.IntPtr(330),
		AlwaysOnTop: astikit.BoolPtr(true),
		WebPreferences: &astilectron.WebPreferences{
			DevTools: astikit.BoolPtr(false),
		},
	})
	if err != nil {
		panic(err)
	}

	menu := app.NewMenu([]*astilectron.MenuItemOptions{
		{
			Label: astikit.StrPtr("Firecrypt"),
		},
	})

	err = menu.Create()
	if err != nil {
		panic(err)
	}

	err = window.Create()
	if err != nil {
		panic(err)
	}

	window.OnMessage(func(message *astilectron.EventMessage) interface{} {
		decodedMessage := decodeMessage(message)

		if decodedMessage.Name == "is-macos" {
			return runtime.GOOS == "darwin"
		} else if decodedMessage.Name == "get-profiles" {
			return profile.GetProfiles()
		} else if decodedMessage.Name == "acquire-profile-lock" {
			return profile.AcquireProfileLock(decodedMessage.Detail[0])
		} else if decodedMessage.Name == "release-profile-lock" {
			profile.ReleaseProfileLock()
		} else if decodedMessage.Name == "set-password" {
			crypt.SetPassword(decodedMessage.Detail[0], decodedMessage.Detail[1])
		} else if decodedMessage.Name == "lock-profile" {
			crypt.LockProfile(decodedMessage.Detail[0])
		} else if decodedMessage.Name == "get-profile-migration-status" {
			migrationStatus := crypt.GetProfileMigrationStatus(decodedMessage.Detail[0])

			if migrationStatus == crypt.ProfileMigrationStatusSupported {
				return "supported"
			} else if migrationStatus == crypt.ProfileMigrationStatusMigratable {
				return "migratable"
			} else {
				return "unsupported"
			}
		} else if decodedMessage.Name == "migrate-profile" {
			return crypt.MigrateProfile(decodedMessage.Detail[0], decodedMessage.Detail[1])
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
