package kuma

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	repoURL     = "https://api.github.com/repos/containous/kuma/releases/latest"
	URLSuffix   = "-linux.tar.gz"
	crdPattern  = "crd(.*)yaml"
	cachePeriod = 1 * time.Hour
)

var (
	localFile                  = path.Join(os.TempDir(), "kuma.tar.gz")
	destinationFolder          = path.Join(os.TempDir(), "kuma")
	basePath                   = path.Join(destinationFolder, "%s")
	installFile                = path.Join(basePath, "install/kubernetes/kuma.yaml")
	installWithmTLSFile        = path.Join(basePath, "install/kubernetes/kuma-demo-auth.yaml")
	bookInfoInstallFile        = path.Join(basePath, "samples/bookinfo/platform/kube/bookinfo.yaml")
	bookInfoGatewayInstallFile = path.Join(basePath, "samples/bookinfo/networking/bookinfo-gateway.yaml")
	crdFolder                  = path.Join(basePath, "install/kubernetes/helm/kuma-init/files/")
)
// APIInfo is used to store individual response from GitHub release call
type APIInfo struct {
	TagName    string   `json:"tag_name,omitempty"`
	PreRelease bool     `json:"prerelease,omitempty"`
	Assets     []*Asset `json:"assets,omitempty"`
}
// APIInfo is used to store individual response from GitHub release call
type Asset struct {
	Name        string `json:"name,omitempty"`
	State       string `json:"state,omitempty"`
	DownloadURL string `json:"browser_download_url,omitempty"`
}

func (iClient *Client) getLatestReleaseURL() error {
	if iClient.kumaReleaseDownloadURL == "" || time.Since(iClient.kumaReleaseUpdatedAt) > cachePeriod {
		logrus.Debugf("API info url: %s", repoURL)
		resp, err := http.Get(repoURL)
		if err != nil {
			err = errors.Wrapf(err, "error getting latest version info")
			logrus.Error(err)
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			err = fmt.Errorf("unable to fetch release info due to an unexpected status code: %d", resp.StatusCode)
			logrus.Error(err)
			return err
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			err = errors.Wrapf(err, "error parsing response body")
			logrus.Error(err)
			return err
		}
		// logrus.Debugf("Raw api info: %s", body)
		result := &APIInfo{}
		err = json.Unmarshal(body, result)
		if err != nil {
			err = errors.Wrapf(err, "error unmarshalling response body")
			logrus.Error(err)
			return err
		}
		logrus.Debugf("retrieved api info: %+#v", result)
		if result != nil && result.Assets != nil && len(result.Assets) > 0 {
			for _, asset := range result.Assets {
				if strings.HasSuffix(asset.Name, URLSuffix) {
					iClient.kumaReleaseVersion = strings.Replace(asset.Name, URLSuffix, "", -1)
					iClient.kumaReleaseDownloadURL = asset.DownloadURL
					iClient.kumaReleaseUpdatedAt = time.Now()
					return nil
				}
			}
		}
		err = errors.New("unable to extract the download URL")
		logrus.Error(err)
		return err
	}
	return nil
}

