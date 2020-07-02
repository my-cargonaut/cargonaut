-- +migrate Up
CREATE TABLE trip (
    id          uuid NOT NULL DEFAULT uuid_generate_v1mc(),
    user_id     uuid NOT NULL,
    vehicle_id  uuid NOT NULL,
    rider_id    uuid,
    start       character varying(128) NOT NULL,
    destination character varying(128) NOT NULL,
    price       numeric NOT NULL, 
    depature    timestamp WITHOUT TIME ZONE,
    arrival     timestamp WITHOUT TIME ZONE,
    created_at  timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    updated_at  timestamp WITHOUT TIME ZONE DEFAULT (now() at time zone 'utc'),
    CONSTRAINT trip_pkey PRIMARY KEY (id),
    CONSTRAINT trip_fkey FOREIGN KEY (user_id) REFERENCES user_account (id) ON DELETE CASCADE,
    CONSTRAINT trip_fkey_2 FOREIGN KEY (rider_id) REFERENCES user_account (id) ON DELETE CASCADE,
    CONSTRAINT trip_fkey_3 FOREIGN KEY (vehicle_id) REFERENCES vehicle (id) ON DELETE CASCADE,
    CONSTRAINT trip_id_key UNIQUE (id)
);
CREATE INDEX trip_id_idx ON trip USING btree (id);
CREATE INDEX trip_user_id_idx ON trip USING btree (user_id);
CREATE INDEX trip_rider_id_idx ON trip USING btree (rider_id);

-- +migrate Down
DROP INDEX trip_id_idx;
DROP INDEX trip_user_id_idx;
DROP INDEX trip_rider_id_idx;
DROP TABLE trip;
