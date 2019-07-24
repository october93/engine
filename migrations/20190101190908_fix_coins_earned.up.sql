-- fix some coins_earned bugs

-- fix cards' coinsearned
UPDATE cards
SET coins_earned = coins_earned - (coins_earned % 10000) + (coins_earned % 10000 * 10000)
WHERE coins_earned % 10000 != 0;
