BEGIN;

CREATE TABLE IF NOT EXISTS users (
    id VARCHAR(26) NOT NULL,
    username VARCHAR(12) NOT NULL,
    email VARCHAR(25) NOT NULL,
    passhash VARCHAR NOT NULL,
    role VARCHAR(25) NOT NULL
);

CREATE TABLE IF NOT EXISTS luxury_items (
    id VARCHAR(26) NOT NULL,
    brand CHAR(26) NOT NULL,
    price INT NOT NULL,
    ownerid VARCHAR(26),
    tokenURI VARCHAR(64) NOT NULL
);

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(26) NOT NULL,
    txHash VARCHAR(64) NOT NULL,
    txType VARCHAR(24) NOT NULL
);

-- CREATE A MANUFACTURER LOGIN
INSERT INTO users(id, username, email, passhash,role ) VALUES ('0ZA0mtj31W3O-U3o7c8_3', 'Sree', 'sree@kaleido.com', '$2a$14$t/8e.sR0PsGtZYFhwqq31ualylWP7JbDWY0mJxn30n0.PzV.Fsobm', 'manufacturer');

-- CREATE LUXURY_ITEMS
INSERT INTO luxury_items(id, brand, price, ownerid,tokenURI ) VALUES ('339313', 'kaleido', 4599, '0ZA0mtj31W3O-U3o7c8_3', 'WBIBGFC5MYQ3170P');
INSERT INTO luxury_items(id, brand, price, ownerid,tokenURI ) VALUES ('417873', 'kaleido', 8999, '0ZA0mtj31W3O-U3o7c8_3', '51H3SD967RWBR8GH');
INSERT INTO luxury_items(id, brand, price, ownerid,tokenURI ) VALUES ('155582', 'kaleido', 3599, '0ZA0mtj31W3O-U3o7c8_3', 'FVXOJW4JXTCTCG0E');
INSERT INTO luxury_items(id, brand, price, ownerid,tokenURI ) VALUES ('826992', 'kaleido', 5999, '0ZA0mtj31W3O-U3o7c8_3', 'IB0WO8RX5CD59ATP');
INSERT INTO luxury_items(id, brand, price, ownerid,tokenURI ) VALUES ('018670', 'kaleido', 6001, '0ZA0mtj31W3O-U3o7c8_3', 'BR8316PD2I9BRVC4');

COMMIT;
