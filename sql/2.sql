SELECT
    u.id AS user_id,
    u.first_name,
    u.last_name,
    p.sku
FROM users u
JOIN purchases p ON p.user_id = u.id
LEFT JOIN ban_list bl ON bl.user_id = u.id
WHERE p.date < bl.date_from OR bl.date_from IS NULL
ORDER BY u.first_name ASC, p.sku ASC
