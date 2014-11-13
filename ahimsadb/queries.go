package ahimsadb

import (
	"database/sql"
	"errors"
	"net/url"

	_ "code.google.com/p/go-sqlite/go1/sqlite3"
	"github.com/NSkelsey/ahimsarest/ahimsajson"
)

var (
	ErrBltnCensored error = errors.New("Bulletin is withheld for some reason")

	// Used by GetJsonBltn
	selectTxid    *sql.Stmt
	selectTxidSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, block, blocks.timestamp, blacklist.reason
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE bulletins.txid = $1
	`
	// Used by GetJsonBlock
	selectBlockHead    *sql.Stmt
	selectBlockHeadSql string = `
		SELECT hash, prevhash, height, timestamp FROM blocks
		WHERE blocks.hash = $1
	`
	selectBlockBltns    *sql.Stmt
	selectBlockBltnsSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, block, blocks.timestamp, blacklist.reason
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE blocks.hash = $1
	`

	// Used by GetJsonAuthor
	selectAuthor    *sql.Stmt
	selectAuthorSql string = `
		SELECT author, count(*), blocks.timestamp
		FROM bulletins LEFT JOIN blocks on bulletins.block = blocks.hash
		WHERE author = $1
		ORDER BY blocks.timestamp ASC
	`

	selectAuthorBltns    *sql.Stmt
	selectAuthorBltnsSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, block, blocks.timestamp, blacklist.reason
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE author = $1
	`

	// Used by GetJsonBlacklist
	selectBlacklist    *sql.Stmt
	selectBlacklistSql string = `
		SELECT txid, reason from blacklist
	`

	// Used by GetWholeBoard
	selectBoardSum    *sql.Stmt
	selectBoardSumSql string = `
		SELECT board, count(*), last_bltn.bltn_ts, first_bltn.blk_ts, author 
		FROM bulletins, 
			(SELECT max(bulletins.timestamp) AS bltn_ts FROM bulletins WHERE board = $1) AS last_bltn,
			(SELECT min(blocks.timestamp)  blk_ts FROM bulletins JOIN blocks on bulletins.block = blocks.hash
				WHERE board = $1	
			) AS first_bltn
		LEFT JOIN blocks ON bulletins.block = blocks.hash
		WHERE board = $1
		ORDER BY bulletins.timestamp ASC
		LIMIT 1
	`

	selectBoardBltns    *sql.Stmt
	selectBoardBltnsSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, block, blocks.timestamp, blacklist.reason
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE board = $1
		ORDER BY blocks.timestamp, bulletins.timestamp
	`

	// Used by GetAllBoards
	selectAllBoards    *sql.Stmt
	selectAllBoardsSql string = `
		SELECT board, count(*), max(bulletins.timestamp), blocks.timestamp, author
		FROM bulletins LEFT JOIN blocks ON bulletins.block = blocks.hash
		GROUP BY board
		ORDER BY blocks.timestamp ASC
	`

	// Used by GetRecentBltns
	selectRecentConf    *sql.Stmt
	selectRecentConfSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, block, blocks.timestamp, blacklist.reason
		FROM bulletins, (
			SELECT max(blocks.height) AS height FROM blocks	
		) AS tip		
		JOIN blocks ON bulletins.block = blocks.hash
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE blocks.height > (tip.height - $1)
		ORDER BY blocks.timestamp DESC
	`

	// Used by GetUnconfirmed
	selectUnconfirmed    *sql.Stmt
	selectUnconfirmedSql string = `
		SELECT bulletins.txid, author, board, message, bulletins.timestamp, NULL, NULL, blacklist.reason
		FROM bulletins
		LEFT JOIN blacklist ON bulletins.txid = blacklist.txid
		WHERE block IS NULL
		ORDER BY bulletins.timestamp
	`

	sqlStmts = map[string]**sql.Stmt{
		selectTxidSql:        &selectTxid,
		selectBlockHeadSql:   &selectBlockHead,
		selectBlockBltnsSql:  &selectBlockBltns,
		selectAuthorSql:      &selectAuthor,
		selectAuthorBltnsSql: &selectAuthorBltns,
		selectBlacklistSql:   &selectBlacklist,
		selectBoardSumSql:    &selectBoardSum,
		selectBoardBltnsSql:  &selectBoardBltns,
		selectAllBoardsSql:   &selectAllBoards,
		selectRecentConfSql:  &selectRecentConf,
		selectUnconfirmedSql: &selectUnconfirmed,
	}
)

// Prepares all of the selects for maximal speediness note that all of the queries
// must be within the sqlStmts map for initialization.
func prepareQueries(db *sql.DB) error {

	for sqlString, sqlStmt := range sqlStmts {
		upStmt, err := db.Prepare(sqlString)
		if err != nil {
			return err
		}
		(*sqlStmt) = upStmt
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
func GetJsonAuthor(db *sql.DB, address string) (*ahimsajson.AuthorResp, error) {

	var numBltns uint64
	var addrstr sql.NullString
	var firstBlockTs sql.NullInt64

	row := selectAuthor.QueryRow(address)
	err := row.Scan(&addrstr, &numBltns, &firstBlockTs)
	if err != nil {
		return nil, err
	}

	// Check to see if query returned a real row indicating that this author
	// acutally exists.
	if !addrstr.Valid {
		return nil, sql.ErrNoRows
	}

	authorSum := &ahimsajson.AuthorSummary{
		Address:  address,
		NumBltns: numBltns,
	}

	if firstBlockTs.Valid {
		authorSum.FirstBlkTs = firstBlockTs.Int64
	}

	rows, err := selectAuthorBltns.Query(address)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	bltns, err := getRelevantBltns(rows)
	if err != nil {
		return nil, err
	}

	authorResp := &ahimsajson.AuthorResp{
		Author:    authorSum,
		Bulletins: bltns,
	}
	return authorResp, nil
}

// Returns the single bulletin in json format, that is identified by txid.
// If the bltn does not exist GetJsonBltn returns sql.ErrNoRows.
func GetJsonBltn(db *sql.DB, txid string) (*ahimsajson.JsonBltn, error) {
	row := selectTxid.QueryRow(txid)
	// If the bulletin is banned withold the bulletin
	withhold := true
	return scanJsonBltn(row, withhold)
}

// Returns the block head
func GetJsonBlock(db *sql.DB, h string) (*ahimsajson.JsonBlkResp, error) {

	var hash, prevhash string
	var timestamp int64
	var height uint64

	row := selectBlockHead.QueryRow(h)
	err := row.Scan(&hash, &prevhash, &height, &timestamp)
	if err != nil {
		return nil, err
	}

	rows, err := selectBlockBltns.Query(h)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	bltns, err := getRelevantBltns(rows)
	if err != nil {
		return nil, err
	}

	blkHead := &ahimsajson.JsonBlkHead{
		Hash:      hash,
		PrevHash:  prevhash,
		Height:    height,
		Timestamp: timestamp,
		NumBltns:  uint64(len(bltns)),
	}

	blkResp := &ahimsajson.JsonBlkResp{
		Head:      blkHead,
		Bulletins: bltns,
	}
	return blkResp, nil
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

// Returns a board summary and the bulletins posted to that board. This works on
// the null board as well!
func GetWholeBoard(db *sql.DB, boardstr string) (*ahimsajson.WholeBoard, error) {

	// Unescape boardstr and consider the string utf-8. After this unescape we
	// must use unescapedboard because that *IS* the value stored in the DB.
	unescapedboard, err := url.QueryUnescape(boardstr)
	if err != nil {
		return nil, err
	}

	row := selectBoardSum.QueryRow(unescapedboard)

	boardSum, err := scanBoardSummary(row)
	if err != nil {
		return nil, err
	}

	rows, err := selectBoardBltns.Query(unescapedboard)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	bltns, err := getRelevantBltns(rows)
	if err != nil {
		return nil, err
	}

	wholeboard := &ahimsajson.WholeBoard{
		Summary:   boardSum,
		Bulletins: bltns,
	}

	return wholeboard, nil
}

// Returns a board summary for every board in the database.
func GetAllBoards(db *sql.DB) ([]*ahimsajson.BoardSummary, error) {
	boards := []*ahimsajson.BoardSummary{}

	rows, err := selectAllBoards.Query()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		boardSum, err := scanBoardSummary(rows)
		if err != nil {
			return nil, err
		}
		boards = append(boards, boardSum)
	}

	return boards, nil
}

// Returns the last num of confirmed bulletins in the order they were mined
func GetRecentConf(db *sql.DB, num int) ([]*ahimsajson.JsonBltn, error) {

	empt := make([]*ahimsajson.JsonBltn, 0, num)

	rows, err := selectRecentConf.Query(num)
	if err != nil {
		return empt, err
	}

	bltns, err := getRelevantBltns(rows)
	if err != nil {
		return empt, err
	}

	return bltns, nil
}

func GetUnconfirmed(db *sql.DB) ([]*ahimsajson.JsonBltn, error) {
	empt := []*ahimsajson.JsonBltn{}

	rows, err := selectUnconfirmed.Query()
	if err != nil {
		return empt, err
	}

	bltns, err := getRelevantBltns(rows)
	if err != nil {
		return empt, err
	}

	return bltns, nil
}
