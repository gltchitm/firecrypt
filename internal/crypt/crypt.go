package crypt

import (
	"archive/zip"
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

const magicVerionPrefix string = "@@5c53512d-FIRECRYPT-VERSION-2-6062fceb@@\n\n\n"

const (
	ProfileMigrationStatusSupported = iota
	ProfileMigrationStatusMigratable
	ProfileMigrationStatusUnsupported
)

const (
	argon2SaltLen = 32
	argon2Time    = 3
	argon2Memory  = 96 * 1024
	argon2Threads = 2
	argon2KeyLen  = 32
)

func GetProfileMigrationStatus(profilePath string) int {
	fileContents, err := ioutil.ReadFile(path.Join(
		filepath.Dir(profilePath),
		filepath.Base(profilePath)+".firecrypt",
	))
	if err != nil {
		panic(err)
	}

	if strings.HasPrefix(string(fileContents), magicVerionPrefix) {
		os.Remove(filepath.Join(profilePath, ".__firecrypt_hash__"))
		return ProfileMigrationStatusSupported
	} else {
		return ProfileMigrationStatusMigratable
	}
}
func MigrateProfile(profilePath, password string) bool {
	hash := sha512.Sum512([]byte(password))
	for i := 0; i < 249999; i++ {
		hash = sha512.Sum512(hash[:])
	}
	key := base64.StdEncoding.EncodeToString(hash[:])

	cmd := exec.Command(
		"openssl",
		"aes-256-cbc",
		"-d",
		"-pbkdf2",
		"-iter",
		"250000",
		"-in",
		path.Join(
			filepath.Dir(profilePath),
			filepath.Base(profilePath)+".firecrypt",
		),
		"-out",
		path.Join(
			filepath.Dir(profilePath),
			filepath.Base(profilePath)+".zip",
		),
		"-k",
		key,
	)

	cmd.Dir = filepath.Dir(profilePath)

	err := cmd.Run()

	if err != nil {
		cmd = exec.Command(
			"rm",
			path.Join(
				filepath.Dir(profilePath),
				filepath.Base(profilePath)+".zip",
			),
		)
		cmd.Run()
		if err != nil {
			panic(err)
		}

		return false
	}

	cmd = exec.Command(
		"unzip",
		path.Join(
			filepath.Dir(profilePath),
			filepath.Base(profilePath)+".zip",
		),
	)
	cmd.Dir = filepath.Dir(profilePath)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command(
		"rm",
		path.Join(
			filepath.Dir(profilePath),
			filepath.Base(profilePath)+".zip",
		),
	)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	cmd = exec.Command(
		"rm",
		path.Join(
			filepath.Dir(profilePath),
			filepath.Base(profilePath)+".firecrypt",
		),
	)
	err = cmd.Run()
	if err != nil {
		panic(err)
	}

	os.Remove(filepath.Join(profilePath, ".__firecrypt_hash__"))
	SetPassword(profilePath, password)

	return true
}
func LockProfile(profilePath string) bool {
	originalWd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	os.Chdir(filepath.Dir(profilePath))

	zipOutput := new(bytes.Buffer)

	zipWriter := zip.NewWriter(zipOutput)

	zipWalker := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			panic(err)
		}
		if info.IsDir() {
			return nil
		}

		file, err := zipWriter.Create(path)
		if err != nil {
			panic(err)
		}

		copySrc, err := os.Open(path)
		defer copySrc.Close()
		if err != nil {
			return err
		}

		_, err = io.Copy(file, copySrc)
		if err != nil {
			panic(err)
		}

		return nil
	}

	err = filepath.Walk(filepath.Base(profilePath), zipWalker)
	if err != nil {
		panic(err)
	}

	zipWriter.Close()

	hashFile, err := os.Open(path.Join(profilePath, ".__firecrypt_key__"))
	defer hashFile.Close()
	if err != nil {
		panic(err)
	}

	for i, v := range readBytesFromFile(*hashFile, len(magicVerionPrefix)) {
		if v != magicVerionPrefix[i] {
			panic("magic version prefix does not match in key file")
		}
	}

	salt := readBytesFromFile(*hashFile, argon2SaltLen)
	key := readBytesFromFile(*hashFile, argon2KeyLen)

	cipher, err := chacha20poly1305.NewX(key)
	if err != nil {
		panic(err)
	}

	nonce := randomBytes(cipher.NonceSize())
	encryptedZipData := cipher.Seal(nonce, nonce, zipOutput.Bytes(), nil)

	output, err := os.Create(filepath.Base(profilePath) + ".firecrypt")

	defer output.Close()

	writeBytesToFile(*output, []byte(magicVerionPrefix))
	writeBytesToFile(*output, salt)
	writeBytesToFile(*output, encryptedZipData)

	err = os.RemoveAll(filepath.Base(profilePath))
	if err != nil {
		panic(err)
	}

	os.Chdir(originalWd)

	return true
}
func UnlockProfile(profilePath, password string) bool {
	originalWd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	os.Chdir(filepath.Dir(profilePath))

	encrypted, err := os.Open(filepath.Base(profilePath) + ".firecrypt")
	if err != nil {
		panic(err)
	}

	readBytesFromFile(*encrypted, len(magicVerionPrefix))
	salt := readBytesFromFile(*encrypted, argon2SaltLen)
	nonce := readBytesFromFile(*encrypted, chacha20poly1305.NonceSizeX)
	encryptedZipData := readBytesFromFileUntilEOF(*encrypted)

	key := argon2.IDKey(
		[]byte(password),
		salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	cipher, err := chacha20poly1305.NewX(key)
	if err != nil {
		panic(err)
	}

	plaintext, err := cipher.Open(nil, nonce, encryptedZipData, nil)
	if err != nil {
		return false
	}

	compressedProfile := bytes.NewReader(plaintext)
	zipReader, err := zip.NewReader(compressedProfile, compressedProfile.Size())
	if err != nil {
		panic(err)
	}

	currentWd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	for _, file := range zipReader.File {
		path := filepath.Join(currentWd, file.Name)
		if !strings.HasPrefix(path, filepath.Clean(currentWd)+string(os.PathSeparator)) {
			panic("file path is invalid!")
		}

		if file.FileInfo().IsDir() {
			os.MkdirAll(path, os.ModePerm)
			continue
		}

		err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
		if err != nil {
			panic(err)
		}

		outputFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		defer outputFile.Close()
		if err != nil {
			panic(err)
		}

		copySrc, err := file.Open()
		defer copySrc.Close()
		if err != nil {
			panic(err)
		}

		_, err = io.Copy(outputFile, copySrc)
		if err != nil {
			panic(err)
		}
	}

	err = os.Remove(filepath.Base(profilePath) + ".firecrypt")
	if err != nil {
		panic(err)
	}

	os.Chdir(originalWd)

	return true
}
func SetPassword(profilePath string, password string) {
	hashFile, err := os.Create(path.Join(profilePath, ".__firecrypt_key__"))
	defer hashFile.Close()
	if err != nil {
		panic(err)
	}

	salt := randomBytes(argon2SaltLen)
	key := argon2.IDKey(
		[]byte(password),
		salt,
		argon2Time,
		argon2Memory,
		argon2Threads,
		argon2KeyLen,
	)

	writeBytesToFile(*hashFile, []byte(magicVerionPrefix))
	writeBytesToFile(*hashFile, salt)
	writeBytesToFile(*hashFile, key)
}

func randomBytes(length int) []byte {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		panic(err)
	}

	return bytes
}
func writeBytesToFile(file os.File, bytes []byte) {
	bytesWritten, err := file.Write(bytes)
	if bytesWritten != len(bytes) {
		panic("too few bytes written!")
	} else if err != nil {
		panic(err)
	}
}
func readBytesFromFile(file os.File, length int) []byte {
	output := make([]byte, length)
	bytesRead, err := file.Read(output)
	if bytesRead != length || err != nil {
		panic(err)
	}
	return output
}
func readBytesFromFileUntilEOF(file os.File) []byte {
	output := make([]byte, 0)

	for true {
		readByte := make([]byte, 1)
		_, err := file.Read(readByte)

		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		output = append(output, readByte[0])
	}

	return output
}
