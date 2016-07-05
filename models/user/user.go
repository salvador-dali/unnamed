package user

import (
	"../../auth"
	"../../mailer"
	"../../misc"
	"../../psql"
	"database/sql"
	"log"
	"reflect"
	"time"
)

// Show user information by Id
func ShowById(userId int) (misc.User, int) {
	if !misc.IsIdValid(userId) {
		log.Println("UserId is not correct")
		return misc.User{}, misc.NoElement
	}

	user := misc.User{}
	var timestamp time.Time
	if err := psql.Db.QueryRow(`
		SELECT nickname, image, about, expertise, followers_num, following_num, purchases_num, questions_num, answers_num, issued_at
		FROM users WHERE id = $1`, userId,
	).Scan(
		&user.Nickname, &user.Image, &user.About, &user.Expertise, &user.Followers_num,
		&user.Following_num, &user.Purchases_num, &user.Questions_num, &user.Answers_num,
		&timestamp,
	); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return misc.User{}, misc.NoElement
		}

		log.Println(err)
		return misc.User{}, misc.NothingToReport
	}
	user.Id = userId
	user.Issued_at = timestamp.Unix()
	return user, misc.NothingToReport
}

// Update information about a user
func Update(userId int, nickname, about string) int {
	if !misc.IsIdValid(userId) {
		log.Println("user was not updated", userId)
		return misc.NothingUpdated
	}

	nickname, ok := misc.ValidateString(nickname, misc.MaxLenS)
	if !ok {
		log.Println("Nickname is not correct", nickname)
		return misc.WrongName
	}

	about, ok = misc.ValidateString(about, misc.MaxLenB)
	if !ok {
		log.Println("About is not correct", about)
		return misc.WrongDescr
	}

	sqlResult, err := psql.Db.Exec(`
		UPDATE users
		SET nickname = $1, about = $2
		WHERE id = $3`, nickname, about, userId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	err, code := psql.IsAffectedOneRow(sqlResult)
	return code
}

// Follow a user by Id
func Follow(whoId, whomId int) int {
	if !misc.IsIdValid(whomId) {
		log.Println("User id is not correct", whomId)
		return misc.NoElement
	}

	if whoId == whomId {
		log.Println("can't follow yourself")
		return misc.FollowYourself
	}

	sqlResult, err := psql.Db.Exec(`
		INSERT INTO followers (who_id, whom_id)
		VALUES ($1, $2)`, whoId, whomId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}
	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = psql.Db.Exec(`
		UPDATE users
		SET followers_num = followers_num + 1
		WHERE id = $1`, whomId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = psql.Db.Exec(`
		UPDATE users
		SET following_num = following_num + 1
		WHERE id = $1`, whoId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return code
	}

	return misc.NothingToReport
}

// Unfollow a user whom you previously followed
func Unfollow(whoId, whomId int) int {
	if !misc.IsIdValid(whomId) {
		log.Println("User id is not correct", whomId)
		return misc.NoElement
	}

	if whoId == whomId {
		log.Println("can't follow yourself")
		return misc.FollowYourself
	}

	sqlResult, err := psql.Db.Exec(`
		DELETE FROM followers
		WHERE who_id = $1 AND whom_id = $2`, whoId, whomId)
	if err != nil {
		log.Println(err)
		return misc.NothingToReport
	}

	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = psql.Db.Exec(`
		UPDATE users
		SET followers_num = followers_num - 1
		WHERE id = $1`, whomId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = psql.Db.Exec(`
		UPDATE users
		SET following_num = following_num - 1
		WHERE id = $1`, whoId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return code
	}

	return misc.NothingToReport
}

// GetFollowing returns a list of users whom a user with Id follows
func GetFollowing(userId int) ([]*misc.User, int) {
	if !misc.IsIdValid(userId) {
		return []*misc.User{}, misc.NothingToReport
	}

	rows, err := psql.Db.Query(`
		SELECT id, nickname, image
		FROM users
		WHERE id IN (
			SELECT whom_id
			FROM followers
			WHERE who_id = $1
		)`, userId)
	if err != nil {
		log.Println(err)
		return []*misc.User{}, misc.NothingToReport
	}
	defer rows.Close()

	users := []*misc.User{}
	for rows.Next() {
		user := misc.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			log.Println(err)
			return []*misc.User{}, misc.NothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return []*misc.User{}, misc.NothingToReport
	}

	return users, misc.NothingToReport
}

// GetFollowers returns a list of users who follow a user with Id
func GetFollowers(userId int) ([]*misc.User, int) {
	if !misc.IsIdValid(userId) {
		return []*misc.User{}, misc.NothingToReport
	}

	rows, err := psql.Db.Query(`
		SELECT id, nickname, image
		FROM users
		WHERE id IN (
			SELECT who_id
			FROM followers
			WHERE whom_id = $1
		)`, userId)
	if err != nil {
		log.Println(err)
		return []*misc.User{}, misc.NothingToReport
	}
	defer rows.Close()

	users := []*misc.User{}
	for rows.Next() {
		user := misc.User{}
		if err := rows.Scan(&user.Id, &user.Nickname, &user.Image); err != nil {
			log.Println(err)
			return []*misc.User{}, misc.NothingToReport
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return []*misc.User{}, misc.NothingToReport
	}

	return users, misc.NothingToReport
}

// Create a new user, sends him a confirmation email
func Create(nickname, email, password string) (int, int) {
	nickname, ok := misc.ValidateString(nickname, misc.MaxLenS)
	if !ok {
		log.Println("Wrong nickname", nickname)
		return 0, misc.WrongName
	}

	email, ok = misc.ValidateEmail(email)
	if !ok {
		log.Println("Wrong email", email)
		return 0, misc.WrongEmail
	}

	if !misc.IsPasswordValid(password) {
		log.Println("Wrong password")
		return 0, misc.WrongPassword
	}

	salt, err := auth.GenerateSalt()
	if err != nil {
		log.Println(err)
		return 0, misc.NoSalt
	}

	hash, err := auth.PasswordHash(password, salt)
	if err != nil {
		log.Println(err)
		return 0, misc.NothingToReport
	}

	userId, confirmationCode := 0, misc.RandomString(misc.ConfCodeLen)
	err = psql.Db.QueryRow(`
		INSERT INTO users (nickname, email, password, salt, confirmation_code)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`, nickname, email, hash, salt, confirmationCode,
	).Scan(&userId)
	if err == nil {
		mailer.EmailConfirmation(email, confirmationCode)
		return userId, misc.NothingToReport
	}

	err, code := psql.CheckSpecificDriverErrors(err)
	log.Println(err)
	return 0, code
}

// VerifyEmail verifies a previously created user
func VerifyEmail(userId int, confCode string) (string, bool) {
	sqlResult, err := psql.Db.Exec(`
		UPDATE users
		SET verified = True, confirmation_code = ''
		WHERE verified = False AND id = $1 AND confirmation_code = $2`, userId, confCode)
	if err, _ := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return "", false
	}

	if err, _ := psql.IsAffectedOneRow(sqlResult); err != nil {
		return "", false
	}

	jwt, err := auth.CreateJWT(userId, true)
	if err != nil {
		return "", false
	}

	return jwt, true
}

// Login a user
func Login(email, password string) (string, bool) {
	email, ok := misc.ValidateEmail(email)
	if !ok || !misc.IsPasswordValid(password) {
		return "", false
	}

	userId, hash, salt, verified := 0, make([]byte, 32), make([]byte, 16), false

	if err := psql.Db.QueryRow(`
		SELECT id, password, salt, verified
		FROM users
		WHERE email = $1`, email,
	).Scan(&userId, &hash, &salt, &verified); err != nil {
		return "", false
	}

	hashAttempt, err := auth.PasswordHash(password, salt)
	if err != nil {
		return "", false
	}

	if !reflect.DeepEqual(hashAttempt, hash) {
		return "", false
	}

	jwt, err := auth.CreateJWT(userId, verified)
	if err != nil {
		return "", false
	}
	return jwt, true
}

// UpdateAvatar updates user's path to avatar
func UpdateAvatar(userId int, path string) int {
	if !misc.IsIdValid(userId) {
		log.Println("user was not updated", userId)
		return misc.NothingUpdated
	}

	sqlResult, err := psql.Db.Exec(`
		UPDATE users
		SET image = $1
		WHERE id = $2`, path, userId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	err, code := psql.IsAffectedOneRow(sqlResult)
	return code
}
