package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/google/go-github/v69/github"
)

type Visibility string
type IssuesState string

const (
	// URLREGEXP = `(\/[\w\-]+)+` // TODO: optimize - remove first slask

	// visibility
	Public   Visibility = "public"
	Private  Visibility = "private"
	Internal Visibility = "internal"

	// issuesState
	Open   IssuesState = "open"
	closed IssuesState = "closed"
)

type Owner struct {
	Id    int64  `json:"id"`
	Type  string `json:"type"`
	Login string `json:"login"`
	Url   string `json:"html_url"`
}

type Permissions struct {
	Admin    bool `json:"admin"`
	Push     bool `json:"push"`
	Pull     bool `json:"pull"`
	Maintain bool `json:"maintain,omitempty"`
}

type License struct {
	Key  string `json:"key"`
	Name string `json:"name"`
	Url  string `json:"html_url,omitempty"`
}

type Label struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Reactions struct {
	Total int `json:"toutal_count"`
	Up    int `json:"+1"`
	Down  int `json:"-1"`
	Laugh int `json:"laugh"`
	Conf  int `json:"confused"`
	Heart int `json:"heart"`
	Hoop  int `json:"hooray"`
	Eyes  int `json:"eyes"`
	Rock  int `json:"rocket"`
}

type Issue struct {
	Id          int64       `json:"id"`
	Number      int         `json:"number"`
	State       IssuesState `json:"state"`
	Title       string      `json:"title"`
	Body        string      `json:"body_text"`
	Labels      []Label
	User        *User
	Association string `json:"author_association"`
	Url         string `json:"html_url"`
	LabelsUrl   string `json:"labels_url"`
	CommentsUrl string `json:"comments_url"`
	Comments    int    `json:"comments"`
	Locked      bool   `json:"locked"`
	Draft       bool   `json:"draft"`
	Reactions   Reactions
}

type Pull struct {
	Id          int64  `json:"id"`
	Number      int    `json:"number"`
	State       string `json:"state"`
	Title       string `json:"title"`
	Body        string `json:"body"`
	Labels      []Label
	User        *User
	Association string  `json:"author_association"`
	Url         string  `json:"html_url"`
	IssueUrl    string  `json:"issue_url"`
	DiffUrl     string  `json:"diff_url"`
	CommitsUrl  string  `json:"commits_url"`
	CommentsUrl string  `json:"comments_url"`
	Comments    int     `json:"comments"`
	Locked      bool    `json:"locked"`
	Draft       bool    `json:"draft"`
	Reviewers   *[]User `json:"requested_reviewers"`
}

type Repo struct {
	// basic information
	Id            int64    `json:"id"`
	Name          string   `json:"name"`
	FullName      string   `json:"full_name"`
	Description   string   `json:"description,omitempty"`
	Language      string   `json:"language,omitempty"`
	Topics        []string `json:"topics,omitempty"`
	DefaultBranch string   `json:"default_branch"`
	Forked        bool     `json:"fork"`
	Owner         Owner
	Issues        []Issue
	Pulls         []Pull

	// urls
	HTMLUrl      string `json:"html_url"`
	IssuesUrl    string `json:"issues_url"`
	PullsUrl     string `json:"pulls_url"`
	DownloadsUrl string `json:"downloads_url"`
	CloneUrl     string `json:"clone_url"`
	SSHUrl       string `json:"ssh_url"`

	// counters
	ForksCount      int   `json:"forks_count"`
	StargazersCount int   `json:"stargazers_count"`
	WatchersCount   int   `json:"watchers_count"`
	OpenIssuesCount int   `json:"open_issues_count"`
	HasIssues       bool  `json:"has_issues"`
	HasProject      bool  `json:"has_projects"`
	HasDisscussions bool  `json:"has_disscussions"`
	HasWiki         bool  `json:"has_wiki"`
	HasPages        bool  `json:"has_pages"`
	HasDownloads    bool  `json:"has_downloads"`
	Archived        bool  `json:"archived"`
	Disabled        bool  `json:"disabled"`
	Size            int64 `json:"size"`

	// permissions & access controls
	Private    bool       `json:"private"`
	Visibility Visibility `json:"visibility"`
	Perms      Permissions
	// license
	License License
}

