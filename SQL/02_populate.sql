-- create a couple of users (everyone has password: password)
INSERT INTO users (nickname, about, email, password, salt, verified) VALUES('Albert Einstein', 'Developed the general theory of relativity.', 'albert@gmail.com', decode('1573b9a8ccfe4bb512bd7e7e50a7693dcdf61633ea791644dcd35453b7837ac52c7c7a95712c9c15a90bf19d26b79d233ccbe45371a390a4f49711e92e82beb65acfed98fff36c645063190ad10908f42addc23f946c57a6cf239323a7d1214e32ca48165f01905153977eaee0391eb278f32cf58ba4aac38e55c9901870698a', 'hex'), decode('215f4c97f4e20852bb395673d1758cce045cde8be21fc3db923b530f5fafda052d56b80df8218c19ed1041ce548daeaf534b9005f137a09fd2f2ed7eaffd52fc', 'hex'), TRUE);
INSERT INTO users (nickname, about, email, password, salt, verified) VALUES('Isaac Newton', 'Mechanics, laws of motion', 'isaac@gmail.com', decode('7f98d95e9064dfb3a6c4d7ef86013ee9ce866b05941f4cdedd4b8a7e107d5f303bdd9cb70c71a0f56fc5a36424e846b312d210250feab6e4e4aa9c6f102a58d8fa0106e601f17cf5287816521e4b84547b7832e28add0388b8b8a3bb4ef57557cbe4c71e517642f68149689ff6bf0ffbecc434f8fe6f4feea80261b38393d662', 'hex'), decode('846c0f8534fd6f6790fdcc989810a9a71df230bb9a26677feb567320d3ffddadc36eb333672a9d60f4819cbae14c4612514911db5cae7610de6c463fb7d37004', 'hex'), TRUE);
INSERT INTO users (nickname, about, email, password, salt, verified) VALUES('Marie Curie', 'Research on radioactivity', 'marie@gmail.com', decode('6b8b91f733bfc06f379516351396a80ebab080f91f768b99669d8ca8958560907db2f7450d11ebac3ed4b0b672e092ff4f7afa9d945b7b3d7b9c0aad13728e83d6aeea0023114f20597efcabb15e6842f115ff1c515601c46af781c46f0500f9e806077e415dfe53a255c940d60a67184cdd3e240893d87171e17b70c2a0e1ac', 'hex'), decode('d5322e4e8bb2df0bed64b172a3f07a38989e91f6d32739aa51940d77587d74e7c0906576da6b45ee886d0e9a397c544b96de80d1b4dec0bf3c7b9e6b1ba0a4f6', 'hex'), TRUE);
INSERT INTO users (nickname, about, email, password, salt, verified) VALUES('Galileo Galilei', 'Astronomy, heliocentrism, dynamics', 'galileo@gmail.com', decode('583007c2a6f09efe63fe77c3aabab92a0fad23355a626d0258366a956bd1f654c32f119a4ac27e3fdcdaf7607491c72a47d582abc4c246fb71d93eab487008b22f4d0ecba8aed1ac30f96d3105bd61edeac3696ce2498f55dac7a68898d73e351ff1cd41cd9f720092b85269e1bc96364fce1de12817c2f38fb20653d6e6dd1e', 'hex'), decode('ac389397945cfd47bddeee4dd7dfbb460268208782a0ddbda11cab033d84975b08266b6c690b12dee007d5874425324710993e608f00e9e41f42b6e857746739', 'hex'), TRUE);
INSERT INTO users (nickname, about, email, password, salt, verified) VALUES('Nikola Tesla', 'Alternating current. Die Edison!', 'nikola@gmail.com', decode('0f57205c0a98c307988a2f48f7943d5fa74d546c0852e728f682750361e7c97771c908aa87265b21be9b6cbef0edc65f3a43829556273dae63c63946efa702f9afa015876806553cf788a66ed5eeee93259b4a1f2b40377cb9482313254cf00d761ed44bce18242194fe0dbec59a5b822595aa04852311de9d1df3404cb0fcba', 'hex'), decode('a700053cd41129aa967e1dc7be51ff721e493c11abb3ca1a45444f510664f45637cb4266b995ee6f1e9c93c17df0fc960d981704c53366b91f13fa64d25e5c08', 'hex'), TRUE);
INSERT INTO users (nickname, about, email, password, salt, verified) VALUES('Louis Pasteur', 'Cool microbiologist. Now you know why the milk is pasteurized', 'louis@gmail.com', decode('1abe80200325d47b31da73447af6882342eec96f96f93838a20293e213b6189aaaf5ee32f29fbcfa1243465582c784d0919f2bd2493f3552ba102d92a2e1d9639fcab7517a61c4ccabafc2ed28476ba593fa63a3be175596b3361fc52004db0f8ef2807ea77ea7768600986c4b8457ac236e2d6397d0652361b0a86f56daf2b4', 'hex'), decode('8c5604c83b325dd46c8f09063a87c4bc78c679062c781c734ef557f3291a5224019e1ece844f93c8eee29245bf66eb2a4e9ce34a4fcb7cc87188e78bc2754ed9', 'hex'), TRUE);
INSERT INTO users (nickname, about, email, password, salt, verified) VALUES('Stephen Hawking', 'Too hard for most people to understand', 'stephen@gmail.com', decode('49ab5ae81e0286086770493c56bb5d633aa47ee8bba2fc1587b2d6769563a1561fb0dc47f886ec4e7d87bfabcd6e92789d6bbb1e0289b343c5c0d8f0f447ac44098274e93568728181c1e6a40a35cf3a4b4b0f28906bf942764c45a51eabb3f4e448d7e7ff9a5a8c45adbb6c0c83acc5855f3cb396b1f821e4254dc2d3b434bb', 'hex'), decode('38e6c8ec4be3274f91b7d78e117644542f97e39b9ca174f32f7b822d5148535f9f05edba2a2c0208e70b4b0a7affa44b9bc660ccb0fd16befd879b3eac2e7071', 'hex'), TRUE);
INSERT INTO users (nickname, about, email, password, salt, verified) VALUES('Charles Darwin', 'Evolution theory says that you are kind of a monkey', 'charles@gmail.com', decode('f48c91551c4bd1177c5b45ff705cec593ed26cdb1f7a403c14023713122eb361c4b614ae7975bde6ba7d9e572faf5f3d7909352586fd8c732f3109079a49241b036b79dc16030822ab07d10b9027506b98456f7f22a81771ac9deab6835ebf1ccb18f35822c10d3eb624ca95c4ada4ef019baad1d9d74cf6e807cf9196ad2d6f', 'hex'), decode('2fba2e6f71f2b1a494f0fe8c9b932f23e7d25894514e10b1ef43c0c5d47467614aa1a8ff70f29c6294f41011bcd34d331d2ff88e6d8d89c7201c8d2a46996f8b', 'hex'), TRUE);
INSERT INTO users (nickname, about, email, password, salt, verified) VALUES('Michael Faraday', 'Electromagnetism, electromagnetic induction and electrolysis', 'michael@gmail.com', decode('00a0ac3698caee1caf91e679ae8fc91f05c0a335b6d2146ae250b16df4cd513380871e2f53477cfdc52a657fd5ec19838ceb26f9a46c0666970bd9cbba39b8930ebf7aa38b76ffc1ac18cf2f32857324bf179863921ca0816311b011fc75e25d4e8f4b4f74697efbc2c11fd31d3d9bb5757c5f8ee144bff3494134cb46de135d', 'hex'), decode('04b4f091d216d8c85a1bdc29f39595bbaa44b1192053970a7611eab945969ef97715c4acd24d0af87116b1ee4f77ed5ff975baca361c658ec9fce85c7e1de258', 'hex'), TRUE);
INSERT INTO users (nickname, about, email, password, salt, confirmation_code) VALUES('Johannes Kepler', 'Mathematician, astronomer. When you speak about motion of planets, you think about me', 'kepler@gmail.com', decode('469653e2b3af2fa52d8be7922344ead2d67f3b200aa5c19e23ed4ab49a7437c40097f78538cba13fe1177bb6119125a8aba41dd7f812826161a04fe0ede0ebace95d200fd1769feaf05310f63d42826b47fba7f859d0f61db222eca6fb19b1a997323615afbd48603c70325755e61ec8fac3c4509838adee05cb4860b1a8af33', 'hex'), decode('42f58bc3edfe4c9cb5fa5d5193ef945402d017e6788aacfbe9f06a10fc746db0a001b3aa88e1210c3493d47c3ccbfc8e43783899a2ea70967a0783558ebe6bae', 'hex'), 'pqaJaBRgAvzLXqzRrrUI');

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