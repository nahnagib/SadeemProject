CREATE TABLE vendor_admins (
    user_id    uuid NOT NULL,
    vendor_id  uuid NOT NULL,

    PRIMARY KEY (user_id, vendor_id),

    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id)
            ON DELETE CASCADE,

    CONSTRAINT fk_vendor_id
        FOREIGN KEY (vendor_id)
            REFERENCES vendors (id)
            ON DELETE CASCADE
);
