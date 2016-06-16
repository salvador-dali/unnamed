-- remove all tables
DROP TABLE IF EXISTS votes_questions;
DROP TABLE IF EXISTS votes_answers;
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS followers;
DROP TABLE IF EXISTS answers;
DROP TABLE IF EXISTS questions;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS purchases;
DROP TABLE IF EXISTS brands;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS timeseries;

-- Users
CREATE TABLE "users" (
    "id" bigserial,
    "nickname" varchar(40) NOT NULL,
    "image" varchar(100) NOT NULL DEFAULT '',
    "about" varchar(1000)  NOT NULL,
    "expertise" int  NOT NULL DEFAULT 0,
    "followers_num" int  NOT NULL DEFAULT 0,
    "following_num" int  NOT NULL DEFAULT 0,
    "purchases_num" int  NOT NULL DEFAULT 0,
    "questions_num" int  NOT NULL DEFAULT 0,
    "answers_num" int  NOT NULL DEFAULT 0,
    "issued_at" timestamp  NOT NULL DEFAULT (now() at time zone 'utc'),
    "tags_like" integer[]  NOT NULL DEFAULT '{}',
    "tags_ignore" integer[]  NOT NULL DEFAULT '{}',
    "brands_like" integer[]  NOT NULL DEFAULT '{}',
    "brands_ignore" integer[]  NOT NULL DEFAULT '{}',
    PRIMARY KEY ("id"),
    UNIQUE ("nickname")
);
COMMENT ON TABLE "users" IS 'All users in the system';
COMMENT ON COLUMN "users"."id" IS 'ID of a user';
COMMENT ON COLUMN "users"."nickname" IS 'Nickname of a user. Does not represent real name';
COMMENT ON COLUMN "users"."image" IS 'Path to the location of the image';
COMMENT ON COLUMN "users"."about" IS 'General field, where person can write anything about himself';
COMMENT ON COLUMN "users"."expertise" IS 'Total number of upvotes/downvotes to all answers';
COMMENT ON COLUMN "users"."followers_num" IS 'Number of people who follows this user';
COMMENT ON COLUMN "users"."following_num" IS 'Number of people whom this user follows';
COMMENT ON COLUMN "users"."purchases_num" IS 'Number of purchases this person posted';
COMMENT ON COLUMN "users"."questions_num" IS 'Number of questions this person asked';
COMMENT ON COLUMN "users"."answers_num" IS 'Number of answers this person provided';
COMMENT ON COLUMN "users"."issued_at" IS 'Time when a user was created';
COMMENT ON COLUMN "users"."tags_like" IS 'Array of tags, the person likes';
COMMENT ON COLUMN "users"."tags_ignore" IS 'Array of tags, the person wishes to ignore';
COMMENT ON COLUMN "users"."brands_like" IS 'Array of brands, the person likes';
COMMENT ON COLUMN "users"."brands_ignore" IS 'Array of brands, the person wishes to ignore';

-- Followers
CREATE TABLE "followers" (
    "who_id" bigint NOT NULL,
    "whom_id" bigint NOT NULL,
    "issued_at" timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
    FOREIGN KEY ("who_id") REFERENCES "users"("id"),
    FOREIGN KEY ("whom_id") REFERENCES "users"("id")
);
CREATE UNIQUE INDEX followers_who_id_whom_id_pkey ON followers (who_id, whom_id);
COMMENT ON TABLE "followers" IS 'Who follows whom';
COMMENT ON COLUMN "followers"."who_id" IS 'ID of a person who follows someone';
COMMENT ON COLUMN "followers"."whom_id" IS 'ID of a person whom a person follows';
COMMENT ON COLUMN "followers"."issued_at" IS 'When who_id started to follow whom_id';

-- Brands
CREATE TABLE "brands" (
    "id" serial,
    "name" varchar(40) NOT NULL,
    "issued_at" timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
    PRIMARY KEY ("id"),
    UNIQUE ("name")
);
COMMENT ON TABLE "brands" IS 'All brands in the system';
COMMENT ON COLUMN "brands"."id" IS 'ID of a brand';
COMMENT ON COLUMN "brands"."name" IS 'Name of a brand. Something like ''Nike''';
COMMENT ON COLUMN "brands"."issued_at" IS 'When a brand was created';

-- Tags
CREATE TABLE "tags" (
    "id" serial,
    "name" varchar(40) NOT NULL,
    "description" varchar(1000) NOT NULL,
    "issued_at" timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
    PRIMARY KEY ("id"),
    UNIQUE ("name")
);
COMMENT ON TABLE "tags" IS 'All tags in the system';
COMMENT ON COLUMN "tags"."id" IS 'ID of a tag';
COMMENT ON COLUMN "tags"."name" IS 'Name of a tag. Something like ''helicopter''';
COMMENT ON COLUMN "tags"."description" IS 'Longer description of a tag';
COMMENT ON COLUMN "tags"."issued_at" IS 'When a tag was created';

