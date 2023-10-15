package utils

type MyError struct {
    Code int64
    Msg  string
}

func (e *MyError) Error() string {
    return e.Msg
}

func NewError(code int64, message string) *MyError {
    return &MyError{
        Code: code,
        Msg:  message,
    }
}