func RetrieveRepos(token string, args Arguments) ([]Repo, error) {
	repos := make([]Repo, 0)
	var idx int // index of current repo
	ctx := context.Background()
	gh := github.NewTokenClient(ctx, token)
	defer ctx.Done()

	reposUrl := BASE_API_URL + "user/starred"
	limit, limitFlag := 0, args.Limit
	page, perPage := 1, 10
	for limit < limitFlag {
		req, err := http.NewRequest(http.MethodGet, reposUrl, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)

		// modify query params
		q := req.URL.Query()
		q.Set("page", strconv.Itoa(page))
		q.Set("per_page", strconv.Itoa(perPage))
		q.Set("sort", "updated")
		req.URL.RawQuery = q.Encode()

		status, body, err := makeRequest(req)
		if status >= 400 || err != nil {
			return nil, fmt.Errorf("HTTPError: %v", err)
		}
		var newRepos []Repo
		if err := json.Unmarshal(body, &newRepos); err != nil {
			return nil, fmt.Errorf("JSONError: %v", err)
		}
		idx = len(repos)
		for _, newRepo := range newRepos { // repo appention
			repos = appendRepo(repos, newRepo, args)
		}

		for ; idx < len(repos); idx++ {
			if !repos[idx].HasIssues { // skip repository
				continue
			}

			// request repository issues
			if strings.Contains(args.Items, "issue") {
				time.Sleep(time.Second)
				opts := github.ListOptions{Page: 1, PerPage: perPage}
				for {
					time.Sleep(time.Millisecond * 200)
					ctx_ := context.Background()
					issues, resp, err := gh.Issues.ListByRepo(
						ctx_, repos[idx].Owner.Login, repos[idx].Name,
						&github.IssueListByRepoOptions{
							ListOptions: opts,
						})
					defer ctx_.Done()
					if err != nil {
						return repos, err
					}
					for _, s := range issues { // issue conversion
						issue := NewIssue(*s)
						repos[idx].Issues = appendIssue(repos[idx].Issues, issue, args)
					}
					issuesLen := len(repos[idx].Issues)
					if resp.NextPage == 0 || issuesLen >= args.Max { // last page
						if issuesLen > args.Max { // max amount
							repos[idx].Issues = repos[idx].Issues[:args.Max]
						}
						limit += len(repos[idx].Issues)
						break
					}
					opts.Page = resp.NextPage
				}
			}

			// request repository pull requests
			if strings.Contains(args.Items, "pr") {
				time.Sleep(time.Second)
				opts := github.ListOptions{Page: 1, PerPage: perPage}
				for {
					time.Sleep(time.Millisecond * 200)
					ctx_ := context.Background()
					pulls, resp, err := gh.PullRequests.List(
						ctx_, repos[idx].Owner.Login, repos[idx].Name,
						&github.PullRequestListOptions{
							ListOptions: opts,
						})
					defer ctx_.Done()
					if err != nil {
						return repos, err
					}
					for _, p := range pulls { // pull conversion
						pull := NewPull(*p)
						repos[idx].Pulls = appendPull(repos[idx].Pulls, pull, args)
					}
					pullsLen := len(repos[idx].Pulls)
					if resp.NextPage == 0 || pullsLen >= args.Max { // last page
						if pullsLen > args.Max { // max amount
							repos[idx].Pulls = repos[idx].Pulls[:args.Max]
						}
						limit += len(repos[idx].Pulls)
						break
					}
					opts.Page = resp.NextPage
				}
			}
		}

		page++ // next page
	}

	// exclude repos without issues, and pulls
	return repos[:idx], nil
}

