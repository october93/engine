package model

import (
	"github.com/october93/engine/kit/globalid"

	"math"
)

const (
	TimeScalingFactor = 45000.0
	PopularRankPower  = 0.7
	OctoberUnixOffset = 1475280000.0

	//these are basically random -- let's revisit these parameters
	UpvoteMult   = 1.0
	DownvoteMult = 2.0
	CommentMult  = 1.5
)

type PopularRankEntry struct {
	CardID                globalid.ID `db:"card_id"`
	Views                 int64       `db:"views"`
	UpvoteCount           int64       `db:"upvote_count"`
	DownvoteCount         int64       `db:"downvote_count"`
	CommentCount          int64       `db:"comment_count"`
	UniqueCommentersCount int64       `db:"unique_commenters_count"`
	ScoreMod              float64     `db:"score_mod"`

	CreatedAtTimestamp int64 `db:"created_at_timestamp"`
	//CardOffset ? Only needed if we ever re-insert cards.
}

func (entry PopularRankEntry) Rank() float64 {
	score := float64(entry.UpvoteCount)*UpvoteMult + float64(entry.CommentCount)*CommentMult - float64(entry.DownvoteCount)*DownvoteMult
	return TimeByNormalizedDerivativePower(score, PopularRankPower) + (float64(entry.CreatedAtTimestamp)-OctoberUnixOffset)/TimeScalingFactor + entry.ScoreMod
}

func TimeByNormalizedDerivativePower(v, power float64) float64 {
	if power == 1 {
		return v
	}
	vp := math.Pow(power, 1/(1-power))
	if v < 0 {
		return -(math.Pow(-v+vp, power) - vp/power)
	}
	return math.Pow(v+vp, power) - vp/power
}

type ConfidenceData struct {
	ID                  globalid.ID `json:"id"`
	UpvoteCount         int64       `json:"upvoteCount"`
	DownvoteCount       int64       `json:"downvoteCount"`
	CommentCount        int64       `json:"commentCount"`
	ScoreMod            float64     `json:"scoreMod"`
	ViewCount           int64       `json:"viewCount"`
	Goodness            float64     `json:"goodness"`
	EngagementScore     float64     `json:"engagementScore"`
	Confidence          float64     `json:"confidence"`
	ProbabilitySurfaced float64     `json:"probabilitySurfaced"`
	Rank                float64     `json:"rank"`
}
