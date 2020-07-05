-- +migrate Up
CREATE TABLE rating (
    id         uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    user_id    uuid NOT NULL,
    author_id  uuid NOT NULL,
    trip_id    uuid NOT NULL,
    comment    character varying(128) NOT NULL,
    value      numeric NOT NULL, 
    created_at timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    CONSTRAINT rating_pkey PRIMARY KEY (id),
    CONSTRAINT rating_fkey FOREIGN KEY (user_id) REFERENCES user_account (id) ON DELETE CASCADE,
    CONSTRAINT rating_fkey_2 FOREIGN KEY (author_id) REFERENCES user_account (id) ON DELETE CASCADE,
    CONSTRAINT rating_fkey_3 FOREIGN KEY (trip_id) REFERENCES trip (id) ON DELETE CASCADE,
    CONSTRAINT rating_id_key UNIQUE (id),
    CONSTRAINT rating_author_id_trip_id_key UNIQUE (author_id, trip_id)
);
CREATE INDEX rating_id_idx ON rating USING btree (id);
CREATE INDEX rating_user_id_idx ON rating USING btree (user_id);
CREATE INDEX rating_trip_id_idx ON rating USING btree (trip_id);

-- +migrate Down
DROP INDEX rating_id_idx;
DROP INDEX rating_user_id_idx;
DROP INDEX rating_trip_id_idx;
DROP TABLE rating;
