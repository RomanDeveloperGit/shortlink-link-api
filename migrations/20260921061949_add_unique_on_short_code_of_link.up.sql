ALTER TABLE links
ADD CONSTRAINT links_short_link_unique UNIQUE (short_code);
