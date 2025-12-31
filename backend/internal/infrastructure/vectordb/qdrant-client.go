package vectordb

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/qdrant/go-client/qdrant"
)

type QdrantClient struct {
	client         *qdrant.Client
	collectionName string
}

func NewQdrantClient() (*QdrantClient, error) {
	host := os.Getenv("QDRANT_HOST")
	portStr := os.Getenv("QDRANT_PORT")
	apiKey := os.Getenv("QDRANT_API_KEY")

	if host == "" || portStr == "" || apiKey == "" {
		return nil, fmt.Errorf("Qdrant configuration not set")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid QDRANT_PORT: %w", err)
	}

	// Create Qdrant client
	client, err := qdrant.NewClient(&qdrant.Config{
		Host:   host,
		Port:   port,
		APIKey: apiKey,
		UseTLS: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Qdrant client: %w", err)
	}

	collectionName := "job_requirements"

	// Check if collection exists
	ctx := context.Background()
	collections, err := client.ListCollections(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list collections: %w", err)
	}

	collectionExists := false
	for _, colName := range collections {
		if colName == collectionName {
			collectionExists = true
			break
		}
	}

	// Create collection if it doesn't exist
	if !collectionExists {
		err = client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: collectionName,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     1536, // OpenAI text-embedding-ada-002 dimension
				Distance: qdrant.Distance_Cosine,
			}),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}
		fmt.Println("✅ Created Qdrant collection:", collectionName)
	}

	return &QdrantClient{
		client:         client,
		collectionName: collectionName,
	}, nil
}

// Document represents a document to be stored
type Document struct {
	ID       string
	Text     string
	Metadata map[string]interface{}
}

// AddDocuments adds documents to Qdrant with their embeddings
func (q *QdrantClient) AddDocuments(ctx context.Context, docs []Document, embeddings [][]float32) error {
	if len(docs) == 0 {
		return nil
	}

	if len(embeddings) != len(docs) {
		return fmt.Errorf("number of embeddings (%d) must match number of documents (%d)", len(embeddings), len(docs))
	}

	points := make([]*qdrant.PointStruct, len(docs))

	for i, doc := range docs {
		// Convert metadata to Qdrant payload
		payload := make(map[string]*qdrant.Value)
		payload["text"] = qdrant.NewValueString(doc.Text)

		for key, val := range doc.Metadata {
			switch v := val.(type) {
			case string:
				payload[key] = qdrant.NewValueString(v)
			case int:
				payload[key] = qdrant.NewValueInt(int64(v))
			case float64:
				payload[key] = qdrant.NewValueDouble(v)
			case bool:
				payload[key] = qdrant.NewValueBool(v)
			}
		}

		points[i] = &qdrant.PointStruct{
			Id:      &qdrant.PointId{PointIdOptions: &qdrant.PointId_Uuid{Uuid: doc.ID}},
			Vectors: qdrant.NewVectors(embeddings[i]...),
			Payload: payload,
		}
	}

	_, err := q.client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: q.collectionName,
		Points:         points,
	})

	return err
}

// SearchResult represents a search result
type SearchResult struct {
	ID       string
	Text     string
	Score    float32
	Metadata map[string]interface{}
}

// Search searches for similar documents using a query embedding
func (q *QdrantClient) Search(ctx context.Context, queryEmbedding []float32, limit uint64) ([]SearchResult, error) {
	searchResult, err := q.client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: q.collectionName,
		Query:          qdrant.NewQuery(queryEmbedding...),
		Limit:          &limit,
		WithPayload:    qdrant.NewWithPayload(true),
	})

	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	results := make([]SearchResult, 0, len(searchResult))

	for _, point := range searchResult {
		metadata := make(map[string]interface{})
		text := ""

		// Extract payload
		if point.Payload != nil {
			for key, val := range point.Payload {
				if key == "text" {
					if stringVal := val.GetStringValue(); stringVal != "" {
						text = stringVal
					}
				} else {
					// Store other metadata
					if stringVal := val.GetStringValue(); stringVal != "" {
						metadata[key] = stringVal
					} else if intVal := val.GetIntegerValue(); intVal != 0 {
						metadata[key] = intVal
					} else if doubleVal := val.GetDoubleValue(); doubleVal != 0 {
						metadata[key] = doubleVal
					} else if boolVal := val.GetBoolValue(); boolVal {
						metadata[key] = boolVal
					}
				}
			}
		}

		var id string
		if point.Id != nil {
			if stringId := point.Id.GetUuid(); stringId != "" {
				id = stringId
			} else if numId := point.Id.GetNum(); numId != 0 {
				id = fmt.Sprintf("%d", numId)
			}
		}

		results = append(results, SearchResult{
			ID:       id,
			Text:     text,
			Score:    point.Score,
			Metadata: metadata,
		})
	}

	return results, nil
}

// DeleteCollection deletes the collection
func (q *QdrantClient) DeleteCollection(ctx context.Context) error {
	return q.client.DeleteCollection(ctx, q.collectionName)
}

// GetCollectionInfo returns collection information
func (q *QdrantClient) GetCollectionInfo(ctx context.Context) (*qdrant.CollectionInfo, error) {
	return q.client.GetCollectionInfo(ctx, q.collectionName)
}
