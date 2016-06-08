-- create a couple of users
INSERT INTO users (nickname, about) VALUES('Albert Einstein', 'Developed the general theory of relativity.');
INSERT INTO users (nickname, about) VALUES('Isaac Newton', 'Mechanics, laws of motion');
INSERT INTO users (nickname, about) VALUES('Marie Curie', 'Research on radioactivity');
INSERT INTO users (nickname, about) VALUES('Galileo Galilei', 'Astronomy, heliocentrism, dynamics');
INSERT INTO users (nickname, about) VALUES('Nikola Tesla', 'Alternating current. Die Edison!');
INSERT INTO users (nickname, about) VALUES('Louis Pasteur', 'Cool microbiologist. Now you know why the milk is pasteurized');
INSERT INTO users (nickname, about) VALUES('Stephen Hawking', 'Too hard for most people to understand');
INSERT INTO users (nickname, about) VALUES('Charles Darwin', 'Evolution theory says that you are kind of a monkey');
INSERT INTO users (nickname, about) VALUES('Michael Faraday', 'Electromagnetism, electromagnetic induction and electrolysis');

-- create a couple of tags
INSERT INTO tags (name, description) VALUES('dress', 'nice dresses');
INSERT INTO tags (name, description) VALUES('drone', 'cool flying machines that do stuff');
INSERT INTO tags (name, description) VALUES('cosmetics', 'Known as make-up, are substances or products used to enhance the appearance or scent of the body');
INSERT INTO tags (name, description) VALUES('car', 'Vehicles that people use to move faster');
INSERT INTO tags (name, description) VALUES('hat', 'Stuff people put on their heads');
INSERT INTO tags (name, description) VALUES('phone', 'People use it to speak with other people');

-- create a couple of tags
INSERT INTO brands (name) VALUES('Apple');
INSERT INTO brands (name) VALUES('BMW');
INSERT INTO brands (name) VALUES('Playstation');
INSERT INTO brands (name) VALUES('Ferrari');
INSERT INTO brands (name) VALUES('Gucci');

-- following each other
INSERT INTO followers ("who_id", "whom_id") VALUES(1, 2);
UPDATE users SET following_num = following_num + 1 WHERE id= 1;
UPDATE users SET followers_num = followers_num + 1 WHERE id= 2;

INSERT INTO followers ("who_id", "whom_id") VALUES(1, 4);
UPDATE users SET following_num = following_num + 1 WHERE id= 1;
UPDATE users SET followers_num = followers_num + 1 WHERE id= 4;

INSERT INTO followers ("who_id", "whom_id") VALUES(1, 7);
UPDATE users SET following_num = following_num + 1 WHERE id= 1;
UPDATE users SET followers_num = followers_num + 1 WHERE id= 7;

INSERT INTO followers ("who_id", "whom_id") VALUES(6, 2);
UPDATE users SET following_num = following_num + 1 WHERE id= 6;
UPDATE users SET followers_num = followers_num + 1 WHERE id= 2;

-- create a few purchases
INSERT INTO purchases (image, description, user_id, tags, brand) VALUES('some_img', 'Look at my new drone', 1, '{2}', Null);
UPDATE users SET purchases_num = purchases_num + 1 WHERE id= 1;

INSERT INTO purchases (image, description, user_id, tags, brand) VALUES('some_img', 'How cool am I?', 4, '{3, 5}', 5);
UPDATE users SET purchases_num = purchases_num + 1 WHERE id= 4;

INSERT INTO purchases (image, description, user_id, tags, brand) VALUES('some_img', 'I really like drones', 1, '{4}', Null);
UPDATE users SET purchases_num = purchases_num + 1 WHERE id= 1;

INSERT INTO purchases (image, description, user_id, tags, brand) VALUES('some_img', 'Now I am fond of cars', 1, '{2}', 4);
UPDATE users SET purchases_num = purchases_num + 1 WHERE id= 1;

-- some people like them
INSERT INTO likes (purchase_id, user_id) VALUES(3, 4);
UPDATE purchases SET likes_num = likes_num + 1 WHERE id= 3;

INSERT INTO likes (purchase_id, user_id) VALUES(3, 6);
UPDATE purchases SET likes_num = likes_num + 1 WHERE id= 3;

INSERT INTO likes (purchase_id, user_id) VALUES(3, 9);
UPDATE purchases SET likes_num = likes_num + 1 WHERE id= 3;

INSERT INTO likes (purchase_id, user_id) VALUES(4, 2);
UPDATE purchases SET likes_num = likes_num + 1 WHERE id= 4;

-- ask questions
INSERT INTO questions (user_id, purchase_id, name) VALUES(7, 4, 'How fast can it go?');
UPDATE users SET questions_num = questions_num + 1 WHERE id= 7;

INSERT INTO questions (user_id, purchase_id, name) VALUES(7, 1, 'What is the maximum distance from the transmitter?');
UPDATE users SET questions_num = questions_num + 1 WHERE id= 7;

-- answer the question
INSERT INTO answers (user_id, question_id, name) VALUES(1, 2, '3 km is the maximum. To be on the safe side, use 2km');
UPDATE users SET answers_num = answers_num + 1 WHERE id= 1;