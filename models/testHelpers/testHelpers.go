// testHelpers is a package with helper functions to test all models
package testHelpers

import (
	"../../config"
	"../../mailer"
	"../../misc"
	"../../psql"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"time"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

var AllPurchases = map[int]misc.Purchase{
	1: {1, "1467954439_isForTests.jpg", "Look at my new drone", 1, 0, []int{2}, 0, 0},
	2: {2, "1467954439_isForTests.jpg", "How cool am I?", 4, 0, []int{3, 5}, 5, 0},
	3: {3, "1467954439_isForTests.jpg", "I really like drones", 1, 0, []int{4}, 0, 3},
	4: {4, "1467954439_isForTests.jpg", "Now I am fond of cars", 1, 0, []int{2}, 4, 1},
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

// RandomString generates a random string of a specific length
func RandomString(length int, isBiggerI, isEdgeCaseI int) string {
	// This length can be bigger or smaller than you predefined. Also you can ask it to be on the edge
	// of the allowed values. For example if you want a value bigger than X, it will generate you
	// some strings of the length X + 1 or bigger (if on the edge it will be only X + 1)
	// If you want smaller than X, it will generate you anything less or equal to X. If on the edge, it
	// will be equal to X
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

func IsSortedArrayEquivalentToArray(arrSorted, arr []int) bool {
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

func IsApproximatelyNow(createdTime int64) bool {
	// sometimes time can be one off. This causes a lot of confusion in the tests.
	timeNow := time.Now().Unix()
	return timeNow == createdTime || timeNow-1 == createdTime
}

func InitAll() {
	// initialize Db connection
	config.Init()
	mailer.Init()
	psql.Init()
}

func CleanUpDb() {
	// prepare database by creating tables and populating it with data
	cmd := exec.Command("../../SQL/set_up_database.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if cmd.Run() != nil {
		log.Fatal("Can't prepare SQL database")
	}
}
