package purchase

import (
	"../../misc"
	"../../psql"
	o "../testHelpers"
	"io/ioutil"
	"log"
	"os"
	"reflect"
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

func TestShowByUserId(t *testing.T) {
	o.CleanUpDb()

	tableSuccess := []struct {
		userId       int
		numPurchases int
	}{
		{1, 3},
		{2, 0},
		{3, 0},
		{4, 1},
		{5, 0},
		{6, 0},
		{7, 0},
		{8, 0},
		{9, 0},
		{-1, 0},
		{10, 0},
	}

	for num, v := range tableSuccess {
		purchases, code := ShowByUserId(v.userId)
		if code != misc.NothingToReport || len(purchases) != v.numPurchases {
			t.Errorf("Case %v. Expect 0 %v. Got %v %v", num, len(purchases), code, v.numPurchases)
		}
	}

	purchases, _ := ShowByUserId(4)
	p, b := purchases[0], o.AllPurchases[2]
	if p.Id != b.Id || p.Image != b.Image || p.Description != b.Description || p.Likes_num != b.Likes_num || p.User_id != b.User_id || p.Brand != b.Brand {
		t.Errorf("Expect %v. Got %v", b, p)
	}

	if !reflect.DeepEqual(p.Tags, b.Tags) {
		t.Errorf("Expect %v. Got %v", b.Tags, p.Tags)
	}
}

func TestShowAll(t *testing.T) {
	o.CleanUpDb()

	purchases, code := ShowAll()
	if code != misc.NothingToReport {
		t.Errorf("Expect %v. Got %v", misc.NothingToReport, code)
	}

	if len(purchases) != len(o.AllPurchases) {
		t.Errorf("Expect %v. Got %v", len(o.AllPurchases), len(purchases))
	}

	for _, v := range purchases {
		p := o.AllPurchases[v.Id]
		if p.Image != v.Image || p.Description != v.Description || p.User_id != v.User_id ||
			p.Brand != v.Brand || p.Likes_num != v.Likes_num {
			t.Errorf("Purcahse %v. Expect %v, got %v", v.Id, v, p)
		}

		if !reflect.DeepEqual(p.Tags, v.Tags) {
			t.Errorf("Purcahse %v. Expect %v, got %v", v.Id, v.Tags, p.Tags)
		}
	}
}

func TestShowById(t *testing.T) {
	o.CleanUpDb()

	tableCorrect := []struct {
		id int
		p  misc.Purchase
	}{
		{1, o.AllPurchases[1]},
		{2, o.AllPurchases[2]},
		{3, o.AllPurchases[3]},
		{4, o.AllPurchases[4]},
	}
	for num, v := range tableCorrect {
		p, code := ShowById(v.id)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect correct execution. Got %v", code)
		}

		if p.Id != v.p.Id || p.Image != v.p.Image || p.Description != v.p.Description ||
			p.User_id != v.p.User_id || p.Brand != v.p.Brand || p.Likes_num != v.p.Likes_num {
			t.Errorf("Case %v. Expect %v. Got %v", num, p, v.p)
		}

		if !reflect.DeepEqual(p.Tags, v.p.Tags) {
			t.Errorf("Purcahse %v. Expect %v, got %v", v.id, v.p.Tags, p.Tags)
		}
	}

	for num, v := range []int{0, -1, 6, 10} {
		p, code := ShowById(v)
		if code != misc.NoElement || p.Id != 0 || p.Image != "" {
			t.Errorf("Case %v. Expectederror. Got %v, %v", num, code, p)
		}
	}
}

