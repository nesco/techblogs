// Package blogs provides domain models and data access for tech blog management.
package blogs

type Kind string

const (
	Individual   Kind = "individual"
	Organization Kind = "organization"
)

var KindByCollection = map[string]Kind{
	"people":        Individual,
	"organizations": Organization,
}

type BlogInfo struct {
	BlogHref          string `json:"blogHref"`
	BlogName          string `json:"blogName"`
	LatestArticleHref string `json:"latestArticleHref"`
	LatestArticleName string `json:"latestArticleName"`
	Kind              Kind   `json:"kind"`
	GitHubHref        string `json:"githubHref"`
}

type BlogConfig struct {
	BlogName            string
	BlogHref            string
	Kind                Kind
	ArticleHrefSelector string
	ArticleNameSelector string
	GitHubHref          string
}
