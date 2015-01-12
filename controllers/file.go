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

var (
	Weedfs *weedo.Client
	//weedfs = weedo.NewClient("localhost:9334")
)

func BindFileApi(m *martini.ClassicMartini) {
	m.Post("/1/file/upload",
		binding.Form(fileUploadForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		fileUploadHandler)
	m.Post("/1/file/delete",
		binding.Json(fileDeleteForm{}, (*Parameter)(nil)),
		ErrorHandler,
		checkTokenHandler,
		fileDeleteHandler)
	//m.Post("/1/file/upload", binding.MultipartForm(fileUploadForm2{}), ErrorHandler, fileUploadHandler2)
	//m.Get(ImageDownloadV1Uri, binding.Form(imageDownloadForm{}), ErrorHandler, imageDownloadHandler)
}

type fileUploadForm struct {
	parameter
}

type fileUploadForm2 struct {
	Token string                `form:"access_token"`
	File  *multipart.FileHeader `form:"filedata"`
}

func fileUploadHandler2(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, form fileUploadForm2) {
	//_, err := form.File.Open()

	log.Println(form.File.Filename)
}

func fileUploadHandler(request *http.Request, resp http.ResponseWriter,
	redis *models.RedisLogger, user *models.Account) {

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
	//log.Println(fid, length, header.Filename, header.Header.Get("Content-Type"))

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

type fileDeleteForm struct {
	Fids []string `json:"fids"`
	parameter
}

func fileDeleteHandler(r *http.Request, w http.ResponseWriter,
	user *models.Account, p Parameter) {

	form := p.(fileDeleteForm)

	for _, fid := range form.Fids {
		file := &models.File{Fid: fid}
		if find, _ := file.OwnedBy(user.Id); !find {
			continue
		}

		if err := file.Delete(); err == nil {
			Weedfs.Delete(fid, 1)
		}
	}

	writeResponse(r.RequestURI, w, nil, nil)
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


*/
