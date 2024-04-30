CREATE TABLE cats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(30) NOT NULL,
    sex VARCHAR(10) NOT NULL,
    age_in_month INT NOT NULL,
    "description" VARCHAR(20) NOT NULL,
    image_urls TEXT [],
    race VARCHAR(50) NOT NULL,
    created_by UUID NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    has_matched BOOLEAN DEFAULT FALSE,
    owner_id UUID NOT NULL,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone
);
ALTER TABLE cats
ADD CONSTRAINT fk_cats_users FOREIGN KEY (owner_id) REFERENCES "users" (id);