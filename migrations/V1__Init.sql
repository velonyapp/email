CREATE TABLE templates (
    id      CHAR(36) PRIMARY KEY,
    alias   VARCHAR(64) NOT NULL,
    subject TEXT NOT NULL,
    html    TEXT NOT NULL,
    text    TEXT NOT NULL
) ENGINE = InnoDB;

CREATE INDEX idx_templates_alias ON templates (alias);
