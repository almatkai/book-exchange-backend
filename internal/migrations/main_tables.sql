-- Create enum types
CREATE TYPE book_condition AS ENUM ('New', 'Like New', 'Very Good', 'Good', 'Fair', 'Poor');
CREATE TYPE exchange_status AS ENUM ('Pending', 'Accepted', 'Rejected', 'Completed', 'Cancelled');

-- Users table
CREATE TABLE users (
                       user_id SERIAL PRIMARY KEY,
                       username VARCHAR(50) UNIQUE NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password_hash VARCHAR(255) NOT NULL,
                       full_name VARCHAR(100),
                       location VARCHAR(255),
                       bio TEXT,
                       rating DECIMAL(3,2) CHECK (rating >= 0 AND rating <= 5),
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       last_login TIMESTAMP WITH TIME ZONE,
                       is_active BOOLEAN DEFAULT true
);

-- Books table
CREATE TABLE books (
                       book_id SERIAL PRIMARY KEY,
                       owner_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
                       title VARCHAR(255) NOT NULL,
                       author VARCHAR(255) NOT NULL,
                       isbn VARCHAR(13),
                       condition book_condition NOT NULL,
                       description TEXT,
                       genre VARCHAR(50),
                       language VARCHAR(50),
                       publication_year INTEGER,
                       is_available BOOLEAN DEFAULT true,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       image_url VARCHAR(255)
);

-- Books table could benefit from:
ALTER TABLE books ADD CONSTRAINT valid_publication_year
    CHECK (publication_year <= EXTRACT(YEAR FROM CURRENT_DATE));

-- Book exchanges table
CREATE TABLE exchanges (
                           exchange_id SERIAL PRIMARY KEY,
                           requester_id INTEGER REFERENCES users(user_id),
                           book_id INTEGER REFERENCES books(book_id),
                           status exchange_status DEFAULT 'Pending',
                           requested_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                           completed_at TIMESTAMP WITH TIME ZONE,
                           meeting_location VARCHAR(255),
                           meeting_time TIMESTAMP WITH TIME ZONE,
                           notes TEXT
);

