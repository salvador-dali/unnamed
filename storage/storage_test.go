// In this file there will be a lot of booleans represented as 0 or 1, also some names
// would be very short. This is done to have tests aligned. (true and false are of different length)
package storage

import (
	"../../unnamed/config"
	"../../unnamed/errorCodes"
	"../../unnamed/structs"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"testing"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	maxLenS     = 40
	maxLenB     = 1000
)

var AllPurchases = map[int]structs.Purchase{
	1: {1, "some_img", "Look at my new drone", 1, nil, []int{}, 0, 0},
	2: {2, "some_img", "How cool am I?", 4, nil, []int{}, 5, 0},
	3: {3, "some_img", "I really like drones", 1, nil, []int{}, 0, 3},
	4: {4, "some_img", "Now I am fond of cars", 1, nil, []int{}, 4, 1},
}

var AllBrands = map[int]structs.Brand{
	1: {1, "Apple", nil},
	2: {2, "BMW", nil},
	3: {3, "Playstation", nil},
	4: {4, "Ferrari", nil},
	5: {5, "Gucci", nil},
}

var AllTags = map[int]structs.Tag{
	1: {1, "dress", "nice dresses", nil},
	2: {2, "drone", "cool flying machines that do stuff", nil},
	3: {3, "cosmetics", "Known as make-up, are substances or products used to enhance the appearance or scent of the body", nil},
	4: {4, "car", "Vehicles that people use to move faster", nil},
	5: {5, "hat", "Stuff people put on their heads", nil},
	6: {6, "phone", "People use it to speak with other people", nil},
}

var AllUsers = map[int]structs.User{
	1: {1, "Albert Einstein", "", "Developed the general theory of relativity.", 0, 0, 3, 3, 0, 1, nil},
	2: {2, "Isaac Newton", "", "Mechanics, laws of motion", 0, 2, 0, 0, 0, 0, nil},
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
	retCode := m.Run()

	defer Db.Close()
	//cleanUpDb()
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

	brands, err, code := GetAllBrands()
	if err != nil || code != errorCodes.DbNothingToReport {
		t.Error("Should finish without no error")
	}

	if len(brands) != len(AllBrands) {
		t.Errorf("Expect %v brands, got %v", len(AllBrands), len(brands))
	}

	for _, brand := range brands {
		b := AllBrands[brand.Id]
		if brand.Id != b.Id || brand.Name != b.Name || brand.Issued_at != nil {
			t.Errorf("Expect %v , got %v", b, brand)
		}
	}
}

func TestGetBrand(t *testing.T) {
	cleanUpDb()

	table := []struct {
		brandId      int
		res_is_error int
		res_code     int
		res          structs.Brand
	}{
		{1, 0, errorCodes.DbNothingToReport, AllBrands[1]},
		{2, 0, errorCodes.DbNothingToReport, AllBrands[2]},
		{3, 0, errorCodes.DbNothingToReport, AllBrands[3]},
		{5, 0, errorCodes.DbNothingToReport, AllBrands[5]},
		{0, 1, errorCodes.DbNoElement, structs.Brand{}},
		{-1, 1, errorCodes.DbNoElement, structs.Brand{}},
		{123, 1, errorCodes.DbNoElement, structs.Brand{}},
		{43, 1, errorCodes.DbNoElement, structs.Brand{}},
	}

	for _, v := range table {
		brand, err, code := GetBrand(v.brandId)
		if v.res_is_error == 1 && err == nil {
			t.Errorf("Wrong result for case %v. Expected error, did not get it", v.brandId)
		}

		if v.res_is_error == 0 && err != nil {
			t.Errorf("Wrong result for case %v. Expected nil, got error", v.brandId)
		}

		if code != v.res_code || brand.Id != v.res.Id || brand.Name != v.res.Name {
			t.Errorf("Wrong result for case %v: \n Expected %v \n Got %v %v %v", v.brandId, v, err, code, brand)
		}
		if brand.Id != 0 && brand.Issued_at == nil {
			t.Errorf("Wrong result for case %v: Real Brand has an Issued_at date", v.brandId)
		}
	}
}

