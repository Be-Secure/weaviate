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

package inverted

import (
	"context"
	"fmt"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/weaviate/weaviate/adapters/repos/db/helpers"
	"github.com/weaviate/weaviate/adapters/repos/db/lsmkv"
	"github.com/weaviate/weaviate/entities/additional"
	enterrors "github.com/weaviate/weaviate/entities/errors"
	"github.com/weaviate/weaviate/entities/models"
	"github.com/weaviate/weaviate/entities/schema"
	"github.com/weaviate/weaviate/entities/searchparams"
	"github.com/weaviate/weaviate/entities/storobj"

	"github.com/weaviate/weaviate/adapters/repos/db/inverted/stopwords"
	"github.com/weaviate/weaviate/adapters/repos/db/inverted/terms"
)

var metrics = lsmkv.BlockMetrics{}

func (b *BM25Searcher) createBlockTerm(N float64, filterDocIds helpers.AllowList, query []string, propName string, propertyBoost float32, duplicateTextBoosts []int, ctx context.Context) ([][]terms.TermInterface, error) {
	bucket := b.store.Bucket(helpers.BucketSearchableFromPropNameLSM(propName))
	return bucket.CreateDiskTerm(N, filterDocIds, query, propName, propertyBoost, duplicateTextBoosts, ctx)
}

