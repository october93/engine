package search

import (
	"encoding/json"
	"fmt"

	"github.com/algolia/algoliasearch-client-go/algoliasearch"
	"github.com/october93/engine/kit/log"
	"github.com/october93/engine/store"
)

const (
	userIndexName    = "users"
	channelIndexName = "channels"
)

// Indexer is used to index cards into a search engine.
type Indexer struct {
	store  *store.Store
	client algoliasearch.Client
	log    log.Logger
	config *Config
	active bool

	userIndex    algoliasearch.Index
	channelIndex algoliasearch.Index
}

// Indexer returns a new instance of Indexer.
func NewIndexer(s *store.Store, l log.Logger, cfg *Config) (*Indexer, error) {
	indexer := &Indexer{
		store:  s,
		log:    l,
		config: cfg,
	}

	if cfg.ApplicationID != "" {
		indexer.client = algoliasearch.NewClient(cfg.ApplicationID, cfg.AlgoliaAPIKey)
		indexer.active = true
	} else {
		indexer.active = false
		return indexer, nil
	}

	// init indexes
	indexer.userIndex = indexer.client.InitIndex(fmt.Sprintf("%v__%v", cfg.IndexName, userIndexName))
	indexer.channelIndex = indexer.client.InitIndex(fmt.Sprintf("%v__%v", cfg.IndexName, channelIndexName))

	return indexer, nil
}

// IndexAll reads all the cards from the database and indexes them to Algolia in
// order to make them queryable for the frontend applications.
//
// All cards are fetched from the datastore. The number of cards is used to
// determine the pagination parameter for the history query. The history query
// is used to determine which cards should be viewable by which user.
//
// A secured API search key is generated per user and attached to the user.
// Finally, the indexed card data is uploaded to Algolia.
//

func (i *Indexer) IndexAll(clearFirst bool) error {
	if !i.active {
		return nil
	}

	err := i.IndexAllChannels(clearFirst)
	if err != nil {
		return nil
	}
	return i.IndexAllUsers(clearFirst)
}

func algoliaObject(iu interface{}) (algoliasearch.Object, error) {
	var object algoliasearch.Object
	b, err := json.Marshal(iu)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &object)
	if err != nil {
		return nil, err
	}
	return object, nil
}
