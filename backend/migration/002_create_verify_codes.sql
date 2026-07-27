CREATE TABLE verify_codes(
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    email VARCHAR(255) NOT NULL,
    code VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL, -- register / login,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at DATETIME NOT NULL,
    apply_times INT NOT NULL DEFAULT 0,
    used BOOLEAN NOT NULL DEFAULT FALSE
        
    )
