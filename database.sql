CREATE TABLE public.users
(
    id         serial  NOT NULL,
    username   varchar NOT NULL,
    "password" varchar NOT NULL,
    CONSTRAINT users_pk PRIMARY KEY (id),
    CONSTRAINT users_username_un UNIQUE (username)
);
CREATE TABLE public.normal_notes
(
    id               serial      NOT NULL,
    for_user         int4        NOT NULL,
    title            text        NOT NULL,
    note_text        text        NOT NULL,
    background_color varchar(6) NULL,
    created_at       timestamptz NOT NULL,
    deadline         timestamptz,
    -- Order of this note in user notes
    note_order       int4        NOT NULL,
    CONSTRAINT normal_notes_pk PRIMARY KEY (id),
    CONSTRAINT normal_notes_user_fk FOREIGN KEY (for_user) REFERENCES public.users (id) ON DELETE CASCADE
);
CREATE TABLE public.normal_note_photos
(
    id       uuid    NOT NULL,
    filename varchar NOT NULL,
    for_note int4    NOT NULL,
    CONSTRAINT normal_notes_photos_pk PRIMARY KEY (id),
    CONSTRAINT normal_notes_photos_note_fk FOREIGN KEY (for_note) REFERENCES public.normal_notes (id) ON DELETE CASCADE
);
CREATE TABLE public.todo
(
    id         serial      NOT NULL,
    for_user   int4        NOT NULL,
    title      text        NOT NULL,
    created_at timestamptz NOT NULL,
    -- Order of this note in user notes
    todo_order int4        NOT NULL,
    CONSTRAINT todo_pk PRIMARY KEY (id),
    CONSTRAINT todo_user_fk FOREIGN KEY (for_user) REFERENCES public.users (id) ON DELETE CASCADE
);
CREATE TABLE public.todo_items
(
    -- We ensure that items of a todo_list are in the same order by deleting everything and adding them again
    id        serial  NOT NULL,
    for_todo  int4    NOT NULL,
    todo_text text    NOT NULL,
    done      boolean NOT NULL,
    deadline  timestamptz,
    CONSTRAINT todo_items_pk PRIMARY KEY (id),
    CONSTRAINT todo_items_todo_fk FOREIGN KEY (for_todo) REFERENCES public.todo (id) ON DELETE CASCADE
);
