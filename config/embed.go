package config

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

//go:embed config-example.yml
var embeddedFiles embed.FS

func GenerateEmbeddedFile(embeddedFile, outputPath string) {
	GenerateFile(embeddedFile, outputPath, embeddedFile)
}

func GenerateFile(embeddedFile, outputPath, outputFile string) {
	// 提取嵌入的文件
	fileContents, err := embeddedFiles.ReadFile(embeddedFile)
	if err != nil {
		log.Fatalf("Error reading embedded file: %v\n", err)
	}

	// 创建输出目录（如果不存在）
	if err := os.MkdirAll(outputPath, os.ModePerm); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
	}
	// 将文件内容写入指定文件夹

	outputPath = filepath.Join(outputPath, outputFile)
	// 将文件内容写入指定文件夹
	err = os.WriteFile(outputPath, fileContents, fs.FileMode(0644))
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
	}
	fmt.Printf("File written successfully to %s\n", outputPath)
}
