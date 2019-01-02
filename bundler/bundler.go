package bundler

import (
	"fmt"
	"github.com/hhrutter/pdfcpu/pkg/api"
	"github.com/hhrutter/pdfcpu/pkg/pdfcpu"
	"github.com/mitchellh/go-homedir"
	"os"
	"strings"
)

// TODO: set in env vars
var home, _ = homedir.Dir()
var pdfOutputPath = home + "/Documents/bundles"
var imageBasePath = home + "/Documents/scanned"

func Bundle(foldername string) (string, error) {
	conf := pdfcpu.NewDefaultConfiguration()
	imp := pdfcpu.DefaultImportConfig()
	out, err := getBundleName(foldername)
	if err != nil {
		return "", err
	}
	names, err := getFileNames(foldername)
	if err != nil {
		return "", err
	}
	clearExistingBundle(foldername)
	comm := api.ImportImagesCommand(names, out, imp, conf)
	_, err = api.ImportImages(comm)
	if err != nil {
		return "", fmt.Errorf("could not bundle images into PDF: %s", err)
	}
	return out, nil
}

func getFileNames(path string) ([]string, error) {
	i, err := os.Stat(imageBasePath + "/" + path)
	if err != nil {
		return nil, fmt.Errorf("requested path does not exist: %s", err)
	}
	isDir := i.IsDir()
	if !isDir {
		return nil, fmt.Errorf("requested path is not a directory: %s", path)
	}
	d, err := os.Open(imageBasePath + "/" + path)
	if err != nil {
		return nil, fmt.Errorf("could not open requested directory: %s", err)
	}
	defer d.Close()
	raw, err := d.Readdirnames(0)
	var names []string
	if err != nil {
		return nil, fmt.Errorf("could not read contents of requested directory: %s", err)
	}
	for _, n := range raw {
		if strings.HasSuffix(n, ".jpg") {
			names = append(names, imageBasePath+ "/" + path + "/" + n)
		}
	}
	return names, nil
}

func getBundleName(path string) (string, error) {
	err := os.MkdirAll(pdfOutputPath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("could not create output path: %s", err)
	}
	s := strings.Split(path, "/")
	return pdfOutputPath + "/" + s[len(s) -1] + ".pdf", nil
}

func clearExistingBundle(bundleName string) {
	_, err := os.Stat(pdfOutputPath + "/" + bundleName)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("no bundle existed at " + pdfOutputPath + "/" + bundleName)
			return
		}
		fmt.Println("could not determine existence of a previous PDF by the name " + bundleName)
		return
	}
	if err = os.Remove(pdfOutputPath + "/" + bundleName); err != nil {
		fmt.Println("could not remove the previous PDF by the name " + bundleName)
	}
	return
}