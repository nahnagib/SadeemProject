CREATE TABLE tables (
    id                      uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    vendor_id               uuid NOT NULL,
    name                    VARCHAR(255) NOT NULL,
    is_available            BOOLEAN DEFAULT TRUE,
    customer_id             uuid DEFAULT NULL,
    is_needs_service        BOOLEAN DEFAULT FALSE,

    CONSTRAINT fk_vendor_id
    FOREIGN KEY (vendor_id)
        REFERENCES vendors (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_customer_id
    FOREIGN KEY (customer_id)
        REFERENCES users (id)
        ON DELETE CASCADE
);