package response

type Response interface {
	// GetStatus returns status code
	GetStatus() int
}

type ErrorResponse struct {
	Status int    `json:"status"`
	Error  string `json:"error"`
}

func (e *ErrorResponse) GetStatus() int {
	return e.Status
}

type OkResponse struct {
	Status int `json:"status"`
	Data   any `json:"data,omitempty"`
}

func (o *OkResponse) GetStatus() int {
	return o.Status
}

// OkResponseEmpty is OkResponse without data, I use it because swagger is shit, or I don't know how to set up it
type OkResponseEmpty struct {
	Status int `json:"status"`
}

func (o *OkResponseEmpty) GetStatus() int {
	return o.Status
}

// Error create new response with error.
// status - status code for error.
// err - error (not string)
func Error(status int, err error) Response {
	return &ErrorResponse{
		Status: status,
		Error:  err.Error(),
	}
}

// OK create new response with status 200, without error
func OK(args ...any) Response {
	if len(args) == 1 {
		return &OkResponse{
			Status: 200,
			Data:   args[0],
		}
	} else if len(args) == 0 {
		return &OkResponseEmpty{
			Status: 200,
		}
	}
	return nil
}
