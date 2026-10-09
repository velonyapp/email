CREATE TABLE templates (
    id           VARCHAR(64) PRIMARY KEY,
    subject      TEXT NOT NULL,
    html         TEXT NOT NULL,
    text         TEXT NOT NULL
) ENGINE = InnoDB;
