CREATE TABLE matches (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	issuer_cat_id UUID NOT NULL,
	receiver_cat_id UUID NOT NULL,
	issuer_id UUID NOT NULL,
	receiver_id UUID NOT NULL,
	"message" VARCHAR(120) NOT NULL,
	"status" VARCHAR(20) NOT NULL,
	created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at timestamp with time zone,
	deleted_at TIMESTAMP WITH TIME ZONE
);
ALTER TABLE matches
ADD CONSTRAINT fk_matches_issuer_id_to_users_id FOREIGN KEY (issuer_id) REFERENCES "users"(id);
ALTER TABLE matches
ADD CONSTRAINT fk_matches_receiver_id_to_users_id FOREIGN KEY (receiver_id) REFERENCES "users"(id);
ALTER TABLE matches
ADD CONSTRAINT fk_matches_receiver_cat_id_to_cats_id FOREIGN KEY (receiver_cat_id) REFERENCES "cats"(id);
ALTER TABLE matches
ADD CONSTRAINT fk_matches_issuer_cat_id_to_cats_id FOREIGN KEY (issuer_cat_id) REFERENCES "cats"(id);
CREATE INDEX idx_matches_id ON matches (id);
CREATE INDEX idx_matches_deleted_at ON matches (deleted_at);
CREATE INDEX idx_matches_receiver_id ON matches (receiver_id);
CREATE INDEX idx_matches_issuer_cat_id ON matches (issuer_cat_id);
CREATE INDEX idx_matches_receiver_cat_id ON matches (receiver_cat_id);