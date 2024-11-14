-- Create users table
CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   username VARCHAR(50) UNIQUE NOT NULL,
   email VARCHAR(100) UNIQUE NOT NULL,
   password_hash VARCHAR(255) NOT NULL,
   created_at TIMESTAMP DEFAULT NOW(),
   updated_at TIMESTAMP
);

-- Create books table
CREATE TABLE books (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   title VARCHAR(255) NOT NULL,
   author VARCHAR(255) NOT NULL,
   description TEXT,
   is_available BOOLEAN DEFAULT TRUE,
   created_at TIMESTAMP DEFAULT NOW(),
   updated_at TIMESTAMP
);

-- Create posts table
CREATE TABLE posts (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   book_id UUID NOT NULL REFERENCES books(id) ON DELETE CASCADE,
   exchange_type VARCHAR(10) CHECK (exchange_type IN ('permanent', 'temporary')) NOT NULL,
   available_until TIMESTAMP,
   location VARCHAR(255) NOT NULL,
   created_at TIMESTAMP DEFAULT NOW(),
   updated_at TIMESTAMP
);

-- Create exchanges table
CREATE TABLE exchanges (
   id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
   post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
   requester_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   status VARCHAR(10) CHECK (status IN ('pending', 'accepted', 'rejected', 'completed')) NOT NULL,
   location VARCHAR(255) NOT NULL,
   exchange_date TIMESTAMP NOT NULL,
   created_at TIMESTAMP DEFAULT NOW(),
   updated_at TIMESTAMP
);

-- Create messages table
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exchange_id UUID NOT NULL REFERENCES exchanges(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    sent_at TIMESTAMP DEFAULT NOW()
);

-- Create ratings table
CREATE TABLE ratings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exchange_id UUID NOT NULL REFERENCES exchanges(id) ON DELETE CASCADE,
    rater_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ratee_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    rating INTEGER CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Additional Indexes for faster lookups
CREATE INDEX idx_books_user_id ON books(user_id);
CREATE INDEX idx_posts_user_id ON posts(user_id);
CREATE INDEX idx_posts_book_id ON posts(book_id);
CREATE INDEX idx_exchanges_post_id ON exchanges(post_id);
CREATE INDEX idx_exchanges_requester_id ON exchanges(requester_id);
CREATE INDEX idx_exchanges_owner_id ON exchanges(owner_id);
CREATE INDEX idx_messages_exchange_id ON messages(exchange_id);
CREATE INDEX idx_ratings_exchange_id ON ratings(exchange_id);