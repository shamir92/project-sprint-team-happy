CREATE TABLE cats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(30) NOT NULL,
    sex VARCHAR(10) NOT NULL,
    age_in_month INT NOT NULL,
    "description" VARCHAR(200) NOT NULL,
    image_urls TEXT [],
    race VARCHAR(50) NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE,
    has_matched BOOLEAN DEFAULT FALSE,
    owner_id UUID NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone
);
ALTER TABLE cats
ADD CONSTRAINT fk_cats_users FOREIGN KEY (owner_id) REFERENCES "users" (id);
CREATE INDEX idx_cats_id ON cats (id);
CREATE INDEX idx_cats_owner_id ON cats (owner_id);
CREATE INDEX idx_cats_deleted_at ON cats (deleted_at);