-- @formatter:off

CREATE TABLE IF NOT EXISTS article
(
    id              uuid            PRIMARY KEY         DEFAULT gen_random_uuid(),
    title           text            NOT NULL,
    description     text            NOT NULL,
    created_at      timestamp       NOT NULL            DEFAULT current_timestamp,
    updated_at      timestamp       NOT NULL            DEFAULT current_timestamp
    );

CREATE INDEX IF NOT EXISTS article_title_idx ON article(title);

CREATE TABLE IF NOT EXISTS outbox_message
(
    id              uuid            PRIMARY KEY         DEFAULT gen_random_uuid(),
    is_sent         boolean                             DEFAULT FALSE,
    message         jsonb           NOT NULL,
    created_at      timestamp                           DEFAULT current_timestamp,
    updated_at      timestamp                           DEFAULT current_timestamp
);

-- @formatter:on