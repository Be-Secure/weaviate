//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2024 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package lsmkv

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

// doReplace parsers all entries into a cache for deduplication first and only
// imports unique entries into the actual memtable as a final step.
func (p *commitloggerParser) doReplace() error {
	nodeCache := make(map[string]segmentReplaceNode)

	var errWhileParsing error

	for {
		var version uint8

		err := binary.Read(p.checksumReader, binary.LittleEndian, &version)
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			errWhileParsing = errors.Wrap(err, "read commit version")
			break
		}

		switch version {
		case 0:
			{
				err = p.doReplaceRecordV0(nodeCache)
			}
		case 1:
			{
				err = p.doReplaceRecordV1(nodeCache)
			}
		default:
			{
				return fmt.Errorf("unsupported commit version %d", version)
			}
		}
		if err != nil {
			errWhileParsing = err
			break
		}
	}

	for _, node := range nodeCache {
		var opts []SecondaryKeyOption
		if p.memtable.secondaryIndices > 0 {
			for i, secKey := range node.secondaryKeys {
				opts = append(opts, WithSecondaryKey(i, secKey))
			}
		}
		if node.tombstone {
			p.memtable.setTombstone(node.primaryKey, opts...)
		} else {
			p.memtable.put(node.primaryKey, node.value, opts...)
		}
	}

	return errWhileParsing
}

func (p *commitloggerParser) doReplaceRecordV0(nodeCache map[string]segmentReplaceNode) error {
	var commitType CommitType

	err := binary.Read(p.reader, binary.LittleEndian, &commitType)
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return errors.Wrap(err, "read commit type")
	}

	if CommitTypeReplace.Is(commitType) {
		if err := p.parseReplaceNode(p.reader, nodeCache); err != nil {
			return errors.Wrap(err, "read replace node")
		}
	} else {
		return errors.Errorf("found a %s commit on a replace bucket", commitType.String())
	}

	return nil
}

func (p *commitloggerParser) doReplaceRecordV1(nodeCache map[string]segmentReplaceNode) error {
	var commitType CommitType

	err := binary.Read(p.checksumReader, binary.LittleEndian, &commitType)
	if err != nil {
		return errors.Wrap(err, "read commit type")
	}
	if !CommitTypeReplace.Is(commitType) {
		return errors.Errorf("found a %s commit on a replace bucket", commitType.String())
	}

	var nodeLen uint32
	err = binary.Read(p.checksumReader, binary.LittleEndian, &nodeLen)
	if err != nil {
		return errors.Wrap(err, "read commit node length")
	}

	p.bufNode.Reset()

	io.CopyN(p.bufNode, p.checksumReader, int64(nodeLen))

	// read checksum directly from the reader
	var checksum [4]byte
	_, err = io.ReadFull(p.reader, checksum[:])
	if err != nil {
		return errors.Wrap(err, "read commit checksum")
	}

	// validate checksum
	if !bytes.Equal(checksum[:], p.checksumReader.Hash()) {
		return errors.Wrap(ErrInvalidChecksum, "read commit entry")
	}

	if err := p.parseReplaceNode(p.bufNode, nodeCache); err != nil {
		return errors.Wrap(err, "read replace node")
	}

	return nil
}

// parseReplaceNode only parses into the deduplication cache, not into the
// final memtable yet. A second step is required to parse from the cache into
// the actual memtable.
func (p *commitloggerParser) parseReplaceNode(r io.Reader, nodeCache map[string]segmentReplaceNode) error {
	n, err := ParseReplaceNode(r, p.memtable.secondaryIndices)
	if err != nil {
		return err
	}

	if !n.tombstone {
		nodeCache[string(n.primaryKey)] = n
	} else {
		if existing, ok := nodeCache[string(n.primaryKey)]; ok {
			existing.tombstone = true
			nodeCache[string(n.primaryKey)] = existing
		} else {
			nodeCache[string(n.primaryKey)] = n
		}
	}

	return nil
}
