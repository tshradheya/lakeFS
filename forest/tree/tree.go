package tree

import (
	cache "github.com/treeverse/lakefs/forest/cache_map"
	gr "github.com/treeverse/lakefs/graveler"
	"github.com/treeverse/lakefs/graveler/committed/sstable"
)

var minimalKey = gr.Key(nil)

type treePartType struct {
	PartName sstable.ID `json:"part_name"`
	MaxKey   gr.Key     `json:"max_path"`
}
type TreeType struct {
	treeSlice []treePartType
}

type treesRepo struct {
	treesMap   cache.CacheMap
	partManger sstable.Manager
}

const (
	CacheMapSize     = 1000
	CacheTrimSize    = 100
	InitialWeight    = 64
	AdditionalWeight = 16
	TrimFactor       = 1
)

func InitTreesRepository(manager sstable.Manager) TreeRepo {
	treesRepository := &treesRepo{
		treesMap:   cache.NewCacheMap(CacheMapSize, CacheTrimSize, InitialWeight, AdditionalWeight, TrimFactor),
		partManger: manager,
	}
	return treesRepository
}

func (t *treesRepo) GetPartManger() sstable.Manager {
	return t.partManger
}

type TreeRepo interface {
	NewTreeWriter(splitFactor int, closeAsync sstable.BatchWriterCloser) TreeWriter
	NewScannerFromID(treeID gr.TreeID, start gr.Key) (gr.ValueIterator, error)
	NewScannerFromTreeParts(treeSlice TreeType, start gr.Key) (gr.ValueIterator, error)
	GetPartManager() sstable.Manager
}

type TreeWriter interface {
	HasOpenWriter() bool
	WriteEntry(record gr.ValueRecord) error
	ForceCloseCurrentPart()
	IsSplitKey(key gr.Key, rowNum int) bool
	FlushIterToTree(iter gr.ValueIterator) error
	SaveTree(reuseTree TreeType) (gr.TreeID, error)
}
