-- deduct all self tips from card amounts
UPDATE cards
SET coins_earned = coins_earned - selftip.amount
FROM (
  SELECT card_id, sum(amount) AS amount
  FROM coin_transactions
  WHERE type = 'tipGiven' AND source_user_id = recipient_user_id
  GROUP BY card_id
) AS selftip
WHERE cards.id = selftip.card_id;

-- delete all self tips from user_tips
DELETE FROM user_tips
WHERE id IN (
  SELECT user_tips.id
  FROM user_tips
    LEFT JOIN cards ON user_tips.card_id = cards.id
  WHERE user_tips.user_id = cards.owner_id
);
