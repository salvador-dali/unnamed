package storage

import (
	"../../unnamed/config"
	"../../unnamed/errorCodes"
	"../../unnamed/structs"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"testing"
	"time"
)

const (
	letterBytes  = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	maxLenSmall  = 40
	maxLenMedium = 100
	maxLenBig    = 1000
)

func randomString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func setup() {
	// prepare database by creating tables and populating it with data
	cmd := exec.Command("../SQL/set_up_database.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if cmd.Run() != nil {
		log.Fatal("Can't prepare SQL database")
	}

	// initialize Db connection
	cnf := config.Init()
	Init(cnf.DbUser, cnf.DbPass, cnf.DbHost, cnf.DbName, cnf.DbPort)

	// initialize randomness
	rand.Seed(time.Now().UTC().UnixNano())
}

// Setup and db.close will be called before and after each test http://stackoverflow.com/a/34102842/1090562
func TestMain(m *testing.M) {
	setup()
	retCode := m.Run()
	defer Db.Close()
	os.Exit(retCode)
}

// --- Brands tests ---
func TestGetAllBrands(t *testing.T) {
	brands, err, code := GetAllBrands()
	if err != nil || code != errorCodes.DbNothingToReport {
		t.Error("Should finish without no error")
	}

	if len(brands) != 5 {
		t.Errorf("Expect 5 brands, got %v", len(brands))
	}

	expected := []structs.Brand{
		structs.Brand{1, "Apple", nil},
		structs.Brand{2, "BMW", nil},
		structs.Brand{3, "Playstation", nil},
		structs.Brand{4, "Ferrari", nil},
		structs.Brand{5, "Gucci", nil},
	}
	for i, brand := range brands {
		if brand.Id != expected[i].Id || brand.Name != expected[i].Name || brand.Issued_at != expected[i].Issued_at {
			t.Errorf("Expect %v , got %v", expected[i], brand)
		}
	}
}

func TestGetBrand(t *testing.T) {
	type testEl struct {
		res_is_error bool
		res_code     int
		res          structs.Brand
	}

	table := map[int]testEl{
		1:   testEl{false, errorCodes.DbNothingToReport, structs.Brand{1, "Apple", nil}},
		2:   testEl{false, errorCodes.DbNothingToReport, structs.Brand{2, "BMW", nil}},
		3:   testEl{false, errorCodes.DbNothingToReport, structs.Brand{3, "Playstation", nil}},
		5:   testEl{false, errorCodes.DbNothingToReport, structs.Brand{5, "Gucci", nil}},
		0:   testEl{true, errorCodes.DbNoElement, structs.Brand{}},
		-1:  testEl{true, errorCodes.DbNoElement, structs.Brand{}},
		123: testEl{true, errorCodes.DbNoElement, structs.Brand{}},
		43:  testEl{true, errorCodes.DbNoElement, structs.Brand{}},
	}

	for id, val := range table {
		brand, err, code := GetBrand(id)
		if val.res_is_error && err == nil {
			t.Errorf("Wrong result for case %v. Expected error, did not get it", id)
		}

		if !val.res_is_error && err != nil {
			t.Errorf("Wrong result for case %v. Expected nil, got error", id)
		}

		if code != val.res_code || brand.Id != val.res.Id || brand.Name != val.res.Name {
			t.Errorf("Wrong result for case %v: \n Expected %v \n Got %v %v %v", id, val, err, code, brand)
		}
		if brand.Id != 0 && brand.Issued_at == nil {
			t.Errorf("Wrong result for case %v: Real Brand has an Issued_at date", id)
		}
	}
}

func TestCreateBrand(t *testing.T) {
	type testEl struct {
		name   string
		res_id int
	}

	correct_table := []testEl{
		testEl{randomString(1), 6},
		testEl{randomString(6), 7},
		testEl{randomString(16), 8},
		testEl{randomString(maxLenSmall), 9},
		testEl{randomString(39), 10},
	}

	wrong_table := []testEl{
		testEl{randomString(maxLenSmall + 1), errorCodes.DbValueTooLong},
		testEl{randomString(56), errorCodes.DbValueTooLong},
		testEl{correct_table[0].name, errorCodes.DbDuplicate},
		testEl{correct_table[1].name, errorCodes.DbDuplicate},
		testEl{correct_table[2].name, errorCodes.DbDuplicate},
		testEl{correct_table[3].name, errorCodes.DbDuplicate},
	}

	for _, val := range correct_table {
		id, err, code := CreateBrand(val.name)
		if id != val.res_id || err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Expected to create a brand. Got %v %v %v", id, err, code)
		}

		brand, err, code := GetBrand(id)
		if brand.Name != val.name {
			t.Errorf("Expected to create a brand with a name %v, got %v", brand.Name, val.name)
		}
	}

	for _, val := range wrong_table {
		id, err, code := CreateBrand(val.name)
		if err == nil || id != 0 || code != val.res_id {
			t.Error("New brand should not be created")
		}
	}

	brands, _, _ := GetAllBrands()
	if len(brands) != 10 {
		t.Errorf("Should have 10 brands, have %v", len(brands))
	}
}

