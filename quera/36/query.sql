show databases;

use mydb;

select * from users;

use quera;

select * from users;


use mydb;

select *
from users;


delete
from products
where id in (2);

select * from users;

show tables;

INSERT INTO products (name, price, color, created_date, create_date)
VALUES ("laptop hp 15 fc", 350, "red",date(now()), date(now()));


select * from products

 SELECT * FROM `products` WHERE `products`.`id` = 3 AND `products`.`deleted_at` IS NULL ORDER BY `products`.`id` LIMIT 1