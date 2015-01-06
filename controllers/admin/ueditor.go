package admin

import (
	"bytes"
	"encoding/json"
	"github.com/ginuerzh/sports/controllers"
	"github.com/ginuerzh/sports/models"
	"github.com/martini-contrib/binding"
	"gopkg.in/go-martini/martini.v1"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	config = make(map[string]interface{})
)

func init() {
	file, err := os.Open("./public/ueditor/php/config.json")
	if err != nil {
		log.Println(err)
		return
	}

	config, err = parseConfig(file)
	if err != nil {
		log.Println(err)
	}
}

func BindUeditorApi(m *martini.ClassicMartini) {
	m.Get("/ueditor/controller",
		binding.Form(ueditorForm{}),
		ueditorHandler)
	m.Post("/ueditor/controller",
		binding.Form(ueditorForm{}),
		ueditorHandler)
	m.Options("/ueditor/controller",
		binding.Form(ueditorForm{}),
		ueditorHandler)
}

func parseConfig(r io.Reader) (map[string]interface{}, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	regex := regexp.MustCompile(`\/\*[\s\S]+?\*\/`)
	data = regex.ReplaceAll(data, []byte(""))

	m := make(map[string]interface{})
	json.NewDecoder(bytes.NewBuffer(data)).Decode(&m)
	return m, err
}

type ueditorForm struct {
	Action   string `form:"action"`
	Callback string `form:"callback"`
}

func ueditorHandler(r *http.Request, w http.ResponseWriter, form ueditorForm) {
	var result interface{}

	switch form.Action {
	case "config":
		result = config

	case "uploadimage":
		if r.Method == "OPTIONS" {
			writeResponse(w, nil)
			return
		}
		filedata, header, err := r.FormFile(config["fileFieldName"].(string))
		if err != nil {
			log.Println(err)
			result = map[string]interface{}{
				"state": "找不到上传文件",
			}
			break
		}

		fid, length, err := controllers.Weedfs.Master().Submit(header.Filename,
			header.Header.Get("Content-Type"), filedata)
		if err != nil {
			log.Println(err)
			result = map[string]interface{}{
				"state": "文件上传时出错",
			}
			break
		}

		var file models.File
		file.Fid = fid
		file.Name = header.Filename
		file.ContentType = header.Header.Get("Content-Type")
		file.Length = length
		file.Owner = "admin"
		file.UploadDate = time.Now()
		if err := file.Save(); err != nil {
			result = map[string]interface{}{
				"state": "文件保存时出错",
			}
			break
		}

		url, _, _ := controllers.Weedfs.GetUrl(fid)
		result = map[string]interface{}{
			"state":    "SUCCESS",
			"url":      url,
			"title":    file.Name,
			"original": file.Name,
			"type":     strings.Split(file.Name, ".")[1],
			"size":     file.Length,
		}
	default:
		result = map[string]interface{}{
			"state": "请求地址出错",
		}
	}

	if len(form.Callback) > 0 {
		if strings.Contains(form.Callback, "_") {
			data, _ := json.Marshal(result)
			writeRawResponse(w, "application/javascript", []byte(form.Callback+"("+string(data)+")"))
			return
		} else {
			result = map[string]interface{}{
				"state": "callback参数不合法",
			}
		}
	}

	data, _ := json.Marshal(result)
	writeRawResponse(w, "text/html; charset=utf-8", data)

}
