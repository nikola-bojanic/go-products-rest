CREATE TABLE categories (
    category_id BIGSERIAL NOT NULL PRIMARY KEY,
    category_name VARCHAR(75) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP 
);

INSERT INTO categories (category_name) VALUES ('Smartphone');  
INSERT INTO categories (category_name) VALUES ('Laptop');
INSERT INTO categories (category_name) VALUES ('Solar Panel');

ALTER TABLE products ADD COLUMN category_id BIGINT REFERENCES categories(category_id);

UPDATE products SET category_id = 1 WHERE id = 1;
UPDATE products SET category_id = 2 WHERE id = 2;
UPDATE products SET category_id = 3 WHERE id = 3;