// `format`: 'short', 'normal', 'verbose'
func (repo Repo) Display(format Format) {
	//TODO: Wrap output properly!
	var tmpl strings.Builder

	if format == Verbose {
		var repoTmpl strings.Builder
		repoStr := fmt.Sprintf("%s: %v `%v` ", repo.FullName, repo.DefaultBranch, repo.License.Key)
		// TODO: display badge for repository.hasSomething
		repoTmpl.WriteString(repoStr)

		urlStr := fmt.Sprintf("%v  %s\n%v  %s\n%v  %s\n%v  %s\n", WebEmoj, repo.HTMLUrl, CloneEmoj, repo.CloneUrl, SSHEmoj, repo.SSHUrl, DownloadEmoj, repo.DownloadsUrl)
		var topicsStr string
		for _, t := range repo.Topics {
			topicsStr = fmt.Sprintf("%s [%s]", topicsStr, t)
		}
		countStr := fmt.Sprintf("%v %d %v %d %v %d  %v %d %v %d \n", IssueEmoj, repo.OpenIssuesCount, StarEmoj, repo.StargazersCount, EyeEmoj, repo.WatchersCount, ForkEmoj, repo.ForksCount, SizeEmoj, repo.Size)
		var issueTmpl strings.Builder
		for _, issue := range repo.Issues {
			headerStr := fmt.Sprintf("  %s  #%d: %s `%v` \n", IssueEmoj, issue.Number, issue.Title, issue.State)
			labelsStr := "  "
			for _, l := range issue.Labels {
				//TODO: make labels color with bash colors
				labelsStr = fmt.Sprintf("%s [%s]", labelsStr, l.Name)
			}
			urlStr := fmt.Sprintf("\n  %s\n", issue.Url)
			bodyStr := fmt.Sprintf("    %s\n    %s\n", issue.Title, issue.Body)
			reactStr := fmt.Sprintf("  %d  %v %d  %v %d  %v %d \n", issue.Reactions.Total, ThumbEmoj, issue.Reactions.Up, DownEmoj, issue.Reactions.Down, LaughEmoj, issue.Reactions.Laugh)
			issueTmpl.WriteString(headerStr)
			issueTmpl.WriteString(labelsStr)
			issueTmpl.WriteString(urlStr)
			issueTmpl.WriteString(bodyStr)
			issueTmpl.WriteString(reactStr)
		}
		var pullTmpl strings.Builder
		for _, pull := range repo.Pulls {
			headerStr := fmt.Sprintf("  %s  #%d: %s `%v` \n", PrEmoj, pull.Number, pull.Title, pull.State)
			labelsStr := "  "
			for _, l := range pull.Labels {
				//TODO: make labels color with bash colors
				labelsStr = fmt.Sprintf("%s [%s]", labelsStr, l.Name)
			}
			urlStr := fmt.Sprintf("\n  %s\n", pull.Url)
			bodyStr := fmt.Sprintf("    %s\n", pull.Body)
			pullTmpl.WriteString(headerStr)
			issueTmpl.WriteString(labelsStr)
			issueTmpl.WriteString(urlStr)
			issueTmpl.WriteString(bodyStr)
		}

		tmpl.WriteString(repoTmpl.String())
		tmpl.WriteString(topicsStr)
		tmpl.WriteString(countStr)
		tmpl.WriteString(urlStr)
		tmpl.WriteString(issueTmpl.String())
		tmpl.WriteString(pullTmpl.String())
		tmpl.WriteRune('\n') // newline

	} else if format == Short {
		repoStr := fmt.Sprintf("%s/%s: ", repo.Owner.Login, repo.Name)
		var issueTmpl strings.Builder
		issueTmpl.WriteString(fmt.Sprintf("%s ", IssueEmoj))
		for _, issue := range repo.Issues {
			issueStr := fmt.Sprintf("#%d ", issue.Number)
			issueTmpl.WriteString(issueStr)
		}
		var pullTmpl strings.Builder
		pullTmpl.WriteString(fmt.Sprintf("%s ", PrEmoj))
		for _, pull := range repo.Pulls {
			pullStr := fmt.Sprintf("#%d ", pull.Number)
			pullTmpl.WriteString(pullStr)
		}

		tmpl.WriteString(repoStr)
		tmpl.WriteString(issueTmpl.String())
		tmpl.WriteString(pullTmpl.String())
		tmpl.WriteRune('\n') // newline

	} else { // default: `normal`
		repoStr := fmt.Sprintf("%s/%s: ", repo.Owner.Login, repo.Name)
		countStr := fmt.Sprintf("%v %d %v %d %v %d \n", StarEmoj, repo.StargazersCount, EyeEmoj, repo.WatchersCount, ForkEmoj, repo.ForksCount)
		urlStr := fmt.Sprintf("%v  %s\n%v  %s\n", CloneEmoj, repo.CloneUrl, SSHEmoj, repo.SSHUrl)
		var issueTmpl strings.Builder
		for _, issue := range repo.Issues {
			issueStr := fmt.Sprintf("  %s  #%d: %s `%v` \n", IssueEmoj, issue.Number, issue.Title, issue.State)
			issueTmpl.WriteString(issueStr)
		}
		var pullTmpl strings.Builder
		for _, pull := range repo.Pulls {
			pullStr := fmt.Sprintf("  %s  #%d: %s `%v` \n", PrEmoj, pull.Number, pull.Title, pull.State)
			pullTmpl.WriteString(pullStr)
		}

		tmpl.WriteString(repoStr)
		tmpl.WriteString(countStr)
		tmpl.WriteString(urlStr)
		tmpl.WriteString(issueTmpl.String())
		tmpl.WriteString(pullTmpl.String())
	}

	fmt.Print(tmpl.String())
}

