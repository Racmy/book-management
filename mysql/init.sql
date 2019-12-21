use book-management;

CREATE TABLE book (
  id int(10) unsigned not null auto_increment,
  title varchar(128) not null,
  author varchar(128),
  latest_issue float unsigned not null default 1,
  front_cover_image_path varchar(128) not null default "static/img/no_img.jpeg",
  primary key (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO book (id, title, author, front_cover_image_path, latest_issue) VALUES (1, 'ジョジョの奇妙な冒険', '荒木比呂彦', 'static/img/jojo.jpg' ,25);
INSERT INTO book (id, title, author, front_cover_image_path, latest_issue) VALUES (2, 'ナルト', '岸本 斉史', 'static/img/naruto.jpg', 45);
INSERT INTO book (id, title, author, front_cover_image_path, latest_issue) VALUES (3, 'テニスの王子様', '許斐剛', '/static/img/tenipuri.jpg',21);
