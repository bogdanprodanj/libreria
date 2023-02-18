CREATE TABLE books
(
    id           SERIAL                         NOT NULL
        CONSTRAINT books_pkey
            PRIMARY KEY,
    title        TEXT                           NOT NULL,
    author       TEXT                           NOT NULL,
    publisher    TEXT                           NOT NULL,
    publish_date TIMESTAMP,
    rating       DOUBLE PRECISION DEFAULT 0,
    status       INTEGER          DEFAULT 0,
    created_at   TIMESTAMP        DEFAULT NOW() NOT NULL,
    updated_at   TIMESTAMP        DEFAULT NOW() NOT NULL,
    deleted_at   TIMESTAMP
);