func TestShowByBrandId(t *testing.T) {
	o.CleanUpDb()

	tableSuccess := []struct {
		brandId     int
		purchaseIds map[int]bool
	}{
		{5, map[int]bool{2: true}},
		{4, map[int]bool{4: true}},
		{0, map[int]bool{}},
		{-1, map[int]bool{}},
		{9, map[int]bool{}},
	}
	for num, v := range tableSuccess {
		purchases, code := ShowByBrandId(v.brandId)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect correct execution. Got %v", num, code)
		}

		if len(purchases) != len(v.purchaseIds) {
			t.Errorf("Case %v. Expect %v. Got %v", num, len(v.purchaseIds), len(purchases))
		}

		if len(purchases) > 0 {
			for _, p := range purchases {
				expected := o.AllPurchases[p.Id]
				if !v.purchaseIds[p.Id] {
					t.Errorf("Not expected purchase with Id %v", p.Id)
				}
				if expected.Id != p.Id || expected.Image != p.Image || expected.Description != p.Description ||
					expected.Likes_num != p.Likes_num || expected.Brand != p.Brand {
					t.Errorf("Expect %v. Got %v", expected, p)
				}

				if !reflect.DeepEqual(p.Tags, expected.Tags) {
					t.Errorf("Expect %v, got %v", expected.Tags, p.Tags)
				}
			}
		}
	}
}

func TestShowByTagId(t *testing.T) {
	o.CleanUpDb()

	tableSuccess := []struct {
		tagId       int
		purchaseIds map[int]bool
	}{
		{1, map[int]bool{}},
		{2, map[int]bool{1: true, 4: true}},
		{3, map[int]bool{2: true}},
		{4, map[int]bool{3: true}},
		{5, map[int]bool{2: true}},
		{9, map[int]bool{}},
		{-1, map[int]bool{}},
	}
	for num, v := range tableSuccess {
		purchases, code := ShowByTagId(v.tagId)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect correct execution. Got %v", num, code)
		}

		if len(purchases) != len(v.purchaseIds) {
			t.Errorf("Case %v. Expect %v purchases. Got %v", num, len(v.purchaseIds), len(purchases))
		}

		if len(purchases) > 0 {
			for _, p := range purchases {
				expected := o.AllPurchases[p.Id]
				if !v.purchaseIds[p.Id] {
					t.Errorf("Not expected purchase with Id %v", p.Id)
				}
				if expected.Id != p.Id || expected.Image != p.Image || expected.Description != p.Description ||
					expected.Likes_num != p.Likes_num || expected.Brand != p.Brand {
					t.Errorf("Expect %v. Got %v", expected, p)
				}

				if !reflect.DeepEqual(p.Tags, expected.Tags) {
					t.Errorf("Expect %v, got %v", expected.Tags, p.Tags)
				}
			}
		}
	}
}

func TestLike(t *testing.T) {
	o.CleanUpDb()

	tableSuccess := []struct {
		purchaseId int
		userId     int
		likesNum   int
	}{
		{1, 2, 1},
		{1, 3, 2},
		{1, 6, 3},
		{4, 3, 2},
		{3, 2, 4},
	}
	for num, v := range tableSuccess {
		code := Like(v.purchaseId, v.userId)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect correct execution. Got %v", num, code)
		}

		p, _ := ShowById(v.purchaseId)
		if p.Likes_num != v.likesNum {
			t.Errorf("Case %v. Expect %v likes. Got %v", num, v.likesNum, p.Likes_num)
		}
	}

	tableFail := []struct {
		purchaseId int
		userId     int
		code       int
		likesNum   int
	}{
		{1, 1, misc.VoteForYourself, 3},
		{3, 4, misc.DbDuplicate, 4},
		{3, 2, misc.DbDuplicate, 4},
		{0, 2, misc.NoPurchase, 0},
		{9, 1, misc.NoPurchase, 0},
		{1, 11, misc.DbForeignKeyViolation, 3},
	}
	for num, v := range tableFail {
		code := Like(v.purchaseId, v.userId)
		if code != v.code {
			t.Errorf("Case %v. Expect to fail. Got %v", num, code)
		}

		p, _ := ShowById(v.purchaseId)
		if p.Likes_num != v.likesNum {
			t.Errorf("Case %v. Expect %v likes. Got %v", num, v.likesNum, p.Likes_num)
		}
	}
}

