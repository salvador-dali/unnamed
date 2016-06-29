// In this file there will be a lot of booleans represented as 0 or 1, also some names
// would be very short. This is done to have tests aligned. (true and false are of different length)
package storage

import (
	"../config"
	"../misc"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"testing"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var AllPurchases = map[int]misc.Purchase{
	1: {1, "some_img", "Look at my new drone", 1, 0, []int{2}, 0, 0},
	2: {2, "some_img", "How cool am I?", 4, 0, []int{3, 5}, 5, 0},
	3: {3, "some_img", "I really like drones", 1, 0, []int{4}, 0, 3},
	4: {4, "some_img", "Now I am fond of cars", 1, 0, []int{2}, 4, 1},
}

var AllBrands = map[int]misc.Brand{
	1: {1, "Apple", 0},
	2: {2, "BMW", 0},
	3: {3, "Playstation", 0},
	4: {4, "Ferrari", 0},
	5: {5, "Gucci", 0},
}

var AllTags = map[int]misc.Tag{
	1: {1, "dress", "nice dresses", 0},
	2: {2, "drone", "cool flying machines that do stuff", 0},
	3: {3, "cosmetics", "Known as make-up, are substances or products used to enhance the appearance or scent of the body", 0},
	4: {4, "car", "Vehicles that people use to move faster", 0},
	5: {5, "hat", "Stuff people put on their heads", 0},
	6: {6, "phone", "People use it to speak with other people", 0},
}

var AllUsers = map[int]misc.User{
	1: {1, "Albert Einstein", "", "Developed the general theory of relativity.", 0, 0, 3, 3, 0, 1, 0},
	2: {2, "Isaac Newton", "", "Mechanics, laws of motion", 0, 2, 0, 0, 0, 0, 0},
	// actually there are more of them
}

// randomString generates a random string of a specific length
// This length can be bigger or smaller than you predefined. Also you can ask it to be on the edge
// of the allowed values. For example if you want a value bigger than X, it will generate you
// some strings of the length X + 1 or bigger (if on the edge it will be only X + 1)
// If you want smaller than X, it will generate you anything less or equal to X. If on the edge, it
// will be equal to X
func randomString(length int, isBiggerI, isEdgeCaseI int) string {
	n := 0
	isEdgeCase := isEdgeCaseI == 1
	isBigger := isBiggerI == 1
	if isEdgeCase {
		if isBigger {
			n = length + 1
		} else {
			n = length
		}
	} else {
		if isBigger {
			n = length + 1 + rand.Intn(10)
		} else {
			n = length - 1 - rand.Intn(length-1)
		}
	}

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func isSortedArrayEquivalentToArray(arrSorted, arr []int) bool {
	if len(arrSorted) != len(arr) {
		return false
	}

	sort.Ints(arr)

	for i, v := range arrSorted {
		if v != arr[i] {
			return false
		}
	}
	return true
}

func initializeDb() {
	// initialize Db connection
	config.Init()
	Init(config.Cfg.DbUser, config.Cfg.DbPass, config.Cfg.DbHost, config.Cfg.DbName, config.Cfg.DbPort)
}

func cleanUpDb() {
	// prepare database by creating tables and populating it with data
	cmd := exec.Command("../SQL/set_up_database.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if cmd.Run() != nil {
		log.Fatal("Can't prepare SQL database")
	}
}

// Setup and db.close will be called before and after each test http://stackoverflow.com/a/34102842/1090562
func TestMain(m *testing.M) {
	initializeDb()
	log.SetOutput(ioutil.Discard)
	retCode := m.Run()

	defer Db.Close()
	cleanUpDb()
	os.Exit(retCode)
}

func TestRandomString(t *testing.T) {
	length := 50

	// bigger, not on the edge
	for i := 0; i < 300; i++ {
		s := randomString(length, 1, 0)
		if len(s) <= length {
			t.Errorf("String should be bigger %v, got %v. IMPORTANT function is random, investigate, do not blindly rerun", length, len(s))
		}
	}

	// bigger, on the edge
	for i := 0; i < 20; i++ {
		s := randomString(length, 1, 1)
		if len(s) != length+1 {
			t.Errorf("String should be equal to %v, got %v. IMPORTANT function is random, investigate, do not blindly rerun", length+1, len(s))
		}
	}

	// smaller, not on the edge
	for i := 0; i < 300; i++ {
		s := randomString(length, 0, 0)
		if len(s) > length {
			t.Errorf("String should be smaller or equal to %v, got %v. IMPORTANT function is random, investigate, do not blindly rerun", length, len(s))
		}
	}

	// smaller, on the edge
	for i := 0; i < 20; i++ {
		s := randomString(length, 0, 1)
		if len(s) != length {
			t.Errorf("String should be equal to %v, got %v. IMPORTANT function is random, investigate, do not blindly rerun", length, len(s))
		}
	}
}

// --- Brands tests ---
func TestGetAllBrands(t *testing.T) {
	cleanUpDb()

	brands, code := GetAllBrands()
	if code != misc.NothingToReport {
		t.Errorf("Expect %v. Got %v", misc.NothingToReport, code)
	}

	if len(brands) != len(AllBrands) {
		t.Errorf("Expect %v. Got %v", len(AllBrands), len(brands))
	}

	for _, brand := range brands {
		b := AllBrands[brand.Id]
		if brand.Id != b.Id || brand.Name != b.Name || brand.Issued_at != 0 {
			t.Errorf("Expect %v. Got %v", b, brand)
		}
	}
}

func TestGetBrand(t *testing.T) {
	cleanUpDb()

	table := []struct {
		brandId int
		code    int
		brand   misc.Brand
	}{
		{1, misc.NothingToReport, AllBrands[1]},
		{2, misc.NothingToReport, AllBrands[2]},
		{3, misc.NothingToReport, AllBrands[3]},
		{5, misc.NothingToReport, AllBrands[5]},
		{0, misc.NoElement, misc.Brand{}},
		{-1, misc.NoElement, misc.Brand{}},
		{12, misc.NoElement, misc.Brand{}},
		{43, misc.NoElement, misc.Brand{}},
	}

	for num, v := range table {
		brand, code := GetBrand(v.brandId)
		if v.code != code || brand.Id != v.brand.Id || brand.Name != v.brand.Name {
			t.Errorf("Case %v. Expect %v, %v. Got %v, %v", num, v.brand, v.code, brand, code)
		}
		if brand.Id != 0 && brand.Issued_at == 0 {
			t.Errorf("Case %v. Expect %v, %v. Got %v, %v", num, v.brand, v.code, brand, code)
		}
	}
}

func TestCreateBrand(t *testing.T) {
	cleanUpDb()

	tableSuccess := []struct {
		name string
		id   int
	}{
		{randomString(misc.MaxLenS, 0, 1), 6},
		{randomString(misc.MaxLenS, 0, 0), 7},
		{randomString(misc.MaxLenS, 0, 0), 8},
	}
	for num, v := range tableSuccess {
		id, code := CreateBrand(v.name)
		if id != v.id || code != misc.NothingToReport {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.id, id)
		}

		brand, code := GetBrand(id)
		if brand.Name != v.name {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.name, brand.Name)
		}
	}

	tableFail := []struct {
		name string
		code int
	}{
		{randomString(misc.MaxLenS, 1, 1), misc.WrongName},
		{randomString(misc.MaxLenS, 1, 0), misc.WrongName},
		{randomString(misc.MaxLenS, 1, 0), misc.WrongName},
		{tableSuccess[0].name, misc.DbDuplicate},
		{tableSuccess[1].name, misc.DbDuplicate},
		{tableSuccess[2].name, misc.DbDuplicate},
	}
	for num, v := range tableFail {
		id, code := CreateBrand(v.name)
		if id != 0 || code != v.code {
			t.Errorf("Case %v. Expect 0 %v. Got %v %v", num, v.code, id, code)
		}
	}

	brands, _ := GetAllBrands()
	brandsNum := len(AllBrands) + len(tableSuccess)
	if len(brands) != brandsNum {
		t.Errorf("Expect %v. Got %v", brandsNum, len(brands))
	}
}

