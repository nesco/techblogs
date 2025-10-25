-- Remove seeded blog configurations
DELETE
FROM blog_configs
WHERE blog_name IN ('Stripe', 'Datadoghq', 'Hillel Wayne', 'Sam Altman');
