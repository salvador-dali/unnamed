package imager

import (
	"../misc"
	"fmt"
	bimg "gopkg.in/h2non/bimg.v1"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	minImgHeight = 400
	minImgWidth  = 600
	avatarBig    = 300
	avatarSmall  = 64

	imgNormalHeight = 600
	imgNormalWidth  = 800
	imgBigHeight    = 900
	imgBigWidth     = 1200

	maxFileSize = 5 * 1024 * 1024 // 5 Mb
)

// have all the mime types that we accept and maps them to file extensions
var mimeToExtension = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

// getTmpLocation is a helper which returns a location of a temporary file
func getTmpLocation(fileName string) string {
	return "images/tmp/" + fileName
}

// SaveTmpFileFromClient checks that the file is below the maximum possible size in Kb and
// saves it on a disk in a temporary folder. It detects the MIME-type of the image and suggests
// an extension based on the MIME-type. If anything is wrong, the file is removed
func SaveTmpFileFromClient(w http.ResponseWriter, r *http.Request) (bool, string, string) {
	// make sure that the file is of correct size
	r.Body = http.MaxBytesReader(w, r.Body, maxFileSize)
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
	fileLoc := getTmpLocation(fileName)
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

// checkTmpFileImgSize makes sure that the dimensions of the temporary image are above min height/width
func checkTmpFileImgSize(fileName string, minHeight, minWidth int) (bool, *bimg.Image) {
	buffer, err := bimg.Read(getTmpLocation(fileName))
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

// findBestDimensions finds the most suitable dimensions for the resize of original image.
// It makes sure that the new dimensions are maximum possible and the aspect ratio is preserved
func findBestDimensions(imgHeight, imgWidth, maxHeight, maxWidth int) (bool, int, int) {
	bestArea, bestHeight, bestWidth := 0, 0, 0

	height := imgHeight * maxWidth / imgWidth
	if height <= imgHeight && bestArea < maxWidth*height {
		bestArea = maxWidth * height
		bestHeight, bestWidth = height, maxWidth
	}

	width := imgWidth * maxHeight / imgHeight
	if width <= imgWidth && bestArea < maxHeight*width {
		bestArea = maxHeight * width
		bestHeight, bestWidth = maxHeight, width
	}

	return bestArea != 0, bestHeight, bestWidth
}

// TmpToAvatar converts a temporary file into a correctly resized avatar. Removes tmp file
func TmpToAvatar(fileName, ext string) (bool, string) {
	ok, img := checkTmpFileImgSize(fileName, avatarBig, avatarBig)
	fullFileName := fileName + ext
	os.Remove(getTmpLocation(fileName))
	if !ok {
		return false, ""
	}

	newImage, err := img.Thumbnail(avatarBig)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	bimg.Write("images/avatars/b/"+fullFileName, newImage)

	newImage, err = img.Thumbnail(avatarSmall)
	if err != nil {
		log.Println(err)
		return false, ""
	}
	bimg.Write("images/avatars/s/"+fullFileName, newImage)

	return true, fullFileName
}

// TmpToAvatar converts a temporary file into a correctly resized purchase. Removes tmp file
func TmpToPurchase(fileName, ext string) (bool, string) {
	ok, img := checkTmpFileImgSize(fileName, minImgHeight, minImgWidth)
	fullFileName := fileName + ext
	os.Remove(getTmpLocation(fileName))
	if !ok {
		os.Remove(getTmpLocation(fileName))
		return false, ""
	}

	sizeInfo, _ := img.Size()
	imgHeight, imgWidth := sizeInfo.Height, sizeInfo.Width

	if ok, h, w := findBestDimensions(imgHeight, imgWidth, imgBigHeight, imgBigWidth); ok {
		if newImage, err := img.Resize(w, h); err != nil {
			log.Println(err)
			os.Remove(getTmpLocation(fileName))
			return false, ""
		} else {
			bimg.Write("images/purchases/b/"+fullFileName, newImage)
		}
	}

	if ok, h, w := findBestDimensions(imgHeight, imgWidth, imgNormalHeight, imgNormalWidth); ok {
		if newImage, err := img.Resize(w, h); err != nil {
			log.Println(err)
			os.Remove(getTmpLocation(fileName))
			return false, ""
		} else {
			bimg.Write("images/purchases/m/"+fullFileName, newImage)
		}
	}

	return true, fullFileName
}

// verifyFile checks that the file exists at a specific location and was created in a right time
func verifyFile(fileName, location string) bool {
	if _, err := os.Stat(location + fileName); err == nil {
		// file exist. Check when it was generated
		timestamp, err := strconv.ParseInt(strings.Split(fileName, "_")[0], 10, 64)
		if err != nil {
			return false
		}

		now := time.Now()
		// make sure that the image was created somewhere in between Now and one day before Now
		return now.AddDate(0, -1, 0).Unix() <= timestamp && timestamp <= now.Unix()
	}

	// file does not exist
	return false
}

// IsAvatarValid makes sure that path to avatar is valid. Empty avatar is also valid
func IsAvatarValid(fileName string) bool {
	if fileName == "" {
		return true
	}

	return verifyFile(fileName, "images/avatars/b/")
}

// IsPurchaseValid makes sure that path to avatar is valid
func IsPurchaseValid(fileName string) bool {
	return verifyFile(fileName, "images/purchases/m/")
}
