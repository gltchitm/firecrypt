package profile

import (
	"os"
	"path"
	"sort"
	"os/exec"
	"strings"
	"strconv"
	"gopkg.in/ini.v1"
)

type Profile struct {
	Id int
	Name string
	Path string
	Configured bool
	CurrentlyEncrypted bool
}

func firefoxPath() string {
	var home, err = os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	return path.Join(home, "Library/Application Support/Firefox")
}
func profilesFromConfig(cfg *ini.File) []Profile {
	var sections = cfg.Sections()
	var profiles = []Profile {}

	for _, section := range sections {
		if (section.HasKey("Name")) {
			var id, err = strconv.Atoi(
				strings.TrimPrefix(
					section.Name(),
					"Profile",
				),
			)

			if err != nil {
				panic(err)
			}

			var profilePath string
			var isRelative int

			if isRelative, err = section.Key("IsRelative").Int(); err != nil {
				panic(err)
			}

			if isRelative == 1 {
				profilePath = path.Join(firefoxPath(), section.Key("Path").String())
			} else {
				profilePath = section.Key("Path").String()
			}

			var configured = false
			var currentlyEncrypted = false

			if _, err = os.ReadFile(profilePath + ".firecrypt"); err == nil {
				currentlyEncrypted = true
			} else if _, err = os.ReadFile(
				path.Join(profilePath, ".__firecrypt_hash__"),
			); err == nil {
				configured = true
			}

            profiles = append(profiles, Profile {
				Id: id,
				Name: section.Key("Name").String(),
				Path: profilePath,
				Configured: configured || currentlyEncrypted,
				CurrentlyEncrypted: currentlyEncrypted,
			})
		}
	}

	sort.SliceStable(profiles, func (i, j int) bool {
		return profiles[i].Id < profiles[j].Id
	})

	return profiles
}

func GetProfiles() []Profile {
	var cfg, err = ini.Load(path.Join(firefoxPath(), "profiles.ini"))

	if err != nil {
		panic(err)
	}

	return profilesFromConfig(cfg)
}
func IsProfileOpen(profilePath string) bool {
	var out, err = exec.Command("lsof").Output()

	if err != nil {
		panic(err)
	}

	return strings.Contains(
		string(out),
		path.Join(profilePath, ".parentlock"),
	)
}
func LaunchProfile(profileName string) {
	var _, err = exec.Command(
		"/Applications/Firefox.app/Contents/MacOS/firefox",
		"-p",
		profileName,
	).Output()

	if err != nil {
		panic(err)
	}

}
func SetPassword(profilePath string, newPassword string) {
    var hashFile, err = os.Create(path.Join(profilePath, ".__firecrypt_hash__"))

    if err != nil {
        panic(err)
    }

    hashFile.WriteString(newPassword)
    hashFile.Close()
}