func TestCreateBrand(t *testing.T) {
	cleanUpDb()

	tableSuccess := []struct {
		name   string
		res_id int
	}{
		{randomString(maxLenS, 0, 1), 6},
		{randomString(maxLenS, 0, 0), 7},
		{randomString(maxLenS, 0, 0), 8},
	}

	tableFail := []struct {
		name string
		code int
	}{
		{randomString(maxLenS, 1, 1), errorCodes.DbValueTooLong},
		{randomString(maxLenS, 1, 0), errorCodes.DbValueTooLong},
		{randomString(maxLenS, 1, 0), errorCodes.DbValueTooLong},
		{tableSuccess[0].name, errorCodes.DbDuplicate},
		{tableSuccess[1].name, errorCodes.DbDuplicate},
		{tableSuccess[2].name, errorCodes.DbDuplicate},
	}

	for _, v := range tableSuccess {
		id, err, code := CreateBrand(v.name)
		if id != v.res_id || err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Expected to create a brand. Got %v %v %v", id, err, code)
		}

		brand, err, code := GetBrand(id)
		if brand.Name != v.name {
			t.Errorf("Expected to create a brand with a name %v, got %v", brand.Name, v.name)
		}
	}

	for _, v := range tableFail {
		id, err, code := CreateBrand(v.name)
		if err == nil || id != 0 || code != v.code {
			t.Error("New brand should not be created")
		}
	}

	brands, _, _ := GetAllBrands()
	brandsNum := len(AllBrands) + len(tableSuccess)
	if len(brands) != brandsNum {
		t.Errorf("Should have %v brands, have %v", brandsNum, len(brands))
	}
}

func TestUpdateBrand(t *testing.T) {
	cleanUpDb()

	randStr := randomString(maxLenS, 0, 0)
	table := []struct {
		id           int
		name         string
		res_is_error int
		res_code     int
	}{
		{2, "Playstation", 1, errorCodes.DbDuplicate},
		{3, "Ferrari", 1, errorCodes.DbDuplicate},
		{3, randStr, 0, errorCodes.DbNothingToReport},
		{4, randStr, 1, errorCodes.DbDuplicate},
		{1, randomString(maxLenS, 0, 0), 0, errorCodes.DbNothingToReport},
		{2, randomString(maxLenS, 0, 1), 0, errorCodes.DbNothingToReport},
		{2, randomString(maxLenS, 1, 0), 1, errorCodes.DbValueTooLong},
		{5, randomString(maxLenS, 1, 1), 1, errorCodes.DbValueTooLong},
		{0, randomString(maxLenS, 0, 0), 1, errorCodes.DbNothingUpdated},
		{-1, randomString(maxLenS, 0, 0), 1, errorCodes.DbNothingUpdated},
		{43, randomString(maxLenS, 0, 0), 1, errorCodes.DbNothingUpdated},
	}

	for _, v := range table {
		err, code := UpdateBrand(v.id, v.name)
		if v.res_is_error == 1 {
			if err == nil || code != v.res_code {
				t.Errorf("The brand %v should not be updated, but it was: %v, %v", v.id, err, code)
			}
		} else {
			if err != nil || code != v.res_code {
				t.Errorf("The brand %v should have been updated, but was not", v.id)
			}

			brand, _, _ := GetBrand(v.id)
			if brand.Name != v.name {
				t.Errorf("Expected value %v after update, got %v", v.name, brand.Name)
			}
		}
	}
}

// --- Tags tests ---
func TestGetAllTags(t *testing.T) {
	cleanUpDb()

	tags, err, code := GetAllTags()
	if err != nil || code != errorCodes.DbNothingToReport {
		t.Error("Should finish without no error")
	}

	if len(tags) != len(AllTags) {
		t.Errorf("Expect %v tags, got %v", len(AllTags), len(tags))
	}

	for _, tag := range tags {
		el := AllTags[tag.Id]
		if tag.Id != el.Id || tag.Name != el.Name || tag.Issued_at != nil || tag.Description != "" {
			t.Errorf("Expect %v , got %v", el, tag)
		}
	}
}

