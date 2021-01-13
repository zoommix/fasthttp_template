CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email character varying(255) NOT NULL,
    password_digest character varying(255) NOT NULL,
    first_name character varying(255) NOT NULL,
    last_name character varying(255) NOT NULL,
    created_at timestamp(6) without time zone NOT NULL,
    updated_at timestamp(6) without time zone NOT NULL
);

CREATE UNIQUE INDEX index_users_on_lowercase_email ON users((lower(email::text)) text_ops);

---- create above / drop below ----

drop table users;
