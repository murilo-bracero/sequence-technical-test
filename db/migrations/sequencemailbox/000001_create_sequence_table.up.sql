CREATE TABLE IF NOT EXISTS sequences(
    id serial primary key,
    external_id uuid not null default gen_random_uuid(),
    sequence_name varchar(255) not null,
    open_tracking_enabled boolean not null,
    click_tracking_enabled boolean not null,
    created timestamp not null default now(),
    updated timestamp
);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE sequences TO sequenceapi;

GRANT USAGE ON SEQUENCE sequences_id_seq TO sequenceapi;