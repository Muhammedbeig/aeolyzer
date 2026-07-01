CREATE TABLE IF NOT EXISTS aeolyzer_sessions (
    app_name VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    session_id VARCHAR(64) NOT NULL,
    title_ciphertext LONGBLOB NULL,
    state_ciphertext LONGBLOB NOT NULL,
    starred BOOLEAN NOT NULL DEFAULT FALSE,
    next_sequence BIGINT UNSIGNED NOT NULL DEFAULT 1,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    PRIMARY KEY (app_name, user_id, session_id),
    INDEX aeolyzer_sessions_updated (app_name, user_id, updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS aeolyzer_events (
    app_name VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    session_id VARCHAR(64) NOT NULL,
    event_id VARCHAR(64) NOT NULL,
    sequence_number BIGINT UNSIGNED NOT NULL,
    event_ciphertext LONGBLOB NOT NULL,
    created_at DATETIME(6) NOT NULL,
    PRIMARY KEY (app_name, user_id, session_id, event_id),
    UNIQUE KEY aeolyzer_events_sequence (app_name, user_id, session_id, sequence_number),
    CONSTRAINT aeolyzer_events_session_fk
        FOREIGN KEY (app_name, user_id, session_id)
        REFERENCES aeolyzer_sessions (app_name, user_id, session_id)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS aeolyzer_attachments (
    app_name VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    session_id VARCHAR(64) NOT NULL,
    attachment_id VARCHAR(64) NOT NULL,
    name_ciphertext BLOB NOT NULL,
    content_type VARCHAR(128) NOT NULL,
    byte_size BIGINT UNSIGNED NOT NULL,
    sha256 BINARY(32) NOT NULL,
    data_ciphertext LONGBLOB NOT NULL,
    created_at DATETIME(6) NOT NULL,
    PRIMARY KEY (app_name, user_id, session_id, attachment_id),
    CONSTRAINT aeolyzer_attachments_session_fk
        FOREIGN KEY (app_name, user_id, session_id)
        REFERENCES aeolyzer_sessions (app_name, user_id, session_id)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

CREATE TABLE IF NOT EXISTS aeolyzer_message_requests (
    app_name VARCHAR(64) NOT NULL,
    user_id VARCHAR(64) NOT NULL,
    session_id VARCHAR(64) NOT NULL,
    request_hash BINARY(32) NOT NULL,
    status VARCHAR(16) NOT NULL,
    response_ciphertext LONGBLOB NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    PRIMARY KEY (app_name, user_id, session_id, request_hash),
    CONSTRAINT aeolyzer_requests_session_fk
        FOREIGN KEY (app_name, user_id, session_id)
        REFERENCES aeolyzer_sessions (app_name, user_id, session_id)
        ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
