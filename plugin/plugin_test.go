package plugin

import (
	"fmt"
	"testing"

	"github.com/franela/goblin"
	"gopkg.in/h2non/gock.v1"
)

func TestPlugin(t *testing.T) {
	g := goblin.Goblin(t)

	g.Describe("NewFromPlugin", func() {
		pl := Plugin{
			BaseURL:     "http://server.com",
			BuildStatus: "failure",
			CommitSHA:   "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			Contexts:    []string{"some/context", "another/context"},
			Description: "cool descrip",
			RepoName:    "test-repo",
			RepoOwner:   "test-owner",
			State:       "success",
			TargetURL:   "https://link.to.com/build",
			Token:       "fake",
		}

		g.It("creates a new/initialized Plugin from a Plugin", func() {
			_, err := NewFromPlugin(pl)
			g.Assert(err == nil).IsTrue(fmt.Sprintf("Received err: %s", err))
		})

		g.It("errors when defined directly", func() {
			defer gock.Off()

			gock.New("http://server.com").
				Post("/repos/test-owner/test-repo/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e").
				Reply(201).
				JSON(map[string]string{})

			defer func() {
				r := recover()
				if r != nil {
					g.Fail("The code should not panic")
				}
			}()

			err := pl.Exec()
			g.Assert(err != nil).IsTrue("should have received error that git client not initialized")
		})
	})

	g.Describe("CalculatedState()", func() {
		g.It("returns valid user defined state", func() {
			pl := Plugin{
				BuildStatus: "failure",
				State:       "success",
			}

			g.Assert(pl.CalculatedState()).Equal("success")

			pl = Plugin{
				BuildStatus: "success",
				State:       "failure",
			}

			g.Assert(pl.CalculatedState()).Equal("failure")
		})

		g.It("returns error for invalid user defined state", func() {
			pl := Plugin{
				State: "asdf",
			}

			g.Assert(pl.CalculatedState()).Equal("error")
		})

		g.It("returns drone build status if no user provided state", func() {
			pl := Plugin{
				BuildStatus: "success",
			}

			g.Assert(pl.CalculatedState()).Equal("success")

			pl = Plugin{
				BuildStatus: "failure",
			}

			g.Assert(pl.CalculatedState()).Equal("failure")
		})

		g.It("returns error if no data is provided", func() {
			pl := Plugin{}

			g.Assert(pl.CalculatedState()).Equal("error")
		})
	})

	g.Describe("Exec()", func() {
		pl := Plugin{
			BaseURL:     "http://server.com",
			BuildStatus: "failure",
			CommitSHA:   "6dcb09b5b57875f334f61aebed695e2e4193db5e",
			Contexts:    []string{"some/context", "another/context"},
			Description: "cool descrip",
			RepoName:    "test-repo",
			RepoOwner:   "test-owner",
			State:       "success",
			TargetURL:   "https://link.to.com/build",
			Token:       "fake",
		}

		g.It("creates a status (single)", func() {
			defer gock.Off()

			pl.Contexts = []string{"some/context"}
			p, err := NewFromPlugin(pl)
			if err != nil {
				g.Fail("Failed to create plugin for testing")
			}

			gock.New("http://server.com").
				Post("/repos/test-owner/test-repo/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e").
				File("../testdata/request/match-first-context.json").
				Reply(201).
				JSON(map[string]string{})

			err = p.Exec()

			g.Assert(err == nil).IsTrue(fmt.Sprintf("Received err: %s", err))
			g.Assert(gock.HasUnmatchedRequest()).IsFalse(fmt.Sprintf("Received unmatched requests: %v\n", gock.GetUnmatchedRequests()))
			if !gock.IsDone() {
				for _, m := range gock.Pending() {
					g.Fail(fmt.Sprintf("Did not make expected request: %s(%s)", m.Request().Method, m.Request().URLStruct))
				}
			}
		})

		g.It("creates statuses (multiple)", func() {
			defer gock.Off()

			pl.Contexts = []string{"some/context", "another/context"}
			p, err := NewFromPlugin(pl)
			if err != nil {
				g.Fail("Failed to create plugin for testing")
			}

			gock.New("http://server.com").
				Post("/repos/test-owner/test-repo/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e").
				File("../testdata/request/match-first-context.json").
				Reply(201).
				JSON(map[string]string{})

			gock.New("http://server.com").
				Post("/repos/test-owner/test-repo/statuses/6dcb09b5b57875f334f61aebed695e2e4193db5e").
				File("../testdata/request/match-second-context.json").
				Reply(201).
				JSON(map[string]string{})

			err = p.Exec()

			g.Assert(err == nil).IsTrue(fmt.Sprintf("Received err: %s", err))
			g.Assert(gock.HasUnmatchedRequest()).IsFalse(fmt.Sprintf("Received unmatched requests: %v\n", gock.GetUnmatchedRequests()))
			if !gock.IsDone() {
				for _, m := range gock.Pending() {
					g.Fail(fmt.Sprintf("Did not make expected request: %s(%s): %s", m.Request().Method, m.Request().URLStruct, string(m.Request().BodyBuffer)))
				}
			}
		})
	})
}
