-- +goose Up
CREATE TABLE IF NOT EXISTS events (
                                      id UUID PRIMARY KEY,
                                      title VARCHAR(255) NOT NULL,
                                      description TEXT,
                                      start_time TIMESTAMP NOT NULL,
                                      end_time TIMESTAMP NOT NULL,
                                      user_id UUID NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS events;
