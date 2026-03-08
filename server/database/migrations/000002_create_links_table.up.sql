CREATE TABLE IF NOT EXISTS tags (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS links (
    id UUID PRIMARY KEY DEFAULT uuidv7(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    link TEXT NOT NULL,
    name VARCHAR(150) NOT NULL,
    link_desc VARCHAR(200) NULL
);

CREATE TABLE IF NOT EXISTS link_tags (
    link_id UUID NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    tag_id  UUID NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
    PRIMARY KEY (link_id, tag_id)
);

-- add indexes
CREATE INDEX tag_id_index ON tags (id);
CREATE INDEX tag_userid_index ON tags (user_id);
CREATE INDEX link_id_index ON links (id);
CREATE INDEX link_userid_index ON links (user_id);
CREATE INDEX linktag_link_id_index ON link_tags (link_id);
CREATE INDEX linktag_tagid_index ON link_tags (tag_id);