package ahimsajson

// Single Items

// Holds all the information available about a given Bulletin
type Bulletin struct {
	Board        string `json:"board,omitempty"`
	Message      string `json:"msg"`
	Txid         string `json:"txid"`
	Block        string `json:"blk,omitempty"`
	Author       string `json:"author"`
	BlkTimestamp uint64 `json:"blkTimestamp,omitempty"`
}

// Holds meta information about a single unique block
type BlockHead struct {
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
type BoardSum struct {
	Name     string `json:"name"`
	NumBltns string `json:"numBltns"`
	// The block timestamp of when this board was started.
	StartedAt string `json:"startedAt,omitempty"`
	// The block timestamp of the latest post.
	LastActive string `json:"lastPost,omitempty"`
}
