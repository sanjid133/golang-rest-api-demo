package repo

type Repo interface {
	Name() string
}

const (
	repoUser = "User"
	repoTag  = "Tag"
)

const (
	msgIsRequired = "is required"
	msgIsInvalid  = "is invalid"
)
