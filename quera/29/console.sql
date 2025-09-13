use mydb;


create table users (
    id bigint auto_increment primary key ,
    level enum('user', 'admin') not null default 'user',
    name varchar(45) not null,
    family varchar(45) not null ,
    birth_date DATETIME

)

INSERT INTO users (level, name, family, birth_date)
VALUES
       (DEFAULT, 'amin', 'pourseyyedy', '2009-08-07 00:00:00');


-- Section1
CREATE TABLE events (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255),
    description TEXT,
    date DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- Section2
CREATE TABLE event_user (
    user_id BIGINT UNSIGNED,
    event_id BIGINT UNSIGNED,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (event_id) REFERENCES events(id)
);