-- +goose Up
create table user (
    id serial primary key,
    name text not null,
    phone_number text not null,
    hashed_password text not null,
    role integer not null default 0;
    created_at timestamp not null default now(),
    updated_at timestamp
);

-- +goose Down
drop table user;