func TestUpdateBrand(t *testing.T) {
	type testEl struct {
		id           int
		name         string
		res_is_error bool
		res_code     int
	}

	randStr := randomString(maxLenSmall)
	table := []testEl{
		testEl{1, randomString(1), false, errorCodes.DbNothingToReport},
		testEl{2, "Playstation", true, errorCodes.DbDuplicate},
		testEl{3, "Ferrari", true, errorCodes.DbDuplicate},
		testEl{2, randomString(6), false, errorCodes.DbNothingToReport},
		testEl{3, randStr, false, errorCodes.DbNothingToReport},
		testEl{4, randStr, true, errorCodes.DbDuplicate},
		testEl{2, randomString(maxLenSmall + 1), true, errorCodes.DbValueTooLong},
		testEl{5, randomString(141), true, errorCodes.DbValueTooLong},
		testEl{0, randomString(10), true, errorCodes.DbNothingUpdated},
		testEl{-1, randomString(10), true, errorCodes.DbNothingUpdated},
		testEl{43, randomString(10), true, errorCodes.DbNothingUpdated},
		testEl{123, randomString(10), true, errorCodes.DbNothingUpdated},
	}

	for _, val := range table {
		err, code := UpdateBrand(val.id, val.name)
		if val.res_is_error {
			if err == nil || code != val.res_code {
				t.Errorf("The brand %v should not be updated, but it was: %v, %v", val.id, err, code)
			}
		} else {
			if err != nil || code != val.res_code {
				t.Errorf("The brand %v should have been updated, but was not", val.id)
			}

			brand, _, _ := GetBrand(val.id)
			if brand.Name != val.name {
				t.Errorf("Expected value %v after update, got %v", val.name, brand.Name)
			}
		}
	}
}

// --- Tags tests ---
func TestGetAllTags(t *testing.T) {
	tags, err, code := GetAllTags()
	if err != nil || code != errorCodes.DbNothingToReport {
		t.Error("Should finish without no error")
	}

	if len(tags) != 6 {
		t.Errorf("Expect 6 tags, got %v", len(tags))
	}

	expected := []structs.Tag{
		structs.Tag{1, "dress", "", nil},
		structs.Tag{2, "drone", "", nil},
		structs.Tag{3, "cosmetics", "", nil},
		structs.Tag{4, "car", "", nil},
		structs.Tag{5, "hat", "", nil},
		structs.Tag{6, "phone", "", nil},
	}
	for i, tag := range tags {
		if tag.Id != expected[i].Id || tag.Name != expected[i].Name || tag.Issued_at != expected[i].Issued_at || tag.Description != expected[i].Description {
			t.Errorf("Expect %v , got %v", expected[i], tag)
		}
	}
}

func TestGetTag(t *testing.T) {
	type testEl struct {
		res_is_error bool
		res_code     int
		res          structs.Tag
	}

	table := map[int]testEl{
		1:   testEl{false, errorCodes.DbNothingToReport, structs.Tag{1, "dress", "nice dresses", nil}},
		2:   testEl{false, errorCodes.DbNothingToReport, structs.Tag{2, "drone", "cool flying machines that do stuff", nil}},
		3:   testEl{false, errorCodes.DbNothingToReport, structs.Tag{3, "cosmetics", "Known as make-up, are substances or products used to enhance the appearance or scent of the body", nil}},
		6:   testEl{false, errorCodes.DbNothingToReport, structs.Tag{6, "phone", "People use it to speak with other people", nil}},
		0:   testEl{true, errorCodes.DbNoElement, structs.Tag{}},
		-1:  testEl{true, errorCodes.DbNoElement, structs.Tag{}},
		123: testEl{true, errorCodes.DbNoElement, structs.Tag{}},
		43:  testEl{true, errorCodes.DbNoElement, structs.Tag{}},
	}

	for id, val := range table {
		tag, err, code := GetTag(id)
		if val.res_is_error && err == nil {
			t.Errorf("Wrong result for case %v. Expected error, did not get it", id)
		}

		if !val.res_is_error && err != nil {
			t.Errorf("Wrong result for case %v. Expected nil, got error", id)
		}

		if code != val.res_code || tag.Id != val.res.Id || tag.Name != val.res.Name || tag.Description != val.res.Description {
			t.Errorf("Wrong result for case %v: \n Expected %v \n Got %v %v %v", id, val, err, code, tag)
		}
		if tag.Id != 0 && tag.Issued_at == nil {
			t.Errorf("Wrong result for case %v: Real Brand has an Issued_at date", id)
		}
	}
}

