CREATE TABLE todo_lists
(
    id          serial       not null primary key,
--     user_id     int          not null, -- внешний ключ логический (fk на users/id в auth service)
    title       varchar(255) not null,
    description varchar(255)
);

CREATE TABLE todo_items
(
    id          serial       not null primary key,
--     list_id     int          not null REFERENCES todo_lists (id) ON DELETE CASCADE,
    title       varchar(255) not null,
    description varchar(255),
    done        boolean      not null default false
);

CREATE TABLE users_lists
(
    id      serial PRIMARY KEY,
    user_id int not null,
    list_id int not null REFERENCES todo_lists (id) ON DELETE CASCADE
);

CREATE TABLE lists_items
(
    id      serial PRIMARY KEY,
    item_id int not null REFERENCES todo_items (id) ON DELETE CASCADE,
    list_id int not null REFERENCES todo_lists (id) ON DELETE CASCADE
);