package repository

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func loadVideosFromFile(t *testing.T, relPath string) []*Video {
	t.Helper()
	p := filepath.Join("..", "..", "testdata", "repository", relPath)
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatalf("read file %s: %v", p, err)
	}

	var list []Video
	if err := json.Unmarshal(b, &list); err != nil {
		t.Fatalf("unmarshal %s: %v", p, err)
	}

	out := make([]*Video, 0, len(list))
	for i := range list {
		// copy to avoid taking address of loop variable
		v := list[i]
		out = append(out, &v)
	}
	return out
}

func TestAddSortedVideos_Basic(t *testing.T) {
	repo := NewVideoRepository(100)

	v2 := loadVideosFromFile(t, "res2.json")
	added := repo.AddSortedVideos(v2)
	if added != len(v2) {
		t.Fatalf("expected %d added from res2, got %d", len(v2), added)
	}
	if len(repo.Videos) != len(v2) {
		t.Fatalf("expected repo.Videos length %d after first add, got %d", len(v2), len(repo.Videos))
	}

	v3 := loadVideosFromFile(t, "res3.json")
	added = repo.AddSortedVideos(v3)
	if added != len(v3) {
		t.Fatalf("expected %d added from res3, got %d", len(v3), added)
	}

	// total should be sum of both files (no duplicates in fixtures)
	if len(repo.Videos) != len(v2)+len(v3) {
		t.Fatalf("expected total videos %d, got %d", len(v2)+len(v3), len(repo.Videos))
	}

	// verify videos are ordered newest (latest) first by StartTime
	for i := 0; i+1 < len(repo.Videos); i++ {
		if repo.Videos[i].StartTime.Before(repo.Videos[i+1].StartTime) {
			t.Fatalf("videos not in newest-first order at index %d: %v before %v", i, repo.Videos[i].StartTime, repo.Videos[i+1].StartTime)
		}
		// fmt.Printf("%s: %s\n", repo.Videos[i].StartTime.Format("01-02 15:04:05"), repo.Videos[i].Title[:15])
	}

	// seenIDs should contain all IDs
	if len(repo.seenIDs) != len(repo.Videos) {
		t.Fatalf("seenIDs size %d does not match videos length %d", len(repo.seenIDs), len(repo.Videos))
	}
}

func TestAddSortedVideos_TrimToCapacity(t *testing.T) {
	// set small capacity so trimming occurs
	repo := NewVideoRepository(5)

	v2 := loadVideosFromFile(t, "res2.json")
	repo.AddSortedVideos(v2)

	v3 := loadVideosFromFile(t, "res3.json")
	repo.AddSortedVideos(v3)

	if len(repo.Videos) > repo.Capacity {
		t.Fatalf("repo.Videos length %d exceeds capacity %d", len(repo.Videos), repo.Capacity)
	}

	// seenIDs should match the kept videos count
	if len(repo.seenIDs) != len(repo.Videos) {
		t.Fatalf("after trim, seenIDs size %d does not match videos length %d", len(repo.seenIDs), len(repo.Videos))
	}

	// ensure no duplicate IDs
	seen := make(map[string]struct{})
	for _, v := range repo.Videos {
		if _, ok := seen[v.ID]; ok {
			t.Fatalf("duplicate id %s in repo.Videos", v.ID)
		}
		seen[v.ID] = struct{}{}
	}
}
