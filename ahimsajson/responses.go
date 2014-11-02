package ahimsajson

// Single Items

// Holds all the information available about a given Bulletin
type JsonBltn struct {
	Txid         string `json:"txid"`
	Board        string `json:"board,omitempty"`
	Author       string `json:"author"`
	Message      string `json:"msg"`
	Timestamp    int64  `json:"timestamp,omitempty"`
	Block        string `json:"blk,omitempty"`
	BlkTimestamp int64  `json:"blkTimestamp,omitempty"`
	BannedReason string `json:"bannedReason,omitempty"`
}

// Holds meta information about a single unique block
type JsonBlkHead struct {
	Hash      string `json:"hash"`
	PrevHash  string `json:"prevHash"`
	Timestamp int64  `json:"timestamp"`
	Height    uint64 `json:"height"`
	NumBltns  uint64 `json:"numBltns"`
}

// Contains a block head and the contained bulletins in that block
type JsonBlkResp struct {
	Head      *JsonBlkHead `json:"head"`
	Bulletins []*JsonBltn  `json:"bltns"`
}

// Holds meta information about a single author
type AuthorSummary struct {
	Address    string `json:"addr"`
	NumBltns   uint64 `json:"numBltns"`
	FirstBlkTs int64  `json:"firstBlkTs,omitempty"`
}

// Contains info about an author and posts by that author
type AuthorResp struct {
	Author    *AuthorSummary `json:"author"`
	Bulletins []*JsonBltn    `json:"bltns"`
}

// Holds meta information about the server
type Info struct {
	Uptime    int64  `json:"uptime"`
	Version   string `json:"version"`
	LatestBlk int64  `json:"latestBlock"`
}

// Holds summary information about a given board
type BoardSummary struct {
	Name       string `json:"name"`
	NumBltns   uint64 `json:"numBltns"`
	CreatedAt  int64  `json:"createdAt"`  // The block timestamp of when this board was started.
	LastActive int64  `json:"lastActive"` // The bulletin timestamp of the latest post.
	CreatedBy  string `json:"createdBy"`
}

// An entire bulletin board that is sorted by default in ascending order
type WholeBoard struct {
	Summary   *BoardSummary `json:"summary"`
	Bulletins []*JsonBltn   `json:"bltns"`
}

// A bulletin that is banned by the administrator for some "Reason."
type BannedBltn struct {
	Txid   string `json:"txid"`
	Reason string `json:"reason"`
}
