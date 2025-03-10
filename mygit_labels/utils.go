package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/google/go-github/v69/github"
)

type Format string
type Emoticon string

const (
	MAX_BYTES   = 128
	access_key  = "Access"
	refresh_key = "Refresh"

	//format
	Short   Format = "short"
	Normal  Format = "normal"
	Verbose Format = "verbose"

	// emoticon
	IssueEmoj    Emoticon = "⚠️"
	PrEmoj       Emoticon = "🔀"
	BranchEmoj   Emoticon = "🪵"
	ForkEmoj     Emoticon = "Ψ"
	LangEmoj     Emoticon = "💻"
	StarEmoj     Emoticon = "⭐"
	ThumbEmoj    Emoticon = "👍"
	DownEmoj     Emoticon = "👎"
	LaughEmoj    Emoticon = "😆"
	ConfuseEmoj  Emoticon = "😳"
	HeartEmoj    Emoticon = "❤️"
	EyeEmoj      Emoticon = "🧐"
	CloneEmoj    Emoticon = "🐑"
	SSHEmoj      Emoticon = "🟩"
	WebEmoj      Emoticon = "🌐"
	DownloadEmoj Emoticon = "⬇️"
	SizeEmoj     Emoticon = "💾"
)

// TODO: validate this with validators
type Arguments struct {
	Items    string   // `issues` & `prs`
	Patterns []string // `[\w.$]+`
	Labels   []string // `[\w]+`
	Topics   []string // `[\w]+`
	Limit    int      // <1024
	Max      int      // <100
	// Sort     string   // :`name` | `star` | `fork`
	// Shuffle  bool     // !~sort
	Output string // :`` | `file.ext`
	Format Format // `^short` | `^normal` | `^verbose`
}

func GetTokens(file_path string) (string, string, error) {
	file, err := os.Open(file_path)
	if err != nil {
		return "", "", err
	}
	buf := make([]byte, MAX_BYTES)
	_, err = file.Read(buf)
	if err != nil {
		return "", "", err
	}
	tokenMap := make(map[string]string)
	data := string(buf)
	tokens := strings.Split(data, "\n")
	for _, token := range tokens {
		KV := strings.Split(token, ":")
		if len(KV) != 2 { // !EOF
			break
		}
		tKey := strings.TrimSpace(KV[0])
		tVal := strings.TrimSpace(KV[1])
		tokenMap[tKey] = tVal
	}
	aToken := tokenMap[access_key]
	rToken := tokenMap[refresh_key]
	return aToken, rToken, nil
}

func SaveTokens(file_path, aToken, rToken string) error {
	file, err := os.Create(file_path)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	buf.Write([]byte(fmt.Sprintf("%s: %s\n", access_key, aToken)))
	buf.Write([]byte(fmt.Sprintf("%s: %s\n", refresh_key, rToken)))
	_, err = file.Write(buf.Bytes())
	if err != nil {
		return err
	}
	return nil
}

func ParseArguments(args *Arguments) *Arguments {
	args = loadOptions(args)
	newargs := new(Arguments)
	// parse arguments from command line
	flag.StringVar(&args.Items, "Items", args.Items, "shows issues, or prs")
	flag.StringVar(&newargs.Patterns[0], "pattern", strings.Join(args.Patterns, ","), "search through title, description, and comments Like RegEx")
	flag.StringVar(&newargs.Labels[0], "labels", strings.Join(args.Labels, ","), "finds labels where part of that")
	flag.StringVar(&newargs.Topics[0], "topics", strings.Join(args.Topics, ","), "find repos relevent to topic")
	flag.IntVar(&newargs.Limit, "limit", args.Limit, "limit number of issues, and prs")
	flag.IntVar(&newargs.Max, "max", args.Max, "maximum number of issues, or prs per repo")
	flag.StringVar(&newargs.Output, "output", args.Output, "the file to save data")
	var maskFormat int8 = 1
	var short, normal, verbose bool
	flag.BoolVar(&short, "short", false, "a compact representation of data")
	flag.BoolVar(&normal, "normal", true, "a normal representation of data")
	flag.BoolVar(&verbose, "verbose", false, "a verbose representation of")
	flag.Parse()
	for _, b := range []bool{verbose, normal, short} {
		maskFormat <<= 1
		if b {
			maskFormat++
		}
	}
	args = validateArgs(args, newargs, maskFormat)

	saveOption(args)
	return args
}

func loadOptions(args *Arguments) *Arguments {
	file, err := os.Open(OPTIONS_FILE)
	if err != nil {
		log.Printf("Unable to open options file [%s]: %v\n", OPTIONS_FILE, err)
		return args
	}
	defer file.Close()
	buf := make([]byte, 1000)
	nl, err := file.Read(buf)
	if err != nil {
		log.Printf("Unable to read file [%v]: %v\n", OPTIONS_FILE, err)
		return args
	}
	newargs := new(Arguments)
	if err = yaml.Unmarshal(buf[:nl], newargs); err != nil {
		log.Fatalf("Unable to deserialize file [%v]: %v\n", OPTIONS_FILE, err)
	}
	args = validateArgs(args, newargs, 0)

	return args
}

