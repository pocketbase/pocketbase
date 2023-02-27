package core_test

import (
	"testing"

	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
	"github.com/pocketbase/pocketbase/tools/list"
)

func TestBaseCollectionEventTags(t *testing.T) {
	c1 := new(models.Collection)

	c2 := new(models.Collection)
	c2.Id = "a"

	c3 := new(models.Collection)
	c3.Name = "b"

	c4 := new(models.Collection)
	c4.Id = "a"
	c4.Name = "b"

	scenarios := []struct {
		collection   *models.Collection
		expectedTags []string
	}{
		{c1, []string{}},
		{c2, []string{"a"}},
		{c3, []string{"b"}},
		{c4, []string{"a", "b"}},
	}

	for i, s := range scenarios {
		event := new(core.BaseCollectionEvent)
		event.Collection = s.collection

		tags := event.Tags()

		if len(s.expectedTags) != len(tags) {
			t.Fatalf("[%d] Expected %v tags, got %v", i, s.expectedTags, tags)
		}

		for _, tag := range s.expectedTags {
			if !list.ExistInSlice(tag, tags) {
				t.Fatalf("[%d] Expected %v tags, got %v", i, s.expectedTags, tags)
			}
		}
	}
}

func TestModelEventTags(t *testing.T) {
	m1 := new(models.Admin)

	c := new(models.Collection)
	c.Id = "a"
	c.Name = "b"
	m2 := models.NewRecord(c)

	scenarios := []struct {
		model        models.Model
		expectedTags []string
	}{
		{m1, []string{"_admins"}},
		{m2, []string{"a", "b"}},
	}

	for i, s := range scenarios {
		event := new(core.ModelEvent)
		event.Model = s.model

		tags := event.Tags()

		if len(s.expectedTags) != len(tags) {
			t.Fatalf("[%d] Expected %v tags, got %v", i, s.expectedTags, tags)
		}

		for _, tag := range s.expectedTags {
			if !list.ExistInSlice(tag, tags) {
				t.Fatalf("[%d] Expected %v tags, got %v", i, s.expectedTags, tags)
			}
		}
	}
}
