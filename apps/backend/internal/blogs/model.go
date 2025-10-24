package blogs

type Kind string

const (
	Person       Kind = "person"
	Organization Kind = "organization"
)

var KindByCollection = map[string]Kind{
	"people":        Person,
	"organizations": Organization,
}

type BlogInfo struct {
	BlogHref          string `json:"blogHref"`
	BlogName          string `json:"blogName"`
	LatestArticleHref string `json:"latestArticleHref"`
	LatestArticleName string `json:"latestArticleName"`
	Kind              Kind   `json:"kind"`
}
