package util

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"log"
	"os/exec"
)

func GetBase64Frame(filename string) string {
	width := 275
	height := 220

	cmd := exec.Command("./libs/bin/ffmpeg.exe", "-i", filename, "-vframes", "1", "-s", fmt.Sprintf("%dx%d", width, height), "-f", "singlejpeg", "-")

	buf := new(bytes.Buffer)

	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		log.Println("could not generate frame, ", err)
	}

	input := buf.Bytes()
	encodeString := base64.StdEncoding.EncodeToString(input)
	return encodeString
}