-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS strategy (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL CONSTRAINT strategy_name_len CHECK (char_length(strategy_name) <= 15)
);

CREATE TABLE IF NOT EXISTS experiment (
    id BIGSERIAL PRIMARY KEY,
    idempotency_key UUID NOT NULL,
    strategy_id INTEGER NOT NULL,
    execution_time INTERVAL NOT NULL
);
CREATE INDEX IF NOT EXISTS idempotency_key_idx ON experiment (idempotency_key);
CREATE TABLE IF NOT EXISTS strategy_task (
    id BIGSERIAL PRIMARY KEY,
    experiment_id INTEGER NOT NULL,
    task_name TEXT NOT NULL CONSTRAINT task_name_len CHECK (char_length(task_name) <= 50),
    task_path TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS experiment_id_idx ON strategy_task (experiment_id);

INSERT INTO strategy (name)
VALUES ("round_robin"), ("fcs"), ("fcn");
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS experiment;
DROP TABLE IF EXISTS strategy_task;
DROP TABLE IF EXISTS strategy;
-- +goose StatementEnd
