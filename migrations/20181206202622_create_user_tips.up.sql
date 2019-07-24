CREATE TABLE user_tips (
  id uuid PRIMARY KEY,
  user_id uuid NOT NULL REFERENCES users(id),
  card_id uuid NOT NULL REFERENCES cards(id),
  alias_id uuid REFERENCES anonymous_aliases(id),
  anonymous bool NOT NULL DEFAULT false,
  amount int NOT NULL DEFAULT 0,

  updated_at timestamp without time zone NOT NULL DEFAULT now(),
  created_at timestamp without time zone NOT NULL DEFAULT now()
);
