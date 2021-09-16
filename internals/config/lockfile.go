package config

import (
	"errors"
	"os"

	"github.com/pakket-project/pakket/util"
	"github.com/pelletier/go-toml/v2"
)

type LockfileMetadata struct {
	Name       string `toml:"name"`
	Version    string `toml:"version"`
	Checksum   string `toml:"checksum"`
	Repository string `toml:"repository"`
}

type LockfileStruct struct {
	Packages []LockfileMetadata `toml:"packages" mapstructure:"packages"`
}

func readLockfile() (err error) {
	file, err := os.ReadFile(util.LockfilePath)
	if errors.Is(err, os.ErrNotExist) {
		_, err = os.Create(util.LockfilePath)
		if err != nil {
			return err
		}
		file, err = os.ReadFile(util.LockfilePath)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = toml.Unmarshal(file, &LockFile)

	return err
}

// Add package information to lockfile
func AddPkgToLockfile(metadata LockfileMetadata) (err error) {
	err = readLockfile()
	if err != nil {
		return err
	}
	LockFile.Packages = append(LockFile.Packages, metadata)

	newLockfile, err := toml.Marshal(&LockFile)
	if err != nil {
		return err
	}

	err = os.WriteFile(util.LockfilePath, newLockfile, 0666)
	return err
}
