// file
package controllers

import (
	"github.com/ginuerzh/sports/errors"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"github.com/nfnt/resize"
	weedo "gopkg.in/ginuerzh/weedo.v0"
	"gopkg.in/go-martini/martini.v1"
	//"io"
	"time"
	//"labix.org/v2/mgo/bson"
	"bytes"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"mime/multipart"
	"net/http"
)

var (
	Weedfs *weedo.Client
	//weedfs = weedo.NewClient("localhost:9334")
)

func init() {
	image.RegisterFormat("jpeg", "\xff\xd8", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("png", "\x89\x50\x4E\x47\x0D\x0A\x1A\x0A", png.Decode, png.DecodeConfig)
	image.RegisterFormat("gif", "\x47\x49\x46\x38\x39\x61", gif.Decode, gif.DecodeConfig)
	image.RegisterFormat("bmp", "\x42\x4D", bmp.Decode, bmp.DecodeConfig)
}

func BindFileApi(m *martini.ClassicMartini) {
	m.Post("/1/file/upload",
		binding.Form(fileUploadForm{}),
		//ErrorHandler,
		//checkTokenHandler,
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
	Width  int `form:"width"`
	Height int `form:"height"`
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
	redis *models.RedisLogger /*, user *models.Account*/, form fileUploadForm) {
	user := &models.Account{}

	if len(form.Token) > 0 {
		id := redis.OnlineUser(form.Token)
		if find, _ := user.FindByUserid(id); !find {
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.NotFoundError))
			return
		}
	}

	var file models.File

	filedata, header, err := request.FormFile("filedata")
	if err != nil {
		log.Println(err)
		writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.FileNotFoundError))
		return
	}

	if form.Width > 0 || form.Height > 0 {
		img, _, err := image.Decode(filedata)
		if err != nil {
			log.Println(err)
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.InvalidFileError))
			return
		}

		fid, err := Weedfs.Master().AssignN(2)
		if err != nil {
			log.Println(err)
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.FileUploadError))
			return
		}
		file.Fid = fid

		thumbnail := resize.Thumbnail(uint(form.Width), uint(form.Height), img, resize.MitchellNetravali)
		vol, err := Weedfs.Volume(fid, "")
		if err != nil {
			log.Println(err)
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.FileUploadError))
			return
		}

		buf := &bytes.Buffer{}
		if err := jpeg.Encode(buf, thumbnail, nil); err != nil {
			log.Println(err)
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.FileUploadError))
			return
		}

		length, err := vol.Upload(fid, 0, header.Filename, "image/jpeg", buf)
		if err != nil {
			log.Println(err)
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.FileUploadError))
			return
		}
		file.Length = length

		filedata.Seek(0, 0)
		if _, err := vol.Upload(fid, 1, header.Filename, header.Header.Get("Content-Type"), filedata); err != nil {
			log.Println(err)
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.FileUploadError))
			return
		}
	} else {
		fid, length, err := Weedfs.Master().Submit(header.Filename, header.Header.Get("Content-Type"), filedata)
		if err != nil {
			writeResponse(request.RequestURI, resp, nil, errors.NewError(errors.FileUploadError))
			return
		}
		//log.Println(fid, length, header.Filename, header.Header.Get("Content-Type"))

		file.Fid = fid
		file.Length = length
	}

	filedata.Seek(0, 0)

	file.Name = header.Filename
	file.ContentType = header.Header.Get("Content-Type")
	file.Md5 = FileMd5(filedata)
	file.Owner = user.Id
	file.UploadDate = time.Now()
	if err := file.Save(); err != nil {
		writeResponse(request.RequestURI, resp, nil, err)
		return
	}

	url, _, _ := Weedfs.GetUrl(file.Fid)
	respData := map[string]interface{}{"fileid": file.Fid, "fileurl": url}

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