func TestCreateTag(t *testing.T) {
	type testEl struct {
		name   string
		descr  string
		res_id int
	}

	correct_table := []testEl{
		testEl{randomString(1), randomString(23), 7},
		testEl{randomString(6), randomString(245), 8},
		testEl{randomString(16), randomString(643), 9},
		testEl{randomString(maxLenSmall), randomString(maxLenBig), 10},
		testEl{randomString(39), randomString(1), 11},
	}

	wrong_table := []testEl{
		testEl{randomString(maxLenSmall + 1), randomString(41), errorCodes.DbValueTooLong},
		testEl{randomString(6), randomString(maxLenBig + 1), errorCodes.DbValueTooLong},
		testEl{randomString(64), randomString(1201), errorCodes.DbValueTooLong},
		testEl{correct_table[0].name, "", errorCodes.DbDuplicate},
		testEl{correct_table[1].name, "", errorCodes.DbDuplicate},
		testEl{correct_table[2].name, "", errorCodes.DbDuplicate},
		testEl{correct_table[3].name, "", errorCodes.DbDuplicate},
	}

	for _, val := range correct_table {
		id, err, code := CreateTag(val.name, val.descr)
		if id != val.res_id || err != nil || code != errorCodes.DbNothingToReport {
			t.Errorf("Expected to create a tag. Got %v %v %v", id, err, code)
		}

		tag, _, _ := GetTag(id)
		if tag.Name != val.name || tag.Description != val.descr {
			t.Errorf("Expected to create a tag with a name %v, got %v", tag.Name, val.name)
		}
	}

	for _, val := range wrong_table {
		id, err, code := CreateTag(val.name, val.descr)
		if err == nil || id != 0 || code != val.res_id {
			t.Errorf("New tag should not be created %v, %v", id, code)
		}
	}

	tags, _, _ := GetAllTags()
	if len(tags) != 11 {
		t.Errorf("Should have 11 tags, have %v", len(tags))
	}
}

func TestUpdateTag(t *testing.T) {
	type testEl struct {
		id           int
		name         string
		descr        string
		res_is_error bool
		res_code     int
	}

	randStr := randomString(maxLenSmall)
	table := []testEl{
		testEl{1, randomString(1), randomString(1), false, errorCodes.DbNothingToReport},
		testEl{2, "car", randomString(159), true, errorCodes.DbDuplicate},
		testEl{3, "phone", randomString(10), true, errorCodes.DbDuplicate},
		testEl{2, randomString(34), randomString(999), false, errorCodes.DbNothingToReport},
		testEl{3, randStr, randomString(989), false, errorCodes.DbNothingToReport},
		testEl{4, randStr, randomString(123), true, errorCodes.DbDuplicate},
		testEl{2, randomString(maxLenSmall + 1), randomString(23), true, errorCodes.DbValueTooLong},
		testEl{5, randomString(31), randomString(maxLenBig + 1), true, errorCodes.DbValueTooLong},
		testEl{5, randomString(maxLenSmall + 1), randomString(maxLenBig + 1), true, errorCodes.DbValueTooLong},
		testEl{0, randomString(10), randomString(10), true, errorCodes.DbNothingUpdated},
		testEl{-1, randomString(10), randomString(10), true, errorCodes.DbNothingUpdated},
		testEl{43, randomString(10), randomString(10), true, errorCodes.DbNothingUpdated},
		testEl{123, randomString(10), randomString(10), true, errorCodes.DbNothingUpdated},
	}

	for _, val := range table {
		err, code := UpdateTag(val.id, val.name, val.descr)
		if val.res_is_error {
			if err == nil || code != val.res_code {
				t.Errorf("The tag %v should not be updated, but it was: %v, %v", val.id, err, code)
			}
		} else {
			if err != nil || code != val.res_code {
				t.Errorf("The tag %v should have been updated, but was not", val.id)
			}

			tag, _, _ := GetTag(val.id)
			if tag.Name != val.name || tag.Description != val.descr {
				t.Errorf("Expected value %v after update, got %v", val.name, tag.Name)
			}
		}
	}
}

