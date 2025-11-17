package scripts

import (
	"context"
	"fmt"
	"log"

	"github.com/adyutaa/parsea/internal/infrastructure/llm"
	"github.com/adyutaa/parsea/internal/infrastructure/vectordb"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	fmt.Println("🚀 Starting Qdrant document ingestion...")

	// Initialize Qdrant
	qdrant, err := vectordb.NewQdrantClient()
	if err != nil {
		log.Fatal("Failed to create Qdrant client:", err)
	}
	fmt.Println("✅ Connected to Qdrant Cloud")

	// Initialize OpenAI for embeddings
	llmClient := llm.NewOpenAIClient()
	if llmClient == nil {
		log.Fatal("Failed to create OpenAI client")
	}
	fmt.Println("✅ Connected to OpenAI")

	// Sample documents (you can expand this later)
	documents := []struct {
		Text     string
		Type     string
		Category string
	}{
		{
			Text: `Backend Engineer Position Requirements:
3+ years experience with backend technologies including Go, Python, or Node.js.
Strong database skills with PostgreSQL, MySQL, or MongoDB.
REST API design and implementation experience.
Experience with Redis or message queues for async processing.
Cloud platform experience (AWS, GCP, Azure).
Understanding of microservices architecture.
Experience with Docker and containerization.`,
			Type:     "job_description",
			Category: "requirements",
		},
		{
			Text: `AI and LLM Skills (Preferred):
Experience with LLM APIs such as OpenAI or Anthropic.
Understanding of prompt engineering best practices.
Knowledge of RAG systems and vector databases.
Experience with embeddings and similarity search.
Familiarity with LLM chaining and workflow orchestration.`,
			Type:     "job_description",
			Category: "ai_skills",
		},
		{
			Text: `Case Study Evaluation Criteria:
Correctness (30%): Proper REST API implementation, async job processing with queues, PDF parsing, LLM integration, RAG implementation.
Code Quality (25%): Clean and readable code, proper separation of concerns, modular components, good error handling.
Resilience (20%): Error handling throughout, retry logic for external APIs, graceful failure handling, input validation.
Documentation (15%): Clear README with setup instructions, architecture explanation, design decisions documented.
Bonus Features (10%): Additional features beyond requirements, performance optimizations, testing.`,
			Type:     "case_study",
			Category: "evaluation",
		},
		{
			Text: `CV Scoring Rubric:
Technical Skills Match (40%): Alignment with job requirements, proficiency in required technologies.
Experience Level (25%): Years of experience, complexity of past projects.
Relevant Achievements (20%): Impact of past work, measurable outcomes.
Cultural Fit (15%): Communication skills, collaboration abilities, learning mindset.`,
			Type:     "cv_rubric",
			Category: "scoring",
		},
		{
			Text: `Project Scoring Rubric:
Implementation Quality: Completeness of features, correctness of implementation.
Code Structure: Modularity, maintainability, use of best practices.
Error Handling: Robustness, retry mechanisms, graceful degradation.
Documentation: README quality, code comments, architecture diagrams.`,
			Type:     "project_rubric",
			Category: "scoring",
		},
	}

	ctx := context.Background()
	allDocs := []vectordb.Document{}
	allTexts := []string{}

	fmt.Println("\n📝 Preparing documents...")
	for i, doc := range documents {
		docID := fmt.Sprintf("%s_%s", doc.Type, uuid.New().String()[:8])

		allDocs = append(allDocs, vectordb.Document{
			ID:   docID,
			Text: doc.Text,
			Metadata: map[string]interface{}{
				"type":     doc.Type,
				"category": doc.Category,
				"index":    i,
			},
		})

		allTexts = append(allTexts, doc.Text)
	}

	fmt.Printf("   Prepared %d documents\n", len(allDocs))

	// Generate embeddings
	fmt.Println("\n🤖 Generating embeddings...")
	embeddings, err := llmClient.GenerateEmbeddings(ctx, allTexts)
	if err != nil {
		log.Fatal("Failed to generate embeddings:", err)
	}
	fmt.Printf("   Generated %d embeddings\n", len(embeddings))

	// Upload to Qdrant
	fmt.Println("\n📥 Uploading to Qdrant...")
	err = qdrant.AddDocuments(ctx, allDocs, embeddings)
	if err != nil {
		log.Fatal("Failed to upload documents:", err)
	}

	fmt.Println("\n✅ Document ingestion complete!")
	fmt.Printf("   Total documents stored: %d\n", len(allDocs))

	// Verify
	info, _ := qdrant.GetCollectionInfo(ctx)
	if info != nil {
		fmt.Printf("   Vectors in collection: %d\n", info.PointsCount)
	}
}
