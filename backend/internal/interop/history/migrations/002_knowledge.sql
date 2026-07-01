ALTER TABLE aeolyzer_sessions
    ADD COLUMN IF NOT EXISTS content_type VARCHAR(32) NOT NULL DEFAULT '' AFTER starred;

CREATE TABLE IF NOT EXISTS aeolyzer_knowledge (
    user_id VARCHAR(64) NOT NULL,
    section VARCHAR(32) NOT NULL,
    version BIGINT UNSIGNED NOT NULL,
    body_ciphertext LONGBLOB NOT NULL,
    summary_ciphertext LONGBLOB NOT NULL,
    created_at DATETIME(6) NOT NULL,
    updated_at DATETIME(6) NOT NULL,
    PRIMARY KEY (user_id, section),
    INDEX aeolyzer_knowledge_updated (user_id, updated_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;