func (b *BM25Searcher) wandBlock(
	ctx context.Context, filterDocIds helpers.AllowList, class *models.Class, params searchparams.KeywordRanking, limit int, additional additional.Properties,
) ([]*storobj.Object, []float32, error) {
	params.AdditionalExplanations = false
	N := float64(b.store.Bucket(helpers.ObjectsBucketLSM).Count())

	var stopWordDetector *stopwords.Detector
	if class.InvertedIndexConfig != nil && class.InvertedIndexConfig.Stopwords != nil {
		var err error
		stopWordDetector, err = stopwords.NewDetectorFromConfig(*(class.InvertedIndexConfig.Stopwords))
		if err != nil {
			return nil, nil, err
		}
	}

	// There are currently cases, for different tokenization:
	// word, lowercase, whitespace and field.
	// Query is tokenized and respective properties are then searched for the search terms,
	// results at the end are combined using WAND

	queryTermsByTokenization := map[string][]string{}
	duplicateBoostsByTokenization := map[string][]int{}
	propNamesByTokenization := map[string][]string{}
	propertyBoosts := make(map[string]float32, len(params.Properties))

	for _, tokenization := range helpers.Tokenizations {
		queryTerms, dupBoosts := helpers.TokenizeAndCountDuplicates(tokenization, params.Query)
		queryTermsByTokenization[tokenization] = queryTerms
		duplicateBoostsByTokenization[tokenization] = dupBoosts

		// stopword filtering for word tokenization
		if tokenization == models.PropertyTokenizationWord {
			queryTerms, dupBoosts = b.removeStopwordsFromQueryTerms(queryTermsByTokenization[tokenization],
				duplicateBoostsByTokenization[tokenization], stopWordDetector)
			queryTermsByTokenization[tokenization] = queryTerms
			duplicateBoostsByTokenization[tokenization] = dupBoosts
		}

		propNamesByTokenization[tokenization] = make([]string, 0)
	}

	averagePropLength := 0.
	for _, propertyWithBoost := range params.Properties {
		property := propertyWithBoost
		propBoost := 1
		if strings.Contains(propertyWithBoost, "^") {
			property = strings.Split(propertyWithBoost, "^")[0]
			boostStr := strings.Split(propertyWithBoost, "^")[1]
			propBoost, _ = strconv.Atoi(boostStr)
		}
		propertyBoosts[property] = float32(propBoost)

		propMean, err := b.GetPropertyLengthTracker().PropertyMean(property)
		if err != nil {
			return nil, nil, err
		}
		averagePropLength += float64(propMean)

		prop, err := schema.GetPropertyByName(class, property)
		if err != nil {
			return nil, nil, err
		}

		switch dt, _ := schema.AsPrimitive(prop.DataType); dt {
		case schema.DataTypeText, schema.DataTypeTextArray:
			if _, exists := propNamesByTokenization[prop.Tokenization]; !exists {
				return nil, nil, fmt.Errorf("cannot handle tokenization '%v' of property '%s'",
					prop.Tokenization, prop.Name)
			}
			propNamesByTokenization[prop.Tokenization] = append(propNamesByTokenization[prop.Tokenization], property)
		default:
			return nil, nil, fmt.Errorf("cannot handle datatype '%v' of property '%s'", dt, prop.Name)
		}
	}

	averagePropLength = averagePropLength / float64(len(params.Properties))

	// 100 is a reasonable expected capacity for the total number of terms to query.
	allResults := make([][][]terms.TermInterface, len(params.Properties))

	for _, tokenization := range helpers.Tokenizations {
		propNames := propNamesByTokenization[tokenization]
		if len(propNames) > 0 {
			queryTerms, duplicateBoosts := helpers.TokenizeAndCountDuplicates(tokenization, params.Query)

			// stopword filtering for word tokenization
			if tokenization == models.PropertyTokenizationWord {
				queryTerms, duplicateBoosts = b.removeStopwordsFromQueryTerms(
					queryTerms, duplicateBoosts, stopWordDetector)
			}
			for i, propName := range propNames {
				results, err := b.createBlockTerm(N, filterDocIds, queryTerms, propName, propertyBoosts[propName], duplicateBoosts, ctx)
				if err != nil {
					return nil, nil, err
				}
				allResults[i] = append(allResults[i], results...)
			}

		}
	}

	// all results. Sum up the length of the results from all terms to get an upper bound of how many results there are
	if limit == 0 {
		for _, res := range allResults {
			for _, res2 := range res {
				for _, res3 := range res2 {
					if res3 != nil {
						limit += res3.Count()
					}
				}
			}
		}
	}

	eg := enterrors.NewErrorGroupWrapper(b.logger)
	eg.SetLimit(_NUMCPU)

	allObjects := make([][][]*storobj.Object, len(allResults))
	allScores := make([][][]float32, len(allResults))

	for i, result1 := range allResults {
		allObjects[i] = make([][]*storobj.Object, len(result1))
		allScores[i] = make([][]float32, len(result1))
		for j, result2 := range result1 {
			i := i
			j := j
			eg.Go(func() (err error) {
				combinedTerms := &terms.Terms{
					T: result2,
				}

				topKHeap := b.doBlockMaxWand(limit, combinedTerms, averagePropLength, params.AdditionalExplanations)
				objects, scores, err := b.getTopKObjects(topKHeap, params.AdditionalExplanations, nil, additional)

				allObjects[i][j] = objects
				allScores[i][j] = scores
				if err != nil {
					return err
				}
				return nil
			})
		}
	}

	if err := eg.Wait(); err != nil {
		return nil, nil, err
	}

	metrics.QueryCount++
	for _, result := range allResults {
		for _, result2 := range result {
			for _, result3 := range result2 {
				if result3 != nil {
					m := result3.(*lsmkv.SegmentBlockMax)
					metrics.BlockCountAdded += m.Metrics.BlockCountAdded
					metrics.BlockCountTotal += m.Metrics.BlockCountTotal
					metrics.BlockCountExamined += m.Metrics.BlockCountExamined
					metrics.DocCountAdded += m.Metrics.DocCountAdded
					metrics.DocCountTotal += m.Metrics.DocCountTotal
					metrics.DocCountExamined += m.Metrics.DocCountExamined
					metrics.LastAddedBlock = m.Metrics.LastAddedBlock

				}
			}
		}
	}

	if metrics.QueryCount%100 == 0 {
		b.logger.Error("BlockMax metrics", "BlockCountTotal ", metrics.BlockCountTotal/metrics.QueryCount, " BlockCountExamined ", metrics.BlockCountExamined/metrics.QueryCount, " BlockCountAdded ", metrics.BlockCountAdded/metrics.QueryCount, " DocCountTotal ", metrics.DocCountTotal/metrics.QueryCount, " DocCountExamined ", metrics.DocCountExamined/metrics.QueryCount, " DocCountAdded ", metrics.DocCountAdded/metrics.QueryCount)
		metrics = lsmkv.BlockMetrics{}
	}

	objects, scores := b.combineResults(allObjects, allScores, limit)

	return objects, scores, nil
}

