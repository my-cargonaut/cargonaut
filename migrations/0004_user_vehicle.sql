-- +migrate Up
CREATE TABLE user_vehicle (
    id                  uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    user_id             uuid NOT NULL,
    brand               character varying(128) NOT NULL,
    model               character varying(128) NOT NULL,
    passengers          smallint, 
    loading_area_length numeric, 
    loading_area_width  numeric, 
    created_at timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    updated_at timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    CONSTRAINT user_vehicle_pkey PRIMARY KEY (id),
    CONSTRAINT user_vehicle_fkey FOREIGN KEY (user_id) REFERENCES user_account (id) ON DELETE CASCADE,
    CONSTRAINT user_vehicle_id_key UNIQUE (brand, model)
);
CREATE INDEX user_vehicle_id_idx ON user_vehicle USING btree (id);
CREATE INDEX user_vehicle_user_id_idx ON user_vehicle USING btree (user_id);

-- +migrate Down
DROP INDEX user_vehicle_id_idx;
DROP INDEX user_vehicle_user_id_idx;
DROP TABLE user_vehicle;
