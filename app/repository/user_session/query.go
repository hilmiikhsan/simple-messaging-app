package user_session

const (
	queryDeleteUserSessionByToken = `
		DELETE FROM user_sessions WHERE token = ?
	`

	queryUpdteUserSessionToken = `
		UPDATE user_sessions 
		SET 
			token = ?, 
			token_expired=? 
		WHERE refresh_token = ?
	`
)
