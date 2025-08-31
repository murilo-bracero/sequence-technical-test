CREATE TABLE IF NOT EXISTS steps(
    id serial primary key,
    external_id uuid not null default gen_random_uuid(),
    mail_subject varchar(255) not null,
    mail_content text,
    step_number integer not null,
    sequence_id integer not null,
    foreign key (sequence_id) references sequences(id) on delete cascade
);

GRANT SELECT, INSERT, UPDATE, DELETE ON TABLE steps TO sequenceapi;

GRANT USAGE ON SEQUENCE steps_id_seq TO sequenceapi;