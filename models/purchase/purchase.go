package purchase

import (
	"../../imager"
	"../../misc"
	"../../psql"
	"../tag"
	"database/sql"
	"log"
	"strconv"
	"strings"
	"time"
)

func getPurchases(rows *sql.Rows, err error) ([]*misc.Purchase, int) {
	if err != nil {
		log.Println(err)
		return []*misc.Purchase{}, misc.NothingToReport
	}
	defer rows.Close()

	purchases, tagString := []*misc.Purchase{}, ""
	var timestamp time.Time
	for rows.Next() {
		p := misc.Purchase{}
		if err := rows.Scan(&p.Id, &p.Image, &p.Description, &p.User_id, &timestamp, &tagString, &p.Brand, &p.Likes_num); err != nil {
			log.Println(err)
			return []*misc.Purchase{}, misc.NothingToReport
		}

		for _, v := range strings.Split(tagString[1:len(tagString)-1], ",") {
			if tagId, err := strconv.Atoi(v); err != nil {
				log.Println(err)
				return []*misc.Purchase{}, misc.NothingToReport
			} else {
				p.Tags = append(p.Tags, tagId)
			}
		}

		p.Issued_at = timestamp.Unix()
		purchases = append(purchases, &p)
	}

	if err = rows.Err(); err != nil {
		log.Println(err)
		return []*misc.Purchase{}, misc.NothingToReport
	}

	return purchases, misc.NothingToReport
}

func getCreatorByPurchaseId(purchaseId int) (int, int) {
	if !misc.IsIdValid(purchaseId) {
		log.Println("purchase ID is wrong", purchaseId)
		return 0, misc.NoPurchase
	}

	whosePurchase := 0
	if err := psql.Db.QueryRow(`
		SELECT user_id
		FROM purchases
		WHERE id = $1`, purchaseId,
	).Scan(&whosePurchase); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return 0, misc.NoPurchase
		}

		log.Println(err)
		return 0, misc.NothingToReport
	}

	return whosePurchase, misc.NothingToReport
}

func getCreatorByQuestionId(questionId int) (int, int) {
	if !misc.IsIdValid(questionId) {
		// if question does not exist, surely there is no purchase for this question
		log.Println("No question ID is wrong", questionId)
		return 0, misc.NoPurchase
	}

	whosePurchase := 0
	if err := psql.Db.QueryRow(`
		SELECT user_id
		FROM purchases
		WHERE id = (
			SELECT purchase_id
			FROM questions
			WHERE id = $1
		)`, questionId).Scan(&whosePurchase); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return 0, misc.NoPurchase
		}

		log.Println(err)
		return 0, misc.NothingToReport
	}

	return whosePurchase, misc.NothingToReport
}

// ShowAll returns all purchases
func ShowAll() ([]*misc.Purchase, int) {
	rows, err := psql.Db.Query(`
		SELECT id, image, description, user_id, issued_at, tag_ids, brand_id, likes_num
		FROM purchases
		ORDER BY issued_at DESC`)

	return getPurchases(rows, err)
}

// ShowById returns one purchase with Id
func ShowById(purchaseId int) (misc.Purchase, int) {
	if !misc.IsIdValid(purchaseId) {
		log.Println("Purchase ID is wrong", purchaseId)
		return misc.Purchase{}, misc.NoElement
	}

	p, tagString := misc.Purchase{}, ""
	var timestamp time.Time
	if err := psql.Db.QueryRow(`
		SELECT image, description, user_id, issued_at, tag_ids, brand_id, likes_num
		FROM purchases
		WHERE id = $1`, purchaseId,
	).Scan(&p.Image, &p.Description, &p.User_id, &timestamp, &tagString, &p.Brand, &p.Likes_num); err != nil {
		if err == sql.ErrNoRows {
			log.Println(err)
			return misc.Purchase{}, misc.NoElement
		}

		log.Println(err)
		return misc.Purchase{}, misc.NothingToReport
	}

	for _, v := range strings.Split(tagString[1:len(tagString)-1], ",") {
		if tagId, err := strconv.Atoi(v); err != nil {
			log.Println(err)
			return misc.Purchase{}, misc.NoElement
		} else {
			p.Tags = append(p.Tags, tagId)
		}
	}
	p.Id = purchaseId
	p.Issued_at = timestamp.Unix()
	return p, misc.NothingToReport
}

