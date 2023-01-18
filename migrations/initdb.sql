CREATE TABLE user_ (
	id_ SERIAL PRIMARY KEY,
	login_ VARCHAR(64) NOT NULL ,
	password_ VARCHAR(64) NOT NULL,
    UNIQUE (login_)
);

CREATE TABLE form_ (
    id_ SERIAL PRIMARY KEY,
    user_id_ INT REFERENCES user_ ON DELETE CASCADE NOT NULL,
    title_ VARCHAR(64) NOT NULL,
    description_ TEXT NOT NULL
);

CREATE TABLE question_ (
    id_ SERIAL PRIMARY KEY,
    form_id_ INT REFERENCES form_ ON DELETE CASCADE NOT NULL,
    header_ TEXT NOT NULL
);

CREATE TABLE pool_answer_ (
    id_ SERIAL PRIMARY KEY,
    user_id_ INT REFERENCES user_ ON DELETE CASCADE NOT NULL,
    form_id_ INT REFERENCES form_ ON DELETE CASCADE NOT NULL
);

CREATE TABLE answer_ (
    id_ SERIAL PRIMARY KEY,
    question_id_ INT REFERENCES question_ ON DELETE CASCADE NOT NULL,
    pool_answer_id_ INT REFERENCES pool_answer_ ON DELETE CASCADE NOT NULL,
    value_ TEXT NOT NULL
);

CREATE ROLE db_readonly;
GRANT CONNECT ON DATABASE quizapp TO db_readonly;
GRANT USAGE ON SCHEMA public TO db_readonly;

GRANT SELECT ON TABLE quizapp.public.form_ TO db_readonly;
GRANT SELECT ON TABLE quizapp.public.question_ TO db_readonly;
GRANT SELECT ON TABLE quizapp.public.pool_answer_ TO db_readonly;
GRANT SELECT ON TABLE quizapp.public.answer_ TO db_readonly;
GRANT SELECT ON TABLE quizapp.public.user_ TO db_readonly;

CREATE USER minotauro_readonly WITH PASSWORD 'Controcarro3_readonly';
GRANT db_readonly TO minotauro_readonly;
