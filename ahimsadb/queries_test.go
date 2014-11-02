package ahimsadb

import "testing"

func TestJsonBlock(t *testing.T) {

	db, err := SetupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	h := "00000000777213b4fd7c5d5a71b9b52608356c4194203b1b63d1bb0e6141d17d"
	jsonBlkResp, err := GetJsonBlock(db, h)

	if err != nil {
		t.Fatal(err)
	}

	respH := jsonBlkResp.Head.Hash
	if respH != h {
		t.Fatalf("Hashes don't match [%s] and returned: [%s]", h, respH)
	}

}

func TestJsonAuthor(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	author := "miUDcP8obUKPhqkrBrQz57sbSg2Mz1kZXH"

	jsonResp, err := GetJsonAuthor(db, author)

	if err != nil {
		t.Fatal(err)
	}

	blkTs := int64(1414017952)
	if jsonResp.Author.NumBltns != 2 || jsonResp.Author.FirstBlkTs != blkTs {
		t.Fatalf("Wrong values:\n [%s]\n", jsonResp)
	}
}

func TestWholeBoard(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	board := "ahimsa-dev"

	wholeBoard, err := GetWholeBoard(db, board)

	if err != nil {
		t.Fatal(err)
	}

	if wholeBoard.Summary.NumBltns != 4 {
		t.Fatalf("Wrong values:\n [%s]\n", wholeBoard)
	}

	expLA := int64(1414624430)
	if wholeBoard.Summary.LastActive != expLA {
		t.Fatalf(
			"Wrong last active time in:\n[%s]\nWanted an LA of: %d\n\tGot: %d",
			wholeBoard.Summary,
			expLA,
			wholeBoard.Summary.LastActive,
		)
	}

}

func TestAllBoards(t *testing.T) {
	db, err := SetupTestDB()
	if err != nil {
		t.Fatal(err)
	}

	allBoards, err := GetAllBoards(db)

	if err != nil {
		t.Fatal(err)
	}

	if len(allBoards) != 2 {
		t.Fatalf("Wrong number of boards returned:\n [%s]\n", allBoards)
	}

}