func (iClient *Client) downloadFile(localFile string) error {
	dFile, err := os.Create(localFile)
	if err != nil {
		err = errors.Wrapf(err, "unable to create a file on the filesystem at %s", localFile)
		logrus.Error(err)
		return err
	}
	defer dFile.Close()

	resp, err := http.Get(iClient.kumaReleaseDownloadURL)
	if err != nil {
		err = errors.Wrapf(err, "unable to download the file from URL: %s", iClient.kumaReleaseDownloadURL)
		logrus.Error(err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("unable to download the file from URL: %s, status: %s", iClient.kumaReleaseDownloadURL, resp.Status)
		logrus.Error(err)
		return err
	}

	_, err = io.Copy(dFile, resp.Body)
	if err != nil {
		err = errors.Wrapf(err, "unable to write the downloaded file to the file system at %s", localFile)
		logrus.Error(err)
		return err
	}
	return nil
}

func (iClient *Client) untarPackage(destination, fileToUntar string) error {
	lFile, err := os.Open(fileToUntar)
	if err != nil {
		err = errors.Wrapf(err, "unable to read the local file %s", fileToUntar)
		logrus.Error(err)
		return err
	}

	gzReader, err := gzip.NewReader(lFile)
	if err != nil {
		err = errors.Wrap(err, "unable to load the file into a gz reader")
		logrus.Error(err)
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		switch {
		case err == io.EOF:
			return nil
		case err != nil:
			err = errors.Wrap(err, "error during untar")
			logrus.Error(err)
			return err
		case header == nil:
			continue
		}

		fileInLoop := filepath.Join(destination, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			if _, err := os.Stat(fileInLoop); err != nil {
				if err := os.MkdirAll(fileInLoop, 0755); err != nil {
					err = errors.Wrapf(err, "error creating directory %s", fileInLoop)
					logrus.Error(err)
					return err
				}
			}
		case tar.TypeReg:
			fileAtLoc, err := os.OpenFile(fileInLoop, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				err = errors.Wrapf(err, "error opening file %s", fileInLoop)
				logrus.Error(err)
				return err
			}

			if _, err := io.Copy(fileAtLoc, tarReader); err != nil {
				err = errors.Wrapf(err, "error writing file %s", fileInLoop)
				logrus.Error(err)
				return err
			}
			fileAtLoc.Close()
		}
	}
}

func (iClient *Client) downloadKuma() (string, error) {
	logrus.Debug("preparing to download the latest kuma release")
	err := iClient.getLatestReleaseURL()
	if err != nil {
		return "", err
	}
	fileName := iClient.kumaReleaseVersion
	downloadURL := iClient.kumaReleaseDownloadURL
	logrus.Debugf("retrieved latest file name: %s and download url: %s", fileName, downloadURL)

	proceedWithDownload := true

	lFileStat, err := os.Stat(localFile)
	if err == nil {
		if time.Since(lFileStat.ModTime()) > cachePeriod {
			proceedWithDownload = true
		} else {
			proceedWithDownload = false
		}
	}

	if proceedWithDownload {
		if err = iClient.downloadFile(localFile); err != nil {
			return "", err
		}
		logrus.Debug("package successfully downloaded, now unzipping . . .")
	}
	if err = iClient.untarPackage(destinationFolder, localFile); err != nil {
		return "", err
	}
	logrus.Debug("successfully unzipped")
	return fileName, nil
}

func (iClient *Client) getKumaComponentYAML(fileName string) (string, error) {
	specificVersionName, err := iClient.downloadKuma()
	if err != nil {
		return "", err
	}
	installFileLoc := fmt.Sprintf(fileName, specificVersionName)
	logrus.Debugf("checking if install file exists at path: %s", installFileLoc)
	_, err = os.Stat(installFileLoc)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Error(err)
			return "", err
		}
		err = errors.Wrap(err, "unknown error")
		logrus.Error(err)
		return "", err
	}
	fileContents, err := ioutil.ReadFile(installFileLoc)
	if err != nil {
		err = errors.Wrap(err, "unable to read file")
		logrus.Error(err)
		return "", err
	}
	return string(fileContents), nil
}

func (iClient *Client) getCRDsYAML() ([]string, error) {
	res := []string{}

	rEx, err := regexp.Compile(crdPattern)
	if err != nil {
		err = errors.Wrap(err, "unable to compile crd pattern")
		logrus.Error(err)
		return nil, err
	}

	specificVersionName, err := iClient.downloadKuma()
	if err != nil {
		return nil, err
	}
	startFolder := fmt.Sprintf(crdFolder, specificVersionName)
	err = filepath.Walk(startFolder, func(currentPath string, info os.FileInfo, err error) error {
		if err == nil && rEx.MatchString(info.Name()) {
			contents, err := ioutil.ReadFile(currentPath)
			if err != nil {
				err = errors.Wrap(err, "unable to read file")
				logrus.Error(err)
				return err
			}
			res = append(res, string(contents))
		}
		return nil
	})
	if err != nil {
		err = errors.Wrap(err, "unable to read the directory")
		logrus.Error(err)
		return nil, err
	}
	return res, nil
}

func (iClient *Client) getLatestKumaYAML(installmTLS bool) (string, error) {
	if installmTLS {
		return iClient.getIstioComponentYAML(installWithmTLSFile)
	}
	return iClient.getIstioComponentYAML(installFile)
}

func (iClient *Client) getBookInfoAppYAML() (string, error) {
	return iClient.getKumaComponentYAML(bookInfoInstallFile)
}

func (iClient *Client) getBookInfoGatewayYAML() (string, error) {
	return iClient.getKumaComponentYAML(bookInfoGatewayInstallFile)
}

func getSMIYamls() (string, error) {
	return "", nil
}
