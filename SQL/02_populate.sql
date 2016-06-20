-- create a couple of users (everyone has password: password)
INSERT INTO users (nickname, about, email, password, salt) VALUES('Albert Einstein', 'Developed the general theory of relativity.', 'albert@gmail.com', decode('59488705fcf1dac7a41eb1641da649d9407e3009e72db57f0da3d67dd75df4c2', 'hex'), decode('43dc7dd85c0ed7afbe7c29f0b4b769ba', 'hex'));
INSERT INTO users (nickname, about, email, password, salt) VALUES('Isaac Newton', 'Mechanics, laws of motion', 'isaac@gmail.com', decode('5efe84f98f592f663e5c59d01cb0827522969de3bad278f5ef6cc98f07ccabdf', 'hex'), decode('3dab762cde874b110763b4e4326830cc', 'hex'));
INSERT INTO users (nickname, about, email, password, salt) VALUES('Marie Curie', 'Research on radioactivity', 'marie@gmail.com', decode('c8c955257a8d2b24347bb8590e080ee4cd0f2a3f844911d18268148ac682f3b9', 'hex'), decode('26f4b81a4f9a777e0f395398aa4b0661', 'hex'));
INSERT INTO users (nickname, about, email, password, salt) VALUES('Galileo Galilei', 'Astronomy, heliocentrism, dynamics', 'galileo@gmail.com', decode('050783af9a98fdc30924d1abf187fe6d67b1323e4f456fb30901295da560c6e6', 'hex'), decode('5e5615bb87eeeb9daf9aed89839640d2', 'hex'));
INSERT INTO users (nickname, about, email, password, salt) VALUES('Nikola Tesla', 'Alternating current. Die Edison!', 'nikola@gmail.com', decode('94bbc94b8a5df0dfc17175c7d7bcd417b029a3e27f7e780453aa868e9367d694', 'hex'), decode('8a377ca8c07b995749b7a71aed3dfcd4', 'hex'));
INSERT INTO users (nickname, about, email, password, salt) VALUES('Louis Pasteur', 'Cool microbiologist. Now you know why the milk is pasteurized', 'louis@gmail.com', decode('e17e05dce423ddcd768b254a228b88a491ed0d1f62a045573c68a6a3a62a9899', 'hex'), decode('aafc0f0021cc1f5d2c31037ff693463a', 'hex'));
INSERT INTO users (nickname, about, email, password, salt) VALUES('Stephen Hawking', 'Too hard for most people to understand', 'stephen@gmail.com', decode('81d2d3f3eca47df640b4531a870821ed30345a6aa4cc0ba090193cca18059756', 'hex'), decode('f8adb5a3506faa1d46e968d6c13d8935', 'hex'));
INSERT INTO users (nickname, about, email, password, salt) VALUES('Charles Darwin', 'Evolution theory says that you are kind of a monkey', 'charles@gmail.com', decode('b4e50b1416dc3b0c340a3ddeb61172f37153aaecf345aecadaecff47fc3b3a9b', 'hex'), decode('8256e8aeaada26c692dc581db5a7e218', 'hex'));
INSERT INTO users (nickname, about, email, password, salt) VALUES('Michael Faraday', 'Electromagnetism, electromagnetic induction and electrolysis', 'michael@gmail.com', decode('1931376ad6311162aae24b34a8cd3db296c37c6778278a41f359fe7284b052ef', 'hex'), decode('1856fca014621b7da833f403357a0e30', 'hex'));

-- create a couple of tags
INSERT INTO tags (name, description) VALUES('dress', 'nice dresses');
INSERT INTO tags (name, description) VALUES('drone', 'cool flying machines that do stuff');
INSERT INTO tags (name, description) VALUES('cosmetics', 'Known as make-up, are substances or products used to enhance the appearance or scent of the body');
INSERT INTO tags (name, description) VALUES('car', 'Vehicles that people use to move faster');
INSERT INTO tags (name, description) VALUES('hat', 'Stuff people put on their heads');
INSERT INTO tags (name, description) VALUES('phone', 'People use it to speak with other people');

-- create a couple of brands
INSERT INTO brands (name) VALUES('Apple');
INSERT INTO brands (name) VALUES('BMW');
INSERT INTO brands (name) VALUES('Playstation');
INSERT INTO brands (name) VALUES('Ferrari');
INSERT INTO brands (name) VALUES('Gucci');
INSERT INTO brands (id, name) VALUES(0, 'EMPTY BRAND');

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
INSERT INTO purchases (image, description, user_id, tag_ids, brand_id) VALUES('some_img', 'Look at my new drone', 1, '{2}', 0);
UPDATE users SET purchases_num = purchases_num + 1 WHERE id= 1;

INSERT INTO purchases (image, description, user_id, tag_ids, brand_id) VALUES('some_img', 'How cool am I?', 4, '{3, 5}', 5);
UPDATE users SET purchases_num = purchases_num + 1 WHERE id= 4;

INSERT INTO purchases (image, description, user_id, tag_ids, brand_id) VALUES('some_img', 'I really like drones', 1, '{4}', 0);
UPDATE users SET purchases_num = purchases_num + 1 WHERE id= 1;

INSERT INTO purchases (image, description, user_id, tag_ids, brand_id) VALUES('some_img', 'Now I am fond of cars', 1, '{2}', 4);
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