func TestGetTag(t *testing.T) {
	cleanUpDb()

	table := []struct {
		tagId        int
		res_is_error int
		res_code     int
		res          structs.Tag
	}{
		{1, 0, errorCodes.DbNothingToReport, AllTags[1]},
		{2, 0, errorCodes.DbNothingToReport, AllTags[2]},
		{3, 0, errorCodes.DbNothingToReport, AllTags[3]},
		{6, 0, errorCodes.DbNothingToReport, AllTags[6]},
		{0, 1, errorCodes.DbNoElement, structs.Tag{}},
		{-1, 1, errorCodes.DbNoElement, structs.Tag{}},
		{23, 1, errorCodes.DbNoElement, structs.Tag{}},
		{43, 1, errorCodes.DbNoElement, structs.Tag{}},
	}

	for _, v := range table {
		tag, err, code := GetTag(v.tagId)
		if v.res_is_error == 1 && err == nil {
			t.Errorf("Wrong result for case %v. Expected error, did not get it", v.tagId)
		}

		if v.res_is_error == 0 && err != nil {
			t.Errorf("Wrong result for case %v. Expected nil, got error", v.tagId)
		}

		if code != v.res_code || tag.Id != v.res.Id || tag.Name != v.res.Name || tag.Description != v.res.Description {
			t.Errorf("Wrong result for case %v: \n Expected %v \n Got %v %v %v", v.tagId, v, err, code, tag)
		}
		if tag.Id != 0 && tag.Issued_at == nil {
			t.Errorf("Wrong result for case %v: Real Brand has an Issued_at date", v.tagId)
		}
	}
}

func TestCreateTag(t *testing.T) {
	cleanUpDb()

	tableSuccess := []struct {
		name   string
		descr  string
		res_id int
	}{
		{randomString(maxLenS, 0, 0), randomString(maxLenB, 0, 0), 7},
		{randomString(maxLenS, 0, 0), randomString(maxLenB, 0, 1), 8},
		{randomString(maxLenS, 0, 1), randomString(maxLenB, 0, 0), 9},
		{randomString(maxLenS, 0, 1), randomString(maxLenB, 0, 1), 10},
	}
	for _, v := range tableSuccess {
		id, err, code := CreateTag(v.name, v.descr)
		if id != v.res_id || err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Expected to create a tag. Got %v %v %v", id, err, code)
		}

		tag, _, _ := GetTag(id)
		if tag.Name != v.name || tag.Description != v.descr {
			t.Errorf("Expected to create a tag with a name %v, got %v", tag.Name, v.name)
		}
	}

	tableFail := []struct {
		name  string
		descr string
		code  int
	}{
		{randomString(maxLenS, 1, 0), randomString(maxLenB, 1, 0), errorCodes.DbValueTooLong},
		{randomString(maxLenS, 1, 0), randomString(maxLenB, 1, 1), errorCodes.DbValueTooLong},
		{randomString(maxLenS, 1, 1), randomString(maxLenB, 1, 0), errorCodes.DbValueTooLong},
		{randomString(maxLenS, 1, 1), randomString(maxLenB, 1, 1), errorCodes.DbValueTooLong},
		{tableSuccess[0].name, "", errorCodes.DbDuplicate},
		{tableSuccess[1].name, "", errorCodes.DbDuplicate},
		{tableSuccess[2].name, "", errorCodes.DbDuplicate},
		{tableSuccess[3].name, "", errorCodes.DbDuplicate},
	}

	for _, v := range tableFail {
		id, err, code := CreateTag(v.name, v.descr)
		if err == nil || id != 0 || code != v.code {
			t.Errorf("New tag should not be created %v, %v", id, code)
		}
	}

	tags, _, _ := GetAllTags()
	tagsNum := len(AllTags) + len(tableSuccess)
	if len(tags) != tagsNum {
		t.Errorf("Should have %v tags, have %v", tagsNum, len(tags))
	}
}

