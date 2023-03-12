package ocr

/*
export CGO_LDFLAGS="-L/opt/homebrew/Cellar/leptonica/1.82.0_1/lib -L/opt/homebrew/Cellar/tesseract/5.3.0_1/lib"
export CGO_CXXFLAGS="-I/opt/homebrew/Cellar/leptonica/1.82.0_1/include -I/opt/homebrew/Cellar/tesseract/5.3.0_1/include"
**/

import (
	"fmt"

	"github.com/otiai10/gosseract/v2"
)

func Detect(path string) error {
	client := gosseract.NewClient()
	defer client.Close()

	client.SetWhitelist("1234567890+-x")
	client.SetPageSegMode(gosseract.PSM_SINGLE_LINE)

	client.SetImage(path)
	text, _ := client.Text()
	fmt.Println(text)
	return nil
}

func Run(token string) error {
	err := Detect("/Users/tim/Downloads/img.png")
	if err != nil {
		return err
	}

	return nil
}
