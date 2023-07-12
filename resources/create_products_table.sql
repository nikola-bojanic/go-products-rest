CREATE TABLE products (
    id BIGSERIAL NOT NULL PRIMARY KEY,
    name VARCHAR(75) NOT NULL,
    short_description TEXT NOT NULL,
    description TEXT NOT NULL,
    price NUMERIC(19, 2) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO products (name, short_description, description, price) 
VALUES ('Sony WF-1000XM4', 'Wireless Earbuds', 'Sony’s WF-1000XM4 have top-notch noise cancellation and lively, enjoyable sound quality. With wireless charging and bonus features like LDAC support, they’re a great overall earphones.', 100.50);
INSERT INTO products (name, short_description, description, price) 
VALUES ('MacBook Pro M2', 'Laptop', 'Apple M2 is a series of ARM-based system on a chip (SoC) designed by Apple Inc. as a central processing unit (CPU) and graphics processing unit (GPU) for its Mac desktops and notebooks, and the iPad Pro tablet.', 2000.50);
INSERT INTO products (name, short_description, description, price) 
VALUES ('Sony WF-1000XM4', 'Solar Panel', 'The Tesla Powerwall is a rechargeable lithium-ion battery stationary home energy storage product manufactured by Tesla Energy. The Powerwall stores electricity for solar self-consumption, time of use load shifting, and backup power. ', 14100.30);