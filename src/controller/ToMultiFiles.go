package controller

import (
	"encoding/json"
	"fmt"
	"img"
	"io/ioutil"
	"os"
	fp "path/filepath"
)

type ImageModel struct {
	Size    uint   // 图片尺寸（以宽为标准）
	ImgType string // 图片格式（png/jpg/jpeg）
	Path    string // 存储路径
}

// 配置图片
type ConfigModel struct {
	Size uint   `json:"size"` // 图片尺寸
	Name string `json:"name"` // 图片名
}

// 配置图片集
type ConfigModelList struct {
	Img       string        `json:"img"`       // 原始图片名
	Dir       string        `json:"dir"`       // 存储目录
	ImageType string        `json:"imagetype"` // 图片格式（png/jpg/jpeg）
	List      []ConfigModel `json:"list"`
}

// 将图片转换成几种尺寸
func ToMultiFiles(filepath string, targets []ImageModel) error {
	originImg, _, err := img.ReadImage(filepath)
	if err != nil {
		return fmt.Errorf("读取图片失败：%v\n", err)
	}

	for _, target := range targets {
		resultImg := img.Resize(target.Size, originImg)
		if err := img.SaveImage(target.ImgType, resultImg, target.Path); err != nil {
			return err
		}
	}

	return nil
}

func ToMultiFilesByConfig(workdir string, configFile string) error {
	f, err := os.Open(configFile)
	if err != nil {
		return err
	}

	cfgData, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	var cfgModelList ConfigModelList
	if err := json.Unmarshal(cfgData, &cfgModelList); err != nil {
		return fmt.Errorf("解析json失败：%v\n", err)
	}

	dir := ""
	if cfgModelList.Dir != "" {
		dir, err = fp.Abs(cfgModelList.Dir)
		if err != nil {
			return fmt.Errorf("获取存储目录全路径失败：%v\n", err)
		}
		fmt.Println("存储目录为：", dir)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("创建存储目录[%s]失败：%v\n", dir, err)
		}
	} else {
		dir = workdir
	}

	fmt.Println("存储目录为：", dir)

	imgModelList := make([]ImageModel, 0, 0)
	for _, cfgModel := range cfgModelList.List {
		imgModel := ImageModel{
			ImgType: cfgModelList.ImageType,
			Size:    cfgModel.Size,
			Path:    fp.Join(dir, cfgModel.Name),
		}

		filedir := fp.Dir(imgModel.Path)
		if filedir != "" {
			if err := os.MkdirAll(filedir, os.ModePerm); err != nil {
				return fmt.Errorf("创建存储目录(%s)失败：%v\n", filedir, err)
			}
		}
		imgModelList = append(imgModelList, imgModel)
	}

	if len(imgModelList) == 0 {
		return fmt.Errorf("图片配置文件错误")
	}

	imgpath := fp.Join(workdir, cfgModelList.Img)

	return ToMultiFiles(imgpath, imgModelList)
}
