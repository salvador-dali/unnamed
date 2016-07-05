package imager

import (
	"../config"
	"../misc"
	"fmt"
	bimg "gopkg.in/h2non/bimg.v1"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const (
	minImgHeight = 400
	minImgWidth  = 600
	avatarBig    = 300
	avatarSmall  = 64
)

// have all the mime types that we accept and maps them to file extensions
var mimeToExtension = map[string]string{
	"image/jpeg": "jpg",
	"image/png":  "png",
	"image/webp": "webp",
}

func SaveTmpFileFromClient(w http.ResponseWriter, r *http.Request) (bool, string, string) {
	// make sure that the file is of correct size
	r.Body = http.MaxBytesReader(w, r.Body, config.Cfg.MaxImgSizeKb)
	clientFile, handler, err := r.FormFile("img")
	if err != nil {
		log.Println(err)
		return false, "", ""
	}
	defer clientFile.Close()

	if handler.Filename == "" {
		log.Println("No filename provided")
		return false, "", ""
	}

	if _, ok := handler.Header["Content-Type"]; !ok {
		log.Println("No content-type provided")
		return false, "", ""
	}

	// save the file locally in a temporary location
	fileName := fmt.Sprintf("%d_%s", time.Now().Unix(), misc.RandomString(10))
	fileLoc := fmt.Sprintf("images/tmp/%s", fileName)
	serverFile, err := os.OpenFile(fileLoc, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Println(err)
		return false, "", ""
	}
	defer serverFile.Close()
	io.Copy(serverFile, clientFile)

	// read that file and check that it is of correct type
	fs, err := os.Open(fileLoc)
	if err != nil {
		log.Println("error reading file")
		os.Remove(fileLoc)
		return false, "", ""
	}
	defer fs.Close()
	buff := make([]byte, 512) // http://golang.org/pkg/net/http/#DetectContentType
	fs.Read(buff)
	mime := http.DetectContentType(buff)
	ext, ok := mimeToExtension[mime]
	if !ok {
		log.Println("File with wrong MIME type", mime)
		os.Remove(fileLoc)
		return false, "", ""
	}

	return true, fileName, ext
}

func CheckTmpFileImgSize(fileName string, minWidth, minHeight int) (bool, *bimg.Image) {
	fileLoc := "images/tmp/" + fileName
	buffer, err := bimg.Read(fileLoc)
	os.Remove(fileLoc)
	if err != nil {
		log.Println(err)
		return false, nil
	}
	img := bimg.NewImage(buffer)

	sizeInfo, err := img.Size()
	if err != nil {
		log.Println(err)
		return false, nil
	}

	if sizeInfo.Width < minWidth || sizeInfo.Height < minHeight {
		log.Println("Image size is too small", sizeInfo)
		return false, nil
	}

	return true, img
}

func TmpToAvatar(fileName, ext string) bool {
	ok, img := CheckTmpFileImgSize(fileName, avatarBig, avatarBig)
	if !ok {
		return false
	}

	newImage, err := img.Thumbnail(avatarBig)
	if err != nil {
		log.Println(err)
		return false
	}
	bimg.Write(fmt.Sprintf("images/avatars/b/%s.%s", fileName, ext), newImage)

	newImage, err = img.Thumbnail(avatarSmall)
	if err != nil {
		log.Println(err)
		return false
	}
	bimg.Write(fmt.Sprintf("images/avatars/s/%s.%s", fileName, ext), newImage)

	return true
}
