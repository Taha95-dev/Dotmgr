package fzf

import (
	"bytes"
	"os"
	"os/exec"
)

func Select(options []string) (string, error) {
	cmd := exec.Command("fzf", "--prompt=👉 Select: ", "--height=40%", "--border")
	var input bytes.Buffer
	for _, opt := range options {
		input.WriteString(opt + "\n")
	}
	cmd.Stdin = &input

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := out.String()
	if len(result) > 0 {
		result = result[:len(result)-1]
	}
	return result, nil
}
