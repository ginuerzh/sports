// weibo
package controllers

import (
	"encoding/json"
	"errors"
	//"fmt"
	"bufio"
	"bytes"
	"io"
	//"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	//"os"
	"crypto/tls"
	"strconv"
	"strings"
)

const (
	WeiboApiUserShow     = "https://api.weibo.com/2/users/show.json"
	WeiboApiStatusUpdate = "https://api.weibo.com/2/statuses/update.json"
	WeiboApiFriends      = "https://api.weibo.com/2/friendships/friends/bilateral.json"
)

var (
	Proxy = "chnpxy01.cn.ta-mp.com:8080"
)

type weiboError struct {
	Request   string `json:"request"`
	ErrorDesc string `json:"error"`
	ErrCode   int    `json:"error_code"`
}

type weiboUser struct {
	Id          uint64 `json:"id"`
	ScreenName  string `json:"screen_name"`
	Gender      string `json:"gender"`
	Url         string `json:"url"`
	Avatar      string `json:"avatar_large"`
	Location    string `json:"location"`
	Description string `json:"description"`
}

type weiboUserShow struct {
	weiboUser
	weiboError
}

type weiboFriends struct {
	Users []weiboUser `json:"users"`
	Total int         `json:"total_number"`
	weiboError
}

func connect(addr string, proxy string) (net.Conn, error) {
	if len(proxy) == 0 {
		return nil, errors.New("No proxy")
	}

	conn, err := net.Dial("tcp", proxy)
	if err != nil {
		return nil, err
	}

	w := bufio.NewWriter(conn)
	w.WriteString("CONNECT " + addr + " HTTP/1.1\r\n")
	w.WriteString("Host: " + addr + "\r\n")
	w.WriteString("Proxy-Connection: keep-alive\r\n\r\n")
	if err = w.Flush(); err != nil {
		conn.Close()
		return nil, err
	}

	b, err := read(conn)
	if err != nil {
		conn.Close()
		return nil, err
	}

	status := string(b)
	//log.Println(status)
	if !strings.Contains(status, "200") {
		conn.Close()
		return nil, errors.New(status)
	}

	config := &tls.Config{InsecureSkipVerify: true}
	cli := tls.Client(conn, config)

	return cli, cli.Handshake()
}

func requestWeiboApi(method string, url string, body io.Reader, result interface{}) (err error) {
	var resp *http.Response

	conn, err := connect("api.weibo.com:443", Proxy)

	switch strings.ToUpper(method) {
	case "GET":
		if err != nil {
			log.Println(err)
			resp, err = http.Get(url)
		} else {
			defer conn.Close()
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Println(err)
				return err
			}

			if err = req.WriteProxy(conn); err != nil {
				log.Println(err)
				return err
			}
			b, err := read(conn)
			if err != nil {
				log.Println(err)
				return err
			}

			resp, err = http.ReadResponse(bufio.NewReader(bytes.NewBuffer(b)), req)
		}

	case "POST":
		resp, err = http.Post(url, "application/json", body)
	default:
		err = errors.New("Unsupported http method: " + method)
	}
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(result)
}

func GetWeiboUserInfo(uid string, appToken string) (weiboUser, error) {
	result := &weiboUserShow{}

	v := url.Values{}
	v.Set("uid", uid)
	v.Set("access_token", appToken)

	if err := requestWeiboApi("GET", WeiboApiUserShow+"?"+v.Encode(), nil, result); err != nil {
		log.Println(err)
		return weiboUser{}, err
	}

	if result.ErrCode > 0 {
		return weiboUser{}, errors.New(strconv.Itoa(result.ErrCode) + ": " + result.ErrorDesc)
	}

	return result.weiboUser, nil
}

func GetWeiboFriends(appkey string, uid string, appToken string) (users []weiboUser, err error) {
	result := &weiboFriends{}
	page := 0
	v := url.Values{}
	v.Set("source", appkey)
	v.Set("access_token", appToken)
	v.Set("uid", uid)

	for {
		page++
		v.Set("page", strconv.Itoa(page))
		//log.Println("get page", WeiboApiFriends+"?"+v.Encode())
		if err = requestWeiboApi("GET", WeiboApiFriends+"?"+v.Encode(), nil, result); err != nil {
			break
		}
		if result.ErrCode > 0 {
			err = errors.New(strconv.Itoa(result.ErrCode) + ": " + result.ErrorDesc)
			break
		}

		users = append(users, result.Users...)
		if len(result.Users) == 0 || len(users) >= result.Total {
			break
		}
	}

	return
}

func read(conn net.Conn) ([]byte, error) {
	var b []byte
	size := 1024
	buf := make([]byte, size)

	for {
		n, err := conn.Read(buf)
		if err != nil {
			return b, err
		}
		b = append(b, buf[:n]...)
		if n < size {
			break
		}
	}

	return b, nil
}
