CREATE TABLE order_items (
    id        uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id  uuid NOT NULL,
    item_id   uuid NOT NULL,
    quantity  INT NOT NULL,
    price     DECIMAL(10,2) NOT NULL,

    CONSTRAINT fk_order_id
    FOREIGN KEY (order_id)
        REFERENCES orders (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_item_id
    FOREIGN KEY (item_id)
        REFERENCES items (id)
        ON DELETE CASCADE
);