func (b *BM25Searcher) combineResults(allObjects [][][]*storobj.Object, allScores [][][]float32, limit int) ([]*storobj.Object, []float32) {
	// combine all results
	combinedObjects := make([]*storobj.Object, 0, limit*len(allObjects))
	combinedScores := make([]float32, 0, limit*len(allObjects))

	// combine all results
	for i := range allObjects {
		singlePropObjects := slices.Concat(allObjects[i]...)
		singlePropScores := slices.Concat(allScores[i]...)
		combinedObjectsProp, combinedScoresProp := b.combineResultsForMultiProp(singlePropObjects, singlePropScores, func(a, b float32) float32 { return b })
		combinedObjects = append(combinedObjects, combinedObjectsProp...)
		combinedScores = append(combinedScores, combinedScoresProp...)
	}

	combinedObjects, combinedScores = b.combineResultsForMultiProp(combinedObjects, combinedScores, func(a, b float32) float32 { return a + b })

	combinedObjects, combinedScores = b.sortResultsByScore(combinedObjects, combinedScores)

	if len(combinedObjects) <= limit {
		return combinedObjects, combinedScores
	}

	return combinedObjects[len(combinedObjects)-limit:], combinedScores[len(combinedObjects)-limit:]
}

type aggregate func(float32, float32) float32

func (b *BM25Searcher) combineResultsForMultiProp(allObjects []*storobj.Object, allScores []float32, aggregateFn aggregate) ([]*storobj.Object, []float32) {
	// if ids are the same, sum the scores
	combinedObjects := make(map[string]*storobj.Object)
	combinedScores := make(map[string]float32)

	for i, obj := range allObjects {
		id := string(obj.ID())
		if _, ok := combinedObjects[id]; !ok {
			combinedObjects[id] = obj
			combinedScores[id] = allScores[i]
		} else {
			combinedScores[id] = aggregateFn(combinedScores[id], allScores[i])
		}
	}

	// sort the combined results
	combinedObjectsSlice := make([]*storobj.Object, 0, len(combinedObjects))
	combinedScoresSlice := make([]float32, 0, len(combinedObjects))

	for id, obj := range combinedObjects {
		combinedObjectsSlice = append(combinedObjectsSlice, obj)
		combinedScoresSlice = append(combinedScoresSlice, combinedScores[id])
	}

	return combinedObjectsSlice, combinedScoresSlice
}

func (b *BM25Searcher) sortResultsByScore(objects []*storobj.Object, scores []float32) ([]*storobj.Object, []float32) {
	sorter := &scoreSorter{
		objects: objects,
		scores:  scores,
	}
	sort.Sort(sorter)
	return sorter.objects, sorter.scores
}

type scoreSorter struct {
	objects []*storobj.Object
	scores  []float32
}

func (s *scoreSorter) Len() int {
	return len(s.objects)
}

func (s *scoreSorter) Less(i, j int) bool {
	if s.scores[i] == s.scores[j] {
		return s.objects[i].ID() < s.objects[j].ID()
	}
	return s.scores[i] < s.scores[j]
}

func (s *scoreSorter) Swap(i, j int) {
	s.objects[i], s.objects[j] = s.objects[j], s.objects[i]
	s.scores[i], s.scores[j] = s.scores[j], s.scores[i]
}
