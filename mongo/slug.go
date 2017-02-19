package mongo

import (
	"fmt"
	"strconv"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type recorder interface {
	getSlug() string
	setSlug(string)
	elem() interface{}
}

func newBoardRecorder(b *Board) *boardRecorder {
	return &boardRecorder{b}
}

type boardRecorder struct {
	*Board
}

func (b *boardRecorder) getSlug() string {
	return b.Slug
}

func (b *boardRecorder) setSlug(slug string) {
	b.Slug = slug
}

func (b *boardRecorder) elem() interface{} {
	return b.Board
}

func newCardRecorder(c *Card) *cardRecorder {
	return &cardRecorder{c}
}

type cardRecorder struct {
	*Card
}

func (c *cardRecorder) getSlug() string {
	return c.Slug
}

func (c *cardRecorder) setSlug(slug string) {
	c.Slug = slug
}

func (c *cardRecorder) elem() interface{} {
	return c.Card
}

func resolveSlugAndDo(col *mgo.Collection, rec recorder, action func(rec recorder) error) (string, error) {
	// try to execute the given action with the generated slug
	err := action(rec)
	if err == nil {
		// success, leaving
		return rec.getSlug(), nil
	}

	e, ok := err.(*mgo.LastError)
	if !ok || e.Code != codeConflict {
		// the action failed because of an unknown error, aborting
		return "", err
	}

	slug := rec.getSlug()

	// the slug already exists, fetch the last record that has the same slug
	// and increment the counter
	var distinctSlugs []string

	err = col.Find(
		bson.M{
			"slug": bson.M{
				"$regex": bson.RegEx{
					Pattern: fmt.Sprintf(`^%s(-\d+)?$`, slug),
					Options: "",
				},
			},
		},
	).Sort("-_id").Select(bson.M{"slug": 1}).Limit(1).Distinct("slug", &distinctSlugs)
	if err != nil {
		return "", err
	}

	lastSlug := distinctSlugs[0]

	// extract the counter from the slug
	var counter int
	if len(lastSlug) > len(slug) {
		counterStr := lastSlug[len(slug)+1 : len(lastSlug)]
		counter, err = strconv.Atoi(counterStr)
		if err != nil {
			return "", err
		}
	}

	// loop until the action succeeds
	for {
		counter++
		rec.setSlug(fmt.Sprintf("%s-%d", slug, counter))
		err := action(rec)
		if err == nil {
			return rec.getSlug(), nil
		}

		e, ok := err.(*mgo.LastError)
		if !ok || e.Code != codeConflict {
			return "", err
		}
	}
}

func getSelector(slugOrID string) bson.M {
	if bson.IsObjectIdHex(slugOrID) {
		return bson.M{"_id": bson.ObjectIdHex(slugOrID)}
	}

	return bson.M{"slug": slugOrID}
}
