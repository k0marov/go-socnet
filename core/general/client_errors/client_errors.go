package client_errors

import "net/http"

type ClientError struct {
	ReadableDetail string `json:"readable_detail"`
	DetailCode     string `json:"detail_code"`
	HTTPCode       int    `json:"-"`
}

func (ce ClientError) Error() string {
	return "An error which will be displayed to the client: " + ce.ReadableDetail
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

var NotFound = ClientError{
	DetailCode:     "not-found",
	ReadableDetail: "The requested entity was not found",
	HTTPCode:       http.StatusNotFound,
}

var AboutTooLong = ClientError{
	DetailCode:     "about-long",
	ReadableDetail: "The about field is too long",
	HTTPCode:       http.StatusBadRequest,
}

var InvalidImage = ClientError{
	DetailCode:     "avatar-non-image",
	ReadableDetail: "The provided image is not in one of supported image codecs.",
	HTTPCode:       http.StatusBadRequest,
}

var NonSquareAvatar = ClientError{
	DetailCode:     "avatar-non-square",
	ReadableDetail: "The provided avatar is a valid image, but it is not square.",
	HTTPCode:       http.StatusBadRequest,
}

var IdNotProvided = ClientError{
	DetailCode:     "no-id",
	ReadableDetail: "The request url was expected to have an 'id' query parameter, but there was none.",
	HTTPCode:       http.StatusBadRequest,
}

var FollowingYourself = ClientError{
	DetailCode:     "following-yourself",
	ReadableDetail: "You cannot follow yourself.",
	HTTPCode:       http.StatusBadRequest,
}

var InsufficientPermissions = ClientError{
	DetailCode:     "insufficient-permissions",
	ReadableDetail: "You are not authorized to perform this action.",
	HTTPCode:       http.StatusUnauthorized,
}

var LikingYourself = ClientError{
	DetailCode:     "liking-yourself",
	ReadableDetail: "You cannot like your own content",
	HTTPCode:       http.StatusBadRequest,
}

var TextTooLong = ClientError{
	DetailCode:     "long-text",
	ReadableDetail: "The provided text is too long.",
	HTTPCode:       http.StatusBadRequest,
}

var EmptyText = ClientError{
	DetailCode:     "empty-text",
	ReadableDetail: "Text cannot be empty.",
	HTTPCode:       http.StatusBadRequest,
}
