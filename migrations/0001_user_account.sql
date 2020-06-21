-- +migrate Up
CREATE TABLE user_account (
    id                 uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    email              character varying(128) NOT NULL,
    password_hash      character varying(128) NOT NULL,
    display_name       character varying(128) NOT NULL,
    birthday           timestamp WITHOUT TIME ZONE NOT NULL,
    avatar             BYTEA NOT NULL,
    created_at         timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    updated_at         timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    CONSTRAINT user_account_pkey PRIMARY KEY (id),
    CONSTRAINT user_account_email_key UNIQUE (email)
);
CREATE INDEX user_account_id_idx ON user_account USING btree (id);
CREATE INDEX user_account_email_idx ON user_account USING btree (email);

-- +migrate Down
DROP INDEX user_account_id_idx;
DROP INDEX user_account_email_idx;
DROP TABLE user_account;
