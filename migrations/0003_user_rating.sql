-- +migrate Up
CREATE TABLE user_rating (
    id         uuid NOT NULL,
    user_id    uuid NOT NULL,
    author_id  uuid NOT NULL,
    comment    character varying(128) NOT NULL,
    value      smallint NOT NULL, 
    created_at timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    CONSTRAINT user_rating_pkey PRIMARY KEY (id),
    CONSTRAINT user_rating_fkey FOREIGN KEY (user_id) REFERENCES user_account (id) ON DELETE CASCADE,
    CONSTRAINT user_rating_fkey_2 FOREIGN KEY (author_id) REFERENCES user_account (id) ON DELETE CASCADE,
    CONSTRAINT user_rating_id_key UNIQUE (id)
);
CREATE INDEX user_rating_id_idx ON user_rating USING btree (id);
CREATE INDEX user_rating_user_id_idx ON user_rating USING btree (user_id);

-- +migrate Down
DROP INDEX user_rating_id_idx;
DROP INDEX user_rating_user_id_idx;
DROP TABLE user_rating;
