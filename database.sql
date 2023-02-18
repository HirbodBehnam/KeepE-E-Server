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
    id               serial    NOT NULL,
    for_user         int4      NOT NULL,
    title            text      NOT NULL,
    note_text        text      NOT NULL,
    background_color varchar(6) NULL,
    created_at       timestamptz NOT NULL,
    deadline         timestamptz,
    -- Order of this note in user notes
    note_order       int4      NOT NULL,
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
