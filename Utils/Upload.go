package Utils

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
)

// 上传单个文件，保存到path指定的位置上
func UploadFile(file *multipart.FileHeader, path string) error {
	// Upload the file to specific dst.
	err := saveToDisk(file, path+file.Filename)
	if err != nil {
		e := fmt.Errorf("文件[%s]上传失败，错误信息：%s", file.Filename, err.Error())
		return e
	}
	return nil
}

//// 上传多个文件，保存到path指定的位置上
//// TODO: 根据前端传递的List创建文件层级结构，然后再保存文件
//func UploadFiles(c *gin.Context, path string) {
//	form, _ := c.MultipartForm()
//	files := form.File["upload[]"]
//	for _, file := range files {
//		log.Println(file.Filename)
//		err := c.SaveUploadedFile(file, file.Filename)
//		if err != nil {
//			e := fmt.Errorf("文件[%s]上传失败，错误信息：%s", file.Filename, err.Error())
//			Models.ResultFail(c, Code.UPLOAD_ERROR, e)
//			return
//		}
//	}
//	Models.ResultSuccessWithMsg(c, nil, fmt.Sprintf("全部文件上传成功，共上传了%d个文件", len(files)))
//}

// 把上传的文件保存到本地
func saveToDisk(file *multipart.FileHeader, dst string) error {
	srcFile, err := file.Open()
	if err != nil {
		return err
	}
	defer srcFile.Close()

	outFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, srcFile)
	return err
}