func TestUpdateBrand(t *testing.T) {
	cleanUpDb()

	randStr := randomString(misc.MaxLenS, 0, 0)
	table := []struct {
		id   int
		name string
		code int
	}{
		{2, "Playstation", misc.DbDuplicate},
		{3, "Ferrari", misc.DbDuplicate},
		{3, randStr, misc.NothingToReport},
		{4, randStr, misc.DbDuplicate},
		{1, randomString(misc.MaxLenS, 0, 0), misc.NothingToReport},
		{2, randomString(misc.MaxLenS, 0, 1), misc.NothingToReport},
		{2, randomString(misc.MaxLenS, 1, 0), misc.WrongName},
		{5, randomString(misc.MaxLenS, 1, 1), misc.WrongName},
		{0, randomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
		{-1, randomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
		{43, randomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
	}
	for num, v := range table {
		code := UpdateBrand(v.id, v.name)
		if code != v.code {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.code, code)

			if v.code == misc.NothingToReport {
				brand, _ := GetBrand(v.id)
				if brand.Name != v.name {
					t.Errorf("Case %v. Expect %v. Got %v", num, v.name, brand.Name)
				}
			}
		}
	}
}

// --- Tags tests ---
func TestGetAllTags(t *testing.T) {
	cleanUpDb()

	tags, code := GetAllTags()
	if code != misc.NothingToReport {
		t.Error("Expect %v. Got %v", misc.NothingToReport, code)
	}

	if len(tags) != len(AllTags) {
		t.Errorf("Expect %v. Got %v", len(AllTags), len(tags))
	}

	for num, tag := range tags {
		el := AllTags[tag.Id]
		if tag.Id != el.Id || tag.Name != el.Name || tag.Issued_at != 0 || tag.Description != "" {
			t.Errorf("Case %v. Expect %v. Got %v", num, el, tag)
		}
	}
}

func TestGetTag(t *testing.T) {
	cleanUpDb()

	table := []struct {
		tagId int
		code  int
		tag   misc.Tag
	}{
		{1, misc.NothingToReport, AllTags[1]},
		{2, misc.NothingToReport, AllTags[2]},
		{3, misc.NothingToReport, AllTags[3]},
		{6, misc.NothingToReport, AllTags[6]},
		{0, misc.NoElement, misc.Tag{}},
		{-1, misc.NoElement, misc.Tag{}},
		{23, misc.NoElement, misc.Tag{}},
		{43, misc.NoElement, misc.Tag{}},
	}
	for num, v := range table {
		tag, code := GetTag(v.tagId)
		if code != v.code || tag.Id != v.tag.Id || tag.Name != v.tag.Name || tag.Description != v.tag.Description {
			t.Errorf("Case %v. Expect %v. Got %v", num, v, tag)
		}
		if tag.Id != 0 && tag.Issued_at == 0 {
			t.Errorf("Case %v. Expect <nil>. Got %v", num, tag)
		}
	}
}

func TestCreateTag(t *testing.T) {
	cleanUpDb()

	tableSuccess := []struct {
		name  string
		descr string
		id    int
	}{
		{randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 0, 0), 7},
		{randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 0, 1), 8},
		{randomString(misc.MaxLenS, 0, 1), randomString(misc.MaxLenB, 0, 0), 9},
		{randomString(misc.MaxLenS, 0, 1), randomString(misc.MaxLenB, 0, 1), 10},
	}
	for num, v := range tableSuccess {
		id, code := CreateTag(v.name, v.descr)
		tag, _ := GetTag(id)
		if id != v.id || code != misc.NothingToReport || tag.Name != v.name || tag.Description != v.descr {
			t.Errorf("Case %v. Expect %v, %v. Got %v %v", num, v.id, misc.NothingToReport, id, code)
		}
	}

	tableFail := []struct {
		name  string
		descr string
		code  int
	}{
		{randomString(misc.MaxLenS, 1, 0), randomString(misc.MaxLenB, 1, 0), misc.WrongName},
		{randomString(misc.MaxLenS, 1, 0), randomString(misc.MaxLenB, 1, 1), misc.WrongName},
		{randomString(misc.MaxLenS, 1, 1), randomString(misc.MaxLenB, 1, 0), misc.WrongName},
		{randomString(misc.MaxLenS, 1, 1), randomString(misc.MaxLenB, 1, 1), misc.WrongName},
		{tableSuccess[0].name, "d", misc.DbDuplicate},
		{tableSuccess[1].name, "d", misc.DbDuplicate},
		{tableSuccess[2].name, "d", misc.DbDuplicate},
		{tableSuccess[3].name, "d", misc.DbDuplicate},
	}
	for num, v := range tableFail {
		id, code := CreateTag(v.name, v.descr)
		if id != 0 || code != v.code {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.code, code)
		}
	}

	tags, _ := GetAllTags()
	tagsNum := len(AllTags) + len(tableSuccess)
	if len(tags) != tagsNum {
		t.Errorf("Expect %v. Got %v", tagsNum, len(tags))
	}
}

