package pkg

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strconv"

	"github.com/pakket-project/pakket/config"
	"github.com/pakket-project/pakket/errors"
	"github.com/pakket-project/pakket/repo"
)

// for use with GetPackage(). Contains all data needed to install a package.
type PkgData struct {
	PkgDef     PackageDefinition
	VerData    VersionMetadata
	PlfData    PlatformData
	TarURL     string
	RepoURL    string
	Version    string
	Repository string
	BinSize    int64
}

func NewPkgData(PkgDef PackageDefinition, VerData VersionMetadata, PlfData PlatformData, Repository string, Version string, TarURL string, RepoURL string, BinSize int64) *PkgData {
	return &PkgData{
		PkgDef:     PkgDef,
		VerData:    VerData,
		PlfData:    PlfData,
		Repository: Repository,
		Version:    Version,
		TarURL:     TarURL,
		RepoURL:    RepoURL,
		BinSize:    BinSize,
	}
}

// One function to get all information needed to install a package. Version should be "latest" for latest version. binSize is the size of the tarball in bytes.
// Probaly going to clean this up and seperate it in seperate functions later.
func GetPackage(pkgName string, pkgVersion *string) (pkgData *PkgData, err error) {
	// search core repository
	resp, err := http.Get(fmt.Sprintf("%s/%s/package.toml", repo.CoreRepositoryURL, pkgName))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// check if response is 200
	if resp.StatusCode != 200 {
		return nil, errors.PackageNotFoundError{Package: pkgName}
	}

	// found, get package definition
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	pkgDef, err := ParsePackage(body)
	if err != nil {
		return nil, err
	}

	var version string

	// get version metadata
	if pkgVersion == nil {
		// latest version
		version = pkgDef.Package.Version
	} else {
		version = *pkgVersion
	}

	resp, err = http.Get(fmt.Sprintf("%s/%s/%s/metadata.toml", repo.CoreRepositoryURL, pkgName, version))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	verData, err := ParseVersion(body)

	var plfData PlatformData

	// get platform data
	if runtime.GOARCH == "arm64" {
		plfData = verData.Arm64
	} else if runtime.GOARCH == "amd64" {
		plfData = verData.Amd64
	}

	pkgUrl := fmt.Sprintf("%s/%s/%s/%s-%s-%s.tar.xz", config.C.Mirrors[0].URL, pkgName, version, pkgName, version, runtime.GOARCH)
	pkgRepoUrl := fmt.Sprintf("%s/%s/%s", repo.CoreRepositoryURL, pkgName, version)
	fmt.Println(pkgUrl, pkgRepoUrl)
	// get pkg size
	size, err := GetPackageSize(pkgUrl)

	return NewPkgData(pkgDef, verData, plfData, "core", version, pkgUrl, pkgRepoUrl, size), err
}

func GetPackageSize(url string) (bytes int64, err error) {
	resp, err := http.Head(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("bad status: %s", resp.Status)
	}

	bytes, err = strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return 0, err
	}

	return bytes, nil
}
