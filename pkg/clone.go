package getcode

import (
	"context"
	"os"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/google/go-github/v57/github"
	"github.com/spf13/viper"
)

type gitRequest struct {
	credentials string
	protocol    string
	targetName  string
	targetType  string
	gitUser     string
	url         string
	projectPath string
}

func Clone(url string) {

	ViperConfig := loadViperConfig()
	requestData := parseInputUrl(url, ViperConfig)
	repos := getRepos(requestData)
	cloneRepos(requestData, repos)
}

func cloneRepos(requestData *gitRequest, repos []*github.Repository) {

	downloadPath := requestData.projectPath + requestData.url + "/" + requestData.targetName + "/"
	cloneUrl := requestData.protocol + "//" + requestData.url + "/" + requestData.targetName + "/"

	for _, repo := range repos {
		fullDownloadPath := downloadPath + *repo.Name
		fullCloneUrl := cloneUrl + *repo.Name
		git.PlainClone(fullDownloadPath, false, &git.CloneOptions{
			URL:      fullCloneUrl,
			Progress: os.Stdout,
		})
	}
}

func parseInputUrl(url string, viperConfig *ViperConfig) *gitRequest {

	var requestData gitRequest
	requestData.protocol = getProtocol(url)
	requestData.targetName = getTargetName(url)
	requestData.url = getGitUrl(url, &requestData)
	requestData.gitUser = getGitUser(requestData.targetName, viperConfig)
	requestData.credentials = getCredentials(requestData.gitUser, viperConfig, requestData.url)
	requestData.targetType = getTargetType(url, viperConfig, &requestData)
	requestData.projectPath = getProjectPath(viperConfig)

	return &requestData
}

func getProjectPath(viperConfig *ViperConfig) string {
	return viperConfig.ProjectPath
}

func getTargetType(url string, viperConfig *ViperConfig, collectedData *gitRequest) string {

	if checkTargetInKeys(collectedData.targetName, viperConfig) {
		return "Authenticated"
	} else {
		return pullTargetType(collectedData)
	}
}

func getRepos(collectedData *gitRequest) []*github.Repository {

	var allRepos []*github.Repository
	client := github.NewClient(nil).WithAuthToken(collectedData.credentials)
	ctx := context.Background()

	listOptions := github.ListOptions{
		Page:    1,
		PerPage: 100,
	}

	if collectedData.targetType == "Organization" {
		opt := &github.RepositoryListByOrgOptions{Type: "owner", Sort: "updated", Direction: "desc", ListOptions: listOptions}

		for {
			repos, resp, err := client.Repositories.ListByOrg(ctx, collectedData.targetName, opt)
			ErrorCheck(err, "Failed to get a list form orgs")
			allRepos = append(allRepos, repos...)

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}
	} else if collectedData.targetType == "User" {

		opt := &github.RepositoryListByUserOptions{Type: "owner", Sort: "updated", Direction: "desc", ListOptions: listOptions}

		for {
			repos, resp, err := client.Repositories.ListByUser(ctx, collectedData.targetName, opt)
			ErrorCheck(err, "Failed to get a list form user")
			allRepos = append(allRepos, repos...)

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}
	} else {

		opt := &github.RepositoryListByAuthenticatedUserOptions{Type: "owner", Sort: "updated", Direction: "desc", ListOptions: listOptions}

		for {
			repos, resp, err := client.Repositories.ListByAuthenticatedUser(ctx, opt)
			ErrorCheck(err, "Failed to get a list from auth user")
			allRepos = append(allRepos, repos...)

			if resp.NextPage == 0 {
				break
			}

			opt.Page = resp.NextPage
		}
	}

	return allRepos
}

func pullTargetType(collectedData *gitRequest) string {

	listOptions := github.ListOptions{
		Page:    1,
		PerPage: 1,
	}

	client := github.NewClient(nil).WithAuthToken(collectedData.credentials)
	opt := &github.RepositoryListByUserOptions{Type: "owner", Sort: "updated", Direction: "desc", ListOptions: listOptions}
	repos, _, err := client.Repositories.ListByUser(context.Background(), collectedData.targetName, opt)
	ErrorCheck(err, "Failed to get a list form user")

	return *repos[0].Owner.Type
}

func getGitUrl(url string, collectedData *gitRequest) string {

	prefix := collectedData.protocol + "//"
	suffix := "/" + collectedData.targetName
	url = strings.TrimPrefix(url, prefix)
	url = strings.TrimSuffix(url, suffix)

	return url
}

func getGitUser(targetName string, viperConfig *ViperConfig) string {

	if checkTargetInKeys(targetName, viperConfig) {
		return targetName
	} else if found, user := checkTargetInValues(targetName, viperConfig); found {
		return user
	} else {
		return viperConfig.DefaultUser
	}
}

func getCredentials(user string, viperConfig *ViperConfig, gitUrl string) string {

	fileData, err := os.ReadFile(viperConfig.GitCredentalPath)
	ErrorCheck(err, "reading .git-gredentials failed")
	var cred string
	prefix := "https://" + user + ":"
	suffix := "@" + gitUrl

	for _, cred = range strings.Split(string(fileData), "\n") {
		if strings.HasPrefix(cred, prefix) {
			cred = strings.TrimPrefix(cred, prefix)
			cred = strings.TrimSuffix(cred, suffix)
			break
		}
	}
	return cred
}

func checkTargetInValues(targetName string, viperConfig *ViperConfig) (bool, string) {

	for user, orgs := range viperConfig.OrgUsers {
		if isValueinList(orgs, targetName) {
			return true, user
		}
	}

	return false, ""
}

func checkTargetInKeys(targetName string, viperConfig *ViperConfig) bool {
	return isMapOfListKey(targetName, viperConfig.OrgUsers)
}

func getTargetName(url string) string {
	return strings.Split(url, "/")[len(strings.Split(url, "/"))-1]
}

func getProtocol(url string) string {
	return strings.Split(url, "//")[0]
}

func loadViperConfig() *ViperConfig {

	var ViperConfig ViperConfig
	ViperConfig.GitCredentalPath = viper.GetString("git_credental_path")
	ViperConfig.ProjectPath = viper.GetString("project_path")
	ViperConfig.OrgUsers = viper.GetStringMapStringSlice("org_users")
	ViperConfig.DefaultUser = viper.GetString("default_user")

	return &ViperConfig
}