func TestUpdateTag(t *testing.T) {
	cleanUpDb()

	randStr := randomString(misc.MaxLenS, 0, 0)
	table := []struct {
		id    int
		name  string
		descr string
		code  int
	}{
		{2, "car", randomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{3, "phone", randomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{3, randStr, randomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{4, randStr, randomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{1, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{2, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 0, 1), misc.NothingToReport},
		{2, randomString(misc.MaxLenS, 0, 1), randomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{2, randomString(misc.MaxLenS, 0, 1), randomString(misc.MaxLenB, 0, 1), misc.NothingToReport},
		{5, randomString(misc.MaxLenS, 1, 1), randomString(misc.MaxLenB, 0, 0), misc.WrongName},
		{5, randomString(misc.MaxLenS, 1, 0), randomString(misc.MaxLenB, 0, 0), misc.WrongName},
		{5, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 1, 1), misc.WrongDescr},
		{5, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 1, 0), misc.WrongDescr},
		{0, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 0, 0), misc.NothingUpdated},
		{-1, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 0, 0), misc.NothingUpdated},
		{43, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 0, 0), misc.NothingUpdated},
	}
	for num, v := range table {
		code := UpdateTag(v.id, v.name, v.descr)
		if code != v.code {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.code, code)
		}

		if v.code == misc.NothingToReport {
			tag, _ := GetTag(v.id)
			if tag.Name != v.name || tag.Description != v.descr {
				t.Errorf("Case %v. Expect %v, %v. Got %v", num, v.name, v.descr, tag)
			}
		}
	}
}