func TestUpdateTag(t *testing.T) {
	cleanUpDb()

	randStr := randomString(maxLenS, 0, 0)
	table := []struct {
		id           int
		name         string
		descr        string
		res_is_error int
		res_code     int
	}{
		{2, "car", randomString(maxLenB, 0, 0), 1, errorCodes.DbDuplicate},
		{3, "phone", randomString(maxLenB, 0, 0), 1, errorCodes.DbDuplicate},
		{3, randStr, randomString(maxLenB, 0, 0), 0, errorCodes.DbNothingToReport},
		{4, randStr, randomString(maxLenB, 0, 0), 1, errorCodes.DbDuplicate},
		{1, randomString(maxLenS, 0, 0), randomString(maxLenB, 0, 0), 0, errorCodes.DbNothingToReport},
		{2, randomString(maxLenS, 0, 0), randomString(maxLenB, 0, 1), 0, errorCodes.DbNothingToReport},
		{2, randomString(maxLenS, 0, 1), randomString(maxLenB, 0, 0), 0, errorCodes.DbNothingToReport},
		{2, randomString(maxLenS, 0, 1), randomString(maxLenB, 0, 1), 0, errorCodes.DbNothingToReport},
		{5, randomString(maxLenS, 1, 1), randomString(maxLenB, 0, 0), 1, errorCodes.DbValueTooLong},
		{5, randomString(maxLenS, 1, 0), randomString(maxLenB, 0, 0), 1, errorCodes.DbValueTooLong},
		{5, randomString(maxLenS, 0, 0), randomString(maxLenB, 1, 1), 1, errorCodes.DbValueTooLong},
		{5, randomString(maxLenS, 0, 0), randomString(maxLenB, 1, 0), 1, errorCodes.DbValueTooLong},
		{0, randomString(maxLenS, 0, 0), randomString(maxLenB, 0, 0), 1, errorCodes.DbNothingUpdated},
		{-1, randomString(maxLenS, 0, 0), randomString(maxLenB, 0, 0), 1, errorCodes.DbNothingUpdated},
		{43, randomString(maxLenS, 0, 0), randomString(maxLenB, 0, 0), 1, errorCodes.DbNothingUpdated},
	}

	for _, v := range table {
		err, code := UpdateTag(v.id, v.name, v.descr)
		if v.res_is_error == 1 {
			if err == nil || code != v.res_code {
				t.Errorf("The tag %v should not be updated, but it was: %v, %v", v.id, err, code)
			}
		} else {
			if err != nil || code != v.res_code {
				t.Errorf("The tag %v should have been updated, but was not", v.id)
			}

			tag, _, _ := GetTag(v.id)
			if tag.Name != v.name || tag.Description != v.descr {
				t.Errorf("Expected value %v after update, got %v", v.name, tag.Name)
			}
		}
	}
}

