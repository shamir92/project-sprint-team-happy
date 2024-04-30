CREATE TABLE matches (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
	issuer_cat_id UUID NOT NULL,
  	receiver_cat_id UUID NOT NULL,
  	issuer_id UUID NOT NULL,
  	receiver_id UUID NOT NULL,
  	"message" VARCHAR(120) NOT NULL,
  	"status" VARCHAR(20) NOT NULL,
  	created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  	deleted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE matches
ADD CONSTRAINT fk_matches_issuer_id_to_users_id
FOREIGN KEY (issuer_id) 
REFERENCES "users"(id);

ALTER TABLE matches
ADD CONSTRAINT fk_matches_receiver_id_to_users_id
FOREIGN KEY (receiver_id) 
REFERENCES "users"(id);

ALTER TABLE matches
ADD CONSTRAINT fk_matches_receiver_cat_id_to_cats_id
FOREIGN KEY (receiver_cat_id) 
REFERENCES "cats"(id);

ALTER TABLE matches
ADD CONSTRAINT fk_matches_issuer_cat_id_to_cats_id
FOREIGN KEY (issuer_cat_id) 
REFERENCES "cats"(id);