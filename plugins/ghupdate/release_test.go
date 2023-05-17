package ghupdate

import "testing"

func TestReleaseFindAssetBySuffix(t *testing.T) {
	r := release{
		Assets: []*releaseAsset{
			{Name: "test1.zip", Id: 1},
			{Name: "test2.zip", Id: 2},
			{Name: "test22.zip", Id: 22},
			{Name: "test3.zip", Id: 3},
		},
	}

	asset, err := r.findAssetBySuffix("2.zip")
	if err != nil {
		t.Fatalf("Expected nil, got err: %v", err)
	}

	if asset.Id != 2 {
		t.Fatalf("Expected asset with id %d, got %v", 2, asset)
	}
}
