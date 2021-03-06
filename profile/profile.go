package profile

import (
	"errors"
	"io/fs"
	"os"
	"path"
	"runtime"
	"sort"
	"strconv"
	"strings"

	"github.com/gltchitm/firecrypt/native"
	"github.com/gofrs/flock"
	"gopkg.in/ini.v1"
)

type Profile struct {
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	Path               string `json:"path"`
	Configured         bool   `json:"configured"`
	CurrentlyEncrypted bool   `json:"currentlyEncrypted"`
}

var profileLock *flock.Flock

func GetProfiles() []Profile {
	cfg, err := ini.Load(path.Join(firefoxPath(), "profiles.ini"))
	if err != nil {
		panic(err)
	}

	return profilesFromConfig(cfg)
}
func AcquireProfileLock(profilePath string) bool {
	profileLock = flock.New(path.Join(profilePath, ".parentlock"))
	locked, err := profileLock.TryLock()
	return locked || err != nil
}
func ReleaseProfileLock() {
	err := profileLock.Unlock()
	if err != nil {
		panic(err)
	}
}
func LaunchProfile(profileName string) {
	lockedPath := profileLock.Path()
	ReleaseProfileLock()
	native.RunFirefox(profileName)
	AcquireProfileLock(lockedPath)
}

func firefoxPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	var relativePath string
	if runtime.GOOS == "darwin" {
		relativePath = "Library/Application Support/Firefox"
	} else if runtime.GOOS == "linux" {
		relativePath = ".mozilla/firefox"
	}

	return path.Join(home, relativePath)
}
func profilesFromConfig(cfg *ini.File) []Profile {
	sections := cfg.Sections()
	profiles := []Profile{}

	for _, section := range sections {
		if section.HasKey("Name") {
			id, err := strconv.Atoi(
				strings.TrimPrefix(
					section.Name(),
					"Profile",
				),
			)
			if err != nil {
				panic(err)
			}

			var profilePath string

			isRelative, err := section.Key("IsRelative").Int()
			if err != nil {
				panic(err)
			}

			if isRelative == 1 {
				profilePath = path.Join(firefoxPath(), section.Key("Path").String())
			} else {
				profilePath = section.Key("Path").String()
			}

			configured := true
			currentlyEncrypted := true

			_, err = os.Stat(profilePath + ".firecrypt")
			if errors.Is(err, fs.ErrNotExist) {
				currentlyEncrypted = false

				_, err = os.Stat(path.Join(profilePath, ".__firecrypt_key__"))
				if errors.Is(err, fs.ErrNotExist) {
					configured = false
				} else if err != nil {
					panic(err)
				}
			} else if err != nil {
				panic(err)
			}

			profiles = append(profiles, Profile{
				Id:                 id,
				Name:               section.Key("Name").String(),
				Path:               profilePath,
				Configured:         configured || currentlyEncrypted,
				CurrentlyEncrypted: currentlyEncrypted,
			})
		}
	}

	sort.SliceStable(profiles, func(i, j int) bool {
		return profiles[i].Id < profiles[j].Id
	})

	return profiles
}
