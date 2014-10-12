package ahimsadb

import (
	"database/sql"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/NSkelsey/ahimsarest/ahimsajson"
)

var (
	// Used by GetJsonBltn
	selectTxid    *sql.Stmt
	selectTxidSql string = `
		SELECT author, board, message, bulletins.timestamp, block, blocks.timestamp 
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		WHERE txid = $1;
	`
	// Used by GetBlock
)

// Prepares all of the selects for maximal speediness
func prepareQueries(db *sql.DB) error {

	var err error
	selectTxid, err = db.Prepare(selectTxidSql)
	if err != nil {
		return err
	}

	return nil
}

// Loads a sqlite db, checks if its reachabale and prepares all the queries.
func LoadDb(dbpath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbpath)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = prepareQueries(db)
	if err != nil {
		return nil, err
	}

	return db, nil
}

// Returns the single bulletin in json format, that is identified by txid.
// If the bltn does not exist GetJsonBltn returns sql.ErrNoRows.
func GetJsonBltn(db *sql.DB, txid string) (*ahimsajson.JsonBltn, error) {

	var author, msg string
	var board, blockH sql.NullString
	var blkTs, bltnTs sql.NullInt64

	row := selectTxid.QueryRow(txid)
	err := row.Scan(&author, &board, &msg, &bltnTs, &blockH, &blkTs)
	if err != nil {
		return nil, err
	}

	bltn := &ahimsajson.JsonBltn{
		Txid:    txid,
		Author:  author,
		Message: msg,
	}

	if bltnTs.Valid {
		bltn.Timestamp = bltnTs.Int64
	}

	if board.Valid {
		bltn.Board = board.String
	}

	// If the response contained a block, fill the optional params
	if blockH.Valid {
		bltn.Block = blockH.String
		bltn.BlkTimestamp = blkTs.Int64
	}
	return bltn, nil
}