// --- Users tests ---
func TestGetUser(t *testing.T) {
	cleanUpDb()

	table := []struct {
		userId int
		code   int
		user   misc.User
	}{
		{1, misc.NothingToReport, AllUsers[1]},
		{2, misc.NothingToReport, AllUsers[2]},
		{0, misc.NoElement, misc.User{}},
		{-1, misc.NoElement, misc.User{}},
		{23, misc.NoElement, misc.User{}},
		{43, misc.NoElement, misc.User{}},
	}
	for num, v := range table {
		user, code := GetUser(v.userId)
		if code != v.code {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.code, code)
		}

		if v.code != code || user.Id != v.user.Id || user.Nickname != v.user.Nickname ||
			user.Image != v.user.Image || user.About != v.user.About || user.Expertise != v.user.Expertise ||
			user.Followers_num != v.user.Followers_num || user.Following_num != v.user.Following_num ||
			user.Purchases_num != v.user.Purchases_num || user.Questions_num != v.user.Questions_num ||
			user.Answers_num != v.user.Answers_num {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.user, user)
		}
		if user.Id != 0 && user.Issued_at == 0 {
			t.Errorf("Case %v. Expect <nil>. Got %v", num, user.Issued_at)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	cleanUpDb()

	randStr := randomString(misc.MaxLenS, 0, 0)
	table := []struct {
		id       int
		nickname string
		about    string
		code     int
	}{
		{2, "Marie Curie", randomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{3, "Nikola Tesla", randomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{3, randStr, randomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{4, randStr, randomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{1, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{2, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 0, 1), misc.NothingToReport},
		{3, randomString(misc.MaxLenS, 0, 1), randomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{4, randomString(misc.MaxLenS, 0, 1), randomString(misc.MaxLenB, 0, 1), misc.NothingToReport},
		{2, randomString(misc.MaxLenS, 1, 0), randomString(misc.MaxLenB, 0, 0), misc.WrongName},
		{5, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenB, 1, 0), misc.WrongDescr},
		{5, randomString(misc.MaxLenS, 1, 0), randomString(misc.MaxLenB, 1, 0), misc.WrongName},
		{0, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
		{-1, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
		{43, randomString(misc.MaxLenS, 0, 0), randomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
	}
	for num, v := range table {
		code := UpdateUser(v.id, v.nickname, v.about)
		if code != v.code {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.code, code)
		}

		if code == misc.NothingToReport {
			user, _ := GetUser(v.id)
			if user.Nickname != v.nickname || user.About != v.about {
				t.Errorf("Case %v. Expect %v %v. Got %v", num, v.nickname, v.about, user)
			}
		}
	}
}

func TestGetFollowers(t *testing.T) {
	cleanUpDb()

	tableSuccess := []struct {
		id        int
		followers []int
	}{
		{1, []int{}},
		{2, []int{1, 6}},
		{3, []int{}},
		{4, []int{1}},
		{7, []int{1}},
	}
	for num, v := range tableSuccess {
		followers, code := GetFollowers(v.id)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect %v. Got %v", num, misc.NothingToReport, code)
		}

		if len(followers) != len(v.followers) {
			t.Errorf("Case %v. Expect %v. Got %v", num, len(v.followers), len(followers))
		}

		followerIds := make([]int, len(followers), len(followers))
		for i, v := range followers {
			followerIds[i] = v.Id
		}

		if !isSortedArrayEquivalentToArray(v.followers, followerIds) {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.followers, followerIds)
		}
	}

	for num, id := range []int{0, 16, 52, -1} {
		followers, code := GetFollowers(id)
		if len(followers) != 0 || code != misc.NothingToReport {
			t.Errorf("Case %v. Expect 0. Got %v", num, len(followers))
		}
	}

	followers, _ := GetFollowers(7)
	u := followers[0]
	if u.Nickname != AllUsers[1].Nickname || u.About != "" || u.Expertise != 0 || u.Followers_num != 0 {
		t.Errorf("Information about follower is not right %v", u)
	}
}

func TestGetFollowing(t *testing.T) {
	cleanUpDb()

	tableSuccess := []struct {
		id        int
		following []int
	}{
		{1, []int{2, 4, 7}},
		{2, []int{}},
		{6, []int{2}},
		{7, []int{}},
	}
	for num, v := range tableSuccess {
		following, code := GetFollowing(v.id)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect %v. Got %v", num, misc.NothingToReport, code)
		}

		if len(following) != len(v.following) {
			t.Errorf("Case %v. Expect %v. Got %v", num, len(v.following), len(following))
		}

		followingIds := make([]int, len(following), len(following))
		for i, v := range following {
			followingIds[i] = v.Id
		}

		if !isSortedArrayEquivalentToArray(v.following, followingIds) {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.following, followingIds)
		}
	}

	for num, id := range []int{0, 16, 52, -1} {
		followers, code := GetFollowing(id)
		if len(followers) != 0 || code != misc.NothingToReport {
			t.Errorf("Case %v. Expect 0. Got %v", num, len(followers))
		}
	}

	following, _ := GetFollowing(6)
	u := following[0]
	if u.Nickname != AllUsers[2].Nickname || u.About != "" || u.Expertise != 0 || u.Followers_num != 0 {
		t.Errorf("Information about following is not right %v", u)
	}
}

func TestFollow(t *testing.T) {
	cleanUpDb()

	table := []struct {
		whoId         int
		whomId        int
		code          int
		followers_num int
		following_num int
	}{
		{1, 1, misc.FollowYourself, 0, 3},
		{1, 2, misc.DbDuplicate, 2, 3},
		{6, 2, misc.DbDuplicate, 2, 1},
		{0, 2, misc.DbForeignKeyViolation, 2, 0},
		{6, -1, misc.NoElement, 0, 1},
		{10, 54, misc.DbForeignKeyViolation, 0, 0},
		{1, 6, misc.NothingToReport, 1, 4},
		{6, 1, misc.NothingToReport, 1, 2},
		{2, 4, misc.NothingToReport, 2, 1},
		{2, 6, misc.NothingToReport, 2, 2},
	}
	for num, v := range table {
		code := Follow(v.whoId, v.whomId)
		if code != v.code {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.code, code)
		}

		followers, _ := GetFollowers(v.whomId)
		following, _ := GetFollowing(v.whoId)

		if len(followers) != v.followers_num || len(following) != v.following_num {
			t.Errorf("Case %v. Expect (%v, %v). Got (%v, %v)", num, v.followers_num, v.following_num, len(followers), len(following))
		}

		u1, _ := GetUser(v.whoId)
		u2, _ := GetUser(v.whomId)
		if u2.Followers_num != v.followers_num || u1.Following_num != v.following_num {
			t.Errorf("Case %v. Expect (%v, %v). Got (%v, %v)", num, v.followers_num, v.following_num, u2.Followers_num, u1.Following_num)
		}
	}
}

func TestUnfollow(t *testing.T) {
	cleanUpDb()

	table := []struct {
		whoId         int
		whomId        int
		code          int
		followers_num int
		following_num int
	}{
		{1, 1, misc.FollowYourself, 0, 3},
		{1, 5, misc.NothingUpdated, 0, 3},
		{6, 3, misc.NothingUpdated, 0, 1},
		{5, 4, misc.NothingUpdated, 1, 0},
		{-1, 4, misc.NothingUpdated, 1, 0},
		{10, 9, misc.NothingUpdated, 0, 0},
		{11, 19, misc.NothingUpdated, 0, 0},
		{1, 6, misc.NothingUpdated, 0, 3},
		{6, 2, misc.NothingToReport, 1, 0},
		{1, 4, misc.NothingToReport, 0, 2},
		{1, 7, misc.NothingToReport, 0, 1},
		{1, 2, misc.NothingToReport, 0, 0},
	}

	for num, v := range table {
		code := Unfollow(v.whoId, v.whomId)
		if code != v.code {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.code, code)
		}

		followers, _ := GetFollowers(v.whomId)
		following, _ := GetFollowing(v.whoId)

		if len(followers) != v.followers_num || len(following) != v.following_num {
			t.Errorf("Case %v. Expect (%v, %v), got (%v, %v)", num, v.followers_num, v.following_num, len(followers), len(following))
		}

		u1, _ := GetUser(v.whoId)
		u2, _ := GetUser(v.whomId)
		if u2.Followers_num != v.followers_num || u1.Following_num != v.following_num {
			t.Errorf("Case %v. Expect (%v, %v), got (%v, %v)", num, v.followers_num, v.following_num, u2.Followers_num, u1.Following_num)
		}
	}
}

