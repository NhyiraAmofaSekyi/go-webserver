CREATE TABLE users (
   id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
   created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
   name TEXT NOT NULL
);
