SELECT 
    node,
    CASE
        WHEN parent IS NULL THEN 'Root'
        WHEN node NOT IN(
            SELECT parent
            FROM tree
            WHERE parent IS NOT NULL
        ) THEN 'Leaf'
        ELSE 'Inner'
    END
    as label
FROM tree