func TestCreateUser(t *testing.T) {
	cleanUpDb()

	tableSuccess := []struct {
		nickname string
		email    string
		password string
		userId   int
	}{
		{"User_134", "email@gmail.com", "just_some_password", 11},
		{"stuff", "email@somemail.com", "anotherPa$$W0rt", 12},
		{"another", "random@yahoo.com", "anotherPa$$W0rt", 13},
	}
	for num, v := range tableSuccess {
		userId, code := CreateUser(v.nickname, v.email, v.password)
		if code != misc.NothingToReport || userId != v.userId {
			t.Errorf("Case %v. Expect 0, %v. Got %v, %v", num, v.userId, code, userId)
		}
	}

	tableFail := []struct {
		nickname string
		email    string
		password string
		code     int
	}{
		{"exist", "albert@gmail.com", "password", misc.DbDuplicate},
		{AllUsers[2].Nickname, "good@mail.com", "password", misc.DbDuplicate},
		{tableSuccess[2].nickname, "amail@mail.com", "password", misc.DbDuplicate},
		{"random", tableSuccess[2].email, "password", misc.DbDuplicate},
	}
	for num, v := range tableFail {
		userId, code := CreateUser(v.nickname, v.email, v.password)
		if code != v.code || userId != 0 {
			t.Errorf("Case %v. Expect 0, %v. Got %v, %v", num, v.code, userId, code)
		}
	}

	user, code := GetUser(tableSuccess[0].userId)
	if user.Nickname != tableSuccess[0].nickname || user.Id != tableSuccess[0].userId {
		t.Errorf("Expected to get user. Got %v, %v", user, code)
	}
}

