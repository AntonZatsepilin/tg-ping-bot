CREATE TABLE users (
    id SERIAL PRIMARY KEY,          
    email VARCHAR(255) NOT NULL UNIQUE,  
    username VARCHAR(50) UNIQUE,     
    password_hash VARCHAR(255) NOT NULL,      
);

CREATE TABLE refresh_tokens (
    id SERIAL PRIMARY KEY,          
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,  
    expires_at TIMESTAMP NOT NULL,        
);
