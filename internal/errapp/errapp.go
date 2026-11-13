package errapp

import "errors"

var (
	ErrContestAlreadyExists = errors.New("Concurso já cadastrado")
	ErrNoDocuments          = errors.New("Jogo Não encontrado")
)

type ErrApp struct {
	Err     error
	Message string
}

func NewErrApp(err error, msg string) *ErrApp {
	return &ErrApp{
		Err:     err,
		Message: msg,
	}
}

type ErrAppReq struct {
	Err     error
	Message string
	Code    int
}

func NewErrAppReq(err error, msg string, code int) *ErrAppReq {
	return &ErrAppReq{
		Err:     err,
		Message: msg,
		Code:    code,
	}
}
