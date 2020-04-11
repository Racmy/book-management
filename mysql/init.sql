use book-management;

CREATE TABLE user (
  id int (10) unsigned not null auto_increment,
  email varchar(256) not null unique,
  password varchar(256) not null,
  name VARCHAR(32) not null unique,
  user_image_path VARCHAR(128) not null default "/static/img/no_img.jpeg",
  active bit(1) not null default 1,
  created_at TIMESTAMP default CURRENT_TIMESTAMP,
  update_at TIMESTAMP default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO user (id, email, password, name,created_at,update_at) VALUES (1, 'aaa@bbb.com', 'p@ssword', 'admin' ,Now(),Now());

CREATE TABLE book (
  id int(10) unsigned not null auto_increment,
  user_id int(10) unsigned not null,
  title varchar(128) not null,
  author varchar(128),
  latest_issue float unsigned not null default 1,
  front_cover_image_path varchar(128) not null default "/static/img/no_img.jpeg",
  active bit(1) not null default 1,
  created_at TIMESTAMP default CURRENT_TIMESTAMP,
  update_at TIMESTAMP default CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  primary key (id),
  foreign key fk_user_id(user_id) REFERENCES user(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT INTO book (id, user_id, title, author, front_cover_image_path, latest_issue) VALUES (1, 1,'ジョジョの奇妙な冒険', '荒木比呂彦', '/static/img/jojo.jpg' ,25);
INSERT INTO book (id, user_id, title, author, front_cover_image_path, latest_issue) VALUES (2, 1,'ナルト', '岸本 斉史', '/static/img/naruto.jpg', 45);
INSERT INTO book (id, user_id, title, author, front_cover_image_path, latest_issue) VALUES (3, 1,'テニスの王子様', '許斐剛', '/static/img/tenipuri.jpg',21);

