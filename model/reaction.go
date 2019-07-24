package model

import (
	"time"

	"github.com/october93/engine/kit/globalid"
)

type ReactionType string

// The Card Reactions, from least to most restrictive
// AllCardReactions is everything the FE can throw into the ReactToCard RPC
// GraphCardReactions is everything the ReactToCard RPC can throw to the graph
// ExplicitCardReactions is one-time explicit button press signals tracked by the graph
// GossipedReactions is reactions which are sent via gossip to all neighbors paying attention
var (
	// AllCardReactions      = []ReactionType{Boost, Bury, HighFive, Pass, Seen, Comment, Like, Love, Haha, Wow, Sad, Angry}
	GraphCardReactions    = []ReactionType{UpVote, DownVote, Comment}
	ExplicitCardReactions = []ReactionType{Boost, Bury}
	GossipedReactions     = []ReactionType{UpVote, DownVote, Comment, Post, UndoComment}
	CardReactions         = []ReactionType{Like, Love, Haha, Wow, Sad, Angry}
	Votes                 = []ReactionType{UpVote, DownVote}
)

const (
	Boost    ReactionType = "boost"
	Bury     ReactionType = "bury"
	HighFive ReactionType = "highfive"
	Pass     ReactionType = "pass"
	Seen     ReactionType = "seen"
	Comment  ReactionType = "comment"

	Like  ReactionType = "like"
	Love  ReactionType = "love"
	Haha  ReactionType = "haha"
	Wow   ReactionType = "wow"
	Sad   ReactionType = "sad"
	Angry ReactionType = "angry"

	UpVote   ReactionType = "up"
	DownVote ReactionType = "down"

	Post        ReactionType = "post"
	ScoreMod    ReactionType = "scoremod"
	Offset      ReactionType = "offset"
	UndoComment ReactionType = "undoComment"
)

type Reaction struct {
	ID        globalid.ID  `db:"id"            json:"id"`
	NodeID    globalid.ID  `db:"node_id"       json:"nodeID"`
	AliasID   globalid.ID  `db:"alias_id"      json:"aliasID,omitempty"`
	CardID    globalid.ID  `db:"card_id"       json:"cardID"`
	Reaction  ReactionType `db:"reaction"      json:"reaction"`
	CreatedAt time.Time    `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time    `db:"updated_at" json:"updatedAt"`
}

func (reaction ReactionType) isAnyOf(others ...ReactionType) bool {
	for _, other := range others {
		if reaction == other {
			return true
		}
	}
	return false
}

// func ValidReaction(reaction ReactionType) bool {
// 	return reaction.isAnyOf(AllCardReactions...)
// }

func GraphCardReaction(reaction ReactionType) bool {
	return reaction.isAnyOf(GraphCardReactions...)
}

func ExplicitReaction(reaction ReactionType) bool {
	return reaction.isAnyOf(ExplicitCardReactions...)
}

func GossipedReaction(reaction ReactionType) bool {
	return reaction.isAnyOf(GossipedReactions...)
}

func CardReaction(reaction ReactionType) bool {
	return reaction.isAnyOf(CardReactions...)
}

func VoteReaction(reaction ReactionType) bool {
	return reaction.isAnyOf(Votes...)
}

func FromVote(vote VoteType) ReactionType {
	if vote == Up {
		return UpVote
	} else {
		return DownVote
	}
}