func saveOption(args *Arguments) {
	file, err := os.Create(OPTIONS_FILE)
	if err != nil {
		log.Printf("Unable to create options file [%s]: %v\n", OPTIONS_FILE, err)
		return
	}
	buf, err := yaml.Marshal(args)
	if err != nil {
		log.Fatalf("Unable to serialize file [%v]: %v\n", OPTIONS_FILE, err)
	}
	_, err = file.Write(buf)
	if err != nil {
		log.Printf("Unable to save file [%v]: %v\n", OPTIONS_FILE, err)
	}
}

func validateArgs(args *Arguments, newargs *Arguments, formatMask int8) *Arguments {
	var short, normal, verbose int8
	short = 9    // 1001
	normal = 10  // 1010
	verbose = 12 // 1100
	var found bool
	for _, v := range []int8{0, verbose, normal, short} {
		if formatMask == v {
			found = true
			if formatMask == verbose {
				args.Format = Verbose
			} else if formatMask == short {
				args.Format = Short
			} else {
				args.Format = Normal
			}
			break
		}
	}
	if !found {
		log.Panicln("Fault: you should choose -short, -normal, or -verbose")
	}

	args.Limit = min(newargs.Limit, 1024)
	args.Max = min(newargs.Max, 100)
	patternsStr := strings.Join(newargs.Patterns, ",")
	patternArr := strings.Split(patternsStr, ",")
	//TODO support expressions; . 1-char, $ n-char
	for _, p := range patternArr {
		p = strings.TrimSpace(p)
		if p != "" {
			args.Patterns = append(args.Patterns, p)
		}
	}
	LabelsStr := strings.Join(newargs.Labels, ",")
	LabelsArr := strings.Split(LabelsStr, ",")
	for _, p := range LabelsArr {
		p = strings.TrimSpace(p)
		if p != "" {
			args.Labels = append(args.Labels, p)
		}
	}
	topicsStr := strings.Join(newargs.Topics, ",")
	topicsArr := strings.Split(topicsStr, ",")
	for _, p := range topicsArr {
		p = removeNonAlpha(p)
		if p != "" {
			args.Topics = append(args.Topics, p)
		}
	}

	return args
}

func removeNonAlpha(s string) string {
	var str strings.Builder
	for _, c := range s {
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			str.WriteRune(c)
		}
	}
	return str.String()
}

func NewIssue(s github.Issue) Issue {
	labels := make([]Label, 0)
	for _, rLabel := range s.Labels {
		label := &Label{
			Id:    *rLabel.ID,
			Name:  *rLabel.Name,
			Color: *rLabel.Color,
		}
		labels = append(labels, *label)
	}
	obj := Issue{
		Id:     *s.ID,
		Number: *s.Number,
		State:  IssuesState(*s.State),
		Labels: labels,
		User: &User{
			Id:    *s.ID,
			Login: *s.User.Login,
			// Name:  *s.User.Name,
			// Email: *s.User.Email,
			Url: *s.URL,
		},
		// Association: *s.AuthorAssociation,
		Url:         *s.URL,
		LabelsUrl:   *s.LabelsURL,
		CommentsUrl: *s.CommentsURL,
		Comments:    *s.Comments,
		//Locked:      *s.Locked,
		//Draft:       *s.Draft,
		Reactions: Reactions{
			Total: *s.Reactions.TotalCount,
			Up:    *s.Reactions.PlusOne,
			Down:  *s.Reactions.MinusOne,
			Laugh: *s.Reactions.Laugh,
			Conf:  *s.Reactions.Confused,
			Heart: *s.Reactions.Heart,
			Hoop:  *s.Reactions.Hooray,
			Eyes:  *s.Reactions.Eyes,
			Rock:  *s.Reactions.Rocket,
		},
	}
	if s.Title != nil {
		obj.Title = *s.Title
	}
	if s.Body != nil {
		obj.Body = *s.Body
	}
	return obj
}

func NewPull(s github.PullRequest) Pull {
	labels := make([]Label, 0)
	for _, rLabel := range s.Labels {
		label := Label{
			Id:    *rLabel.ID,
			Name:  *rLabel.Name,
			Color: *rLabel.Color,
		}
		labels = append(labels, label)
	}
	reviewrs := make([]User, 0)
	for _, r := range s.RequestedReviewers {
		user := User{
			Id:    *r.ID,
			Login: *r.Login,
			// Name:  *r.Name,
			// Email: *r.Email,
			Url: *r.URL,
		}
		reviewrs = append(reviewrs, user)
	}
	obj := Pull{
		Id:     *s.ID,
		Number: *s.Number,
		State:  *s.State,
		Labels: labels,
		User: &User{
			Id:    *s.ID,
			Login: *s.User.Login,
			// Name:  *s.User.Name,
			// Email: *s.User.Email,
			Url: *s.URL,
		},
		// Association: *s.AuthorAssociation,
		Url:         *s.URL,
		IssueUrl:    *s.IssueURL,
		DiffUrl:     *s.DiffURL,
		CommitsUrl:  *s.CommitsURL,
		CommentsUrl: *s.CommentsURL,
		//Locked:      *s.Locked,
		//Draft:       *s.Draft,
		// Reviewers:   &reviewrs,
	}
	if s.Title != nil {
		obj.Title = *s.Title
	}
	if s.Body != nil {
		obj.Body = *s.Body
	}
	if s.Comments != nil {
		obj.Comments = *s.Comments
	}
	return obj
}