func TestUnlike(t *testing.T) {
	o.CleanUpDb()

	tableSuccess := []struct {
		purchaseId int
		userId     int
		likesNum   int
	}{
		{3, 4, 2},
		{4, 2, 0},
		{3, 9, 1},
	}
	for num, v := range tableSuccess {
		code := Unlike(v.purchaseId, v.userId)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect correct execution. Got %v", num, code)
		}

		p, _ := ShowById(v.purchaseId)
		if p.Likes_num != v.likesNum {
			t.Errorf("Case %v. Expect %v likes. Got %v", num, v.likesNum, p.Likes_num)
		}
	}

	tableFail := []struct {
		purchaseId int
		userId     int
		code       int
		likesNum   int
	}{
		{3, 1, misc.VoteForYourself, 1},
		{3, 4, misc.NothingUpdated, 1},
		{1, 7, misc.NothingUpdated, 0},
		{3, 9, misc.NothingUpdated, 1},
		{2, 3, misc.NothingUpdated, 0},
		{-2, -1, misc.NoPurchase, 0},
		{9, 2, misc.NoPurchase, 0},
		{3, -1, misc.NothingUpdated, 1},
	}
	for num, v := range tableFail {
		code := Unlike(v.purchaseId, v.userId)
		if code != v.code {
			t.Errorf("Case %v. Expect to fail. Got %v", num, code)
		}

		p, _ := ShowById(v.purchaseId)
		if p.Likes_num != v.likesNum {
			t.Errorf("Case %v. Expect %v likes. Got %v", num, v.likesNum, p.Likes_num)
		}
	}
}

func TestCreate(t *testing.T) {
	o.CleanUpDb()

	tableSuccess := []struct {
		userId  int
		descr   string
		brandId int
		tagIds  []int
	}{
		{7, o.RandomString(misc.MaxLenB, 0, 0), 1, []int{4}},
		{2, o.RandomString(misc.MaxLenB, 0, 1), 4, []int{2}},
		{1, o.RandomString(misc.MaxLenB, 0, 0), 0, []int{2}},
		{4, o.RandomString(misc.MaxLenB, 0, 0), 0, []int{2, 1}},
		{6, o.RandomString(misc.MaxLenB, 0, 0), 1, []int{2, 4, 1}},
		{5, o.RandomString(misc.MaxLenB, 0, 0), 1, []int{2, 4, 5, 3}},
	}
	for num, v := range tableSuccess {
		id, code := Create(v.userId, v.descr, v.brandId, v.tagIds)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect correct execution. Got %v", num, code)
		}

		if id != num+len(o.AllPurchases)+1 {
			t.Errorf("Case %v. Expect correct ID %v. Got %v", num, num+len(o.AllPurchases)+1, id)
		}

		p, _ := ShowById(id)
		if p.Id != id || p.Description != v.descr || p.Brand != v.brandId {
			t.Errorf("Case %v. Expect %v %v %v. Got %v %v %v", num, p.Id, len(p.Description), p.Brand, id, len(v.descr), v.brandId)
		}

		if !reflect.DeepEqual(p.Tags, v.tagIds) {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.tagIds, p.Tags)
		}

		if !o.IsApproximatelyNow(p.Issued_at) {
			t.Errorf("Case %v. Time %v is not approximately now", num, p.Issued_at)
		}
	}

	tableFail := []struct {
		userId  int
		descr   string
		brandId int
		tagIds  []int
		code    int
	}{
		{19, o.RandomString(misc.MaxLenB, 0, 0), 1, []int{4}, misc.DbForeignKeyViolation},
		{8, o.RandomString(misc.MaxLenB, 0, 0), 1, []int{}, misc.NoTags},
		{5, o.RandomString(misc.MaxLenB, 0, 0), 1, []int{1, 3, 3}, misc.WrongTags},
		{1, o.RandomString(misc.MaxLenB, 0, 0), 1, []int{1, 3, 9}, misc.WrongTags},
		{2, o.RandomString(misc.MaxLenB, 0, 0), 1, []int{1, 3, 2, 5, 1}, misc.WrongTagsNum},
		{3, o.RandomString(misc.MaxLenB, 0, 0), 9, []int{1, 3}, misc.DbForeignKeyViolation},
		{3, o.RandomString(misc.MaxLenB, 1, 0), 2, []int{1, 3}, misc.WrongDescr},
		{3, o.RandomString(misc.MaxLenB, 1, 1), 2, []int{1, 3}, misc.WrongDescr},
	}
	for num, v := range tableFail {
		id, code := Create(v.userId, v.descr, v.brandId, v.tagIds)
		if id != 0 || code != v.code {
			t.Errorf("Case %v. Expect failing. Got %v", num, code)
		}
	}
}
