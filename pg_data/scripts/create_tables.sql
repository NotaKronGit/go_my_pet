-- create users table
CREATE TABLE IF NOT EXISTS USERS (
    id INT PRIMARY KEY,
    name VARCHAR(10)
);

-- insert initial default users
INSERT INTO USERS (id,name) VALUES (0,'USER'),(1,'ADMIN');