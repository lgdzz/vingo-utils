// *****************************************************************************
// 作者: lgdz
// 创建时间: 2024/8/5
// 描述：
// *****************************************************************************

package soffice

import (
	"fmt"
	"github.com/lgdzz/vingo-utils/vingo"
	"github.com/lgdzz/vingo-utils/vingo/request"
	"os"
	"os/exec"
	"path/filepath"
)

func Word2Pdf(fileUrl string) (pdf []byte) {
	var tmpPath = filepath.Join(".", "tmp")
	vingo.Mkdir(tmpPath)

	wordPath := request.DownloadFile(fileUrl, tmpPath, false)
	pdfPath := vingo.ReplaceFilePathExt(wordPath, ".pdf")
	wordPathAbs, _ := filepath.Abs(wordPath)
	tmpAbsPath, _ := filepath.Abs(tmpPath)
	defer func() {
		vingo.FileDelete(wordPath, false)
		vingo.FileDelete(pdfPath, false)
	}()
	// libreoffice --headless --convert-to txt /data/test.doc --outdir /data
	// soffice --headless --convert-to txt /Users/lgdz/Downloads/2e446484eb9945b9bcdc2d204a7d35ac.docx --outdir /Users/lgdz/Downloads
	//cmd := exec.Command("libreoffice", "--headless", "--convert-to", "txt", wordPathAbs, "--outdir", tmpAbsPath)
	cmd := exec.Command("soffice", "--headless", "--convert-to", "pdf", wordPathAbs, "--outdir", tmpAbsPath)
	err := cmd.Run()
	if err != nil {
		vingo.LogError(fmt.Sprintf("执行libreoffice错误：%v", err))
		panic(err)
	}
	pdf, err = os.ReadFile(pdfPath)
	if err != nil {
		vingo.LogError(fmt.Sprintf("读取pdf文件错误：%v", err))
		panic(err)
	}
	return
}
