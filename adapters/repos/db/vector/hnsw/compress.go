//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package hnsw

import (
	"context"
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
	"github.com/weaviate/weaviate/adapters/repos/db/helpers"
	"github.com/weaviate/weaviate/adapters/repos/db/lsmkv"
	ssdhelpers "github.com/weaviate/weaviate/adapters/repos/db/vector/ssdhelpers"
	entlsmkv "github.com/weaviate/weaviate/entities/lsmkv"
	ent "github.com/weaviate/weaviate/entities/vectorindex/hnsw"
)

func (h *hnsw) initCompressedStore() error {
	store, err := lsmkv.New(h.compressedStoreLSMPath(), "", h.logger, nil,
		h.shardCompactionCallbacks, h.shardFlushCallbacks)
	if err != nil {
		return errors.Wrap(err, "Init lsmkv (compressed vectors store)")
	}
	err = store.CreateOrLoadBucket(context.Background(), helpers.CompressedObjectsBucketLSM)
	if err != nil {
		return errors.Wrapf(err, "Create or load bucket (compressed vectors store)")
	}
	h.compressedStore = store
	return nil
}

func (h *hnsw) compressedStoreLSMPath() string {
	return fmt.Sprintf("%s/%s/%s", h.rootPath, h.className, h.shardName)
}

func (h *hnsw) Compress(cfg ent.PQConfig) error {
	if h.isEmpty() {
		return errors.New("Compress command cannot be executed before inserting some data. Please, insert your data first.")
	}
	err := h.initCompressedStore()
	if err != nil {
		return errors.Wrap(err, "Initializing compressed vector store")
	}

	// get representative vector for inferring dimensions, use entrypoint
	h.RLock()
	epID := h.entryPointID
	h.RUnlock()

	vec, err := h.vectorForID(context.Background(), epID)
	if err != nil {
		return errors.Wrap(err, "Inferring data dimensions")
	}
	dims := len(vec)

	// segments == 0 (default value) means use as many segments as dimensions
	if cfg.Segments <= 0 {
		cfg.Segments = dims
	}

	h.pq, err = ssdhelpers.NewProductQuantizer(cfg, h.distancerProvider, dims)
	if err != nil {
		return errors.Wrap(err, "Compressing vectors.")
	}

	data := h.cache.all()
	cleanData := make([][]float32, 0, len(data))
	for _, point := range data {
		if point == nil {
			continue
		}
		cleanData = append(cleanData, point)
	}
	h.compressedVectorsCache.grow(uint64(len(data)))
	h.pq.Fit(cleanData)

	h.compressActionLock.Lock()
	defer h.compressActionLock.Unlock()
	ssdhelpers.Concurrently(uint64(len(data)),
		func(index uint64) {
			if data[index] == nil {
				return
			}

			encoded := h.pq.Encode(data[index])
			h.storeCompressedVector(index, encoded)
			h.compressedVectorsCache.preload(index, encoded)
		})
	if err := h.commitLog.AddPQ(h.pq.ExposeFields()); err != nil {
		return errors.Wrap(err, "Adding PQ to the commit logger")
	}

	h.compressed.Store(true)
	h.cache.drop()
	return nil
}

//nolint:unused
func (h *hnsw) encodedVector(id uint64) ([]byte, error) {
	return h.compressedVectorsCache.get(context.Background(), id)
}

func (h *hnsw) storeCompressedVector(id uint64, vector []byte) {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, id)
	h.compressedStore.Bucket(helpers.CompressedObjectsBucketLSM).Put(key, vector)
}

func (h *hnsw) getCompressedVectorForID(ctx context.Context, id uint64) ([]byte, error) {
	key := make([]byte, 8)
	binary.LittleEndian.PutUint64(key, id)

	if vec, err := h.compressedStore.Bucket(helpers.CompressedObjectsBucketLSM).Get(key); err == nil {
		return vec, nil
	} else if err != entlsmkv.NotFound {
		return nil, fmt.Errorf("getting vector '%d' from compressed store: %w", id, err)
	}

	// not found, fallback to uncompressed source
	h.logger.
		WithField("action", "compress").
		WithField("vector", id).
		Warnf("Vector not found in compressed store")

	vec, err := h.VectorForIDThunk(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("getting vector '%d' from store: %w", id, err)
	}

	return h.pq.Encode(h.normalizeVec(vec)), nil
}
