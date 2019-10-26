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

// Tag interface represents the tag repository
type Tag interface {
	Repo

	Add(v *model.Tag) error
	Get(id string) (*model.Tag, error)
	ListTags(tags []string) ([]model.Tag, error)
	ListByUserID(uid string) ([]model.Tag, error)
	//Update(id string, v *model.Tag) error
}

// bsonTag holds bson Tag
type bsonTag struct {
	ID       primitive.ObjectID `bson:"_id"`
	Tag      string             `bson:"tag"`
	UserID   string             `bson:"user_id"`
	ExpireAt time.Time          `bson:"expire_at"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// Valid checks if a bsonTag is valid
func (u *bsonTag) valid() (bool, error) {
	return true, nil
}

// prepareBsonTag prepare a bson tag from model tag
func prepareBsonTag(u *model.Tag) (*bsonTag, error) {
	tag := bsonTag{
		Tag:       u.Tag,
		UserID:    u.UserID,
		ExpireAt:  u.ExpireAt,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
	if u.ID != "" {
		id, err := primitive.ObjectIDFromHex(u.ID)
		if err != nil {
			return nil, errors.Errorf("Invalied ID %v", u.ID)
		}
		tag.ID = id
	}
	return &tag, nil
}

// formTag forms a model tag from bson tag
func formTag(u *bsonTag) *model.Tag {
	return &model.Tag{
		ID:        u.ID.Hex(),
		Tag:       u.Tag,
		UserID:    u.UserID,
		ExpireAt:  u.ExpireAt,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// MgoTag implements Tag repository with mongodb
type MgoTag struct {
	db    *database.Client
	table string
}

// NewMgoTag takes the collection name for tag repository and a mongo db client
// returns a MgoTag repository
func NewMgoTag(col string, db *database.Client) *MgoTag {
	return &MgoTag{
		db:    db,
		table: col,
	}
}

// Name of the repository "Tag"
func (r *MgoTag) Name() string {
	return repoTag
}

// Add adds a new tag to the repository
func (r *MgoTag) Add(v *model.Tag) error {
	bTag, err := prepareBsonTag(v)
	if err != nil {
		return err
	}

	now := util.Now()

	if bTag.ID.IsZero() {
		bTag.ID = primitive.NewObjectID()
	}
	bTag.CreatedAt = now
	bTag.UpdatedAt = now

	if ok, err := bTag.valid(); !ok {
		return errors.Errorf(repoTag, err)
	}

	if _, err := r.db.Insert(context.Background(), r.table, bTag); err != nil {
		return err
	}

	*v = *formTag(bTag)

	return nil
}

// Get returns a tag from the repository by its id
// if no tag found it returns nil, nil
func (r *MgoTag) Get(id string) (*model.Tag, error) {
	if _, err := primitive.ObjectIDFromHex(id); err != nil {
		return nil, errors.Errorf("Invalid id", id)
	}
	row, err := r.db.FindID(context.Background(), r.table, id)
	if err != nil {
		return nil, err
	}
	defer row.Close()

	tag := bsonTag{}

	if row.Next() {
		if err := row.Scan(&tag); err != nil {
			return nil, err
		}
	}

	if err := row.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return formTag(&tag), nil
}

func (r *MgoTag) ListTags(tags []string) ([]model.Tag, error) {
	query := bson.M{
		"tag": bson.M{
			"$in": tags,
		},
		"expire_at": bson.M{
			"$gte": util.Now().UTC(),
		},
	}
	rows, err := r.db.Find(context.Background(), r.table, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	respTags := []model.Tag{}
	for rows.Next() {
		bTag := bsonTag{}
		if err := rows.Scan(&bTag); err != nil {
			return nil, err
		}
		respTags = append(respTags, *formTag(&bTag))

	}
	if rows.Err(); err != nil {
		return nil, err
	}
	return respTags, nil

}

func (r *MgoTag) ListByUserID(uid string) ([]model.Tag, error) {
	query := bson.M{
		"user_id": uid,
		"expire_at": bson.M{
			"$gte": util.Now().UTC(),
		},
	}
	rows, err := r.db.Find(context.Background(), r.table, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	respTags := []model.Tag{}
	for rows.Next() {
		bTag := bsonTag{}
		if err := rows.Scan(&bTag); err != nil {
			return nil, err
		}
		respTags = append(respTags, *formTag(&bTag))

	}
	if rows.Err(); err != nil {
		return nil, err
	}
	return respTags, nil

}
