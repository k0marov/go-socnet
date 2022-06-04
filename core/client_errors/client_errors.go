package client_errors

import "net/http"

type ClientError struct {
	ReadableDetail string
	DetailCode     string
	HTTPCode       int
}

func (ce ClientError) Error() string {
	return "An error which will be displayed to the client: " + ce.ReadableDetail
}

var NoError = ClientError{
	DetailCode:     "no-error",
	ReadableDetail: "This error shouldn't be displayed. It is used as a zero value.",
	HTTPCode:       0,
}

var InvalidJsonError = ClientError{
	DetailCode:     "invalid-json",
	ReadableDetail: "The provided request body is not valid JSON.",
	HTTPCode:       http.StatusBadRequest,
}

var AvatarNotProvidedError = ClientError{
	DetailCode:     "no-avatar",
	ReadableDetail: "You should provide an image in 'avatar' field.",
	HTTPCode:       http.StatusBadRequest,
}

var AvatarTooBigError = ClientError{
	DetailCode:     "big-avatar",
	ReadableDetail: "The avatar image you provided is too big.",
	HTTPCode:       http.StatusBadRequest,
}

var BodyIsNotMultipartForm = ClientError{
	DetailCode:     "not-multipartform",
	ReadableDetail: "Post data for this endpoint should be provided as a multipart form.",
	HTTPCode:       http.StatusBadRequest,
}

var ProfileNotFound = ClientError{
	DetailCode:     "profile-not-found",
	ReadableDetail: "The requested profile was not found",
	HTTPCode:       http.StatusNotFound,
}

var AboutTooLong = ClientError{
	DetailCode:     "about-long",
	ReadableDetail: "The about field is too long",
	HTTPCode:       http.StatusBadRequest,
}

var NonImageAvatar = ClientError{
	DetailCode:     "avatar-non-image",
	ReadableDetail: "The provided avatar is not in one of supported image codecs.",
	HTTPCode:       http.StatusBadRequest,
}

var NonSquareAvatar = ClientError{
	DetailCode:     "avatar-non-square",
	ReadableDetail: "The provided avatar is a valid image, but it is not square.",
	HTTPCode:       http.StatusBadRequest,
}
