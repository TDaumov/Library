-- Запустить один раз:
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    email text NOT NULL UNIQUE,
    password text NOT NULL,
    role text NOT NULL DEFAULT 'reader',
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS books (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    title text NOT NULL,
    author text NOT NULL,
    description text,
    published_at timestamptz,
    available boolean DEFAULT true,
    created_at timestamptz DEFAULT now(),
    updated_at timestamptz DEFAULT now()
    );

CREATE TABLE IF NOT EXISTS borrows (
    id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id uuid REFERENCES users(id),
    book_id uuid REFERENCES books(id),
    borrowed_at timestamptz DEFAULT now(),
    returned_at timestamptz
    );