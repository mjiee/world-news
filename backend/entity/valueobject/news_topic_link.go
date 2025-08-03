package valueobject

// NewsTopicLink represents a link associated with a news category.
type NewsTopicLink struct {
	Topic string `json:"topic"`
	URL   string `json:"url"`
}

// NewNewsTopicLink creates a new news topic link.
func NewNewsTopicLink(topic string, url string) *NewsTopicLink {
	return &NewsTopicLink{
		Topic: topic,
		URL:   url,
	}
}

// Compare compare two news topic link.
func (n *NewsTopicLink) Compare(other *NewsTopicLink) bool {
	return n.URL == other.URL && n.Topic == other.Topic
}