func TestLogin(t *testing.T) {
	cleanUpDb()

	email, pass := "some_strange_mail@gmail.com", "very_new_password"
	CreateUser("username", email, pass)

	tableSuccess := []struct {
		email    string
		password string
	}{
		{"albert@gmail.com", "password"},
		{"isaac@gmail.com", "password"},
		{"michael@gmail.com", "password"},
		{email, pass},
	}
	for num, v := range tableSuccess {
		jwt, ok := Login(v.email, v.password)
		if !ok || len(jwt) < 10 {
			t.Errorf("Case %v. Expect to log in. Got %v, %v", num, ok, jwt)
		}
	}

	tableFail := []struct {
		email    string
		password string
	}{
		{tableSuccess[0].email, tableSuccess[0].email},
		{"a" + tableSuccess[1].email, tableSuccess[1].password},
		{tableSuccess[2].email, tableSuccess[2].password + "a"},
	}
	for num, v := range tableFail {
		jwt, ok := Login(v.email, v.password)
		if ok || jwt != "" {
			t.Errorf("Case %v. Expect to fail. Got %v, %v", num, ok, jwt)
		}
	}
}

func TestVerifyEmail(t *testing.T) {
	cleanUpDb()

	tableFail := []struct {
		userId     int
		verifyCode string
	}{
		{11, "pqaJaBRgAvzLXqzRrrUI"},
		{1, ""},
		{5, "pqaJaBRgAvzLXqzRrrUI"},
		{500, "pqaJaBRgAvzLXqzRrrUIsafasdfsad"},
		{10, "pqaJaBRgAvzLXqzRrrUIsafasdfsad"},
	}
	for num, v := range tableFail {
		if VerifyEmail(v.userId, v.verifyCode) {
			t.Errorf("Case %v. Expect to fail. Got True", num)
		}
	}

	if !VerifyEmail(10, "pqaJaBRgAvzLXqzRrrUI") {
		t.Errorf("Expect to verify email. Got False")
	}
}

