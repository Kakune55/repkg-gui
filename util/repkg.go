package util

import (
	"os/exec"
)



func RePkgVersion() string {
	cmd := exec.Command("./RePKG.exe", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "ERROR"
	}
	return string(out)
}

func RePkgExtract(source string, target string) (err error) {
	cmd := exec.Command("./RePKG.exe", "extract", source, "-o "+target)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if string(out) != "" {
		return nil
	}
	return nil
}


