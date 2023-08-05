CREATE TABLE IF NOT EXISTS books(
    id SERIAL,
    title VARCHAR(250) NOT NULL,
    description VARCHAR(250) NOT NULL,
    price NUMERIC(2) NOT NULL,
    PRIMARY KEY (id)
);
