package db

const (
	initRequest = `
		CREATE TABLE IF NOT EXISTS links (
			short_suffix TEXT PRIMARY KEY,
			link TEXT,
			secret_key TEXT UNIQUE,
			clicks INTEGER NOT NULL DEFAULT 0,
			expiration_date TIMESTAMP WITH TIME ZONE,
			is_deleted BOOLEAN DEFAULT FALSE
		);
	`

	dropRequest = `
		DROP TABLE IF EXISTS links;
	`

	cleanRequest = `
		DELETE FROM links;
	`

	saveRequest = `
		INSERT INTO links(short_suffix, link, secret_key, expiration_date, is_deleted) 
			VALUES ($1, $2, $3, $4, $5);
	`

	selectBySuffixRequest = `
		SELECT * FROM links 
			WHERE short_suffix=$1;
	`

	selectByLinkRequest = `
		SELECT * FROM links 
			WHERE link=$1;
	`

	selectBySecretKeyRequest = `
		SELECT * FROM links 
			WHERE secret_key=$1;
	`

	deleteBySecretKeyRequest = `
		DELETE FROM links 
			WHERE secret_key=$1;
	`

	incrementClicksBySuffixRequest = `
		UPDATE links
			SET clicks = clicks+1
			WHERE short_suffix=$1;
	`

	updateDeletedBySuffixRequest = `
		UPDATE links
			SET is_deleted=true
			WHERE short_suffix=$1;
	`
	deleteRowsByIsDeleted = `
		DELETE FROM links
			WHERE is_deleted=true;
	`

	updateExpirationDate = `
		UPDATE links
			SET is_deleted=true
			WHERE expiration_date < $1`
)