-- Purchases
CREATE TABLE "purchases" (
    "id" bigserial,
    "image" varchar(100) NOT NULL,
    "description" varchar(1000) NOT NULL,
    "user_id" bigint NOT NULL,
    "issued_at" timestamp  NOT NULL DEFAULT (now() at time zone 'utc'),
    "tags" integer[] NOT NULL,
    "brand" integer NOT NULL DEFAULT 0,
    "likes_num" integer NOT NULL DEFAULT 0,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("brand") REFERENCES "brands"("id") ON DELETE SET NULL,
    FOREIGN KEY ("user_id") REFERENCES "users"("id")
    -- TODO add foreign key for tags
);
COMMENT ON TABLE "purchases" IS 'All purchases in the system';
COMMENT ON COLUMN "purchases"."id" IS 'ID of a purchase';
COMMENT ON COLUMN "purchases"."image" IS 'Path to the location of the image';
COMMENT ON COLUMN "purchases"."description" IS 'Short description of what exactly was bought and why is it so exciting for a user';
COMMENT ON COLUMN "purchases"."user_id" IS 'Who posted this purchase';
COMMENT ON COLUMN "purchases"."issued_at" IS 'When was the purchase posted';
COMMENT ON COLUMN "purchases"."tags" IS 'Array of tags associated with the purchase';
COMMENT ON COLUMN "purchases"."brand" IS 'A brand associated with the purchase. A purchase can have no brand';
COMMENT ON COLUMN "purchases"."likes_num" IS 'Number of likes a purchase received';

-- Likes
CREATE TABLE "likes" (
    "purchase_id" bigint NOT NULL,
    "user_id" bigint NOT NULL,
    "issued_at" timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
    FOREIGN KEY ("user_id") REFERENCES "users"("id"),
    FOREIGN KEY ("purchase_id") REFERENCES "purchases"("id")
);
CREATE UNIQUE INDEX likes_user_id_purchase_id_pkey ON likes (user_id, purchase_id);
COMMENT ON TABLE "likes" IS 'All likes/unlikes related to purchases';
COMMENT ON COLUMN "likes"."purchase_id" IS 'ID of a purchase that was liked';
COMMENT ON COLUMN "likes"."user_id" IS 'ID of a user who issued a like';
COMMENT ON COLUMN "likes"."issued_at" IS 'When a like was issued';

-- Questions
CREATE TABLE "questions" (
    "id" bigserial,
    "user_id" bigint NOT NULL,
    "purchase_id" bigint NOT NULL,
    "name" varchar(100) NOT NULL,
    "issued_at" timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
    "votes_num" int DEFAULT 0 NOT NULl,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("purchase_id") REFERENCES "purchases"("id"),
    FOREIGN KEY ("user_id") REFERENCES "users"("id")
);
COMMENT ON TABLE "questions" IS 'All questions in the system';
COMMENT ON COLUMN "questions"."id" IS 'ID of a question';
COMMENT ON COLUMN "questions"."user_id" IS 'Who asked a question';
COMMENT ON COLUMN "questions"."purchase_id" IS 'What purchase is the question related to';
COMMENT ON COLUMN "questions"."name" IS 'The actual question';
COMMENT ON COLUMN "questions"."issued_at" IS 'When was the question asked';
COMMENT ON COLUMN "questions"."votes_num" IS 'Number of votes that this question received';

-- Answers
CREATE TABLE "answers" (
    "id" bigserial,
    "user_id" bigint NOT NULL,
    "question_id" bigint NOT NULL,
    "name" varchar(1000) NOT NULL,
    "issued_at" timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
    "votes_num" int NOT NULL DEFAULT 0,
    PRIMARY KEY ("id"),
    FOREIGN KEY ("user_id") REFERENCES "users"("id"),
    FOREIGN KEY ("question_id") REFERENCES "questions"("id")
);
COMMENT ON TABLE "answers" IS 'All the answers in the system';
COMMENT ON COLUMN "answers"."id" IS 'ID of an answer';
COMMENT ON COLUMN "answers"."user_id" IS 'ID of a user who wrote an answer';
COMMENT ON COLUMN "answers"."question_id" IS 'ID of a question which this answer tries to answer';
COMMENT ON COLUMN "answers"."name" IS 'Actual answer';
COMMENT ON COLUMN "answers"."issued_at" IS 'When was the answer given';
COMMENT ON COLUMN "answers"."votes_num" IS 'Number of votes, the answer received';

-- Votes
CREATE TABLE "votes_questions" (
    "user_id" bigint NOT NULL,
    "question_id" bigint NOT NULL,
    "issued_at" timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
    "is_voting_up" boolean NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users"("id"),
    FOREIGN KEY ("question_id") REFERENCES "questions"("id")
);
COMMENT ON TABLE "votes_questions" IS 'All votes for questions in the system';
COMMENT ON COLUMN "votes_questions"."user_id" IS 'ID of a user who voted';
COMMENT ON COLUMN "votes_questions"."question_id" IS 'ID of a question which was voted';
COMMENT ON COLUMN "votes_questions"."issued_at" IS 'When was the vote issued';
COMMENT ON COLUMN "votes_questions"."is_voting_up" IS 'Is person voting up or down';

