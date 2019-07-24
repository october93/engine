package search

import (
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/october93/engine/model"
)

// IndexedUser is
//
// ObjectID is used to unique identify objects in Algolia.
//
type IndexedUser struct {
	ObjectID string `json:"objectID"`

	DisplayName      string `json:"displayName"`
	FirstName        string `json:"firstName"`
	LastName         string `json:"lastName"`
	Username         string `json:"username"`
	ProfileImagePath string `json:"profileImagePath"`
	CoverImagePath   string `json:"coverImagePath"`
	Bio              string `json:"bio"`
}

// NewIndexedUser returns a new instance of IndexedUser by copying over
// relevant fields.
func NewIndexedUser(m *model.User) *IndexedUser {
	iU := &IndexedUser{
		DisplayName:      m.DisplayName,
		FirstName:        m.FirstName,
		LastName:         m.LastName,
		Username:         m.Username,
		ProfileImagePath: m.ProfileImagePath,
		CoverImagePath:   m.CoverImagePath,
		Bio:              m.Bio,
		ObjectID:         m.ID.String(),
	}
	return iU
}

func (i *Indexer) IndexAllUsers(clearFirst bool) error {
	if !i.active {
		return nil
	}

	if clearFirst {
		_, err := i.userIndex.Clear()

		if err != nil {
			return err
		}
	}

	allUsers, err := i.store.GetUsers()
	if err != nil {
		return err
	}

	return i.IndexUsers(allUsers)
}

// UpdateIndex updates the index for a given card by making sure all added
// viewers will be available on the index for the given card.
func (i *Indexer) IndexUsers(users []*model.User) error {
	if !i.active {
		return nil
	}

	objs := make([]algoliasearch.Object, len(users))
	// fetch existing object from index
	for i, user := range users {
		obj, err := algoliaObject(NewIndexedUser(user))
		if err != nil {
			return err
		}
		objs[i] = obj
	}

	// upload updated object to Algolia
	_, err := i.userIndex.AddObjects(objs)
	return err
}

func (i *Indexer) IndexUser(m *model.User) error {
	if !i.active {
		return nil
	}
	// Make a new object from the user (this works even on updates because of the objectID)
	obj, err := algoliaObject(NewIndexedUser(m))

	if err != nil {
		return err
	}
	// upload updated object to Algolia
	_, err = i.userIndex.AddObject(obj)
	return err
}

// UpdateIndex updates the index for a given card by making sure all added
// viewers will be available on the index for the given card.
func (i *Indexer) RemoveIndexForUser(m *model.User) error {
	if !i.active {
		return nil
	}
	// upload updated object to Algolia
	_, err := i.userIndex.DeleteObject(m.ID.String())
	return err
}
