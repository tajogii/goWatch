CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE room (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    size SMALLINT ,
    password VARCHAR(255) NOT NULL
);