// --- Users tests ---
func TestGetUser(t *testing.T) {
	cleanUpDb()

	table := []struct {
		userId       int
		res_is_error int
		res_code     int
		res          structs.User
	}{
		{1, 0, errorCodes.DbNothingToReport, AllUsers[1]},
		{2, 0, errorCodes.DbNothingToReport, AllUsers[2]},
		{0, 1, errorCodes.DbNoElement, structs.User{}},
		{-1, 1, errorCodes.DbNoElement, structs.User{}},
		{23, 1, errorCodes.DbNoElement, structs.User{}},
		{43, 1, errorCodes.DbNoElement, structs.User{}},
	}

	for _, v := range table {
		user, err, code := GetUser(v.userId)
		if v.res_is_error == 1 && err == nil {
			t.Errorf("Wrong result for case %v. Expected error, did not get it", v.userId)
		}

		if v.res_is_error == 0 && err != nil {
			t.Errorf("Wrong result for case %v. Expected nil, got error", v.userId)
		}

		if code != v.res_code || user.Id != v.res.Id || user.Nickname != v.res.Nickname ||
			user.Image != v.res.Image || user.About != v.res.About || user.Expertise != v.res.Expertise ||
			user.Followers_num != v.res.Followers_num || user.Following_num != v.res.Following_num ||
			user.Purchases_num != v.res.Purchases_num || user.Questions_num != v.res.Questions_num ||
			user.Answers_num != v.res.Answers_num {
			t.Errorf("Wrong result for case %v: \n Expected %v \n Got %v %v %v", v.userId, v, err, code, user)
		}
		if user.Id != 0 && user.Issued_at == nil {
			t.Errorf("Wrong result for case %v: Real User has an Issued_at date", v.userId)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	cleanUpDb()

	randStr := randomString(maxLenS, 0, 0)
	table := []struct {
		id           int
		nickname     string
		about        string
		res_is_error int
		res_code     int
	}{
		{2, "Marie Curie", randomString(maxLenB, 0, 0), 1, errorCodes.DbDuplicate},
		{3, "Nikola Tesla", randomString(maxLenB, 0, 0), 1, errorCodes.DbDuplicate},
		{3, randStr, randomString(maxLenB, 0, 0), 0, errorCodes.DbNothingToReport},
		{4, randStr, randomString(maxLenB, 0, 0), 1, errorCodes.DbDuplicate},
		{1, randomString(maxLenS, 0, 0), randomString(maxLenB, 0, 0), 0, errorCodes.DbNothingToReport},
		{2, randomString(maxLenS, 0, 0), randomString(maxLenB, 0, 1), 0, errorCodes.DbNothingToReport},
		{3, randomString(maxLenS, 0, 1), randomString(maxLenB, 0, 0), 0, errorCodes.DbNothingToReport},
		{4, randomString(maxLenS, 0, 1), randomString(maxLenB, 0, 1), 0, errorCodes.DbNothingToReport},
		{2, randomString(maxLenS, 1, 0), randomString(maxLenB, 0, 0), 1, errorCodes.DbValueTooLong},
		{5, randomString(maxLenS, 0, 0), randomString(maxLenB, 1, 0), 1, errorCodes.DbValueTooLong},
		{5, randomString(maxLenS, 1, 0), randomString(maxLenB, 1, 0), 1, errorCodes.DbValueTooLong},
		{0, randomString(maxLenS, 0, 0), randomString(maxLenS, 0, 0), 1, errorCodes.DbNothingUpdated},
		{-1, randomString(maxLenS, 0, 0), randomString(maxLenS, 0, 0), 1, errorCodes.DbNothingUpdated},
		{43, randomString(maxLenS, 0, 0), randomString(maxLenS, 0, 0), 1, errorCodes.DbNothingUpdated},
	}

	for _, v := range table {
		err, code := UpdateUser(v.id, v.nickname, v.about)
		if v.res_is_error == 1 {
			if err == nil || code != v.res_code {
				t.Errorf("User %v should not be updated, but it was: %v, %v", v.id, err, code)
			}
		} else {
			if err != nil || code != v.res_code {
				t.Errorf("User %v should have been updated, but was not", v.id)
			}

			user, _, _ := GetUser(v.id)
			if user.Nickname != v.nickname || user.About != v.about {
				t.Errorf("Expected value %v after update, got %v", v.nickname, user.Nickname)
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
	for _, v := range tableSuccess {
		followers, err, code := GetFollowers(v.id)
		if err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Expected to get followers, got a mistake %v %v", err, code)
		}

		if len(followers) != len(v.followers) {
			t.Errorf("Expected to get %v followers, got %v", len(v.followers), len(followers))
		}

		followerIds := make([]int, len(followers), len(followers))
		for i, v := range followers {
			followerIds[i] = v.Id
		}

		if !isSortedArrayEquivalentToArray(v.followers, followerIds) {
			t.Errorf("Followers are not equal %v %v", v.followers, followerIds)
		}
	}

	for _, id := range []int{0, 16, 52, -1} {
		followers, err, code := GetFollowers(id)
		if len(followers) != 0 || err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Should receive empty array with no errors. Received %v %v %v", followers, err, code)
		}
	}

	followers, _, _ := GetFollowers(7)
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
	for _, v := range tableSuccess {
		following, err, code := GetFollowing(v.id)
		if err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Expected to get following, got a mistake %v %v", err, code)
		}

		if len(following) != len(v.following) {
			t.Errorf("Expected to get %v following, got %v", len(v.following), len(following))
		}

		followingIds := make([]int, len(following), len(following))
		for i, v := range following {
			followingIds[i] = v.Id
		}

		if !isSortedArrayEquivalentToArray(v.following, followingIds) {
			t.Errorf("Followers are not equal %v %v", v.following, followingIds)
		}
	}

	for _, id := range []int{0, 16, 52, -1} {
		followers, err, code := GetFollowing(id)
		if len(followers) != 0 || err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Should receive empty array with no errors. Received %v %v %v", followers, err, code)
		}
	}

	following, _, _ := GetFollowing(6)
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
		res_is_error  int
		res_code      int
		followers_num int
		following_num int
	}{
		{1, 1, 1, errorCodes.FollowYourself, 0, 3},
		{1, 2, 1, errorCodes.DbDuplicate, 2, 3},
		{6, 2, 1, errorCodes.DbDuplicate, 2, 1},
		{0, 2, 1, errorCodes.DbForeignKeyViolation, 2, 0},
		{6, -1, 1, errorCodes.DbForeignKeyViolation, 0, 1},
		{10, 54, 1, errorCodes.DbForeignKeyViolation, 0, 0},
		{1, 6, 0, errorCodes.DbNothingToReport, 1, 4},
		{6, 1, 0, errorCodes.DbNothingToReport, 1, 2},
		{2, 4, 0, errorCodes.DbNothingToReport, 2, 1},
		{2, 6, 0, errorCodes.DbNothingToReport, 2, 2},
	}

	for _, v := range table {
		err, code := Follow(v.whoId, v.whomId)
		if v.res_is_error == 1 {
			if err == nil || code != v.res_code {
				t.Errorf("Expect follow to fail, got %v, %v", err, code)
			}
		} else {
			if err != nil || code != v.res_code {
				t.Errorf("Expect follow to happen, got %v, %v", err, code)
			}
		}

		followers, _, _ := GetFollowers(v.whomId)
		following, _, _ := GetFollowing(v.whoId)

		if len(followers) != v.followers_num || len(following) != v.following_num {
			t.Errorf("Number of followers and following in FOLLOWERS table is not right. Expect (%v, %v), got (%v, %v)", v.followers_num, v.following_num, len(followers), len(following))
		}

		u1, _, _ := GetUser(v.whoId)
		u2, _, _ := GetUser(v.whomId)
		if u2.Followers_num != v.followers_num || u1.Following_num != v.following_num {
			t.Errorf("Number of followers and following in USERS table is not right. Expect (%v, %v), got (%v, %v)", v.followers_num, v.following_num, u2.Followers_num, u1.Following_num)
		}
	}
}

