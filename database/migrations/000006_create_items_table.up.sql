CREATE TABLE items (
    id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id     uuid NOT NULL,
    name          VARCHAR(255) NOT NULL,
    price         DECIMAL(10,2) NOT NULL,
    img           VARCHAR(255),
    created_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_vendor_id
    FOREIGN KEY (vendor_id)
        REFERENCES vendors (id)
        ON DELETE CASCADE
);
