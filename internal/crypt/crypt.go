package crypt

import (
	"path"
	"os/exec"
	"path/filepath"
)

func LockProfile(profilePath, password string) bool {
	var zipPath = path.Join(
		filepath.Dir(profilePath),
		filepath.Base(profilePath) + ".zip",
	)
	var cmd = exec.Command(
		"zip",
		"-mr",
		zipPath,
		filepath.Base(profilePath),
	)
	cmd.Dir = filepath.Dir(profilePath)

	var err = cmd.Run()

	if err != nil {
		panic(err)
	}

	cmd = exec.Command(
		"openssl",
		"aes-256-cbc",
		"-pbkdf2",
		"-iter",
		"250000",
		"-in",
		zipPath,
		"-out",
		path.Join(
			filepath.Dir(profilePath),
			filepath.Base(profilePath) + ".firecrypt",
		),
		"-k",
		password,
	)

	if err = cmd.Run(); err != nil {
		panic(err)
	}

	cmd = exec.Command(
		"rm",
		zipPath,
	)

	if err = cmd.Run(); err != nil {
		panic(err)
	}

	return false
}
func UnlockProfile(profilePath, password string) bool {
	var cmd = exec.Command(
		"openssl",
		"aes-256-cbc",
		"-d",
		"-pbkdf2",
		"-iter",
		"250000",
		"-in",
		path.Join(
			filepath.Dir(profilePath),
			filepath.Base(profilePath) + ".firecrypt",
		),
		"-out",
		path.Join(
			filepath.Dir(profilePath),
			filepath.Base(profilePath) + ".zip",
		),
		"-k",
		password,
	)

	cmd.Dir = filepath.Dir(profilePath)

	var err = cmd.Run()

	if err != nil {
		cmd = exec.Command(
			"rm",
			path.Join(
				filepath.Dir(profilePath),
				filepath.Base(profilePath) + ".zip",
			),
		)

		if err = cmd.Run(); err != nil {
			panic(err)
		}

		return false
	}

	cmd = exec.Command(
		"unzip",
		path.Join(
			filepath.Dir(profilePath),
			filepath.Base(profilePath) + ".zip",
		),
	)
	cmd.Dir = filepath.Dir(profilePath)

	if err = cmd.Run(); err != nil {
		panic(err)
	}

	cmd = exec.Command(
		"rm",
		path.Join(
			filepath.Dir(profilePath),
			filepath.Base(profilePath) + ".zip",
		),
	)

	if err = cmd.Run(); err != nil {
		panic(err)
	}

	cmd = exec.Command(
		"rm",
		path.Join(
			filepath.Dir(profilePath),
			filepath.Base(profilePath) + ".firecrypt",
		),
	)

	if err = cmd.Run(); err != nil {
		panic(err)
	}

	return true
}