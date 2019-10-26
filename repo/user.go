package repo

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sanjid133/rest-user-store/model"
	"github.com/sanjid133/rest-user-store/util"
	"time"


	"github.com/sanjid133/rest-user-store/database"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

// User interface represents the user repository
type User interface {
	Repo

	Add(v *model.User) error
	Get(id string) (*model.User, error)
	//Update(id string, v *model.User) error
}

// bsonUser holds bson User
type bsonUser struct {
	ID               primitive.ObjectID `bson:"_id"`
	FirstName string `bson:"first_name"`
	LastName string `bson:"last_name"`
	Password string `bson:"password"`
	//Tags []Tag

	CreatedAt        time.Time          `bson:"created_at"`
	UpdatedAt        time.Time          `bson:"updated_at"`
}

// Valid checks if a bsonUser is valid
func (u *bsonUser) valid() (bool, error) {
	return true, nil
}

// prepareBsonUser prepare a bson user from model user
func prepareBsonUser(u *model.User) (*bsonUser, error) {
	usr := bsonUser{
		FirstName: u.FirstName,
		LastName: u.LastName,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}
	if u.ID != "" {
		id, err := primitive.ObjectIDFromHex(u.ID)
		if err != nil {
			return nil, errors.Errorf("Invalied ID %v", u.ID)
		}
		usr.ID = id
	}
	return &usr, nil
}

// formUser forms a model user from bson user
func formUser(u *bsonUser) *model.User {
	return &model.User{
		ID:               u.ID.Hex(),
		FirstName: u.FirstName,
		LastName: u.LastName,
		CreatedAt:        u.CreatedAt,
		UpdatedAt:        u.UpdatedAt,
	}
}

// MgoUser implements User repository with mongodb
type MgoUser struct {
	db    *database.Client
	table string
}

// NewMgoUser takes the collection name for user repository and a mongo db client
// returns a MgoUser repository
func NewMgoUser(col string, db *database.Client) *MgoUser {
	return &MgoUser{
		db:    db,
		table: col,
	}
}

// Name of the repository "User"
func (r *MgoUser) Name() string {
	return repoUser
}

/*// Indices returns mongo database indexes for MgoUser
func (r *MgoUser) Indices() []mongo.Index {
	return []mongo.Index{}
}
*/
// EnsureIndices ensures indices in database
func (r *MgoUser) EnsureIndices() error {
	return nil
}

// DropIndices drops indices from database
func (r *MgoUser) DropIndices() error {
	return nil
}

// Add adds a new user to the repository
func (r *MgoUser) Add(v *model.User) error {
	bUser, err := prepareBsonUser(v)
	if err != nil {
		return err
	}

	now := util.Now()

	if bUser.ID.IsZero() {
		bUser.ID = primitive.NewObjectID()
	}
	bUser.CreatedAt = now
	bUser.UpdatedAt = now

	if ok, err := bUser.valid(); !ok {
		return errors.Errorf(repoUser, err)
	}

	if _, err := r.db.Insert(context.Background(),r.table, bUser); err != nil {
		return err
	}

	*v = *formUser(bUser)

	return nil
}

// Get returns a user from the repository by its id
// if no user found it returns nil, nil
func (r *MgoUser) Get(id string) (*model.User, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, errors.Errorf("Invalid id", id)
	}
	row, err := r.db.FindID(context.Background(), r.table, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	usr := bsonUser{}

	if row.Next() {
		if err := row.Scan(&usr); err != nil {
			return nil, err
		}
	}

	if err := row.Err(); err != nil {
		return nil, errors.PrepareDocumentError(err)
	}

	return formUser(&usr), nil
}
