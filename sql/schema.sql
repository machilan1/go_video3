create table app_user (
    id bigserial primary key ,
    name text not null,
    email text not null unique,
    password text not null,
    created_at timestamptz not null default now(),
    last_login timestamptz,
    is_admin boolean not null default false,
    is_active boolean not null default true,
    can_update boolean not null default false
);


create table course(
    id bigserial primary key,
    name text not null,
    description text,
    instructor_name text not null,
    created_at timestamptz not null default now(),
    created_by bigint references app_user(id) on delete set null,
    deleted_at timestamptz,
    click_count int not null default 0,
    updated_at timestamptz not null default now()
);

CREATE TABLE sessions (
	token TEXT PRIMARY KEY,
	data BYTEA NOT NULL,
	expiry TIMESTAMPTZ NOT NULL
);

CREATE INDEX sessions_expiry_idx ON sessions (expiry);


create table video(
    id bigserial primary key,
    file_name text not null,
    description text,
    created_at timestamptz not null default now(),
    deleted_at timestamptz,
    updated_by bigint references app_user(id) on delete set null
);



create table chapter(
    id bigserial primary key,
    title text not null,
    description text,
    chap_num int not null default 1 check (chap_num > 0),
    created_at timestamptz not null default now(),
    course_id bigint references course(id) on delete cascade,
    video_id bigint references video(id) on delete cascade
);

create table tag(
    id bigserial primary key,
    label text not null,
    created_at timestamptz not null default now(),
    created_by bigint references app_user(id) on delete set null
);

create table course_tag(
    course_id bigint references course(id) on delete cascade,
    tag_id bigint references tag(id) on delete cascade,
    primary key (course_id, tag_id)
) ;

create table favorite(
    user_id bigint references app_user(id) on delete cascade,
    course_id bigint references course(id) on delete cascade,
    primary key (user_id, course_id)
);

create table history(
    chapter_id bigint references chapter(id) on delete cascade,
    user_id bigint references app_user(id) on delete cascade,
    updated_at timestamptz not null default now(),
    primary key (chapter_id, user_id)
);


CREATE OR REPLACE FUNCTION course_update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    IF (NEW.* IS DISTINCT FROM OLD.*) AND (NEW.click_count = OLD.click_count) THEN
        NEW.updated_at = NOW();
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;
