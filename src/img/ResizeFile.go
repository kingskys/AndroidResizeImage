package img

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

func ReadImage(filepath string) (image.Image, string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, "", err
	}

	defer f.Close()

	return image.Decode(f)
}

func Resize(size uint, img image.Image) image.Image {
	return resize.Resize(size, 0, img, resize.Bilinear)
}

func SaveImage(imgeType string, m image.Image, filepath string) error {

	f, err := os.Create(filepath)
	if err != nil {
		return err
	}

	if imgeType == "png" {
		return png.Encode(f, m)
	} else if imgeType == "jpg" || imgeType == "jpeg" {
		return jpeg.Encode(f, m, nil)
	}

	return fmt.Errorf("不支持格式：%s\n", imgeType)
}

// 缩放图片的主要方法
func DecodeImage(filepath string, outpath string, size uint) error {

	img, t, err := ReadImage(filepath)
	if err != nil {
		return err
	}

	resultImg := Resize(size, img)

	return SaveImage(t, resultImg, outpath)

}
