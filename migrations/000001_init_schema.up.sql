
CREATE TABLE songs (
                       id VARCHAR(255) PRIMARY KEY not null,
                       group_name VARCHAR(255) NOT NULL,
                       song_name VARCHAR(255) NOT NULL,
                       release_date DATE,
                       lyrics TEXT,
                       link VARCHAR(255),
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