func SaveFile(repos []Repo, file string) {
	out, err := os.Create(file)
	if err != nil {
		log.Fatalf("Unable to create file: %v", file)
	}
	buff := make([]byte, 1)
	if buff, err = json.MarshalIndent(repos, "", "  "); err != nil {
		log.Fatalf("Unable to serialize data: %v", err)
	}
	n, err := out.Write(buff)
	if err != nil {
		log.Fatalf("Unable to save file: %v", file)
	}
	log.Printf("%v chars wrote to file %v \n", n, file)
}

func appendRepo(repos []Repo, repo Repo, args Arguments) []Repo {
	var found bool
	for _, topic := range repo.Topics {
		if len(args.Topics) == 0 {
			found = true
			break
		}
		if slices.Contains(args.Topics, topic) {
			found = true
		}
	}
	//TODO Consider implementing patterns;; `FullName`, `Description`
	if found {
		newRepo := repo
		newRepo.Issues = make([]Issue, 0)
		newRepo.Pulls = make([]Pull, 0)
		repos = append(repos, newRepo)
	}
	return repos
}

func appendIssue(issues []Issue, issue Issue, args Arguments) []Issue {
	if len(args.Patterns) == 0 && len(args.Labels) == 0 {
		issues = append(issues, issue)
		return issues // early returns
	}
	//TODO Consider implementing patterns;; `title`, `body`, `comments`
	var found bool
	for _, text := range []string{issue.Title, issue.Body, issue.CommentsUrl} {
		found = findPattern(args.Patterns, text)
		if found {
			break
		}
	}
	for _, label := range args.Labels {
		for _, l := range issue.Labels {
			lbl := removeNonAlpha(l.Name)
			found = strings.Contains(lbl, label)
			if found {
				break
			}
		}
	}
	if found {
		issues = append(issues, issue)
	}
	return issues
}

func appendPull(pulls []Pull, pull Pull, args Arguments) []Pull {
	if len(args.Patterns) == 0 && len(args.Labels) == 0 {
		pulls = append(pulls, pull)
	}
	//TODO Consider implementing patterns;; `title`, `body`, `comments`
	var found bool
	for _, text := range []string{pull.Title, pull.Body, pull.CommentsUrl} {
		found = findPattern(args.Patterns, text)
		if found {
			break
		}
	}
	for _, label := range args.Labels {
		for _, l := range pull.Labels {
			lbl := removeNonAlpha(l.Name)
			found = strings.Contains(lbl, label)
			if found {
				break
			}
		}
	}
	if found {
		pulls = append(pulls, pull)
	}
	return pulls
}

func findPattern(patterns []string, text string) bool {
	var matched bool
	for _, pattern := range patterns {
		textRe := regexp.MustCompile(pattern)
		matched = textRe.MatchString(text)
	}
	return matched
}
