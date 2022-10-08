begin;

create table users (
    id serial,
    username varchar(100) primary key
);

create table userdata (
    userid int not null,
    firstName varchar(100),
    lastName varchar(100),
    description varchar(200)
);

commit;
