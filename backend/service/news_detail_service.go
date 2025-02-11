package service

type NewsDetailService interface {
}

type newsDetailService struct {
}

func NewNewsDetailService() NewsDetailService {
	return &newsDetailService{}
}
