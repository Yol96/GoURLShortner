package apiserver

import (
	"fmt"
	"image/png"
	"os"

	"github.com/Yol96/GoURLShortner/internal/app/model"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func createQRCodeImage(sl *model.Link) error {
	baseURL := "http://localhost:5000/api/" // Move to const or config file
	directoryPath := "./static/"

	qrCode, err := qr.Encode(baseURL+sl.ShortLink, qr.M, qr.Auto)
	if err != nil {
		return err
	}

	qrCode, err = barcode.Scale(qrCode, 200, 200)
	if err != nil {
		return err
	}

	imageName := fmt.Sprintf("%s.png", sl.ShortLink)
	file, err := os.Create(directoryPath + imageName)
	if err != nil {
		return err
	}

	defer file.Close()
	png.Encode(file, qrCode)
	return nil
}
