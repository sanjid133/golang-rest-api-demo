package repo

type Repo interface {
	Name() string
	EnsureIndices() error
	DropIndices() error
}

const (
	repoUser = "User"
)

const (
	msgIsRequired = "is required"
	msgIsInvalid = "is invalid"
)