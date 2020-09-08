create table boards (
    id int primary key auto_increment,
    title varchar(150) not null,
    code varchar(32) not null unique,
    created datetime not null,
    updated datetime not null,
    index boards_idx_1 (code)
);

create table roles (
    role int primary key,
    description text not null
);

insert into roles
    (role, description)
    values (1, "admin");
insert into roles
    (role, description)
    values (2, "user");

create table users (
    uid varchar(128) primary key,
    name text
);

create table board_members (
    id int primary key auto_increment,
    board_id int not null,
    uid varchar(128) not null,
    role int not null,
    foreign key fk_board (board_id) references boards(id),
    foreign key fk_role (role) references roles(role)
);

create table items (
    id int primary key auto_increment,
    board_id int not null,
    y text,
    w text,
    t text,
    author_uid varchar(128),
    created datetime not null,
    updated datetime not null,
    foreign key fk_board (board_id) references boards(id)
);