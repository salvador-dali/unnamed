package user

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

func TestShowById(t *testing.T) {
	o.CleanUpDb()

	table := []struct {
		userId int
		code   int
		user   misc.User
	}{
		{1, misc.NothingToReport, o.AllUsers[1]},
		{2, misc.NothingToReport, o.AllUsers[2]},
		{0, misc.NoElement, misc.User{}},
		{-1, misc.NoElement, misc.User{}},
		{23, misc.NoElement, misc.User{}},
		{43, misc.NoElement, misc.User{}},
	}
	for num, v := range table {
		user, code := ShowById(v.userId)
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

func TestUpdate(t *testing.T) {
	o.CleanUpDb()

	randStr := o.RandomString(misc.MaxLenS, 0, 0)
	table := []struct {
		id       int
		nickname string
		about    string
		code     int
	}{
		{2, "Marie Curie", o.RandomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{3, "Nikola Tesla", o.RandomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{3, randStr, o.RandomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{4, randStr, o.RandomString(misc.MaxLenB, 0, 0), misc.DbDuplicate},
		{1, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{2, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 0, 1), misc.NothingToReport},
		{3, o.RandomString(misc.MaxLenS, 0, 1), o.RandomString(misc.MaxLenB, 0, 0), misc.NothingToReport},
		{4, o.RandomString(misc.MaxLenS, 0, 1), o.RandomString(misc.MaxLenB, 0, 1), misc.NothingToReport},
		{2, o.RandomString(misc.MaxLenS, 1, 0), o.RandomString(misc.MaxLenB, 0, 0), misc.WrongName},
		{5, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenB, 1, 0), misc.WrongDescr},
		{5, o.RandomString(misc.MaxLenS, 1, 0), o.RandomString(misc.MaxLenB, 1, 0), misc.WrongName},
		{0, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
		{-1, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
		{43, o.RandomString(misc.MaxLenS, 0, 0), o.RandomString(misc.MaxLenS, 0, 0), misc.NothingUpdated},
	}
	for num, v := range table {
		code := Update(v.id, v.nickname, v.about)
		if code != v.code {
			t.Errorf("Case %v. Expect %v. Got %v", num, v.code, code)
		}

		if code == misc.NothingToReport {
			user, _ := ShowById(v.id)
			if user.Nickname != v.nickname || user.About != v.about {
				t.Errorf("Case %v. Expect %v %v. Got %v", num, v.nickname, v.about, user)
			}
		}
	}
}

func TestGetFollowers(t *testing.T) {
	o.CleanUpDb()

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

		if !o.IsSortedArrayEquivalentToArray(v.followers, followerIds) {
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
	if u.Nickname != o.AllUsers[1].Nickname || u.About != "" || u.Expertise != 0 || u.Followers_num != 0 {
		t.Errorf("Information about follower is not right %v", u)
	}
}

func TestGetFollowing(t *testing.T) {
	o.CleanUpDb()

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

		if !o.IsSortedArrayEquivalentToArray(v.following, followingIds) {
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
	if u.Nickname != o.AllUsers[2].Nickname || u.About != "" || u.Expertise != 0 || u.Followers_num != 0 {
		t.Errorf("Information about following is not right %v", u)
	}
}

func TestFollow(t *testing.T) {
	o.CleanUpDb()

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

		u1, _ := ShowById(v.whoId)
		u2, _ := ShowById(v.whomId)
		if u2.Followers_num != v.followers_num || u1.Following_num != v.following_num {
			t.Errorf("Case %v. Expect (%v, %v). Got (%v, %v)", num, v.followers_num, v.following_num, u2.Followers_num, u1.Following_num)
		}
	}
}

func TestUnfollow(t *testing.T) {
	o.CleanUpDb()

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

		u1, _ := ShowById(v.whoId)
		u2, _ := ShowById(v.whomId)
		if u2.Followers_num != v.followers_num || u1.Following_num != v.following_num {
			t.Errorf("Case %v. Expect (%v, %v), got (%v, %v)", num, v.followers_num, v.following_num, u2.Followers_num, u1.Following_num)
		}
	}
}

func TestCreate(t *testing.T) {
	o.CleanUpDb()

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
		userId, code := Create(v.nickname, v.email, v.password)
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
		{o.AllUsers[2].Nickname, "good@mail.com", "password", misc.DbDuplicate},
		{tableSuccess[2].nickname, "amail@mail.com", "password", misc.DbDuplicate},
		{"random", tableSuccess[2].email, "password", misc.DbDuplicate},
	}
	for num, v := range tableFail {
		userId, code := Create(v.nickname, v.email, v.password)
		if code != v.code || userId != 0 {
			t.Errorf("Case %v. Expect 0, %v. Got %v, %v", num, v.code, userId, code)
		}
	}

	user, code := ShowById(tableSuccess[0].userId)
	if user.Nickname != tableSuccess[0].nickname || user.Id != tableSuccess[0].userId {
		t.Errorf("Expected to get user. Got %v, %v", user, code)
	}
}

func TestLogin(t *testing.T) {
	o.CleanUpDb()

	email, pass := "some_strange_mail@gmail.com", "very_new_password"
	Create("username", email, pass)

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
	o.CleanUpDb()

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
		if _, ok := VerifyEmail(v.userId, v.verifyCode); ok {
			t.Errorf("Case %v. Expect to fail. Got True", num)
		}
	}

	if _, ok := VerifyEmail(10, "pqaJaBRgAvzLXqzRrrUI"); !ok {
		t.Errorf("Expect to verify email. Got False")
	}
}
