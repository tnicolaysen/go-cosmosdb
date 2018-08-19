package cosmosdb

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
)

type Collection struct {
	Resource
	IndexingPolicy *IndexingPolicy `json:"indexingPolicy,omitempty"`
	Docs           string          `json:"_docs,omitempty"`
	Udf            string          `json:"_udfs,omitempty"`
	Sprocs         string          `json:"_sprocs,omitempty"`
	Triggers       string          `json:"_triggers,omitempty"`
	Conflicts      string          `json:"_conflicts,omitempty"`
	PartitionKey   *PartitionKey   `json:"partitionKey,omitempty"`
}

type DocumentCollection struct {
	Rid                 string       `json:"_rid,omitempty"`
	Count               int32        `json:"_count,omitempty"`
	DocumentCollections []Collection `json:"DocumentCollections"`
}

type IndexingPolicy struct {
	IndexingMode IndexingMode   `json:"indexingMode,omitempty"`
	Automatic    bool           `json:"automatic"`
	Included     []IncludedPath `json:"includedPaths,omitempty"`
	Excluded     []ExcludedPath `json:"excludedPaths,omitempty"`
}

type IndexingMode string
type OfferThroughput int32
type OfferType string

//const (
//	Consistent = IndexingMode("Consistent")
//	Lazy       = IndexingMode("Lazy")
//)
//
//const (
//	OfferTypeS1 = OfferType("S1")
//	OfferTypeS2 = OfferType("S2")
//	OfferTypeS3 = OfferType("S3")
//)

type PartitionKey struct {
	Paths []string `json:"paths"`
	Kind  string   `json:"kind"`
}

// https://docs.microsoft.com/en-us/rest/api/cosmos-db/create-a-collection
type CollectionCreateOptions struct {
	Id             string          `json:"id"`
	IndexingPolicy *IndexingPolicy `json:"indexingPolicy,omitempty"`
	PartitionKey   *PartitionKey   `json:"partitionKey,omitempty"`

	// RTUs [400 - 250000]. Do not use in combination with OfferType
	OfferThroughput OfferThroughput `json:"offerThroughput,omitempty"`
	// S1,S2,S3. Do not use in combination with OfferThroughput
	OfferType OfferType `json:"offerType,omitempty"`
}

type CollectionReplaceOptions struct {
	Resource
	Id             string          `json:"id"`
	IndexingPolicy *IndexingPolicy `json:"indexingPolicy,omitempty"`
	PartitionKey   *PartitionKey   `json:"partitionKey,omitempty"`
}

// https://docs.microsoft.com/en-us/rest/api/cosmos-db/create-a-collection
func (c *Client) CreateCollection(ctx context.Context, dbName string,
	colOps CollectionCreateOptions, ops *RequestOptions) (*Collection, error) {

	headers := make(map[string]string)

	if colOps.OfferThroughput > 0 {
		headers[HEADER_OFFER_THROUGHPUT] = fmt.Sprintf("%d", colOps.OfferThroughput)
	}

	if colOps.OfferThroughput >= 10000 && colOps.PartitionKey == nil {
		return nil, errors.New(fmt.Sprintf("Must specify PartitionKey for collection '%s' when OfferThroughput is >= 10000", colOps.Id))
	}

	if colOps.OfferType != "" {
		headers[HEADER_OFFER_TYPE] = fmt.Sprintf("%s", colOps.OfferType)
	}

	collection := &Collection{}
	link := CreateCollLink(dbName, "")

	err := c.create(ctx, link, colOps, collection, headers)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

// https://docs.microsoft.com/en-us/rest/api/cosmos-db/list-collections
func (c *Client) ListCollections(ctx context.Context, dbName string,
	ops *RequestOptions) (*DocumentCollection, error) {
	url := createDatabaseLink(dbName) + "/colls"

	docCol := &DocumentCollection{}
	err := c.get(ctx, url, docCol, nil)
	if err != nil {
		return nil, err
	}

	return docCol, nil
}

func (c *Client) GetCollection(ctx context.Context, dbName, colName string,
	ops *RequestOptions) (*Collection, error) {
	return nil, ErrorNotImplemented
}

func (c *Client) DeleteCollection(ctx context.Context, dbName, colName string,
	ops *RequestOptions) error {
	return ErrorNotImplemented
}

// https://docs.microsoft.com/en-us/rest/api/cosmos-db/replace-a-collection
func (c *Client) ReplaceCollection(ctx context.Context, dbName string,
	colOps CollectionReplaceOptions, ops *RequestOptions) (*Collection, error) {

	collection := &Collection{}
	link := CreateCollLink(dbName, colOps.Id)

	err := c.replace(ctx, link, colOps, collection, nil)
	if err != nil {
		return nil, err
	}

	return collection, nil
}

// TODO: add model for partition key ranges
func (c *Client) GetPartitionKeyRanges(ctx context.Context, dbName, colName string,
	ops *RequestOptions) error {
	return ErrorNotImplemented
}
