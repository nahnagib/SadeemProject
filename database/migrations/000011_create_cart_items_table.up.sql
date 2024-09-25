CREATE TABLE cart_items (
    cart_id  uuid NOT NULL,
    item_id   uuid NOT NULL,
    quantity INT NOT NULL,

    PRIMARY KEY (cart_id, item_id),

    CONSTRAINT fk_cart_id
    FOREIGN KEY (cart_id)
        REFERENCES carts (id)
        ON DELETE CASCADE,

    CONSTRAINT fk_item_id
    FOREIGN KEY (item_id)
        REFERENCES items (id)
        ON DELETE CASCADE
);
