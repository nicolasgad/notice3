package common

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/storage/memory"
	"regexp"
	"strings"
)

func CheckGitUrl(url string) bool {
	reg := regexp.MustCompile(`^((git|ssh|http(s){0,1})|(git@[\w\.]{1,}))(:(//){0,1})([\w\.@:/\-~]{1,})(/){0,1}$`)
	return reg.MatchString(url)
}

func CheckHttpGitUrl(url string) bool {
	reg := regexp.MustCompile(`^(http(s){0,1})(:(//){0,1})([\w\.@:/\-~]{1,})(/){0,1}$`)
	return reg.MatchString(url)
}

func CheckSshGitUrl(url string) bool {
	reg := regexp.MustCompile(`^((git|ssh)|(git@[\w\.]{1,}))(:(//){0,1})([\w\.@:/\-~]{1,})(/){0,1}$`)
	return reg.MatchString(url)
}

func NormalizeGitUrlToSsh(url string) string {
	replaceArr := [][]string{
		{"https://github.com/", "git@github.com:"},
		{"https://gitlab.com/", "git@gitlab.com:"},
		{"https://git.coding.net/", "git@git.coding.net:"},
		{"https://gitee.com/", "git@gitee.com:"},
		{"https://bitbucket.org/", "git@bitbucket.org:"},
	}
	result := strings.ToLower(url)
	for _, pattern := range replaceArr {
		result = strings.Replace(result, pattern[0], pattern[1], 1)
	}

	result = strings.TrimSuffix(result, "/")

	if !strings.HasSuffix(result, ".git") {
		result += ".git"
	}

	return result
}

func NormalizeGitUrlToHttp(url string) string {
	replaceArr := [][]string{
		{"git@github.com:", "https://github.com/"},
		{"git@gitlab.com:", "https://gitlab.com/"},
		{"git@git.coding.net:", "https://git.coding.net/"},
		{"git@gitee.com:", "https://gitee.com/"},
		{"git@bitbucket.org:", "https://bitbucket.org/"},
	}
	result := strings.ToLower(url)
	for _, pattern := range replaceArr {
		result = strings.Replace(result, pattern[0], pattern[1], 1)
	}

	result = strings.TrimSuffix(result, "/")

	if !strings.HasSuffix(result, ".git") {
		result += ".git"
	}

	return result
}

func GetListFromGitRemote(gitUrl string, auth transport.AuthMethod) (branchs []string, authError error) {
	// Create a new repository
	logrus.Info("git init in memory")
	repository, err := git.Init(memory.NewStorage(), nil)
	AssertOrInterrupt(err)

	// Add a new remote, with the default fetch refspec
	remote, err := repository.CreateRemote(&config.RemoteConfig{
		Name: "origin",
		URLs: []string{gitUrl},
	})
	AssertOrInterrupt(err)

	logrus.Info("git ls-remote ", gitUrl)
	rfs, authError := remote.List(&git.ListOptions{
		Auth: auth,
	})
	if authError != nil {
		return
	}

	// get branchs
	for _, rf := range rfs {
		if rf.Type() == plumbing.HashReference && rf.Name().IsBranch() {
			branchs = append(branchs, string(rf.Name()))
		}
	}
	return
}
