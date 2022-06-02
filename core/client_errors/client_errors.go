package client_errors

type ClientError struct {
	ReadableDetail string
	DetailCode     string
}

func (ce ClientError) Error() string {
	return "An error which will be displayed to the client: " + ce.ReadableDetail
}

var InvalidJsonError = ClientError{
	DetailCode:     "invalid-json",
	ReadableDetail: "The provided request body is not valid JSON.",
}

var AvatarNotProvidedError = ClientError{
	DetailCode:     "no-avatar",
	ReadableDetail: "You should provide an image in 'avatar' field.",
}

var AvatarTooBigError = ClientError{
	DetailCode:     "big-avatar",
	ReadableDetail: "The avatar image you provided is too big.",
}

var BodyIsNotMultipartForm = ClientError{
	DetailCode:     "not-multipartform",
	ReadableDetail: "Post data for this endpoint should be provided as a multipart form.",
}
