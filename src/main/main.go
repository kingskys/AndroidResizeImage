package main

import (
	"controller"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	fp, err := filepath.Abs(os.Args[0])
	if err != nil {
		log.Fatal(err)
	}

	srcdir := filepath.Dir(fp)
	fmt.Println("工作目录：", srcdir)

	cfgpath := filepath.Join(srcdir, "config.json")

	if err := controller.ToMultiFilesByConfig(srcdir, cfgpath); err != nil {
		log.Fatal(err)
	}

	log.Println("完成.")

}