CREATE TABLE "votes_answers" (
    "user_id" bigint NOT NULL,
    "answer_id" bigint NOT NULL,
    "issued_at" timestamp NOT NULL DEFAULT (now() at time zone 'utc'),
    "is_voting_up" boolean NOT NULL,
    FOREIGN KEY ("user_id") REFERENCES "users"("id"),
    FOREIGN KEY ("answer_id") REFERENCES "answers"("id")
);
COMMENT ON TABLE "votes_answers" IS 'All votes for answers in the system';
COMMENT ON COLUMN "votes_answers"."user_id" IS 'ID of a user who voted';
COMMENT ON COLUMN "votes_answers"."answer_id" IS 'ID of an answer which was voted';
COMMENT ON COLUMN "votes_answers"."issued_at" IS 'When was the vote issued';
COMMENT ON COLUMN "votes_answers"."is_voting_up" IS 'Is person voting up or down';

-- information about all events in the system
CREATE TABLE "timeseries" (
    "id" bigserial,
    "issued_at" timestamp,
    "user_id" bigint,
    "fingerprint" bit(128),
    "lat" real,
    "lng" real,
    "os" varchar(100),
    "browser" varchar(100),
    "cpu" varchar(100),
    "is_charging" boolean,
    "battery" smallint,
    "ip" cidr,
    "action_type" int,
    "experiment_type" smallint,
    "user_id2" bigint,
    "brand_id" int,
    "tag_id" int,
    "purchase_id" bigint,
    "like_id" bigint,
    "question_id" bigint,
    "answer_id" bigint,
    "vote_id" bigint,
    "is_backend" boolean,
    PRIMARY KEY ("id")
);
COMMENT ON TABLE "timeseries" IS 'This table is for data analysis purposes. It stores information about all possible events that happened in the system';
COMMENT ON COLUMN "timeseries"."id" IS 'Id of the event. The main purpose is to give a link to an event quickly.';
COMMENT ON COLUMN "timeseries"."issued_at" IS 'When an event happened';
COMMENT ON COLUMN "timeseries"."user_id" IS 'Id of the user who has done something. Null if it is an anonymous person.';
COMMENT ON COLUMN "timeseries"."fingerprint" IS 'Binary string which attempts to identify the device of anonymous user. Implementation: https://github.com/Valve/fingerprintjs2';
COMMENT ON COLUMN "timeseries"."lat" IS 'Latitude http://gis.stackexchange.com/a/8674/9972';
COMMENT ON COLUMN "timeseries"."lng" IS 'Longitude';
COMMENT ON COLUMN "timeseries"."os" IS 'User''s operation system with minor numbers. Like Mac OS 10.10.1';
COMMENT ON COLUMN "timeseries"."browser" IS 'User''s browser with build numbers. Like Chrome 50.0.2661.102';
COMMENT ON COLUMN "timeseries"."cpu" IS 'Information about user''s CPU hardware: MacIntel, 4 Cores';
COMMENT ON COLUMN "timeseries"."is_charging" IS 'Whether the battery is charging at the moment. http://webkay.robinlinus.com/';
COMMENT ON COLUMN "timeseries"."battery" IS 'Battery level in percentage. Actual values are from 0 to 100.';
COMMENT ON COLUMN "timeseries"."ip" IS 'IP address of a user';
COMMENT ON COLUMN "timeseries"."action_type" IS 'What was the action, user performed. Clicked submit button, followed someone. Would be enum';
COMMENT ON COLUMN "timeseries"."experiment_type" IS 'Some of the events are permanent. Information about them is always collected. The value for them is 0. Some other are temporary experiments that run for a small period of time. This allows to find out what experiment was it.';
COMMENT ON COLUMN "timeseries"."user_id2" IS 'If user_id does something with another user (follows, flags, etc), this field stores info about that second user.';
COMMENT ON COLUMN "timeseries"."brand_id" IS 'Id of a brand, if it is relevant to an event';
COMMENT ON COLUMN "timeseries"."tag_id" IS 'Id of a tag, if it is relevant to an event';
COMMENT ON COLUMN "timeseries"."purchase_id" IS 'Id of a purchase, if it is relevant to an event';
COMMENT ON COLUMN "timeseries"."like_id" IS 'Id of a like, if it is relevant to an event';
COMMENT ON COLUMN "timeseries"."question_id" IS 'Id of a question, if it is relevant to an event';
COMMENT ON COLUMN "timeseries"."answer_id" IS 'Id of an answer, if it is relevant to an event';
COMMENT ON COLUMN "timeseries"."vote_id" IS 'Id of a vote, if it is relevant to an event';
COMMENT ON COLUMN "timeseries"."is_backend" IS 'Whether the change was verified on a backend. For example when user started to follow another user, a backend event is fired, thus verifying the event. When a user hovered an item, this is a purely frontend event and it can''t be verified.';
