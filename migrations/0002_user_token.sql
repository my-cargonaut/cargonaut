-- +migrate Up
CREATE TABLE user_token (
    id         uuid NOT NULL,
    user_id    uuid NOT NULL,
    expires_at timestamp WITHOUT TIME ZONE,
    created_at timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    CONSTRAINT user_token_pkey PRIMARY KEY (id),
    CONSTRAINT user_token_fkey FOREIGN KEY (user_id) REFERENCES user_account (id) ON DELETE CASCADE,
    CONSTRAINT user_token_id_key UNIQUE (id)
);
CREATE INDEX user_token_id_idx ON user_token USING btree (id);
CREATE INDEX user_token_user_id_idx ON user_token USING btree (user_id);
CREATE INDEX user_token_expires_at_idx ON user_token USING btree (expires_at);

-- +migrate Down
DROP INDEX user_token_id_idx;
DROP INDEX user_token_user_id_idx;
DROP INDEX user_token_expires_at_idx;
DROP TABLE user_token;
