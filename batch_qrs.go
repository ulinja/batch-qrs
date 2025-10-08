package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	"io"
	"os"
)

func main() {
	var args = os.Args[1:]
	if len(args) < 1 {
		errorExit("at least one argument is required")
	}

	for _, arg := range args {
		bytes := generateQrCode(arg)
		fmt.Println(base64Encode(bytes))
	}
}

func errorExit(message string) {
	fmt.Fprintf(os.Stderr, "ERROR: %s\n", message)
	os.Exit(1)
}

func generateQrCode(payload string) []byte {
	qrc, err := qrcode.NewWith(
		payload,
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium),
	)
	if err != nil {
		errorExit(err.Error())
	}

	buf := bytes.NewBuffer(nil)
	wr := nopCloser{Writer: buf}
	w2 := standard.NewWithWriter(
		wr,
		standard.WithQRWidth(10),
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
		standard.WithBorderWidth(20),
	)
	if err = qrc.Save(w2); err != nil {
		errorExit(err.Error())
	}

	return buf.Bytes()
}

func base64Encode(payload []byte) string {
	return base64.StdEncoding.EncodeToString(payload)
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
