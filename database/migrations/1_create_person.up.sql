CREATE TYPE sex AS ENUM ('male', 'female');
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS person (
    id SERIAL,
    uuid uuid default uuid_generate_v4 (),
    name text, 
    surname text,
    patronymic text,
    age integer,
    gender sex,
    nation text
);