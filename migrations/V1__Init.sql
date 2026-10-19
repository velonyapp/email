CREATE TABLE templates (
    id      CHAR(36) PRIMARY KEY,
    alias   VARCHAR(64) NOT NULL,
    subject TEXT NOT NULL,
    html    TEXT NOT NULL,
    text    TEXT NOT NULL,

    INDEX idx_templates_alias (alias)
) ENGINE = InnoDB;
