-- Remove seeded blog configurations
DELETE FROM blog_configs
WHERE blog_name IN (
    'Datadoghq',
    'Google',
    'Jane Street',
    'Stripe',
    'Alex Edwards',
    'Dan Luu',
    'Martin Kleppman',
    'Max Bernstein',
    'Michael Stapelberg',
    'Neal Krawetz',
    'Evan Hahn',
    'John D. Cook',
    'Josh W. Comeau',
    'Julia Evans',
    'Robert C. Martin',
    'Hasen Judy',
    'Hillel Wayne',
    'Rain',
    'Robin Ward',
    'Sam Altman',
    'Scott Aaronson',
    'Steve Klabnik'
);
