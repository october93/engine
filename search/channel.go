package search

import (
	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/october93/engine/model"
)

// IndexedUser is
//
// ObjectID is used to unique identify objects in Algolia.
//
type IndexedChannel struct {
	ObjectID string `json:"objectID"`

	Name string `json:"name"`
}

// NewIndexedUser returns a new instance of IndexedUser by copying over
// relevant fields.
func NewIndexedChannel(m *model.Channel) *IndexedChannel {
	i := &IndexedChannel{
		Name:     m.Name,
		ObjectID: m.ID.String(),
	}
	return i
}

func (i *Indexer) IndexAllChannels(clearFirst bool) error {
	if !i.active {
		return nil
	}

	if clearFirst {
		_, err := i.channelIndex.Clear()

		if err != nil {
			return err
		}
	}

	allChannels, err := i.store.GetChannels()
	if err != nil {
		return err
	}

	return i.IndexChannels(allChannels)
}

// UpdateIndex updates the index for a given card by making sure all added
// viewers will be available on the index for the given card.
func (i *Indexer) IndexChannels(channels []*model.Channel) error {
	if !i.active {
		return nil
	}

	objs := make([]algoliasearch.Object, len(channels))
	// fetch existing object from index
	for i, channel := range channels {
		obj, err := algoliaObject(NewIndexedChannel(channel))
		if err != nil {
			return err
		}
		objs[i] = obj
	}

	// upload updated object to Algolia
	_, err := i.channelIndex.AddObjects(objs)
	return err
}

// UpdateIndex updates the index for a given card by making sure all added
// viewers will be available on the index for the given card.
func (i *Indexer) IndexChannel(m *model.Channel) error {
	if !i.active {
		return nil
	}

	// Make a new object from the user (this works even on updates because of the objectID)
	obj, err := algoliaObject(NewIndexedChannel(m))

	if err != nil {
		return err
	}
	// upload updated object to Algolia
	_, err = i.channelIndex.AddObject(obj)
	return err
}

// UpdateIndex updates the index for a given card by making sure all added
// viewers will be available on the index for the given card.
func (i *Indexer) RemoveIndexForChannel(m *model.Channel) error {
	// upload updated object to Algolia
	_, err := i.channelIndex.DeleteObject(m.ID.String())
	return err
}