func TestUnfollow(t *testing.T) {
	cleanUpDb()

	table := []struct {
		whoId         int
		whomId        int
		res_is_error  int
		res_code      int
		followers_num int
		following_num int
	}{
		{1, 1, 1, errorCodes.FollowYourself, 0, 3},
		{1, 5, 1, errorCodes.DbNothingUpdated, 0, 3},
		{6, 3, 1, errorCodes.DbNothingUpdated, 0, 1},
		{5, 4, 1, errorCodes.DbNothingUpdated, 1, 0},
		{-1, 4, 1, errorCodes.DbNothingUpdated, 1, 0},
		{10, 9, 1, errorCodes.DbNothingUpdated, 0, 0},
		{11, 19, 1, errorCodes.DbNothingUpdated, 0, 0},
		{1, 6, 1, errorCodes.DbNothingUpdated, 0, 3},
		{6, 2, 0, errorCodes.DbNothingToReport, 1, 0},
		{1, 4, 0, errorCodes.DbNothingToReport, 0, 2},
		{1, 7, 0, errorCodes.DbNothingToReport, 0, 1},
		{1, 2, 0, errorCodes.DbNothingToReport, 0, 0},
	}

	for _, v := range table {
		err, code := Unfollow(v.whoId, v.whomId)
		if v.res_is_error == 1 {
			if err == nil || code != v.res_code {
				t.Errorf("Expect unfollow to fail, got %v, %v", err, code)
			}
		} else {
			if err != nil || code != v.res_code {
				t.Errorf("Expect unfollow to happen, got %v, %v", err, code)
			}
		}

		followers, _, _ := GetFollowers(v.whomId)
		following, _, _ := GetFollowing(v.whoId)

		if len(followers) != v.followers_num || len(following) != v.following_num {
			t.Errorf("Number of followers and following in FOLLOWERS table is not right. Expect (%v, %v), got (%v, %v)", v.followers_num, v.following_num, len(followers), len(following))
		}

		u1, _, _ := GetUser(v.whoId)
		u2, _, _ := GetUser(v.whomId)
		if u2.Followers_num != v.followers_num || u1.Following_num != v.following_num {
			t.Errorf("Number of followers and following in USERS table is not right. Expect (%v, %v), got (%v, %v)", v.followers_num, v.following_num, u2.Followers_num, u1.Following_num)
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
		{"User_134", "email@gmail.com", "just_some_password", 10},
		{"stuff", "email@somemail.com", "anotherPa$$W0rt", 11},
		{"another", "random@yahoo.com", "anotherPa$$W0rt", 12},
	}
	for _, v := range tableSuccess {
		userId, err, code := CreateUser(v.nickname, v.email, v.password)
		if err != nil || code != errorCodes.DbNothingToReport || userId != v.userId {
			t.Errorf("Expected to create users. Got %v, %v, %v", userId, err, code)
		}
	}

	tableFail := []struct {
		nickname string
		email    string
		password string
		code     int
	}{
		{"exist", "albert@gmail.com", "password", errorCodes.DbDuplicate},
		{AllUsers[2].Nickname, "good@mail.com", "password", errorCodes.DbDuplicate},
		{tableSuccess[2].nickname, "amail@mail.com", "password", errorCodes.DbDuplicate},
		{"random", tableSuccess[2].email, "password", errorCodes.DbDuplicate},
	}
	for _, v := range tableFail {
		userId, err, code := CreateUser(v.nickname, v.email, v.password)
		if err == nil || code != v.code || userId != 0 {
			t.Errorf("Expected to fail. Got %v, %v, %v", userId, err, code)
		}
	}

	user, err, _ := GetUser(tableSuccess[0].userId)
	if err != nil || user.Nickname != tableSuccess[0].nickname || user.Id != tableSuccess[0].userId {
		t.Errorf("Expected to get user. Got %v, %v", err, user)
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
	for _, v := range tableSuccess {
		jwt, ok := Login(v.email, v.password)
		if !ok || len(jwt) < 10 {
			t.Errorf("Expected to log in. Got %v, %v", ok, jwt)
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
	for _, v := range tableFail {
		jwt, ok := Login(v.email, v.password)
		if ok || jwt != "" {
			t.Errorf("Expected to fail. Got %v, %v", ok, jwt)
		}
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

	for _, v := range tableSuccess {
		purchases, err, code := GetUserPurchases(v.userId)
		if err != nil || code != errorCodes.DbNothingToReport || len(purchases) != v.numPurchases {
			t.Errorf("Expected to see no errors and %v purchases. Got %v %v %v", v.numPurchases, err, code, len(purchases))
		}
	}

	purchases, _, _ := GetUserPurchases(4)
	p, o := purchases[0], AllPurchases[2]
	if p.Id != o.Id || p.Image != o.Image || p.Description != o.Description || p.Likes_num != o.Likes_num || p.User_id != o.User_id || p.Brand != o.Brand {
		t.Errorf("Purchase does not look right %v", p)
	}
}

func TestGetAllPurchases(t *testing.T) {
	cleanUpDb()

	purchases, err, code := GetAllPurchases()
	if err != nil || code != errorCodes.DbNothingToReport {
		t.Errorf("GetAllPurchases should succeed. Got %v %v", err, code)
	}

	if len(purchases) != len(AllPurchases) {
		t.Errorf("Expect to see %v purchases. Got %v", len(AllPurchases), len(purchases))
	}

	for _, v := range purchases {
		p := AllPurchases[v.Id]
		if p.Image != v.Image || p.Description != v.Description || p.User_id != v.User_id ||
			p.Brand != v.Brand || p.Likes_num != v.Likes_num {
			t.Errorf("Purchase %v does not look right. Expect %v, got %v", v.Id, v, p)
		}
	}
}

func TestGetPurchase(t *testing.T) {
	cleanUpDb()

	tableCorrect := []struct {
		id int
		p  structs.Purchase
	}{
		{1, AllPurchases[1]},
		{2, AllPurchases[2]},
		{3, AllPurchases[3]},
		{4, AllPurchases[4]},
	}

	for _, v := range tableCorrect {
		p, err, code := GetPurchase(v.id)
		if err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Expected correct execution. Got %v %v", err, code)
		}

		if p.Id != v.p.Id || p.Image != v.p.Image || p.Description != v.p.Description ||
			p.User_id != v.p.User_id || p.Brand != v.p.Brand || p.Likes_num != v.p.Likes_num {
			t.Errorf("Purchase looks different %v, %v", p, v.p)
		}
	}

	for _, v := range []int{0, -1, 6, 10} {
		p, err, code := GetPurchase(v)
		if err == nil || code != errorCodes.DbNoElement || p.Id != 0 || p.Image != "" {
			t.Errorf("Expected to get error. Got %v, %v, %v", err, code, p)
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

	for _, v := range tableSuccess {
		purchases, err, code := GetAllPurchasesWithBrand(v.brandId)
		if err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Expect correct execution. Got %v %v", err, code)
		}

		if len(purchases) != len(v.purchaseIds) {
			t.Errorf("Expected %v purchases. Got %v", len(v.purchaseIds), len(purchases))
		}

		if len(purchases) > 0 {
			for _, p := range purchases {
				expected := AllPurchases[p.Id]
				if !v.purchaseIds[p.Id] {
					t.Errorf("Not expected purchase with Id %v", p.Id)
				}
				if expected.Id != p.Id || expected.Image != p.Image || expected.Description != p.Description ||
					expected.Likes_num != p.Likes_num || expected.Brand != p.Brand {
					t.Errorf("Expected %v. Got %v", expected, p)
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

	for _, v := range tableSuccess {
		purchases, err, code := GetAllPurchasesWithTag(v.tagId)
		if err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Expect correct execution. Got %v %v", err, code)
		}

		if len(purchases) != len(v.purchaseIds) {
			t.Errorf("Expected %v purchases. Got %v", len(v.purchaseIds), len(purchases))
		}

		if len(purchases) > 0 {
			for _, p := range purchases {
				expected := AllPurchases[p.Id]
				if !v.purchaseIds[p.Id] {
					t.Errorf("Not expected purchase with Id %v", p.Id)
				}
				if expected.Id != p.Id || expected.Image != p.Image || expected.Description != p.Description ||
					expected.Likes_num != p.Likes_num || expected.Brand != p.Brand {
					t.Errorf("Expected %v. Got %v", expected, p)
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
	for _, v := range tableSuccess {
		err, code := LikePurchase(v.purchaseId, v.userId)
		if err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Expect correct execution. Got %v %v", err, code)
		}

		p, _, _ := GetPurchase(v.purchaseId)
		if p.Likes_num != v.likesNum {
			t.Errorf("Expect to see %v likes. Got %v", v.likesNum, p.Likes_num)
		}
	}

	tableFail := []struct {
		purchaseId int
		userId     int
		code       int
		likesNum   int
	}{
		{1, 1, errorCodes.DbVoteForOwnStuff, 3},
		{3, 4, errorCodes.DbDuplicate, 4},
		{3, 2, errorCodes.DbDuplicate, 4},
		{0, 2, errorCodes.DbNoPurchase, 0},
		{9, 1, errorCodes.DbNoPurchase, 0},
		{1, 10, errorCodes.DbForeignKeyViolation, 3},
	}
	for _, v := range tableFail {
		err, code := LikePurchase(v.purchaseId, v.userId)
		if err == nil || code != v.code {
			t.Errorf("Expect to fail. Got %v %v", err, code)
		}

		p, _, _ := GetPurchase(v.purchaseId)
		if p.Likes_num != v.likesNum {
			t.Errorf("Expect to see %v likes. Got %v", v.likesNum, p.Likes_num)
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
	for _, v := range tableSuccess {
		err, code := UnlikePurchase(v.purchaseId, v.userId)
		if err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Expect correct execution. Got %v %v", err, code)
		}

		p, _, _ := GetPurchase(v.purchaseId)
		if p.Likes_num != v.likesNum {
			t.Errorf("Expect to see %v likes. Got %v", v.likesNum, p.Likes_num)
		}
	}

	tableFail := []struct {
		purchaseId int
		userId     int
		code       int
		likesNum   int
	}{
		{3, 1, errorCodes.DbVoteForOwnStuff, 1},
		{3, 4, errorCodes.DbNothingUpdated, 1},
		{1, 7, errorCodes.DbNothingUpdated, 0},
		{3, 9, errorCodes.DbNothingUpdated, 1},
		{2, 3, errorCodes.DbNothingUpdated, 0},
		{-2, -1, errorCodes.DbNoPurchase, 0},
		{9, 2, errorCodes.DbNoPurchase, 0},
		{3, -1, errorCodes.DbNothingUpdated, 1},
	}
	for _, v := range tableFail {
		err, code := UnlikePurchase(v.purchaseId, v.userId)
		if err == nil || code != v.code {
			t.Errorf("Expect to fail. Got %v %v", err, code)
		}

		p, _, _ := GetPurchase(v.purchaseId)
		if p.Likes_num != v.likesNum {
			t.Errorf("Expect to see %v likes. Got %v", v.likesNum, p.Likes_num)
		}
	}
}
