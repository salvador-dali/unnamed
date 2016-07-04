package tag

import (
	"../../misc"
	"../../psql"
	o "../testHelpers"
	"io/ioutil"
	"log"
	"os"
	"testing"
)

// Setup and db.close will be called before and after each test http://stackoverflow.com/a/34102842/1090562
func TestMain(m *testing.M) {
	o.InitAll()
	log.SetOutput(ioutil.Discard)
	retCode := m.Run()

	defer psql.Db.Close()
	o.CleanUpDb()
	os.Exit(retCode)
}

func TestShowAll(t *testing.T) {
	o.CleanUpDb()

	tags, code := ShowAll()
	if code != misc.NothingToReport {
		t.Error("Expect %v. Got %v", misc.NothingToReport, code)
	}

	if len(tags) != len(o.AllTags) {
		t.Errorf("Expect %v. Got %v", len(o.AllTags), len(tags))
	}

	for num, tag := range tags {
		el := o.AllTags[tag.Id]
		if tag.Id != el.Id || tag.Name != el.Name || tag.Issued_at != 0 || tag.Description != "" {
			t.Errorf("Case %v. Expect %v. Got %v", num, el, tag)
		}
	}
}

func TestShowById(t *testing.T) {
	o.CleanUpDb()

	table := []struct {
		tagId int
		code  int
		tag   misc.Tag
	}{
		{1, misc.NothingToReport, o.AllTags[1]},
		{2, misc.NothingToReport, o.AllTags[2]},
		{3, misc.NothingToReport, o.AllTags[3]},
		{6, misc.NothingToReport, o.AllTags[6]},
		{0, misc.NoElement, misc.Tag{}},
		{-1, misc.NoElement, misc.Tag{}},
		{23, misc.NoElement, misc.Tag{}},
		{43, misc.NoElement, misc.Tag{}},
	}
	for num, v := range table {
		tag, code := ShowById(v.tagId)
		if code != v.code || tag.Id != v.tag.Id || tag.Name != v.tag.Name || tag.Description != v.tag.Description {
			t.Errorf("Case %v. Expect %v. Got %v", num, v, tag)
		}
		if tag.Id != 0 && tag.Issued_at == 0 {
			t.Errorf("Case %v. Expect <nil>. Got %v", num, tag)
		}
	}
}

func TestCreate(t *testing.T) {
	o.CleanUpDb()

	tableSuccess := []struct {
		name  string
		descr string
		id    int
	}{
		{o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 0, 0), 7},
		{o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 0, 1), 8},
		{o.RandomString(misc.MaxLenS, 0, 1), o.RandomString(misc.MaxLenB, 0, 0), 9},
		{o.RandomString(misc.MaxLenS, 0, 1), o.RandomString(misc.MaxLenB, 0, 1), 10},
	}
	for num, v := range tableSuccess {
		id, code := Create(v.name, v.descr)
		tag, _ := ShowById(id)
		if id != v.id || code != misc.NothingToReport || tag.Name != v.name || tag.Description != v.descr {
			t.Errorf("Case %v. Expect %v, %v. Got %v %v", num, v.id, misc.NothingToReport, id, code)
		}
	}

	tableFail := []struct {
		name  string
		descr string
		code  int
	}{
		{o.RandomString(misc.MaxLenS, 1, 0), o.RandomString(misc.MaxLenB, 1, 0), misc.WrongName},
		{o.RandomString(misc.MaxLenS, 1, 0), o.RandomString(misc.MaxLenB, 1, 1), misc.WrongName},
		{o.RandomString(misc.MaxLenS, 1, 1), o.RandomString(misc.MaxLenB, 1, 0), misc.WrongName},
		{o.RandomString(misc.MaxLenS, 1, 1), o.RandomString(misc.MaxLenB, 1, 1), misc.WrongName},
		{tableSuccess[0].name, "d", misc.DbDuplicate},
		{tableSuccess[1].name, "d", misc.DbDuplicate},
		{tableSuccess[2].name, "d", misc.DbDuplicate},
		{tableSuccess[3].name, "d", misc.DbDuplicate},
	}
	for num, v := range tableFail {
		id, code := Create(v.name, v.descr)
		if id != 0 || code != v.code {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.code, code)
		}
	}

	tags, _ := ShowAll()
	tagsNum := len(o.AllTags) + len(tableSuccess)
	if len(tags) != tagsNum {
		t.Errorf("Expect %v. Got %v", tagsNum, len(tags))
	}
}

func TestUpdate(t *testing.T) {
	o.CleanUpDb()

	randStr := o.RandomString(misc.MaxLenS, 0, 0)
	table := []struct {
		id    int
		name  string
		descr string
		code  int
	}{
		{2, "car", o.RandomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{3, "phone", o.RandomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{3, randStr, o.RandomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{4, randStr, o.RandomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{1, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{2, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 0, 1), misc.NothingToReport},
		{2, o.RandomString(misc.MaxLenS, 0, 1), o.RandomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{2, o.RandomString(misc.MaxLenS, 0, 1), o.RandomString(misc.MaxLenB, 0, 1), misc.NothingToReport},
		{5, o.RandomString(misc.MaxLenS, 1, 1), o.RandomString(misc.MaxLenB, 0, 0), misc.WrongName},
		{5, o.RandomString(misc.MaxLenS, 1, 0), o.RandomString(misc.MaxLenB, 0, 0), misc.WrongName},
		{5, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 1, 1), misc.WrongDescr},
		{5, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 1, 0), misc.WrongDescr},
		{0, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 0, 0), misc.NothingUpdated},
		{-1, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 0, 0), misc.NothingUpdated},
		{43, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 0, 0), misc.NothingUpdated},
	}
	for num, v := range table {
		code := Update(v.id, v.name, v.descr)
		if code != v.code {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.code, code)
		}

		if v.code == misc.NothingToReport {
			tag, _ := ShowById(v.id)
			if tag.Name != v.name || tag.Description != v.descr {
				t.Errorf("Case %v. Expect %v, %v. Got %v", num, v.name, v.descr, tag)
			}
		}
	}
}
