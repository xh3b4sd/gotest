package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/giantswarm/microerror"
)

const path = "/Users/xh3b4sd/go/src/github.com/giantswarm/api"

const search = `var ([a-zA-Z0-9]+) = microerror.New\("[ a-zA-Z0-9]+"\)`

const replace = `var $1 = &microerror.Error{
	Kind: "$1",
}`

func main() {
	walkFunc := func(p string, i os.FileInfo, err error) error {
		if err != nil {
			return microerror.Mask(err)
		}

		if i.IsDir() && i.Name() == ".circleci" {
			return filepath.SkipDir
		}
		if i.IsDir() && i.Name() == ".git" {
			return filepath.SkipDir
		}
		if i.IsDir() && i.Name() == "vendor" {
			return filepath.SkipDir
		}
		if i.IsDir() {
			return nil
		}
		if !i.IsDir() && i.Name() != "error.go" {
			return nil
		}

		c, err := ioutil.ReadFile(p)
		if err != nil {
			return microerror.Mask(err)
		}

		r, err := regexp.Compile(search)
		if err != nil {
			return microerror.Mask(err)
		}

		n := r.ReplaceAllString(string(c), replace)

		err = ioutil.WriteFile(p, []byte(n), i.Mode())
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	err := filepath.Walk(path, walkFunc)
	if err != nil {
		panic(err)
	}
}