// --- Purchases tests ---
func TestGetUserPurchases(t *testing.T) {
	cleanUpDb()

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
		purchases, code := GetUserPurchases(v.userId)
		if code != misc.NothingToReport || len(purchases) != v.numPurchases {
			t.Errorf("Case %v. Expect 0 %v. Got %v %v", num, len(purchases), code, v.numPurchases)
		}
	}

	purchases, _ := GetUserPurchases(4)
	p, o := purchases[0], AllPurchases[2]
	if p.Id != o.Id || p.Image != o.Image || p.Description != o.Description || p.Likes_num != o.Likes_num || p.User_id != o.User_id || p.Brand != o.Brand {
		t.Errorf("Expect %v. Got %v", o, p)
	}

	if !reflect.DeepEqual(p.Tags, o.Tags) {
		t.Errorf("Expect %v. Got %v", o.Tags, p.Tags)
	}
}

func TestGetAllPurchases(t *testing.T) {
	cleanUpDb()

	purchases, code := GetAllPurchases()
	if code != misc.NothingToReport {
		t.Errorf("Expect %v. Got %v", misc.NothingToReport, code)
	}

	if len(purchases) != len(AllPurchases) {
		t.Errorf("Expect %v. Got %v", len(AllPurchases), len(purchases))
	}

	for _, v := range purchases {
		p := AllPurchases[v.Id]
		if p.Image != v.Image || p.Description != v.Description || p.User_id != v.User_id ||
			p.Brand != v.Brand || p.Likes_num != v.Likes_num {
			t.Errorf("Purcahse %v. Expect %v, got %v", v.Id, v, p)
		}

		if !reflect.DeepEqual(p.Tags, v.Tags) {
			t.Errorf("Purcahse %v. Expect %v, got %v", v.Id, v.Tags, p.Tags)
		}
	}
}

func TestGetPurchase(t *testing.T) {
	cleanUpDb()

	tableCorrect := []struct {
		id int
		p  misc.Purchase
	}{
		{1, AllPurchases[1]},
		{2, AllPurchases[2]},
		{3, AllPurchases[3]},
		{4, AllPurchases[4]},
	}
	for num, v := range tableCorrect {
		p, code := GetPurchase(v.id)
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
		p, code := GetPurchase(v)
		if code != misc.NoElement || p.Id != 0 || p.Image != "" {
			t.Errorf("Case %v. Expectederror. Got %v, %v", num, code, p)
		}
	}
}

