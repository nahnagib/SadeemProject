CREATE TYPE order_status AS ENUM ('completed', 'preparing');

CREATE TABLE orders (
    id                uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    total_order_cost  DECIMAL(10,2) NOT NULL,
    customer_id       uuid NOT NULL,
    vendor_id         uuid NOT NULL,
    status            order_status NOT NULL,
    created_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_vendor_id
    FOREIGN KEY (vendor_id)
        REFERENCES vendors (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_customer_id
    FOREIGN KEY (customer_id)
        REFERENCES users (id)
        ON DELETE CASCADE
);
