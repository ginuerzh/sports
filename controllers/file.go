// file
package controllers

import (
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	weedo "gopkg.in/ginuerzh/weedo.v0"
	"gopkg.in/go-martini/martini.v1"
	//"io"
	"time"
	//"labix.org/v2/mgo/bson"
	"log"
	"mime/multipart"
	"net/http"
)

const (
	ImageDownloadV1Uri = "/1/image/get"
	FileDeleteV1Uri    = "/1/file/del"
)

var (
	Weedfs *weedo.Client
	//weedfs = weedo.NewClient("localhost:9334")
)

func BindFileApi(m *martini.ClassicMartini) {
	m.Post("/1/file/upload", binding.Form(fileUploadForm{}), ErrorHandler, fileUploadHandler)
	//m.Post("/1/file/upload", binding.MultipartForm(fileUploadForm2{}), ErrorHandler, fileUploadHandler2)
	//m.Get(ImageDownloadV1Uri, binding.Form(imageDownloadForm{}), ErrorHandler, imageDownloadHandler)
	//m.Post(FileDeleteV1Uri, binding.Json(fileDeleteForm{}), ErrorHandler, fileDeleteHandler)
}

type fileUploadForm struct {
	AccessToken string `form:"access_token" binding:"required"`
	//user        models.User `form:"-"`
}

func (form *fileUploadForm) Validate(e *binding.Errors, req *http.Request) {
	//log.Println(form.AccessToken)
	//form.user = userAuth(form.AccessToken, e)
}

type fileUploadForm2 struct {
	Token string                `form:"access_token"`
	File  *multipart.FileHeader `form:"filedata"`
}

func fileUploadHandler2(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form fileUploadForm2) {
	//_, err := form.File.Open()

	log.Println(form.File.Filename)
}

func fileUploadHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form fileUploadForm) {
	user := redis.OnlineUser(form.AccessToken)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.AccessError))
		return
	}

	filedata, header, err := request.FormFile("filedata")
	if err != nil {
		log.Println(err)
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.FileNotFoundError))
		return
	}

	fid, length, err := Weedfs.Master().Submit(header.Filename, header.Header.Get("Content-Type"), filedata)
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.FileUploadError))
		return
	}
	log.Println(fid, length, header.Filename, header.Header.Get("Content-Type"))

	filedata.Seek(0, 0)

	var file models.File
	file.Fid = fid
	file.Name = header.Filename
	file.ContentType = header.Header.Get("Content-Type")
	file.Length = length
	file.Md5 = FileMd5(filedata)
	file.Owner = user.Id
	file.UploadDate = time.Now()
	if err := file.Save(); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	url, _, _ := Weedfs.GetUrl(fid)
	respData := map[string]interface{}{"fileid": fid, "fileurl": url}

	writeResponse(request.RequestURI, resp, respData, nil)
}

/*
type imageDownloadForm struct {
	ImageId   string `form:"image_id" binding:"required"`
	ImageSize string `form:"image_size_type"`
}


func imageDownloadHandler(request *http.Request, resp http.ResponseWriter, form imageDownloadForm) {
	var file models.File

	if exist, err := file.FindByFid(form.ImageId); !exist {
		if err == errors.NoError {
			err = errors.FileNotFoundError
		}
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	fileData, err := weedo.Download(form.ImageId)
	if err != nil {
		writeResponse(request.RequestURI, resp, nil, errors.FileNotFoundError)
		return
	}
	defer fileData.Close()

	resp.Header().Set("Content-Length", strconv.FormatInt(file.Size, 10))
	io.Copy(resp, fileData)
}


func imageDownloadHandler(request *http.Request, resp http.ResponseWriter, form imageDownloadForm) {
	url := imageUrl(form.ImageId, ImageOriginal)

	respData := map[string]string{"image_url": url}
	writeResponse(request.RequestURI, resp, respData, errors.NoError)
}

type fileDeleteForm struct {
	Fid         string `json:"image_id" binding:"required"`
	AccessToken string `json:"access_token" binding:"required"`
	//user        models.User `json:"-"`
}

func (form *fileDeleteForm) Validate(e *binding.Errors, req *http.Request) {
	//form.user = userAuth(form.AccessToken, e)
}

func fileDeleteHandler(request *http.Request, resp http.ResponseWriter, redis *models.RedisLogger, form fileDeleteForm) {
	var file models.File

	user := redis.OnlineUser(form.AccessToken)
	if user == nil {
		writeResponse(request.RequestURI, resp, nil, errors.AccessError)
		return
	}

	if find, err := file.FindByFid(form.Fid); !find {
		if err == errors.NoError {
			err = errors.FileNotFoundError
		}
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	if file.Owner != user.Userid {
		writeResponse(request.RequestURI, resp, nil, errors.FileNotFoundError)
		return
	}

	err := file.Delete()
	writeResponse(request.RequestURI, resp, nil, err)
}
*/