// ShowByUserId returns all purchases done by user Id
func ShowByUserId(userId int) ([]*misc.Purchase, int) {
	// userId is the current user and is always valid
	rows, err := psql.Db.Query(`
		SELECT id, image, description, user_id, issued_at, tag_ids, brand_id, likes_num
		FROM purchases
		WHERE user_id = $1
		ORDER BY issued_at DESC`, userId)

	return getPurchases(rows, err)
}

// ShowByBrandId returns all purchases with a brand Id
func ShowByBrandId(brandId int) ([]*misc.Purchase, int) {
	if !misc.IsIdValid(brandId) {
		log.Println("Brand Id is wrong", brandId)
		return []*misc.Purchase{}, misc.NothingToReport
	}

	rows, err := psql.Db.Query(`
		SELECT id, image, description, user_id, issued_at, tag_ids, brand_id, likes_num
		FROM purchases
		WHERE brand_id = $1
		ORDER BY issued_at DESC`, brandId)

	return getPurchases(rows, err)
}

// ShowByTagId returns all purchases with a tag Id
func ShowByTagId(tagId int) ([]*misc.Purchase, int) {
	if !misc.IsIdValid(tagId) {
		log.Println("Tag ID is wrong", tagId)
		return []*misc.Purchase{}, misc.NothingToReport
	}

	rows, err := psql.Db.Query(`
		SELECT id, image, description, user_id, issued_at, tag_ids, brand_id, likes_num
		FROM purchases
		WHERE $1 = ANY (tag_ids)
		ORDER BY issued_at DESC`, tagId)

	return getPurchases(rows, err)
}

// Create a new purchase
func Create(userId int, description, image string, brandId int, tagsId []int) (int, int) {
	// userID is the current user and should be valid
	description, ok := misc.ValidateString(description, misc.MaxLenB)
	if !ok {
		log.Println("description is wrong", description)
		return 0, misc.WrongDescr
	}

	if !imager.IsPurchaseValid(image) {
		log.Println("Purchase is not valid", image)
		return 0, misc.WrongImg
	}

	if brandId < 0 {
		log.Println("BrandID is wrong", brandId)
		return 0, misc.NoElement
	}

	if err, code := tag.ValidateTags(tagsId); err != nil {
		log.Println(err)
		return 0, code
	}

	stringTagIds, id := make([]string, len(tagsId), len(tagsId)), 0
	for k, v := range tagsId {
		stringTagIds[k] = strconv.Itoa(v)
	}

	tagsToInsert := "{" + strings.Join(stringTagIds, ",") + "}"
	err := psql.Db.QueryRow(`
		INSERT INTO purchases (image, description, user_id, tag_ids, brand_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`, image, description, userId, tagsToInsert, brandId).Scan(&id)
	if err != nil {
		if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
			log.Println(err)
			return 0, code
		}
	}

	sqlResult, err := psql.Db.Exec(`
		UPDATE users
		SET purchases_num = purchases_num + 1
		WHERE id=$1`, userId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return 0, code
	}

	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return 0, code
	}

	return id, misc.NothingToReport
}

