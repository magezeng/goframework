package Utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"tipu.com/go-framework/Code"
	"tipu.com/go-framework/Models"
)

// 上传单个文件，保存到path指定的位置上
func UploadFile(c *gin.Context, path string) {
	file, _ := c.FormFile("file")
	// Upload the file to specific dst.
	err := c.SaveUploadedFile(file, path+file.Filename)
	if err != nil {
		e := fmt.Errorf("文件[%s]上传失败，错误信息：%s", file.Filename, err.Error())
		Models.ResultFail(c, Code.UPLOAD_ERROR, e)
		return
	}
	Models.ResultSuccessWithMsg(c, nil, "上传成功！")
}

// 上传多个文件，保存到path指定的位置上
// TODO: 根据webuploader传递的List创建文件层级结构
func UploadFiles(c *gin.Context, path string) {
	form, _ := c.MultipartForm()
	files := form.File["upload[]"]
	for _, file := range files {
		log.Println(file.Filename)
		err := c.SaveUploadedFile(file, file.Filename)
		if err != nil {
			e := fmt.Errorf("文件[%s]上传失败，错误信息：%s", file.Filename, err.Error())
			Models.ResultFail(c, Code.UPLOAD_ERROR, e)
			return
		}
	}
	// TODO: 这个长度不一定正确
	Models.ResultSuccessWithMsg(c, nil, fmt.Sprintf("全部文件上传成功，共上传了%d个文件", len(files)))
}
