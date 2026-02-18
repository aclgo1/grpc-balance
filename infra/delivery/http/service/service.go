package service

import "net/http"

type HttpService struct {
}

func NewHttpService() *HttpService {
	return &HttpService{}
}

func CreditPix(w http.ResponseWriter, r *http.Request) {

}
