CREATE TABLE carts (
    id               uuid PRIMARY KEY NOT NULL UNIQUE,
    total_price      DECIMAL(10,2) NOT NULL DEFAULT 0,
    quantity         INT NOT NULL DEFAULT 0,
    vendor_id        uuid DEFAULT NULL,
    created_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_user_id
    FOREIGN KEY (id)
        REFERENCES users (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_vendor_id
    FOREIGN KEY (vendor_id)
        REFERENCES vendors (id)
        ON DELETE CASCADE
);
