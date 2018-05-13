create table page (
    id serial,
    uri text unique,
    title text,
    tags text[],
    body text
);
