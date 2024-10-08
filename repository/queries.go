package repository

var (
	queryInsertUsers         = "INSERT INTO \"user\"(full_name, phone_number, password, balance, created_time) VALUES"
	valuesInsertUsersF       = "($%d, $%d, $%d, $%d, $%d),"
	returnLastInsertedUserID = "RETURNING id"
)

var (
	querySelectUsers     = "SELECT id, full_name, phone_number, balance, password, created_time, updated_time FROM \"user\" WHERE true"
	whereUserPhoneNumber = " AND phone_number = $%d"
	whereUserID          = " AND id = $%d"
)

var (
	queryUpdateUserF      = "UPDATE \"user\" SET %s WHERE TRUE"
	incrementUserBalanceF = "balance = balance + $%d"
	decrementUserBalanceF = "balance = balance - $%d"
	setUserUpdatedTimeF   = "updated_time = $%d"
)

var (
	queryInsertTransaction = "INSERT INTO transaction(id, user_id, amount, type, recipient_id, status, description, created_time) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id"
)

var (
	queryLockUser = "SELECT balance from \"user\" WHERE id = $1 FOR UPDATE"
)
