package helpers

import (
	"bytes"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/makiuchi-d/gozxing/qrcode/decoder"
	"image"
	"image/png"
	"io"
	"io/ioutil"
)

func EncodeQrcode(content string, size int) []byte {
	encodeHints := map[gozxing.EncodeHintType]interface{}{
		gozxing.EncodeHintType_ERROR_CORRECTION: decoder.ErrorCorrectionLevel_L,
	}
	writer := qrcode.NewQRCodeWriter()
	encode, _ := writer.Encode(content, gozxing.BarcodeFormat_QR_CODE, size, size, encodeHints)

	var b bytes.Buffer
	ioWriter := io.Writer(&b)
	_ = png.Encode(ioWriter, encode)
	return b.Bytes()
}

func DecodeQrcode(qr []byte) gozxing.Result {
	reader := qrcode.NewQRCodeReader()
	newReader := bytes.NewReader(qr)
	img, _, _ := image.Decode(newReader)
	fromImage, _ := gozxing.NewBinaryBitmapFromImage(img)

	decodeHints := map[gozxing.DecodeHintType]interface{}{
		gozxing.DecodeHintType_TRY_HARDER: true,
	}
	decoded, _ := reader.Decode(fromImage, decodeHints)
	return *decoded
}

func DecodeQrTextFromFile(path string) gozxing.Result {
	return DecodeQrcode(CheckErrorWithResult(ioutil.ReadFile(path)).([]byte))
}

func EncodeQrTextToFile(path string, content string, size int) {
	encodeQrcode := EncodeQrcode(content, size)
	CheckError(ioutil.WriteFile(path, encodeQrcode, 0x644))
}
