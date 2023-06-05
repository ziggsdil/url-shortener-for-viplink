package db

const (
	initRequest = `
		CREATE TABLE IF NOT EXISTS links (
			short_suffix TEXT PRIMARY KEY,
			link TEXT,
			secret_key TEXT UNIQUE,
			clicks INTEGER NOT NULL DEFAULT 0
		);
	`

	saveRequest = `
		INSERT INTO links(short_suffix, link, secret_key) 
			VALUES ($1, $2, $3);
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
)
