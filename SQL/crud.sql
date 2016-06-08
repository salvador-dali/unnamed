-- create new user
INSERT INTO users (nickname, about) VALUES('user_1', 'about user 1');

-- update avatar
UPDATE users SET image = 'new_path' WHERE id = 10;

-- create/update tag/brand
INSERT INTO tags (name, description) VALUES('tag_6', 'descr_tag_6');
UPDATE tags SET name='tag_6', description='descr_tag_6' WHERE id=5;
INSERT INTO brands (name) VALUES('brand_4');
UPDATE brands SET name='brand_4' WHERE id=3;

-- follow
INSERT INTO followers ("who_id", "whom_id") VALUES(6, 10);
UPDATE users SET followers_num = followers_num + 1 WHERE id= 10;
UPDATE users SET following_num = following_num + 1 WHERE id= 6;

-- unfollow
DELETE FROM followers WHERE who_id=6 AND whom_id=10;
UPDATE users SET followers_num = followers_num - 1 WHERE id= 10;
UPDATE users SET following_num = following_num - 1 WHERE id= 6;

-- create a purchase
INSERT INTO purchases (image, description, user_id, tags, brand) VALUES('some_img', 'descr', 8, '{1, 3}', 2); -- verify tags
UPDATE users SET purchases_num = purchases_num + 1 WHERE id= 8;

-- like a purchase
INSERT INTO likes (purchase_id, user_id) VALUES(4, 10);
UPDATE purchases SET likes_num = likes_num + 1 WHERE id= 4;

-- unlike a purchase
DELETE FROM likes WHERE purchase_id = 4 AND user_id = 10;
UPDATE purchases SET likes_num = likes_num - 1 WHERE id= 4;

-- ask a question
INSERT INTO questions (user_id, purchase_id, name) VALUES(7, 6, 'So what exactly have you bought?');
UPDATE users SET questions_num = questions_num + 1 WHERE id= 7;

-- answer a question
INSERT INTO answers (user_id, question_id, name) VALUES(7, 3, 'Not sure');
UPDATE users SET answers_num = answers_num + 1 WHERE id= 7;

-- vote question/answer up/down
INSERT INTO votes_questions (user_id, question_id, is_voting_up) VALUES(7, 1, 1);
-- find whose question it is and update his stats