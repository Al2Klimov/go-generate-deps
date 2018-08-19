package go_generate_deps

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

var errBadOut = errors.New("got bad output from 'go list -json .'")

func GenDeps() error {
	cmd := exec.Command("go", "list", "-json", ".")
	outBuf := bytes.Buffer{}

	cmd.Stdout = &outBuf
	cmd.Stderr = os.Stderr

	if errCR := cmd.Run(); errCR != nil {
		return errCR
	}

	var goList interface{}

	if errJUM := json.Unmarshal(outBuf.Bytes(), &goList); errJUM != nil {
		return errJUM
	}

	var pkgs map[string]struct{}
	var pkgName string

	if goListMap, goListIsMap := goList.(map[string]interface{}); goListIsMap {
		if name, hasName := goListMap["Name"]; hasName {
			if nameString, nameIsString := name.(string); nameIsString {
				pkgName = nameString
			} else {
				return errBadOut
			}
		} else {
			return errBadOut
		}

		if deps, hasDeps := goListMap["Deps"]; hasDeps {
			if depsArray, depsIsArray := deps.([]interface{}); depsIsArray {
				pkgs = make(map[string]struct{}, len(depsArray))

				for _, pkg := range depsArray {
					if pkgString, pkgIsString := pkg.(string); pkgIsString {
						pkgs[pkgString] = struct{}{}
					} else {
						return errBadOut
					}
				}
			} else {
				return errBadOut
			}
		} else {
			pkgs = map[string]struct{}{}
		}
	} else {
		return errBadOut
	}

	return ioutil.WriteFile(
		"GithubcomAl2klimovGogeneratedeps.go",
		[]byte(fmt.Sprintf("package %s\nvar GithubcomAl2klimovGogeneratedeps = %#v", pkgName, pkgs)),
		0666,
	)
}
