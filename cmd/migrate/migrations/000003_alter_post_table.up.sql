ALTER TABLE Posts
    ADD CONSTRAINT fk_user FOREIGN KEY(
        user_id
    )
        REFERENCES Users(
            id
        );