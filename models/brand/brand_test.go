package brand

import (
	"testing"
	o "../testHelpers"
	"../../misc"
	"../../psql"
	"log"
	"io/ioutil"
	"os"
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

func TestShowAll(t *testing.T){
	o.CleanUpDb()

	brands, code := ShowAll()
	if code != misc.NothingToReport {
		t.Errorf("Expect %v. Got %v", misc.NothingToReport, code)
	}

	if len(brands) != len(o.AllBrands) {
		t.Errorf("Expect %v. Got %v", len(o.AllBrands), len(brands))
	}

	for _, brand := range brands {
		b := o.AllBrands[brand.Id]
		if brand.Id != b.Id || brand.Name != b.Name || brand.Issued_at != 0 {
			t.Errorf("Expect %v. Got %v", b, brand)
		}
	}
}

func TestShowById(t *testing.T) {
	o.CleanUpDb()

	table := []struct {
		brandId int
		code    int
		brand   misc.Brand
	}{
		{1, misc.NothingToReport, o.AllBrands[1]},
		{2, misc.NothingToReport, o.AllBrands[2]},
		{3, misc.NothingToReport, o.AllBrands[3]},
		{5, misc.NothingToReport, o.AllBrands[5]},
		{0, misc.NoElement, misc.Brand{}},
		{-1, misc.NoElement, misc.Brand{}},
		{12, misc.NoElement, misc.Brand{}},
		{43, misc.NoElement, misc.Brand{}},
	}

	for num, v := range table {
		brand, code := ShowById(v.brandId)
		if v.code != code || brand.Id != v.brand.Id || brand.Name != v.brand.Name {
			t.Errorf("Case %v. Expect %v, %v. Got %v, %v", num, v.brand, v.code, brand, code)
		}
		if brand.Id != 0 && brand.Issued_at == 0 {
			t.Errorf("Case %v. Expect %v, %v. Got %v, %v", num, v.brand, v.code, brand, code)
		}
	}
}

func TestCreate(t *testing.T) {
	o.CleanUpDb()

	tableSuccess := []struct {
		name string
		id   int
	}{
		{o.RandomString(misc.MaxLenS, 0, 1), 6},
		{o.RandomString(misc.MaxLenS, 0, 0), 7},
		{o.RandomString(misc.MaxLenS, 0, 0), 8},
	}
	for num, v := range tableSuccess {
		id, code := Create(v.name)
		if id != v.id || code != misc.NothingToReport {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.id, id)
		}

		brand, code := ShowById(id)
		if brand.Name != v.name {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.name, brand.Name)
		}
	}

	tableFail := []struct {
		name string
		code int
	}{
		{o.RandomString(misc.MaxLenS, 1, 1), misc.WrongName},
		{o.RandomString(misc.MaxLenS, 1, 0), misc.WrongName},
		{o.RandomString(misc.MaxLenS, 1, 0), misc.WrongName},
		{tableSuccess[0].name, misc.DbDuplicate},
		{tableSuccess[1].name, misc.DbDuplicate},
		{tableSuccess[2].name, misc.DbDuplicate},
	}
	for num, v := range tableFail {
		id, code := Create(v.name)
		if id != 0 || code != v.code {
			t.Errorf("Case %v. Expect 0 %v. Got %v %v", num, v.code, id, code)
		}
	}

	brands, _ := ShowAll()
	brandsNum := len(o.AllBrands) + len(tableSuccess)
	if len(brands) != brandsNum {
		t.Errorf("Expect %v. Got %v", brandsNum, len(brands))
	}
}

func TestUpdateBrand(t *testing.T) {
	o.CleanUpDb()

	randStr := o.RandomString(misc.MaxLenS, 0, 0)
	table := []struct {
		id   int
		name string
		code int
	}{
		{2, o.AllBrands[3].Name, misc.DbDuplicate},
		{3, o.AllBrands[4].Name, misc.DbDuplicate},
		{3, randStr, misc.NothingToReport},
		{4, randStr, misc.DbDuplicate},
		{1, o.RandomString(misc.MaxLenS, 0, 0), misc.NothingToReport},
		{2, o.RandomString(misc.MaxLenS, 0, 1), misc.NothingToReport},
		{2, o.RandomString(misc.MaxLenS, 1, 0), misc.WrongName},
		{5, o.RandomString(misc.MaxLenS, 1, 1), misc.WrongName},
		{0, o.RandomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
		{-1, o.RandomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
		{43, o.RandomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
	}
	for num, v := range table {
		code := Update(v.id, v.name)
		if code != v.code {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.code, code)

			if v.code == misc.NothingToReport {
				brand, _ := ShowById(v.id)
				if brand.Name != v.name {
					t.Errorf("Case %v. Expect %v. Got %v", num, v.name, brand.Name)
				}
			}
		}
	}
}