-- Reviews table
CREATE TABLE reviews (
                         review_id SERIAL PRIMARY KEY,
                         reviewer_id INTEGER REFERENCES users(user_id),
                         reviewed_user_id INTEGER REFERENCES users(user_id),
                         exchange_id INTEGER REFERENCES exchanges(exchange_id),
                         rating INTEGER CHECK (rating >= 1 AND rating <= 5),
                         comment TEXT,
                         created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Wishlists table
CREATE TABLE wishlists (
                           wishlist_id SERIAL PRIMARY KEY,
                           user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
                           title VARCHAR(255) NOT NULL,
                           author VARCHAR(255),
                           created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Messages table for communication between users
CREATE TABLE messages (
                          message_id SERIAL PRIMARY KEY,
                          sender_id INTEGER REFERENCES users(user_id),
                          receiver_id INTEGER REFERENCES users(user_id),
                          exchange_id INTEGER REFERENCES exchanges(exchange_id),
                          content TEXT NOT NULL,
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          is_read BOOLEAN DEFAULT false
);

-- Create indexes for better query performance
CREATE INDEX idx_books_owner ON books(owner_id);
CREATE INDEX idx_books_available ON books(is_available);
CREATE INDEX idx_exchanges_status ON exchanges(status);
CREATE INDEX idx_exchanges_requester ON exchanges(requester_id);
CREATE INDEX idx_reviews_reviewed_user ON reviews(reviewed_user_id);
CREATE INDEX idx_messages_receiver ON messages(receiver_id);
CREATE INDEX idx_messages_exchange ON messages(exchange_id);

-- Trigger to update the updated_at timestamp for books
CREATE OR REPLACE FUNCTION update_updated_at_column()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_books_updated_at
    BEFORE UPDATE ON books
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Function to calculate user rating
CREATE OR REPLACE FUNCTION update_user_rating()
    RETURNS TRIGGER AS $$
BEGIN
    UPDATE users
    SET rating = (
        SELECT ROUND(AVG(rating)::numeric, 2)
        FROM reviews
        WHERE reviewed_user_id = NEW.reviewed_user_id
    )
    WHERE user_id = NEW.reviewed_user_id;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_user_rating_trigger
    AFTER INSERT OR UPDATE ON reviews
    FOR EACH ROW
EXECUTE FUNCTION update_user_rating();


-- Create enum for different comment target types
CREATE TYPE commentable_type AS ENUM ('book', 'user', 'exchange', 'review');

-- Comments table with polymorphic association
CREATE TABLE comments (
                          comment_id SERIAL PRIMARY KEY,
                          user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
                          target_type commentable_type NOT NULL,
                          target_id INTEGER NOT NULL,  -- ID of the entity being commented on
                          content TEXT NOT NULL,
                          parent_comment_id INTEGER REFERENCES comments(comment_id), -- For nested comments
                          created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                          is_edited BOOLEAN DEFAULT false,
                          is_deleted BOOLEAN DEFAULT false
);

-- Create indexes for better query performance
CREATE INDEX idx_comments_user ON comments(user_id);
CREATE INDEX idx_comments_target ON comments(target_type, target_id);
CREATE INDEX idx_comments_parent ON comments(parent_comment_id);

-- Comment reactions table (likes, etc.)
CREATE TABLE comment_reactions (
                                   reaction_id SERIAL PRIMARY KEY,
                                   comment_id INTEGER REFERENCES comments(comment_id) ON DELETE CASCADE,
                                   user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
                                   reaction_type VARCHAR(20) NOT NULL, -- 'like', 'dislike', etc.
                                   created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                                   UNIQUE(comment_id, user_id) -- One reaction type per user per comment
);

-- Create index for reactions
CREATE INDEX idx_comment_reactions_comment ON comment_reactions(comment_id);

-- Trigger to update the updated_at timestamp for comments
CREATE TRIGGER update_comments_updated_at
    BEFORE UPDATE ON comments
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Function to get comment count for any entity
CREATE OR REPLACE FUNCTION get_comment_count(
    target_type_param commentable_type,
    target_id_param INTEGER
)
    RETURNS INTEGER AS $$
BEGIN
    RETURN (
        SELECT COUNT(*)
        FROM comments
        WHERE target_type = target_type_param
          AND target_id = target_id_param
          AND NOT is_deleted
    );
END;
$$ LANGUAGE plpgsql;

-- Example queries for different use cases:

-- Get all comments for a specific book
CREATE OR REPLACE FUNCTION get_book_comments(book_id_param INTEGER)
    RETURNS TABLE (
                      comment_id INTEGER,
                      user_id INTEGER,
                      username VARCHAR(50),
                      content TEXT,
                      created_at TIMESTAMP WITH TIME ZONE,
                      reaction_count INTEGER
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            c.comment_id,
            c.user_id,
            u.username,
            c.content,
            c.created_at,
            COUNT(cr.reaction_id)::INTEGER as reaction_count
        FROM comments c
                 JOIN users u ON c.user_id = u.user_id
                 LEFT JOIN comment_reactions cr ON c.comment_id = cr.comment_id
        WHERE c.target_type = 'book'
          AND c.target_id = book_id_param
          AND NOT c.is_deleted
        GROUP BY c.comment_id, c.user_id, u.username, c.content, c.created_at
        ORDER BY c.created_at DESC;
END;
$$ LANGUAGE plpgsql;

-- Get threaded comments (with replies)
CREATE OR REPLACE FUNCTION get_threaded_comments(
    target_type_param commentable_type,
    target_id_param INTEGER,
    max_depth INTEGER DEFAULT 5,
    page_size INTEGER DEFAULT 20,
    page_number INTEGER DEFAULT 1
)
    RETURNS TABLE (
                      comment_id INTEGER,
                      parent_id INTEGER,
                      user_id INTEGER,
                      username VARCHAR(50),
                      content TEXT,
                      created_at TIMESTAMP WITH TIME ZONE,
                      level INTEGER,
                      reaction_count INTEGER,
                      reply_count INTEGER,
                      path INTEGER[],
                      sort_order BIGINT
                  ) AS $$
BEGIN
    RETURN QUERY
        WITH RECURSIVE comment_tree AS (
            -- Base case: top-level comments
            SELECT
                c.comment_id,
                c.parent_comment_id AS parent_id,
                c.user_id,
                u.username,
                c.content,
                c.created_at,
                1 AS level,
                ARRAY[c.comment_id] AS path,
                c.created_at::BIGINT AS sort_order,
                COUNT(cr.reaction_id)::INTEGER as reaction_count,
                (SELECT COUNT(*)
                 FROM comments replies
                 WHERE replies.parent_comment_id = c.comment_id
                   AND NOT replies.is_deleted)::INTEGER as reply_count
            FROM comments c
                     JOIN users u ON c.user_id = u.user_id
                     LEFT JOIN comment_reactions cr ON c.comment_id = cr.comment_id
            WHERE c.target_type = target_type_param
              AND c.target_id = target_id_param
              AND c.parent_comment_id IS NULL
              AND NOT c.is_deleted
            GROUP BY c.comment_id, c.parent_comment_id, c.user_id, u.username,
                     c.content, c.created_at

            UNION ALL

            -- Recursive case: replies
            SELECT
                c.comment_id,
                c.parent_comment_id AS parent_id,
                c.user_id,
                u.username,
                c.content,
                c.created_at,
                ct.level + 1,
                ct.path || c.comment_id,
                ct.sort_order + ROW_NUMBER() OVER (
                    PARTITION BY c.parent_comment_id
                    ORDER BY c.created_at
                    )::BIGINT,
                COUNT(cr.reaction_id)::INTEGER,
                (SELECT COUNT(*)
                 FROM comments replies
                 WHERE replies.parent_comment_id = c.comment_id
                   AND NOT replies.is_deleted)::INTEGER
            FROM comments c
                     JOIN users u ON c.user_id = u.user_id
                     JOIN comment_tree ct ON ct.comment_id = c.parent_comment_id
                     LEFT JOIN comment_reactions cr ON c.comment_id = cr.comment_id
            WHERE NOT c.is_deleted
              AND ct.level < max_depth
            GROUP BY c.comment_id, c.parent_comment_id, c.user_id, u.username,
                     c.content, c.created_at, ct.level, ct.path, ct.sort_order
        )
        SELECT
            comment_id,
            parent_id,
            user_id,
            username,
            content,
            created_at,
            level,
            reaction_count,
            reply_count,
            path,
            sort_order
        FROM comment_tree
        ORDER BY path, created_at
        LIMIT page_size
            OFFSET (page_number - 1) * page_size;
END;
$$ LANGUAGE plpgsql;

-- Create enum for post types
CREATE TYPE post_type AS ENUM ('short_review', 'long_review', 'discussion', 'recommendation');

-- Posts table
CREATE TABLE posts (
                       post_id SERIAL PRIMARY KEY,
                       user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
                       book_id INTEGER REFERENCES books(book_id),
                       type post_type NOT NULL,
                       title VARCHAR(255) NOT NULL,
                       content TEXT NOT NULL,
                       summary TEXT, -- Short description or excerpt
                       rating INTEGER CHECK (rating >= 1 AND rating <= 5),
                       likes_count INTEGER DEFAULT 0,
                       views_count INTEGER DEFAULT 0,
                       created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                       published_at TIMESTAMP WITH TIME ZONE,
                       is_published BOOLEAN DEFAULT true,
                       is_featured BOOLEAN DEFAULT false,
                       is_edited BOOLEAN DEFAULT false
);

-- Tags for posts
CREATE TABLE tags (
                      tag_id SERIAL PRIMARY KEY,
                      name VARCHAR(50) UNIQUE NOT NULL,
                      description TEXT,
                      created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Post tags relation
CREATE TABLE post_tags (
                           post_id INTEGER REFERENCES posts(post_id) ON DELETE CASCADE,
                           tag_id INTEGER REFERENCES tags(tag_id) ON DELETE CASCADE,
                           PRIMARY KEY (post_id, tag_id)
);

-- Post likes
CREATE TABLE post_likes (
                            post_id INTEGER REFERENCES posts(post_id) ON DELETE CASCADE,
                            user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            PRIMARY KEY (post_id, user_id)
);

-- Post saves (bookmarks)
CREATE TABLE post_saves (
                            post_id INTEGER REFERENCES posts(post_id) ON DELETE CASCADE,
                            user_id INTEGER REFERENCES users(user_id) ON DELETE CASCADE,
                            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                            PRIMARY KEY (post_id, user_id)
);

-- Create indexes
CREATE INDEX idx_posts_user ON posts(user_id);
CREATE INDEX idx_posts_book ON posts(book_id);
CREATE INDEX idx_posts_type ON posts(type);
CREATE INDEX idx_posts_published ON posts(is_published, published_at);
CREATE INDEX idx_posts_featured ON posts(is_featured);
CREATE INDEX idx_post_tags_tag ON post_tags(tag_id);

-- Update triggers
CREATE TRIGGER update_posts_updated_at
    BEFORE UPDATE ON posts
    FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

-- Function to update post likes count
CREATE OR REPLACE FUNCTION update_post_likes_count()
    RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        UPDATE posts
        SET likes_count = likes_count + 1
        WHERE post_id = NEW.post_id;
    ELSIF TG_OP = 'DELETE' THEN
        UPDATE posts
        SET likes_count = likes_count - 1
        WHERE post_id = OLD.post_id;
    END IF;
    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_post_likes_count_trigger
    AFTER INSERT OR DELETE ON post_likes
    FOR EACH ROW
EXECUTE FUNCTION update_post_likes_count();

-- Function to get post with tags
CREATE OR REPLACE FUNCTION get_post_with_tags(post_id_param INTEGER)
    RETURNS TABLE (
                      post_id INTEGER,
                      user_id INTEGER,
                      username VARCHAR(50),
                      book_id INTEGER,
                      book_title VARCHAR(255),
                      post_type post_type,
                      title VARCHAR(255),
                      content TEXT,
                      summary TEXT,
                      rating INTEGER,
                      likes_count INTEGER,
                      views_count INTEGER,
                      created_at TIMESTAMP WITH TIME ZONE,
                      tags TEXT[]
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            p.post_id,
            p.user_id,
            u.username,
            p.book_id,
            b.title as book_title,
            p.type,
            p.title,
            p.content,
            p.summary,
            p.rating,
            p.likes_count,
            p.views_count,
            p.created_at,
            ARRAY_AGG(t.name) as tags
        FROM posts p
                 JOIN users u ON p.user_id = u.user_id
                 LEFT JOIN books b ON p.book_id = b.book_id
                 LEFT JOIN post_tags pt ON p.post_id = pt.post_id
                 LEFT JOIN tags t ON pt.tag_id = t.tag_id
        WHERE p.post_id = post_id_param
        GROUP BY
            p.post_id, p.user_id, u.username, p.book_id,
            b.title, p.type, p.title, p.content, p.summary,
            p.rating, p.likes_count, p.views_count, p.created_at;
END;
$$ LANGUAGE plpgsql;

-- Function to get trending posts
CREATE OR REPLACE FUNCTION get_trending_posts(
    days_param INTEGER DEFAULT 7,
    limit_param INTEGER DEFAULT 10
)
    RETURNS TABLE (
                      post_id INTEGER,
                      title VARCHAR(255),
                      summary TEXT,
                      likes_count INTEGER,
                      views_count INTEGER,
                      comment_count BIGINT,
                      created_at TIMESTAMP WITH TIME ZONE
                  ) AS $$
BEGIN
    RETURN QUERY
        SELECT
            p.post_id,
            p.title,
            p.summary,
            p.likes_count,
            p.views_count,
            COUNT(c.comment_id) as comment_count,
            p.created_at
        FROM posts p
                 LEFT JOIN comments c ON c.target_type = 'post' AND c.target_id = p.post_id
        WHERE p.is_published
          AND p.created_at >= CURRENT_TIMESTAMP - (days_param || ' days')::INTERVAL
        GROUP BY p.post_id, p.title, p.summary, p.likes_count, p.views_count, p.created_at
        ORDER BY (p.likes_count + p.views_count + COUNT(c.comment_id)) DESC
        LIMIT limit_param;
END;
$$ LANGUAGE plpgsql;

-- Reviews table could benefit from:
ALTER TABLE reviews ADD CONSTRAINT no_self_reviews
    CHECK (reviewer_id != reviewed_user_id);