// --- Users tests ---
func TestGetUser(t *testing.T) {
	type testEl struct {
		res_is_error bool
		res_code     int
		res          structs.User
	}

	table := map[int]testEl{
		1:   testEl{false, errorCodes.DbNothingToReport, structs.User{1, "Albert Einstein", "", "Developed the general theory of relativity.", 0, 0, 3, 3, 0, 1, nil}},
		2:   testEl{false, errorCodes.DbNothingToReport, structs.User{2, "Isaac Newton", "", "Mechanics, laws of motion", 0, 2, 0, 0, 0, 0, nil}},
		0:   testEl{true, errorCodes.DbNoElement, structs.User{}},
		-1:  testEl{true, errorCodes.DbNoElement, structs.User{}},
		123: testEl{true, errorCodes.DbNoElement, structs.User{}},
		43:  testEl{true, errorCodes.DbNoElement, structs.User{}},
	}

	for id, val := range table {
		user, err, code := GetUser(id)
		if val.res_is_error && err == nil {
			t.Errorf("Wrong result for case %v. Expected error, did not get it", id)
		}

		if !val.res_is_error && err != nil {
			t.Errorf("Wrong result for case %v. Expected nil, got error", id)
		}

		if code != val.res_code || user.Id != val.res.Id || user.Nickname != val.res.Nickname ||
			user.Image != val.res.Image || user.About != val.res.About || user.Expertise != val.res.Expertise ||
			user.Followers_num != val.res.Followers_num || user.Following_num != val.res.Following_num ||
			user.Purchases_num != val.res.Purchases_num || user.Questions_num != val.res.Questions_num ||
			user.Answers_num != val.res.Answers_num {
			t.Errorf("Wrong result for case %v: \n Expected %v \n Got %v %v %v", id, val, err, code, user)
		}
		if user.Id != 0 && user.Issued_at == nil {
			t.Errorf("Wrong result for case %v: Real User has an Issued_at date", id)
		}
	}
}

func TestUpdateUser(t *testing.T) {
	type testEl struct {
		id           int
		nickname     string
		about        string
		res_is_error bool
		res_code     int
	}

	randStr := randomString(maxLenSmall)
	table := []testEl{
		testEl{1, randomString(10), randomString(531), false, errorCodes.DbNothingToReport},
		testEl{2, "Marie Curie", randomString(159), true, errorCodes.DbDuplicate},
		testEl{3, "Nikola Tesla", randomString(10), true, errorCodes.DbDuplicate},
		testEl{2, randomString(34), randomString(999), false, errorCodes.DbNothingToReport},
		testEl{3, randStr, randomString(989), false, errorCodes.DbNothingToReport},
		testEl{4, randStr, randomString(123), true, errorCodes.DbDuplicate},
		testEl{2, randomString(maxLenSmall + 1), randomString(23), true, errorCodes.DbValueTooLong},
		testEl{5, randomString(31), randomString(maxLenBig + 1), true, errorCodes.DbValueTooLong},
		testEl{0, randomString(10), randomString(10), true, errorCodes.DbNothingUpdated},
		testEl{-1, randomString(10), randomString(10), true, errorCodes.DbNothingUpdated},
		testEl{43, randomString(10), randomString(10), true, errorCodes.DbNothingUpdated},
		testEl{123, randomString(10), randomString(10), true, errorCodes.DbNothingUpdated},
	}

	for _, val := range table {
		err, code := UpdateUser(val.id, val.nickname, val.about)
		if val.res_is_error {
			if err == nil || code != val.res_code {
				t.Errorf("User %v should not be updated, but it was: %v, %v", val.id, err, code)
			}
		} else {
			if err != nil || code != val.res_code {
				t.Errorf("User %v should have been updated, but was not", val.id)
			}

			user, _, _ := GetUser(val.id)
			if user.Nickname != val.nickname || user.About != val.about {
				t.Errorf("Expected value %v after update, got %v", val.nickname, user.Nickname)
			}
		}
	}
}
