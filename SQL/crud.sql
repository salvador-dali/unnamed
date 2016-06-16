-- create a purchase
INSERT INTO purchases (image, description, user_id, tags, brand) VALUES('some_img', 'descr', 8, '{1, 3}', 2); -- verify tags
UPDATE users SET purchases_num = purchases_num + 1 WHERE id= 8;

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