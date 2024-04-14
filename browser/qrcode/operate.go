package qrcode

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func CreateBarcodeAsBase64(msg string) string {
	// Create the barcode
	qrCode, _ := qr.Encode(msg, qr.M, qr.Auto)
	qrCode, _ = barcode.Scale(qrCode, 200, 200)

	// encode the barcode as png and save it to buffer memory
	buf := new(bytes.Buffer)
	err := png.Encode(buf, qrCode)
	if err != nil {
		fmt.Println("Error: ", err)
		panic(err)
	}

	// convert but fo base64
	base64img := base64.StdEncoding.EncodeToString(buf.Bytes())
	pngImgString := "data:image/png;base64," + base64img

	return pngImgString
}
