package Utils

import "testing"

func TestUnzipToCurrentPath(t *testing.T) {
	path := "./src.zip"
	err := UnzipToCurrentPath(path)
	if err != nil{
		t.Fatal(err)
	}
	t.Log("解压缩完成!")
}