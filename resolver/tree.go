package resolver

import (
	"context"

	"github.com/gitpods/gitpods/repository"
)

type TreeResolver struct {
	repositories repository.Service
}

func NewTree(rs repository.Service) *TreeResolver {
	return &TreeResolver{
		repositories: rs,
	}
}

type treeArgs struct {
	Owner string
	Name  string
}

type treeObjectResolver struct {
	mode   string
	typ    string
	object string
	file   string
	commit *commitResolver
}

func (r *TreeResolver) Tree(ctx context.Context, args treeArgs) ([]*treeObjectResolver, error) {
	objects, err := r.repositories.Tree(ctx, &repository.Owner{Username: args.Owner}, args.Name)
	if err != nil {
		return nil, err
	}

	var tor []*treeObjectResolver
	for _, obj := range objects {
		tor = append(tor, &treeObjectResolver{
			mode:   obj.Mode,
			typ:    obj.Type,
			object: obj.Object,
			file:   obj.File,
			commit: &commitResolver{
				hash:           obj.Commit.Hash,
				tree:           obj.Commit.Tree,
				parent:         obj.Commit.Parent,
				subject:        obj.Commit.Subject,
				author:         obj.Commit.Author,
				authorEmail:    obj.Commit.AuthorEmail,
				authorDate:     int32(obj.Commit.AuthorDate.Unix()),
				committer:      obj.Commit.Committer,
				committerEmail: obj.Commit.CommitterEmail,
				committerDate:  int32(obj.Commit.CommitterDate.Unix()),
			},
		})
	}

	return tor, nil
}

func (r *treeObjectResolver) Mode() string {
	return r.mode
}

func (r *treeObjectResolver) Type() string {
	return r.typ
}

func (r *treeObjectResolver) Object() string {
	return r.object
}

func (r *treeObjectResolver) File() string {
	return r.file
}

func (r *treeObjectResolver) Commit() *commitResolver {
	return &commitResolver{
		hash:    r.commit.hash,
		tree:    r.commit.tree,
		parent:  r.commit.parent,
		subject: r.commit.subject,

		author:      r.commit.author,
		authorEmail: r.commit.authorEmail,
		authorDate:  r.commit.authorDate,

		committer:      r.commit.committer,
		committerEmail: r.commit.committerEmail,
		committerDate:  r.commit.committerDate,
	}
}

type commitResolver struct {
	hash    string
	tree    string
	parent  string
	subject string

	author      string
	authorEmail string
	authorDate  int32

	committer      string
	committerEmail string
	committerDate  int32
}

func (r *commitResolver) Hash() string {
	return r.hash
}

func (r *commitResolver) Tree() string {
	return r.tree
}

func (r *commitResolver) Parent() string {
	return r.parent
}

func (r *commitResolver) Subject() string {
	return r.subject
}

func (r *commitResolver) Author() *commitAuthorResolver {
	return &commitAuthorResolver{
		name:  r.author,
		email: r.authorEmail,
		date:  r.authorDate,
	}
}

type commitAuthorResolver struct {
	name  string
	email string
	date  int32
}

func (r *commitAuthorResolver) Name() string {
	return r.name
}
func (r *commitAuthorResolver) Email() string {
	return r.email
}
func (r *commitAuthorResolver) Date() int32 {
	return r.date
}

func (r *commitResolver) Committer() *commitCommitterResolver {
	return &commitCommitterResolver{
		name:  r.committer,
		email: r.committerEmail,
		date:  r.committerDate,
	}
}

type commitCommitterResolver struct {
	name  string
	email string
	date  int32
}

func (r *commitCommitterResolver) Name() string {
	return r.name
}
func (r *commitCommitterResolver) Email() string {
	return r.email
}
func (r *commitCommitterResolver) Date() int32 {
	return r.date
}
