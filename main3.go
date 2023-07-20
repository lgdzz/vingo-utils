package main

import (
	"fmt"
	"path/filepath"
)

func main() {
	filePath := "uploads/2023-07-19/81208dbd-4ae4-4358-afd1-1c5c1df41851.doc"
	fileName := filepath.Base(filePath)
	fileNameWithoutExt := fileName[:len(fileName)-len(filepath.Ext(fileName))]
	newFilePath := filepath.Join(filepath.Dir(filePath), fileNameWithoutExt+".txt")

	fmt.Println(newFilePath)
}
