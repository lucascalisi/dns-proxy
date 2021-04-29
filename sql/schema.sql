CREATE TABLE `resolved` (
    `uid` INTEGER PRIMARY KEY AUTOINCREMENT,
    `domain` VARCHAR(200) NOT NULL,
    `type` VARCHAR(10) NOT NULL,
    `class` VARCHAR(15) NOT NULL,
    `action` VARCHAR(20) NOT NULL,
    `date` DATE NOT NULL
);
