package service

import (
	"context"

	"github.com/adyutaa/parsea/internal/infrastructure/llm"
	"github.com/adyutaa/parsea/internal/infrastructure/vectordb"
)

type ContextService struct {
	qdrant    *vectordb.QdrantClient
	llmClient *llm.OpenAIClient
}

func NewContextService(qdrant *vectordb.QdrantClient, llmClient *llm.OpenAIClient) *ContextService {
	return &ContextService{
		qdrant:    qdrant,
		llmClient: llmClient,
	}
}

// GetJobRequirementsContext retrieves relevant context from vector DB
func (s *ContextService) GetJobRequirementsContext(ctx context.Context, jobTitle string) (string, error) {
	// TODO: Implement vector search once embedding functions are ready
	// For now, return hardcoded context
	return GetHardcodedJobContext(), nil
}

// GetCaseStudyContext retrieves case study requirements from vector DB
func (s *ContextService) GetCaseStudyContext(ctx context.Context) (string, error) {
	// TODO: Implement vector search once embedding functions are ready
	return GetHardcodedCaseStudyContext(), nil
}

// GetCVScoringContext retrieves CV scoring rubric
func (s *ContextService) GetCVScoringContext(ctx context.Context) (string, error) {
	return GetCVScoringRubric(), nil
}

// GetProjectScoringContext retrieves project scoring rubric
func (s *ContextService) GetProjectScoringContext(ctx context.Context) (string, error) {
	return GetProjectScoringRubric(), nil
}

// Hardcoded fallbacks (untuk development atau jika RAG gagal)

func GetHardcodedJobContext() string {
	return `Backend Engineer Position Requirements:

Technical Skills:
- 3+ years experience with backend technologies
- Proficiency in Go, Python, or Node.js
- Strong database skills (PostgreSQL, MySQL, MongoDB)
- REST API design and implementation
- Experience with Redis or message queues
- Cloud platform experience (AWS, GCP, Azure)
- Understanding of microservices architecture

AI/LLM Skills (Preferred):
- Experience with LLM APIs
- Understanding of prompt engineering
- Knowledge of RAG systems and vector databases

Evaluation Focus:
- Technical skills match with requirements
- Years of relevant experience
- Quality of past projects
- Communication and collaboration skills`
}

func GetHardcodedCaseStudyContext() string {
	return `Case Study Evaluation Criteria:

1. Correctness (30%): REST API implementation, async job processing, LLM integration, RAG implementation
2. Code Quality (25%): Clean, modular, well-documented code
3. Resilience (20%): Error handling, retries, graceful failures
4. Documentation (15%): Clear README, setup instructions, design explanations
5. Bonus Features (10%): Additional features beyond requirements

Scoring: 1-5 scale (1=Insufficient, 5=Exceptional)`
}

func GetCVScoringRubric() string {
	return `CV Scoring Rubric:
- Technical Skills Match (40%): 1-5 based on alignment with job requirements
- Experience Level (25%): 1-5 based on years and project complexity
- Achievements (20%): 1-5 based on measurable impact
- Cultural Fit (15%): 1-5 based on communication and collaboration`
}

func GetProjectScoringRubric() string {
	return `Project Scoring Rubric:
- Implementation Correctness (30%): 1-5 based on feature completeness
- Code Quality (25%): 1-5 based on structure and maintainability
- Resilience (20%): 1-5 based on error handling
- Documentation (15%): 1-5 based on clarity
- Bonus Features (10%): 1-5 based on extra features`
}
