CREATE INDEX user_reactions_card_id_type_idx ON user_reactions (card_id, type);
CREATE INDEX cards_thread_root_id_idx ON cards (thread_root_id);
CREATE INDEX cards_owner_id_idx ON cards (owner_id);
