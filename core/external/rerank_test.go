package external

import (
	"testing"

	"github.com/navidrome/navidrome/core/agents"
	"github.com/navidrome/navidrome/model"
	. "github.com/onsi/gomega"
)

func TestReRankByPopularity(t *testing.T) {
	g := NewWithT(t)

	songs := []agents.Song{
	    {Name: "low scrobbles", Artist: "artist a", Match: 1.0, Scrobbles: 200000},
	    {Name: "mid scrobbles", Artist: "artist b", Match: 0.8, Scrobbles: 1000000},
	    {Name: "high scrobbles", Artist: "artist c", Match: 0.5, Scrobbles: 5000000},
	}

	matched := model.MediaFiles{
		{ID: "1", Title: "low scrobbles", Artist: "artist a"},
		{ID: "2", Title: "mid scrobbles", Artist: "artist b"},
		{ID: "3", Title: "high scrobbles", Artist: "artist c"},
	}

	t.Run("weight=0 preserves original order", func(t *testing.T) {
		result := reRankByPopularity(songs, matched, 0.0)
		g.Expect(result[0].ID).To(Equal("1"))
		g.Expect(result[1].ID).To(Equal("2"))
		g.Expect(result[2].ID).To(Equal("3"))
	})

	t.Run("weight=1 sorts purely by scrobbles descending", func(t *testing.T) {
		result := reRankByPopularity(songs, matched, 1.0)
		g.Expect(result[0].ID).To(Equal("3")) // highest scrobbles
		g.Expect(result[1].ID).To(Equal("2"))
		g.Expect(result[2].ID).To(Equal("1")) // lowest scrobbles
	})

	t.Run("weight=0.5 balances similarity and popularity", func(t *testing.T) {
	        result := reRankByPopularity(songs, matched, 0.5)
	        g.Expect(result[2].ID).To(Equal("3")) // lowest combined score should be last
	        g.Expect(result).To(HaveLen(3))
	})

	t.Run("all zero scrobbles returns matched unchanged", func(t *testing.T) {
		zeroSongs := []agents.Song{
			{Name: "low scrobbles", Artist: "artist a", Match: 1.0, Scrobbles: 0},
			{Name: "mid scrobbles", Artist: "artist b", Match: 0.8, Scrobbles: 0},
		}
		zeroMatched := model.MediaFiles{
			{ID: "1", Title: "low scrobbles", Artist: "artist a"},
			{ID: "2", Title: "mid scrobbles", Artist: "artist b"},
		}
		result := reRankByPopularity(zeroSongs, zeroMatched, 0.5)
		g.Expect(result[0].ID).To(Equal("1"))
		g.Expect(result[1].ID).To(Equal("2"))
	})

	t.Run("empty matched returns empty", func(t *testing.T) {
		result := reRankByPopularity(songs, model.MediaFiles{}, 0.5)
		g.Expect(result).To(BeEmpty())
	})
	t.Run("scrobbles below floor score 0 on popularity", func(t *testing.T) {
        	floorSongs := []agents.Song{
                	{Name: "below floor", Artist: "artist a", Match: 0.5, Scrobbles: 1000},
                	{Name: "above floor", Artist: "artist b", Match: 0.4, Scrobbles: 500000},
        	}
        	floorMatched := model.MediaFiles{
                	{ID: "1", Title: "below floor", Artist: "artist a"},
                	{ID: "2", Title: "above floor", Artist: "artist b"},
        	}
        	// with weight=1.0, above-floor track should win despite lower match
        	result := reRankByPopularity(floorSongs, floorMatched, 1.0)
        	g.Expect(result[0].ID).To(Equal("2"))
        	g.Expect(result[1].ID).To(Equal("1"))
	})
}
