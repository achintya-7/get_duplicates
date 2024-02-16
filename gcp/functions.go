package gcp

import (
	"context"
	db "duplicates-finder/db/generated"
	"sort"
	"sync"

	"log"

	"cloud.google.com/go/storage"
	"go.uber.org/zap"
)

const (
	BUCKET = "editing_userdata"
)

func (c *Client) GetUserProfiles(user_id string) ([]string, error) {
	paths := []string{}

	data, err := c.Store.GetProfilesFoldersCatalogs(context.TODO(), user_id)
	if err != nil {
		return nil, err
	}

	for _, it := range data {
		path := it.UserID + "/" + it.ProfileKey + "/" + it.FolderKey + "/" + it.CatalogKey
		paths = append(paths, path)
	}

	return paths, nil
}

func (c *Client) ListZipObjects(ctx context.Context, path string, wg *sync.WaitGroup) {
	defer wg.Done()

	objects := []ObjectAttrs{}

	it := c.Client.Bucket(BUCKET).Objects(ctx, &storage.Query{Prefix: path})
	for {
		attrs, err := it.Next()
		if err != nil {
			if err.Error() == "no more items in iterator" {
				break
			} else {
				return
			}
		}

		obj := ObjectAttrs{
			path:         attrs.Name,
			size:         attrs.Size,
			storageClass: attrs.StorageClass,
			lastModified: attrs.Updated,
			checksum:     attrs.CRC32C,
		}

		objects = append(objects, obj)
	}

	checkLastTwoZips(objects)

}

func checkLastTwoZips(objects []ObjectAttrs) {
	if len(objects) <= 2 {
		return
	}

	sort.Slice(objects, func(i, j int) bool {
		return objects[i].path < objects[j].path
	})

	lastZip := objects[len(objects)-1]
	secondLastZip := objects[len(objects)-2]

	if lastZip.size == secondLastZip.size && lastZip.checksum == secondLastZip.checksum {
		zap.L().Warn(
			"Duplicate zip found\n",
			zap.String("secondLastZip", secondLastZip.path),
			zap.String("lastZip", lastZip.path),
			zap.String("----", "----"),
		)
	}
}

func (c *Client) Start() {
	count, err := c.Store.GetProfilesCount(context.TODO())
	if err != nil {
		log.Println("Error getting profiles count", err)
		return
	}

	batchSize := int64(10)
	batches := count / int64(batchSize)

	for i := int64(0); i < batches; i++ {
		params := db.GetProfilesWithOffsetParams{
			Limit:  int32(batchSize),
			Offset: int32(i * batchSize),
		}

		data, err := c.Store.GetProfilesWithOffset(context.TODO(), params)
		if err != nil {
			log.Println("Error getting profiles with limit and offser", params.Limit, params.Offset, err)
			continue
		}

		var batch_wg sync.WaitGroup

		for _, it := range data {
			batch_wg.Add(1)
			
		}
	}

}
