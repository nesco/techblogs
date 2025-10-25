-- Seed blog configurations with CSS selectors
INSERT INTO blog_configs (
    blog_name, blog_href, kind, article_href_selector, article_name_selector
) VALUES
(
    'Datadoghq',
    'https://www.datadoghq.com/blog/',
    'organization',
    'main article:first-of-type a:first-of-type',
    'main article:first-of-type .card-header'
),
(
    'Google',
    'https://research.google/blog/',
    'organization',
    '#page-content .blog-index a[href^="/blog"]',
    '#page-content .blog-index a[href^="/blog"] .headline-5'
),
(
    'Jane Street',
    'https://blog.janestreet.com/archive/',
    'organization',
    '.archive .table .cell.title > a',
    '.archive .table .cell.title > a'
),
(
    'Stripe',
    'https://stripe.com/blog',
    'organization',
    'article.BlogIndexPost:first-of-type '
    || '.BlogIndexPost__title a.BlogIndexPost__titleLink',
    'article.BlogIndexPost:first-of-type '
    || '.BlogIndexPost__title a.BlogIndexPost__titleLink'
),
(
    'Alex Edwards',
    'https://www.alexedwards.net/blog',
    'individual',
    '.articles li:first-of-type a',
    '.articles li:first-of-type a'
),
(
    'Dan Luu',
    'https://danluu.com/',
    'individual',
    'ul a[href^="https://danluu.com"]',
    'ul a[href^="https://danluu.com"]'
),
(
    'Martin Kleppman',
    'https://martin.kleppmann.com/archive.html',
    'individual',
    '#content ul > li > a',
    '#content ul > li > a'
),
(
    'Max Bernstein',
    'https://bernsteinbear.com/blog/',
    'individual',
    '.container ul:first-of-type li:first-of-type > a',
    '.container ul:first-of-type li:first-of-type > a'
),
(
    'Michael Stapelberg',
    'https://michael.stapelberg.ch/posts/',
    'individual',
    'main .ArticleList li:first-of-type a',
    'main .ArticleList li:first-of-type a'

),
(
    'Neal Krawetz',
    'https://www.hackerfactor.com/blog/index.php',
    'individual',
    '',
    ''
),
(
    'Evan Hahn',
    'https://evanhahn.com/blog/',
    'individual',
    '.post-list li:first-of-type > a',
    '.post-list li:first-of-type > a'
),
(
    'John D. Cook',
    'https://www.johndcook.com/blog/',
    'individual',
    '#content .entry-title > a',
    '#content .entry-title > a'
),
(
    'Josh W. Comeau',
    'https://www.joshwcomeau.com/',
    'individual',
    '.w124ae9d > a',
    '.w124ae9d > a > span'
),
(
    'Julia Evans',
    'https://jvns.ca/',
    'individual',
    '#content .article-list > a',
    '#content .article-list > a'
),
(
    'Robert C. Martin',
    'https://blog.cleancoder.com/',
    'individual',
    'aside ul li:first-of-type > a',
    'aside ul li:first-of-type > a'
),
(
    'Hasen Judy',
    'https://hasen.substack.com/',
    'individual',
    '.portable-archive-list a[href^="https://hasen.substack.com/p/"]',
    '.portable-archive-list a[href^="https://hasen.substack.com/p/"]'
),
(
    'Hillel Wayne',
    'https://buttondown.com/hillelwayne/archive/',
    'individual',
    '.email-list a',
    '.email-list a .email > div:first-of-type > div:first-of-type'
),
(
    'Rain',
    'https://sunshowers.io/',
    'individual',
    '.posts .index-post.on-list h2 span a',
    '.posts .index-post.on-list h2 span a'
),
(
    'Robin Ward',
    'https://eviltrout.com/blog/',
    'individual',
    '.container ul li:first-of-type > a',
    '.container ul li:first-of-type > a'
),
(
    'Sam Altman',
    'https://blog.samaltman.com/',
    'individual',
    '#main article:first-child h2 a',
    '#main article:first-child h2 a'
),
(
    'Scott Aaronson',
    'https://scottaaronson.blog/',
    'individual',
    '#content .post h2 > a',
    '#content .post h2 > a'
),
(
    'Steve Klabnik',
    'https://steveklabnik.com/writing/',
    'individual',
    '#main-content section:first-of-type li:first-of-type > a',
    '#main-content section:first-of-type li:first-of-type > a'
);
