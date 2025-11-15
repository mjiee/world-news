package service

// PodcastService represents the podcast service.
type PodcastService interface {
}

type podcastService struct {
}

func NewPodcastService() PodcastService {
	return &podcastService{}
}
