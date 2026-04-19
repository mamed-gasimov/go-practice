-- First version
SELECT city FROM city_population
ORDER BY population ASC
LIMIT 1

-- Second version
SELECT city FROM city_population
WHERE population = (SELECT MIN(population) FROM city_population)
