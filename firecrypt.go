package main

import (
	"github.com/gltchitm/firecrypt/crypt"
	"github.com/gltchitm/firecrypt/native"
	"github.com/gltchitm/firecrypt/profile"
)

func main() {
	native.StartFirecrypt(func(name string, detail []string) interface{} {
		if name == "get-profiles" {
			return profile.GetProfiles()
		} else if name == "acquire-profile-lock" {
			return profile.AcquireProfileLock(detail[0])
		} else if name == "release-profile-lock" {
			profile.ReleaseProfileLock()
		} else if name == "set-password" {
			crypt.SetPassword(detail[0], detail[1])
		} else if name == "lock-profile" {
			crypt.LockProfile(detail[0])
		} else if name == "get-profile-migration-status" {
			migrationStatus := crypt.GetProfileMigrationStatus(detail[0])

			if migrationStatus == crypt.ProfileMigrationStatusSupported {
				return "supported"
			} else if migrationStatus == crypt.ProfileMigrationStatusMigratable {
				return "migratable"
			} else {
				return "unsupported"
			}
		} else if name == "migrate-profile" {
			return crypt.MigrateProfile(detail[0], detail[1])
		} else if name == "unlock-profile" {
			return crypt.UnlockProfile(detail[0], detail[1])
		} else if name == "launch-profile" {
			go profile.LaunchProfile(detail[0])
		}

		return nil
	})
}
