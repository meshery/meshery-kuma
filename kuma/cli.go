package kuma

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
)

func GetKumactl(version string) error {
	if version == "latest" {
		version = "1.0.0"
	}
	distro := os.Getenv("GODISTRO")
	arch := os.Getenv("GOARCH")

	url := fmt.Sprintf("https://kong.bintray.com/kuma/kuma-%s-%s-%s.tar.gz", version, distro, arch)
	fmt.Println(url)
	response, err := http.Get(url)
	if err != nil {
		return ErrGetKumactl(err)
	}

	err = untar(response.Body)
	if err != nil {
		return ErrGetKumactl(err)
	}

	err = moveBinary(fmt.Sprintf("kuma-%s/bin/kumactl", version), "./kumactl")
	if err != nil {
		return ErrGetKumactl(err)
	}

	return nil
}

func untar(gzipStream io.Reader) error {

	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return ErrUntar(err)
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return ErrUntar(err)
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil && !os.IsExist(err) {
				return ErrUntar(err)
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return ErrUntar(err)
			}
			defer outFile.Close()
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return ErrUntar(err)
			}
		default:
			return ErrUntarDefault
		}
	}
	return nil
}

func moveBinary(sourcePath, destPath string) error {
	inputFile, err := os.Open(sourcePath)
	if err != nil {
		return ErrMoveBinary(err)
	}
	outputFile, err := os.Create(destPath)
	if err != nil {
		inputFile.Close()
		return ErrMoveBinary(err)
	}
	defer outputFile.Close()

	err = outputFile.Chmod(0755)
	if err != nil {
		return ErrMoveBinary(err)
	}
	_, err = io.Copy(outputFile, inputFile)
	inputFile.Close()
	if err != nil {
		return ErrMoveBinary(err)
	}
	// The copy was successful, so now delete the original file
	err = os.Remove(sourcePath)
	if err != nil {
		return ErrMoveBinary(err)
	}
	return nil
}
