package repo

import (
	"context"
	"github.com/pkg/errors"
	"github.com/sanjid133/rest-user-store/model"
	"github.com/sanjid133/rest-user-store/util"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"time"

	"github.com/sanjid133/rest-user-store/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User interface represents the user repository
type User interface {
	Repo

	Add(v *model.User) error
	Get(id string) (*model.User, error)
	ListUsers(IDs []string) ([]model.User, error)
	//Update(id string, v *model.User) error
}

// bsonUser holds bson User
type bsonUser struct {
	ID        primitive.ObjectID `bson:"_id"`
	FirstName string             `bson:"first_name"`
	LastName  string             `bson:"last_name"`
	PassHash  []byte             `bson:"pass_hash"`
	PassSalt  []byte             `bson:"pass_salt"`
	PassIter  int                `bson:"pass_iter"`
	PassAlgo  string             `bson:"pass_algo"`
	//Tags []Tag

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// Valid checks if a bsonUser is valid
func (u *bsonUser) valid() (bool, error) {
	return true, nil
}

// prepareBsonUser prepare a bson user from model user
func prepareBsonUser(u *model.User) (*bsonUser, error) {
	usr := bsonUser{
		FirstName: u.FirstName,
		LastName:  u.LastName,
		PassHash:  u.PassHash,
		PassSalt:  u.PassSalt,
		PassIter:  u.PassIter,
		PassAlgo:  u.PassAlgo,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
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
		ID:        u.ID.Hex(),
		FirstName: u.FirstName,
		LastName:  u.LastName,
		PassHash:  u.PassHash,
		PassSalt:  u.PassSalt,
		PassIter:  u.PassIter,
		PassAlgo:  u.PassAlgo,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
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

	if _, err := r.db.Insert(context.Background(), r.table, bUser); err != nil {
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
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return formUser(&usr), nil
}

func (r *MgoUser) ListUsers(IDs []string) ([]model.User, error) {
	usrIds := []primitive.ObjectID{}
	for _, id := range IDs {
		objID, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, errors.Errorf("Invalid id", id)
		}
		usrIds = append(usrIds, objID)
	}

	query := bson.M{
		"_id": bson.M{
			"$in": usrIds,
		},
	}
	rows, err := r.db.Find(context.Background(), r.table, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	respUsers := []model.User{}
	for rows.Next() {
		bUser := bsonUser{}
		if err := rows.Scan(&bUser); err != nil {
			return nil, err
		}
		respUsers = append(respUsers, *formUser(&bUser))

	}
	if rows.Err(); err != nil {
		return nil, err
	}
	return respUsers, nil

}
