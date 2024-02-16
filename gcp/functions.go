package gcp

import (
	"context"
	db "duplicates-finder/db/generated"
	"fmt"
	"math"
	"sort"
	"sync"

	"log"

	"cloud.google.com/go/storage"
	"go.uber.org/zap"
)

const (
	BUCKET = "editing_userdata"
	BATCH  = 90
)

func (c *Client) Start() {
	count, err := c.Store.GetProfilesCount(context.TODO())
	if err != nil {
		log.Println("Error getting profiles count", err)
		return
	}

	batches := int64(math.Ceil(float64(count) / float64(BATCH)))

	log.Println("Total profiles", count)
	log.Println("Batch size", BATCH)
	log.Println("Total batches", batches)

	for i := int64(0); i < batches; i++ {
		log.Println("Batch", i+1)

		params := db.GetProfilesWithOffsetParams{
			Limit:  int32(BATCH),
			Offset: int32(i * BATCH),
		}

		data, err := c.Store.GetProfilesWithOffset(context.TODO(), params)
		if err != nil {
			log.Println("Error getting profiles with limit and offser", params.Limit, params.Offset, err)
			continue
		}

		var batch_wg sync.WaitGroup

		for _, it := range data {
			batch_wg.Add(1)
			go c.startCollectingDuplicates(it, &batch_wg)
		}

		batch_wg.Wait()
	}
}

func (c *Client) startCollectingDuplicates(userid string, wg *sync.WaitGroup) {
	defer wg.Done()

	paths := c.getUserPaths(userid)

	var path_wg sync.WaitGroup

	for _, path := range paths {
		path_wg.Add(1)
		go c.listZipObjects(context.Background(), path, &path_wg)
	}

	path_wg.Wait()
}

func (c *Client) getUserPaths(user_id string) []string {
	paths := []string{}

	data, err := c.Store.GetProfilesFoldersCatalogs(context.TODO(), user_id)
	if err != nil {
		return nil
	}

	for _, it := range data {
		path := it.UserID + "/" + it.ProfileKey + "/" + it.FolderKey + "/" + it.CatalogKey
		paths = append(paths, path)
	}

	return paths
}

func (c *Client) listZipObjects(ctx context.Context, path string, wg *sync.WaitGroup) {
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
			"Duplicate zip found",
			zap.String("second_last_path", secondLastZip.path),
			zap.String("last_path", lastZip.path),
			zap.String("second_last_zip_checksum", fmt.Sprint(secondLastZip.checksum)),
			zap.String("last_zip_checksum", fmt.Sprint(lastZip.checksum)),
			zap.String("second_last_zip_size", fmt.Sprint(secondLastZip.size)),
			zap.String("last_zip_size", fmt.Sprint(lastZip.size)),
			zap.String("second_last_zip_storage_class", secondLastZip.storageClass),
			zap.String("last_zip_storage_class", lastZip.storageClass),
			zap.String("second_last_zip_last_modified", fmt.Sprint(secondLastZip.lastModified)),
			zap.String("last_zip_last_modified", fmt.Sprint(lastZip.lastModified)),
			zap.String("second_last_dump", fmt.Sprint(secondLastZip)),
			zap.String("last_dump", fmt.Sprint(lastZip)),
		)
	}
}
