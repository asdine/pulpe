package mongo

import (
	"fmt"
	"strconv"
	"strings"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func resolveSlugAndDo(col *mgo.Collection, slugField, slug, sep string, action func(string) error) (string, error) {
	// try to execute the given action with the generated slug
	err := action(slug)
	if err == nil {
		// success, leaving
		return slug, nil
	}

	if !mgo.IsDup(err) || !strings.Contains(err.Error(), slugField) {
		// the action failed because of an unknown error, aborting
		return "", err
	}

	// the slug already exists, fetch the last record that has the same slug
	// and increment the counter
	var distinctSlugs []string

	err = col.Find(
		bson.M{
			slugField: bson.M{
				"$regex": bson.RegEx{
					Pattern: fmt.Sprintf(`^%s(%s\d+)?$`, slug, sep),
					Options: "",
				},
			},
		},
	).Sort("-_id").Select(bson.M{slugField: 1}).Limit(1).Distinct(slugField, &distinctSlugs)
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
		currentSlug := fmt.Sprintf("%s%s%d", slug, sep, counter)
		err := action(currentSlug)
		if err == nil {
			return currentSlug, nil
		}

		if !mgo.IsDup(err) || !strings.Contains(err.Error(), slugField) {
			return "", err
		}
	}
}

func getSelector(slugField, slugOrID string) bson.M {
	if bson.IsObjectIdHex(slugOrID) {
		return bson.M{"_id": bson.ObjectIdHex(slugOrID)}
	}

	return bson.M{slugField: slugOrID}
}