func TestGetAllPurchasesWithBrand(t *testing.T) {
	cleanUpDb()

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
		purchases, code := GetAllPurchasesWithBrand(v.brandId)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect correct execution. Got %v", num, code)
		}

		if len(purchases) != len(v.purchaseIds) {
			t.Errorf("Case %v. Expect %v. Got %v", num, len(v.purchaseIds), len(purchases))
		}

		if len(purchases) > 0 {
			for _, p := range purchases {
				expected := AllPurchases[p.Id]
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

func TestGetAllPurchasesWithTag(t *testing.T) {
	cleanUpDb()

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
		purchases, code := GetAllPurchasesWithTag(v.tagId)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect correct execution. Got %v", num, code)
		}

		if len(purchases) != len(v.purchaseIds) {
			t.Errorf("Case %v. Expect %v purchases. Got %v", num, len(v.purchaseIds), len(purchases))
		}

		if len(purchases) > 0 {
			for _, p := range purchases {
				expected := AllPurchases[p.Id]
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

func TestLikePurchase(t *testing.T) {
	cleanUpDb()

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
		code := LikePurchase(v.purchaseId, v.userId)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect correct execution. Got %v", num, code)
		}

		p, _ := GetPurchase(v.purchaseId)
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
		code := LikePurchase(v.purchaseId, v.userId)
		if code != v.code {
			t.Errorf("Case %v. Expect to fail. Got %v", num, code)
		}

		p, _ := GetPurchase(v.purchaseId)
		if p.Likes_num != v.likesNum {
			t.Errorf("Case %v. Expect %v likes. Got %v", num, v.likesNum, p.Likes_num)
		}
	}
}

func TestUnlikePurchase(t *testing.T) {
	cleanUpDb()

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
		code := UnlikePurchase(v.purchaseId, v.userId)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect correct execution. Got %v", num, code)
		}

		p, _ := GetPurchase(v.purchaseId)
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
		code := UnlikePurchase(v.purchaseId, v.userId)
		if code != v.code {
			t.Errorf("Case %v. Expect to fail. Got %v", num, code)
		}

		p, _ := GetPurchase(v.purchaseId)
		if p.Likes_num != v.likesNum {
			t.Errorf("Case %v. Expect %v likes. Got %v", num, v.likesNum, p.Likes_num)
		}
	}
}

func TestCreatePurchase(t *testing.T) {
	cleanUpDb()

	tableSuccess := []struct {
		userId  int
		descr   string
		brandId int
		tagIds  []int
	}{
		{7, randomString(misc.MaxLenB, 0, 0), 1, []int{4}},
		{2, randomString(misc.MaxLenB, 0, 1), 4, []int{2}},
		{1, randomString(misc.MaxLenB, 0, 0), 0, []int{2}},
		{4, randomString(misc.MaxLenB, 0, 0), 0, []int{2, 1}},
		{6, randomString(misc.MaxLenB, 0, 0), 1, []int{2, 4, 1}},
		{5, randomString(misc.MaxLenB, 0, 0), 1, []int{2, 4, 5, 3}},
	}
	for num, v := range tableSuccess {
		timeNow := time.Now().Unix()
		id, code := CreatePurchase(v.userId, v.descr, v.brandId, v.tagIds)
		if code != misc.NothingToReport {
			t.Errorf("Case %v. Expect correct execution. Got %v", num, code)
		}

		if id != num+len(AllPurchases)+1 {
			t.Errorf("Case %v. Expect correct ID %v. Got %v", num, num+len(AllPurchases)+1, id)
		}

		p, _ := GetPurchase(id)
		if p.Id != id || p.Description != v.descr || p.Brand != v.brandId {
			t.Errorf("Case %v. Expect %v %v %v. Got %v %v %v", num, p.Id, len(p.Description), p.Brand, id, len(v.descr), v.brandId)
		}

		if !reflect.DeepEqual(p.Tags, v.tagIds) {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.tagIds, p.Tags)
		}

		if timeNow != p.Issued_at {
			t.Errorf("Case %v. Expect %v. Got %v", num, timeNow, p.Issued_at)
		}
	}

	tableFail := []struct {
		userId  int
		descr   string
		brandId int
		tagIds  []int
		code    int
	}{
		{19, randomString(misc.MaxLenB, 0, 0), 1, []int{4}, misc.DbForeignKeyViolation},
		{8, randomString(misc.MaxLenB, 0, 0), 1, []int{}, misc.NoTags},
		{5, randomString(misc.MaxLenB, 0, 0), 1, []int{1, 3, 3}, misc.WrongTags},
		{1, randomString(misc.MaxLenB, 0, 0), 1, []int{1, 3, 9}, misc.WrongTags},
		{2, randomString(misc.MaxLenB, 0, 0), 1, []int{1, 3, 2, 5, 1}, misc.WrongTagsNum},
		{3, randomString(misc.MaxLenB, 0, 0), 9, []int{1, 3}, misc.DbForeignKeyViolation},
		{3, randomString(misc.MaxLenB, 1, 0), 2, []int{1, 3}, misc.WrongDescr},
		{3, randomString(misc.MaxLenB, 1, 1), 2, []int{1, 3}, misc.WrongDescr},
	}
	for num, v := range tableFail {
		id, code := CreatePurchase(v.userId, v.descr, v.brandId, v.tagIds)
		if id != 0 || code != v.code {
			t.Errorf("Case %v. Expect failing. Got %v", num, code)
		}
	}
}
