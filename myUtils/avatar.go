package myUtils

import (
	"bytes"
	"encoding/base64"
	"github.com/o1egl/govatar"
	"image"
	"image/jpeg"
	"unsafe"
)

//GetBase64AvatarByStr 通过str获取base64编码的头像
func GetBase64AvatarByStr(s string) (imgCode *string, err error) {
	if img, err1 := govatar.GenerateForUsername(govatar.MALE, s); err1 != nil {
		err = err1
	} else {
		imgCode, err = img2Base64(img)
	}
	return
}

func img2Base64(img image.Image) (imgBase64 *string, err error) {
	emptyBuff := bytes.NewBuffer(nil)                       //开辟一个新的空buff
	if err = jpeg.Encode(emptyBuff, img, nil); err != nil { //img写入到buff
		return
	}
	dist := make([]byte, 50000)                        //开辟存储空间
	base64.StdEncoding.Encode(dist, emptyBuff.Bytes()) //buff转成base64
	index := bytes.IndexByte(dist, 0)                  //这里要注意，因为申请的固定长度数组，所以没有被填充完的部分需要去掉，负责输出可能出错
	baseImage := append([]byte("data:image/bmp;base64,"), dist[0:index]...)

	imgBase64 = (*string)(unsafe.Pointer(&baseImage))
	return
}
