package gql

import (
	"net/http"
	"time"

	"github.com/vektah/gqlgen/graphql"
	"github.com/vektah/gqlgen/handler"
)

// SetupHandler returns a new http.Handler capable of processing GraphQL
// queries.
func (g *GraphQL) SetupHandler(exec graphql.ExecutableSchema) http.Handler {
	return g.LoaderMiddleware(handler.GraphQL(exec))
}

const (
	START_OF_WEEK_DAY = time.Sunday
	DATE_FORMAT       = "2006-01-02 MST"
	RETURN_SUCCESS    = "success"
	RETURN_NO_CHANGE  = "no changes made"
)

func truncateDateToDays(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

func containingWeek(t time.Time, startingDayOfWeek time.Weekday) (time.Time, time.Time) {
	diff := int(t.Weekday()) - int(startingDayOfWeek)
	before := (7 + diff) % 7
	after := (7 - diff) % 7

	if after == 0 {
		after = 7
	}

	return truncateDateToDays(t.AddDate(0, 0, -before)),
		truncateDateToDays(t.AddDate(0, 0, after)).Add(-(time.Nanosecond))
}
