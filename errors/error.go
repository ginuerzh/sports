// error
package errors

import (
	"fmt"
)

const (
	NoError   = iota
	AuthError = 1000 + iota
	UserExistError
	AccessError
	DbError
	_
	JsonError
	NotFoundError
	PasswordError
	InvalidFileError
	HttpError
	FileNotFoundError
	_
	NotExistsError
	InvalidAddrError
	InvalidMsgError
	DeviceTokenError
	ReviewNotFoundError
	InviteCodeError
	FileTooLargeError
	FileUploadError
	UnimplementedError
)

var errMap map[int]string = map[int]string{
	NoError:             "success",
	AuthError:           "auth error",
	UserExistError:      "user exists",
	AccessError:         "access token error",
	DbError:             "database error",
	JsonError:           "json data error",
	NotFoundError:       "not found",
	PasswordError:       "password invalid",
	InvalidFileError:    "file invalid",
	HttpError:           "http error",
	FileNotFoundError:   "file not found",
	NotExistsError:      "not exists",
	InvalidAddrError:    "address invalid",
	InvalidMsgError:     "message invalid",
	DeviceTokenError:    "device token invalid",
	ReviewNotFoundError: "review not found",
	InviteCodeError:     "invite code invalid",
	FileTooLargeError:   "file too large",
	FileUploadError:     "file upload error",
	UnimplementedError:  "unimplemented",
}

type Error struct {
	Id   int    `json:"error_id"`
	Desc string `json:"error_desc"`
}

func NewError(id int, desc ...string) *Error {
	s := errMap[id]
	if len(desc) > 0 {
		s = desc[0]
	}
	return &Error{Id: id, Desc: s}
}

func (e *Error) Error() string {
	return fmt.Sprintf("%d %s", e.Id, e.Desc)
}