// Like a purchase with some Id
func Like(purchaseId, userId int) int {
	if !misc.IsIdValid(purchaseId) {
		log.Println("Purchase Id is not valid", purchaseId)
		return misc.NoPurchase
	}

	// check whose purchase is it
	whosePurchase, code := getCreatorByPurchaseId(purchaseId)
	if whosePurchase == 0 {
		return code
	}

	if whosePurchase == userId {
		log.Println("can't vote for own purchase")
		return misc.VoteForYourself
	}

	// now allow the person to vote for someones else purchase
	sqlResult, err := psql.Db.Exec(`
		INSERT INTO likes (purchase_id, user_id)
		VALUES ($1, $2)`, purchaseId, userId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}
	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = psql.Db.Exec(`
		UPDATE purchases
		SET likes_num = likes_num + 1
		WHERE id = $1`, purchaseId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return code
	}

	return misc.NothingToReport
}

// Unlike a purchase which a user previously liked
func Unlike(purchaseId, userId int) int {
	if !misc.IsIdValid(purchaseId) {
		log.Println("Purchase Id is not possitive", purchaseId)
		return misc.NoPurchase
	}

	// check whose purchase is it
	whosePurchase, code := getCreatorByPurchaseId(purchaseId)
	if whosePurchase == 0 {
		return code
	}

	if whosePurchase == userId {
		log.Println("can't vote for own purchase")
		return misc.VoteForYourself
	}

	sqlResult, err := psql.Db.Exec(`
		DELETE FROM likes
		WHERE purchase_id = $1 AND user_id = $2`, purchaseId, userId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}
	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return code
	}

	sqlResult, err = psql.Db.Exec(`
		UPDATE purchases
		SET likes_num = likes_num - 1
		WHERE id = $1`, purchaseId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return code
	}

	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return code
	}

	return misc.NothingToReport
}

// AskQuestion about a specific purchase
func AskQuestion(purchaseId, userId int, question string) (int, int) {
	whosePurchase, code := getCreatorByPurchaseId(purchaseId)
	if whosePurchase == 0 {
		return 0, code
	}

	if whosePurchase == userId {
		log.Println("can't vote for own purchase")
		return 0, misc.AskYourself
	}

	question, ok := misc.ValidateString(question, misc.MaxLenB)
	if !ok {
		log.Println("Wrong question", question)
		return 0, misc.WrongName
	}

	questionId := 0
	err := psql.Db.QueryRow(`
		INSERT INTO questions (user_id, purchase_id, name)
		VALUES ($1, $2, $3)
		RETURNING id`, userId, purchaseId, question,
	).Scan(&questionId)
	if err != nil {
		err, code := psql.CheckSpecificDriverErrors(err)
		log.Println(err)
		return 0, code
	}

	sqlResult, err := psql.Db.Exec(`
		UPDATE users
		SET questions_num = questions_num + 1
		WHERE id = $1`, userId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return 0, code
	}

	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return 0, code
	}

	return questionId, misc.NothingToReport
}

// AnswerQuestion answers previously asked question
func AnswerQuestion(questionId, userId int, answer string) (int, int) {
	whosePurchase, code := getCreatorByQuestionId(questionId)
	if whosePurchase == 0 {
		return 0, code
	}

	if whosePurchase != userId {
		log.Println("can asnwer only questions regarding your purchase")
		return 0, misc.AnswerOtherPurchase
	}

	answer, ok := misc.ValidateString(answer, misc.MaxLenB)
	if !ok {
		log.Println("Wrong answer", answer)
		return 0, misc.WrongName
	}

	answerId := 0
	err := psql.Db.QueryRow(`
		INSERT INTO answers (user_id, question_id, name)
		VALUES ($1, $2, $3)
		RETURNING id`, userId, questionId, answer,
	).Scan(&answerId)
	if err != nil {
		err, code := psql.CheckSpecificDriverErrors(err)
		log.Println(err)
		return 0, code
	}

	sqlResult, err := psql.Db.Exec(`
		UPDATE users
		SET answers_num = answers_num + 1
		WHERE id = $1`, userId)
	if err, code := psql.CheckSpecificDriverErrors(err); err != nil {
		log.Println(err)
		return 0, code
	}

	if err, code := psql.IsAffectedOneRow(sqlResult); err != nil {
		return 0, code
	}

	return answerId, misc.NothingToReport
}
