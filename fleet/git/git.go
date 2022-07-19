package git

import (
	"context"

	"github.com/google/go-github/v45/github"
	"github.com/tgs266/fleet/config"
	"github.com/tgs266/fleet/rest-gen/generated/com/fleet/common"
	"golang.org/x/oauth2"
)

type GitClient struct {
	AccessToken  string
	GithubClient *github.Client
}

var gc *GitClient

func NewClient(config *config.Config) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Git.GithubAccessToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	gc = &GitClient{
		AccessToken:  config.Git.GithubAccessToken,
		GithubClient: client,
	}

}

func GetConfig(source common.GitSource) (string, error) {
	fc, _, _, err := gc.GithubClient.Repositories.GetContents(context.Background(), source.Owner, source.Repo, "fleet/deployment/config.yaml", &github.RepositoryContentGetOptions{})
	if err != nil {
		return "", err
	}
	return fc.GetContent()
}
