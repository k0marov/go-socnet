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
