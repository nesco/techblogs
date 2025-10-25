-- Seed blog configurations with CSS selectors
INSERT INTO blog_configs (
    blog_name, blog_href, kind, article_selector
) VALUES
(
    'Stripe',
    'https://stripe.com/blog',
    'organization',
    '.BlogHero__layout article h1 a'
),
(
    'Datadoghq',
    'https://www.datadoghq.com/blog/',
    'organization',
    'main article a'
),
(
    'Hillel Wayne',
    'https://buttondown.com/hillelwayne/archive/',
    'person',
    '.email-list a:nth-child(2)'
),
(
    'Sam Altman',
    'https://blog.samaltman.com/',
    'person',
    '#main article:first-child h2 a'
);
