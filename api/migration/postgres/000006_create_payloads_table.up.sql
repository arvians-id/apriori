CREATE TABLE IF NOT EXISTS payloads (
    id_payload SERIAL,
    payload VARCHAR(256) NOT NULL,
    PRIMARY KEY (id_payload)
)