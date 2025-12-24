package repository

import (
	"net/url"
	"strings"
	"time"
)

type Video struct {
	ID               string    `json:"contentId"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	StartTime        time.Time `json:"startTime"`
	ThumbnailURL     string    `json:"thumbnailUrl"`
	ThumbnailType    string    `json:"-"`
	ThumbnailLength  int64     `json:"-"`
	TagsConnectedStr string    `json:"tags"`
}

// VideoRepository 動画情報をメモリに保持する
// 基本的にAddSortedVideos()で追加を行うことを想定、その際Capacity超過分のTrimも行われる。
type VideoRepository struct {
	Videos   []*Video
	seenIDs  map[string]struct{}
	Capacity int
}

// URL 動画視聴URLを返す
func (v Video) URL() string {
	return "https://nico.ms/" + v.ID
}

// TagSearchURL タグ検索用のURLを返す
func TagSearchURL(tag string) string {
	return "https://www.nicovideo.jp/tag/" + url.PathEscape(tag)
}

// Tags タグをスライス形式で返す
func (v Video) Tags() []string {
	return strings.Split(v.TagsConnectedStr, " ")
}

func NewVideoRepository(capacity int) *VideoRepository {
	return &VideoRepository{
		Videos:   make([]*Video, 0, capacity),
		seenIDs:  make(map[string]struct{}),
		Capacity: capacity,
	}
}

// AddSortedVideos ソート済み動画スライスをマージし、重複を排除して格納する。重複を省いた後の追加数を返す。
// ソート済みであることを前提としている。APIから-startTimeで取ったデータを入れるなら問題なし
func (r *VideoRepository) AddSortedVideos(newVideos []*Video) int {
	if len(newVideos) == 0 {
		return 0
	}

	toMerge := make([]*Video, 0, len(newVideos))
	for _, v := range newVideos {
		if _, exists := r.seenIDs[v.ID]; !exists {
			toMerge = append(toMerge, v)
			r.seenIDs[v.ID] = struct{}{}
		}
	}
	if len(toMerge) == 0 {
		return 0
	}

	r.Videos = mergeSortedVideos(r.Videos, toMerge)

	r.TrimToCapacity()
	return len(toMerge)
}

func mergeSortedVideos(existing, newVideos []*Video) []*Video {
	result := make([]*Video, 0, len(existing)+len(newVideos))
	i, j := 0, 0
	for i < len(existing) && j < len(newVideos) {
		// merge so that newest (latest StartTime) comes first
		if existing[i].StartTime.After(newVideos[j].StartTime) {
			result = append(result, existing[i])
			i++
		} else {
			result = append(result, newVideos[j])
			j++
		}
	}

	result = append(result, existing[i:]...)
	result = append(result, newVideos[j:]...)
	return result
}

// TrimToCapacity 規定数を超えた動画を削除しその数を返す
func (r *VideoRepository) TrimToCapacity() int {
	if len(r.Videos) <= r.Capacity {
		return 0
	}

	removed := r.Videos[r.Capacity:]
	r.Videos = r.Videos[:r.Capacity]
	for _, v := range removed {
		delete(r.seenIDs, v.ID)
	}

	return len(removed)
}
