-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE fcm_tokens (
                            id int(11) NOT NULL AUTO_INCREMENT,
                            token text NOT NULL,
                            user_id int(11) NOT NULL,
                            os varchar(20) NOT NULL,
                            created_at date DEFAULT current_timestamp(),
                            updated_at date DEFAULT current_timestamp(),
                            deleted_at date DEFAULT NULL,
                            PRIMARY KEY (id)
);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
drop table fcm_tokens;