package storage

import (
	"../../unnamed/config"
	"../../unnamed/errorCodes"
	"../../unnamed/structs"
	"log"
	"os"
	"os/exec"
	"testing"
	"math/rand"
	"time"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

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
		43: testEl{true, errorCodes.DbNoElement, structs.Brand{}},
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
		name string
		res_id int
	}

	correct_table := []testEl{
		testEl{randomString(1), 6},
		testEl{randomString(6), 7},
		testEl{randomString(16), 8},
		testEl{randomString(40), 9},
		testEl{randomString(39), 10},
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

	wrong_table := []testEl{
		testEl{randomString(41), errorCodes.DbValueTooLong},
		testEl{randomString(56), errorCodes.DbValueTooLong},
		testEl{correct_table[0].name, errorCodes.DbDuplicate},
		testEl{correct_table[1].name, errorCodes.DbDuplicate},
		testEl{correct_table[2].name, errorCodes.DbDuplicate},
		testEl{correct_table[3].name, errorCodes.DbDuplicate},
	}

	for _, val := range wrong_table {
		id, err, code := CreateBrand(val.name)
		if err == nil || id != 0 || code !=val.res_id{
			t.Error("New brand should not be created")
		}
	}

	brands, _, _ := GetAllBrands()
	if len(brands) != 10 {
		t.Error("It looks like a couple of brands were created, but they should not")
	}
}