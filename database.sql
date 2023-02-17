CREATE TABLE public.users
(
    id         serial  NOT NULL,
    username   varchar NOT NULL,
    "password" varchar NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY (id),
    CONSTRAINT users_username_un UNIQUE (username)
);
