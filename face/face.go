package face

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"sort"

	pigo "github.com/esimov/pigo/core"
)

type Coordinate [2]int

type Solution struct {
	FaceTiles []Coordinate `json:"face_tiles"`
}

func DetectFaces(image io.Reader) (ret []pigo.Detection, err error) {
	cascadeFile, err := ioutil.ReadFile("face/cascade/facefinder")
	if err != nil {
		return
	}

	src, err := pigo.DecodeImage(image)
	if err != nil {
		return
	}

	pixels := pigo.RgbToGrayscale(src)
	cols, rows := src.Bounds().Max.X, src.Bounds().Max.Y

	cParams := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     500,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,

		ImageParams: pigo.ImageParams{
			Pixels: pixels,
			Rows:   rows,
			Cols:   cols,
			Dim:    cols,
		},
	}

	pigo := pigo.NewPigo()

	classifier, err := pigo.Unpack(cascadeFile)
	if err != nil {
		return
	}

	angle := 0.0

	return classifier.RunCascade(cParams, angle), nil
}

func Run(token string) error {
	file, err := os.Open("/Users/tim/Downloads/face.png")
	if err != nil {
		return err
	}

	dets, err := DetectFaces(file)
	if err != nil {
		return err
	}

	file, err = os.Open("/Users/tim/Downloads/face.png")
	if err != nil {
		return err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	n := image.NewRGBA(img.Bounds())

	for i := 0; i < n.Bounds().Dx(); i++ {
		for j := 0; j < n.Bounds().Dy(); j++ {
			n.Set(i, j, img.At(i, j))
		}
	}

	resultSet := make(map[Coordinate]bool)
	result := []Coordinate{}

	for _, det := range dets {
		if det.Q > 5 {
			for i := 0; i < 20; i++ {
				for j := 0; j < 20; j++ {
					n.Set(det.Col+i, det.Row+j, color.RGBA{0, 255, 0, 1})
				}
			}

			x := int(math.Floor(float64(det.Row) / 100.0))
			y := int(math.Floor(float64(det.Col) / 100.0))

			resultSet[Coordinate{x, y}] = true
		}
	}

	for k := range resultSet {
		result = append(result, k)
	}

	sort.Slice(result, func(i, j int) bool {
		if result[i][0] != result[j][0] {
			return result[i][0] < result[j][0]
		} else {
			return result[i][1] < result[j][1]
		}

	})

	out, err := os.Create("/Users/tim/Downloads/annotated.png")
	if err != nil {
		return err
	}

	err = jpeg.Encode(out, n, nil)
	if err != nil {
		return err
	}

	solutionData, err := json.Marshal(Solution{FaceTiles: result})
	if err != nil {
		return err
	}

	response, err := http.Post(fmt.Sprintf("https://hackattic.com/challenges/basic_face_detection/solve?access_token=%s", token), "application/json", bytes.NewBuffer(solutionData))
	if err != nil {
		return err
	}

	fmt.Println(response.Status)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}
