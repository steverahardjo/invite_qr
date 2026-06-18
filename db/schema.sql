
--name: InitParticipants: one
CREATE TABLE participants (
	id SERIAL PRIMARY KEY,
	name VARCHAR(255) NOT NULL,
	email VARCHAR(255) NOT NULL,
	wa_number VARCHAR(20) NOT NULL,
	accessed BOOLEAN NOT NULL DEFAULT FALSE,

)

-- name: InsertParticipant: one
INSERT INTO participants (name, email, wa_number)
VALUES ('', '', '');

-- ListParticipants: many
SELECT * FROM participants
where owner = $1
ORDER BY id
LIMIT $5 OFFSET $3;


-- name: UpdateParticipant: one
UPDATE participants
SET accessed = TRUE
WHERE id = (SELECT id FROM participants WHERE email = '' or wa_number = '');
