package ahimsadb

import (
	"database/sql"
	"errors"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/NSkelsey/ahimsarest/ahimsajson"
)

var (
	ErrBltnCensored error = errors.New("Bulletin is withheld for some reason")

	// Used by GetJsonBltn
	selectTxid    *sql.Stmt
	selectTxidSql string = `
		SELECT author, board, message, bulletins.timestamp, block, blocks.timestamp, blacklist.reason
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE bulletins.txid = $1;
	`
	// Used by GetJsonBlock

	// Used by GetJsonAuthor
	selectAuthor    *sql.Stmt
	selectAuthorSql string = `
		SELECT author, count(*), blocks.timestamp
		FROM bulletins LEFT JOIN blocks on bulletins.block = blocks.hash
		WHERE author = $1
		ORDER BY blocks.timestamp ASC
		LIMIT 1;
	`

	// Used by GetJsonBlacklist
	selectBlacklist    *sql.Stmt
	selectBlacklistsql string = `
		SELECT txid, reason from blacklist;
	`
)

// Prepares all of the selects for maximal speediness
func prepareQueries(db *sql.DB) error {

	var err error
	if selectTxid, err = db.Prepare(selectTxidSql); err != nil {
		return err
	}

	if selectAuthor, err = db.Prepare(selectAuthorSql); err != nil {
		return err
	}

	if selectBlacklist, err = db.Prepare(selectBlacklistsql); err != nil {
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

// Returns information about a single author
func GetJsonAuthor(db *sql.DB, address string) (*ahimsajson.Author, error) {

	var numBltns uint64
	var firstBlockTs sql.NullInt64

	row := selectAuthor.QueryRow(address)
	err := row.Scan(&numBltns, &firstBlockTs)
	if err != nil {
		return nil, err
	}

	author := &ahimsajson.Author{
		Address:  address,
		NumBltns: numBltns,
	}

	if firstBlockTs.Valid {
		author.FirstBlkTs = firstBlockTs.Int64
	}

	return author, nil
}

// Returns the single bulletin in json format, that is identified by txid.
// If the bltn does not exist GetJsonBltn returns sql.ErrNoRows.
func GetJsonBltn(db *sql.DB, txid string) (*ahimsajson.JsonBltn, error) {

	var author, msg string
	var board, blockH, bannedReason sql.NullString
	var blkTs, bltnTs sql.NullInt64

	row := selectTxid.QueryRow(txid)
	err := row.Scan(&author, &board, &msg, &bltnTs, &blockH, &blkTs, &bannedReason)
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

	// If the bulletin was banned, still return the bltn, but provide
	// the error for applications to handle.
	if bannedReason.Valid {
		return bltn, ErrBltnCensored
	}

	return bltn, nil
}

func GetJsonBlock(db *sql.DB, hash string) (*ahimsajson.JsonBlkHead, error) {

	blkHead := &ahimsajson.JsonBlkHead{}

	return blkHead, nil
}

func GetJsonBlacklist(db *sql.DB) ([]*ahimsajson.BannedBltn, error) {

	blacklist := []*ahimsajson.BannedBltn{}
	empt := []*ahimsajson.BannedBltn{}
	rows, err := selectBlacklist.Query()
	defer rows.Close()
	if err != nil {
		return empt, err
	}
	for rows.Next() {
		var txid, reason string
		if err := rows.Scan(&txid, &reason); err != nil {
			return empt, err
		}
		bannedBltn := &ahimsajson.BannedBltn{
			Txid:   txid,
			Reason: reason,
		}
		blacklist = append(blacklist, bannedBltn)

	}

	return blacklist, nil
}
