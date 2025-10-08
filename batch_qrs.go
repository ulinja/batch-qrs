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
		fmt.Fprintf(os.Stderr, "ERROR: at least one argument is required\n")
		os.Exit(1)
	}

	for _, arg := range args {
		qrc, err := qrcode.NewWith(
			arg,
			qrcode.WithEncodingMode(qrcode.EncModeByte),
			qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionMedium),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}

		fmt.Printf("%s\n", base64.StdEncoding.EncodeToString(buf.Bytes()))
	}
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
