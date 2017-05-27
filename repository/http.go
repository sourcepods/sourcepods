package repository

import (
	"fmt"
	"net/http"
	"time"

	"github.com/google/jsonapi"
	"github.com/pressly/chi"
)

type response struct {
	ID            string    `jsonapi:"primary,repositories"`
	Name          string    `jsonapi:"attr,name"`
	Description   string    `jsonapi:"attr,description"`
	Website       string    `jsonapi:"attr,website"`
	DefaultBranch string    `jsonapi:"attr,default_branch"`
	Private       bool      `jsonapi:"attr,private"`
	Bare          bool      `jsonapi:"attr,bare"`
	Created       time.Time `jsonapi:"attr,created_at"`
	Updated       time.Time `jsonapi:"attr,updated_at"`

	Stars            int                       `jsonapi:"attr,stars"`
	Forks            int                       `jsonapi:"attr,forks"`
	IssueStats       *responseIssueStats       `jsonapi:"attr,issue_stats"`
	PullRequestStats *responsePullRequestStats `jsonapi:"attr,pull_request_stats"`

	Owner *ResponseOwner `jsonapi:"relation,owner"`
}

type responseIssueStats struct {
	TotalCount  int `json:"total_count"`
	OpenCount   int `json:"open_count"`
	ClosedCount int `json:"closed_count"`
}

type responsePullRequestStats struct {
	TotalCount  int `json:"total_count"`
	OpenCount   int `json:"open_count"`
	ClosedCount int `json:"closed_count"`
}

type ResponseOwner struct {
	ID string `jsonapi:"primary,user"`
}

// NewUsersHandler returns a RESTful http router interacting with the Service.
func NewUsersHandler(s Service) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", listByOwner(s))

	return r
}

func listByOwner(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username := chi.URLParam(r, "username")

		repositories, stats, owner, err := s.ListByOwnerUsername(username)
		if err != nil {
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  http.StatusText(http.StatusNotFound),
				Detail: "Can't find repositories for this owner",
				Status: fmt.Sprintf("%d", http.StatusNotFound),
			}})
			return
		}

		resRepos := make([]interface{}, len(repositories))
		for i, repo := range repositories {
			resRepos[i] = &response{
				ID:            repo.ID,
				Name:          repo.Name,
				Description:   repo.Description,
				Website:       repo.Website,
				DefaultBranch: repo.DefaultBranch,
				Private:       repo.Private,
				Bare:          repo.Bare,
				Created:       repo.Created,
				Updated:       repo.Updated,

				Stars: stats[i].Stars,
				Forks: stats[i].Forks,
				IssueStats: &responseIssueStats{
					TotalCount:  stats[i].IssueTotalCount,
					OpenCount:   stats[i].IssueOpenCount,
					ClosedCount: stats[i].IssueClosedCount,
				},
				PullRequestStats: &responsePullRequestStats{
					TotalCount:  stats[i].PullRequestTotalCount,
					OpenCount:   stats[i].PullRequestOpenCount,
					ClosedCount: stats[i].PullRequestClosedCount,
				},

				Owner: &ResponseOwner{ID: owner.ID},
			}
		}

		w.Header().Set("Content-Type", jsonapi.MediaType)
		if err := jsonapi.MarshalManyPayload(w, resRepos); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func NewHandler(s Service) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/", get(s))

	return r
}

func get(s Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ownerUsername := chi.URLParam(r, "owner")
		name := chi.URLParam(r, "name")

		repository, stats, owner, err := s.Find(ownerUsername, name)
		if err != nil {
			jsonapi.MarshalErrors(w, []*jsonapi.ErrorObject{{
				Title:  http.StatusText(http.StatusNotFound),
				Detail: "Repository with this owner and name does not exist",
				Status: fmt.Sprintf("%d", http.StatusNotFound),
			}})
			return
		}

		res := &response{
			ID:            repository.ID,
			Name:          repository.Name,
			Description:   repository.Description,
			Website:       repository.Website,
			DefaultBranch: repository.DefaultBranch,
			Private:       repository.Private,
			Bare:          repository.Bare,
			Created:       repository.Created,
			Updated:       repository.Updated,

			Stars: stats.Stars,
			Forks: stats.Forks,

			IssueStats: &responseIssueStats{
				TotalCount:  stats.IssueTotalCount,
				OpenCount:   stats.IssueOpenCount,
				ClosedCount: stats.IssueClosedCount,
			},

			PullRequestStats: &responsePullRequestStats{
				TotalCount:  stats.PullRequestTotalCount,
				OpenCount:   stats.PullRequestOpenCount,
				ClosedCount: stats.PullRequestClosedCount,
			},

			Owner: &ResponseOwner{ID: owner.ID},
		}

		w.Header().Set("Content-Type", jsonapi.MediaType)
		jsonapi.MarshalOnePayload(w, res)
	}
}
