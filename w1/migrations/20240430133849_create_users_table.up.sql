CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    fullname VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE,
    password CHAR(72) NOT NULL
);