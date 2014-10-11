package ahimsajson

// Single Items

// Holds all the information available about a given Bulletin
type JsonBltn struct {
	Txid         string `json:"txid"`
	Board        string `json:"board,omitempty"`
	Author       string `json:"author"`
	Message      string `json:"msg"`
	Block        string `json:"blk,omitempty"`
	BlkTimestamp uint64 `json:"blkTimestamp,omitempty"`
}

// Holds meta information about a single unique block
type JsonBlkHead struct {
	Hash      string `json:"hash"`
	prevHash  string `json:"prevHash"`
	Timestamp uint64 `json:"timestamp"`
	Height    uint64 `json:"height"`
	NumBltns  uint64 `json:"numBltns"`
}

// Holds meta information about the server
type Info struct {
	Uptime    uint64 `json:"uptime"`
	Version   string `json:"version"`
	LatestBlk uint64 `json:"latestblock"`
}

// Holds summary information about a given board
type BoardSummary struct {
	Name       string `json:"name"`
	NumBltns   string `json:"numBltns"`
	StartedAt  string `json:"startedAt,omitempty"` // The block timestamp of when this board was started.
	LastActive string `json:"lastPost,omitempty"`  // The block timestamp of the latest post.
}
