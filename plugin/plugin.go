package plugin

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/google/go-github/github"
	"github.com/urfave/cli"
	"golang.org/x/oauth2"
)

type (
	Plugin struct {
		BaseURL     string
		BuildStatus string
		CommitSHA   string
		Contexts    []string
		Description string
		Password    string
		RepoName    string
		RepoOwner   string
		Username    string
		State       string
		TargetURL   string
		Token       string

		gitClient  *github.Client
		gitContext context.Context
	}
)

// NewFromCLI creates new Plugin instance from CLI flags
func NewFromCLI(c *cli.Context) (*Plugin, error) {
	p := Plugin{
		BaseURL:     c.String("base-url"),
		BuildStatus: c.String("build-status"),
		CommitSHA:   c.String("commit-sha"),
		Contexts:    c.StringSlice("context"),
		Description: c.String("description"),
		Password:    c.String("password"),
		RepoName:    c.String("repo-name"),
		RepoOwner:   c.String("repo-owner"),
		State:       c.String("state"),
		TargetURL:   c.String("target-url"),
		Token:       c.String("api-key"),
		Username:    c.String("username"),
	}

	err := p.init()

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// NewFromPlugin allows someone to populate Plugin struct and then initializes it
func NewFromPlugin(p Plugin) (*Plugin, error) {
	err := p.init()

	if err != nil {
		return nil, err
	}

	return &p, nil
}

// Exec executes the plugin
func (p Plugin) Exec() error {
	if p.gitClient == nil {
		return fmt.Errorf("Exec(): git client not initialized")
	}

	state := p.CalculatedState()

	status := &github.RepoStatus{
		Description: &p.Description,
		State:       &state,
		TargetURL:   &p.TargetURL,
	}

	for _, c := range p.Contexts {
		status.Context = &c

		logrus.WithFields(logrus.Fields{
			"build-status":     p.BuildStatus,
			"calculated-state": state,
			"context":          c,
			"description":      p.Description,
			"repo-name":        p.RepoName,
			"repo-owner":       p.RepoOwner,
			"sha":              p.CommitSHA,
			"state":            p.State,
			"target-url":       p.TargetURL,
		}).Debug("creating status")
		_, _, err := p.gitClient.Repositories.CreateStatus(p.gitContext, p.RepoOwner, p.RepoName, p.CommitSHA, status)

		if err != nil {
			return err
		}
	}

	return nil
}

// CalculatedState to determine state for the status
func (p Plugin) CalculatedState() string {
	if p.State != "" {
		switch p.State {
		case "success", "error", "failure", "pending":
			return p.State
		default:
			return "error"
		}
	}

	if p.BuildStatus == "" {
		return "error"
	}

	return p.BuildStatus
}

func (p *Plugin) init() error {
	err := p.validate()

	if err != nil {
		return err
	}

	err = p.initGitClient()

	if err != nil {
		return err
	}

	return nil
}

func (p *Plugin) initGitClient() error {
	if !strings.HasSuffix(p.BaseURL, "/") {
		p.BaseURL = p.BaseURL + "/"
	}

	baseURL, err := url.Parse(p.BaseURL)

	if err != nil {
		return fmt.Errorf("Failed to parse base URL. %s", err)
	}

	p.gitContext = context.Background()

	if p.Token != "" {
		ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: p.Token})
		tc := oauth2.NewClient(p.gitContext, ts)
		p.gitClient = github.NewClient(tc)
	} else {
		tp := github.BasicAuthTransport{
			Username: strings.TrimSpace(p.Username),
			Password: strings.TrimSpace(p.Password),
		}
		p.gitClient = github.NewClient(tp.Client())
	}
	p.gitClient.BaseURL = baseURL

	return nil
}

func (p Plugin) validate() error {
	if p.Token == "" && (p.Username == "" || p.Password == "") {
		return fmt.Errorf("You must provide an API key or Username and Password")
	}

	return nil
}
