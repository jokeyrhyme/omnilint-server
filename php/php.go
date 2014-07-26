package php

import (
  "bytes"
  "os/exec"
  "io"
)

func ParseSyntax(code io.Reader) (string, error) {
	path, err := exec.LookPath("php")
	if err != nil {
		return "", err
	}
	cmd := exec.Command(path, "-l")
	cmd.Stdin = code
	var (
    out bytes.Buffer
  )
	cmd.Stdout = &out
	cmd.Stderr = &out
	err = cmd.Run()
	if err != nil {
		return out.String(), nil
	}
	return "", nil
}
