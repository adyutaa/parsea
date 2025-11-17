package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/adyutaa/parsea/internal/infrastructure/llm"
	"github.com/adyutaa/parsea/internal/infrastructure/vectordb"
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
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		log.Fatal("OPENAI_API_KEY not set in environment")
	}
	llmClient := llm.NewOpenAIClient(openaiKey)
	if llmClient == nil {
		log.Fatal("Failed to create OpenAI client")
	}
	fmt.Println("✅ Connected to OpenAI")

	// Real documents from the case study brief
	documents := []struct {
		Text     string
		Type     string
		Category string
	}{
		{
			Text: `Job Description - Product Engineer (Backend) 2025

Rakamin is hiring a Product Engineer (Backend) to work on Rakamin. We're looking for dedicated engineers who write code they're proud of and who are eager to keep scaling and improving complex systems, including those powered by AI.

About the Job:
You'll be building new product features alongside a frontend engineer and product manager using our Agile methodology, as well as addressing issues to ensure our apps are robust and our codebase is clean. As a Product Engineer, you'll write clean, efficient code to enhance our product's codebase in meaningful ways.

In addition to classic backend work, this role also touches on building AI-powered systems, where you'll design and orchestrate how large language models (LLMs) integrate into Rakamin's product ecosystem.

Key Responsibilities:
- Collaborating with frontend engineers and 3rd parties to build robust backend solutions that support highly configurable platforms and cross-platform integration
- Developing and maintaining server-side logic for central database, ensuring high performance throughput and response time
- Designing and fine-tuning AI prompts that align with product requirements and user contexts
- Building LLM chaining flows, where the output from one model is reliably passed to and enriched by another
- Implementing Retrieval-Augmented Generation (RAG) by embedding and retrieving context from vector databases, then injecting it into AI prompts to improve accuracy and relevance
- Handling long-running AI processes gracefully — including job orchestration, async background workers, and retry mechanisms
- Designing safeguards for uncontrolled scenarios: managing failure cases from 3rd party APIs and mitigating the randomness/nondeterminism of LLM outputs

Required Skills:
- Experience with backend languages and frameworks (Node.js, Django, Rails)
- Database management (MySQL, PostgreSQL, MongoDB)
- RESTful APIs
- Security compliance
- Cloud technologies (AWS, Google Cloud, Azure)
- Server-side languages (Java, Python, Ruby, or JavaScript)
- Understanding of frontend technologies
- User authentication and authorization between multiple systems, servers, and environments
- Scalable application design principles
- Creating database schemas that represent and support business processes
- Implementing automated testing platforms and unit tests
- Familiarity with LLM APIs, embeddings, vector databases and prompt design best practices`,
			Type:     "job_description",
			Category: "requirements",
		},
		{
			Text: `Case Study Brief - Backend Developer Evaluation System

Objective:
Your mission is to build a backend service that automates the initial screening of a job application. The service will receive a candidate's CV and a project report, evaluate them against a specific job description and a case study brief, and produce a structured, AI-generated evaluation report.

Core Logic & Data Flow:
The system operates with a clear separation of inputs and reference documents:

Candidate-Provided Inputs (The Data to be Evaluated):
1. Candidate CV: The candidate's resume (PDF)
2. Project Report: The candidate's project report to our take-home case study (PDF)

System-Internal Documents (The "Ground Truth" for Comparison):
1. Job Description: A document detailing the requirements and responsibilities for the role
2. Case Study Brief: This document. Used as ground truth for Project Report
3. Scoring Rubric: A predefined set of parameters for evaluating CV and Report

Required API Endpoints:
- POST /upload: Accepts multipart/form-data containing the Candidate CV and Project Report (PDF)
- POST /evaluate: Triggers the asynchronous AI evaluation pipeline
- GET /result/{id}: Retrieves the status and result of an evaluation job

Evaluation Pipeline Components:
1. RAG (Context Retrieval): Ingest all System-Internal Documents into a vector database
2. Prompt Design & LLM Chaining: CV Evaluation, Project Report Evaluation, Final Analysis
3. Long-Running Process Handling: Asynchronous processing with job tracking
4. Error Handling & Randomness Control: Simulate edge cases and API failures`,
			Type:     "case_study_brief",
			Category: "requirements",
		},
		{
			Text: `CV Match Evaluation Rubric (1–5 scale per parameter)

Technical Skills Match (Weight: 40%)
Description: Alignment with job requirements (backend, databases, APIs, cloud, AI/LLM)
Scoring Guide: 1 = Irrelevant skills, 2 = Few overlaps, 3 = Partial match, 4 = Strong match, 5 = Excellent match + AI/LLM exposure

Experience Level (Weight: 25%)
Description: Years of experience and project complexity
Scoring Guide: 1 = <1 yr / trivial projects, 2 = 1–2 yrs, 3 = 2–3 yrs with mid-scale projects, 4 = 3–4 yrs solid track record, 5 = 5+ yrs / high-impact projects

Relevant Achievements (Weight: 20%)
Description: Impact of past work (scaling, performance, adoption)
Scoring Guide: 1 = No clear achievements, 2 = Minimal improvements, 3 = Some measurable outcomes, 4 = Significant contributions, 5 = Major measurable impact

Cultural / Collaboration Fit (Weight: 15%)
Description: Communication, learning mindset, teamwork/leadership
Scoring Guide: 1 = Not demonstrated, 2 = Minimal, 3 = Average, 4 = Good, 5 = Excellent and well-demonstrated

Final CV Match Rate: Weighted Average (1–5) → Convert to 0-1 decimal (×0.2)`,
			Type:     "cv_rubric",
			Category: "scoring",
		},
		{
			Text: `Project Deliverable Evaluation Rubric (1–5 scale per parameter)

Correctness (Prompt & Chaining) (Weight: 30%)
Description: Implements prompt design, LLM chaining, RAG context injection
Scoring Guide: 1 = Not implemented, 2 = Minimal attempt, 3 = Works partially, 4 = Works correctly, 5 = Fully correct + thoughtful

Code Quality & Structure (Weight: 25%)
Description: Clean, modular, reusable, tested
Scoring Guide: 1 = Poor, 2 = Some structure, 3 = Decent modularity, 4 = Good structure + some tests, 5 = Excellent quality + strong tests

Resilience & Error Handling (Weight: 20%)
Description: Handles long jobs, retries, randomness, API failures
Scoring Guide: 1 = Missing, 2 = Minimal, 3 = Partial handling, 4 = Solid handling, 5 = Robust, production-ready

Documentation & Explanation (Weight: 15%)
Description: README clarity, setup instructions, trade-off explanations
Scoring Guide: 1 = Missing, 2 = Minimal, 3 = Adequate, 4 = Clear, 5 = Excellent + insightful

Creativity / Bonus (Weight: 10%)
Description: Extra features beyond requirements
Scoring Guide: 1 = None, 2 = Very basic, 3 = Useful extras, 4 = Strong enhancements, 5 = Outstanding creativity

Final Project Score: Weighted Average (1–5)`,
			Type:     "project_rubric",
			Category: "scoring",
		},
		{
			Text: `AI and LLM Integration Requirements for Backend Engineers

Specific AI/LLM Skills Expected:
- Experience with LLM APIs such as OpenAI, Anthropic, or similar providers
- Understanding of prompt engineering best practices and temperature control
- Knowledge of RAG (Retrieval-Augmented Generation) systems and vector databases
- Experience with embeddings and semantic similarity search
- Familiarity with LLM chaining and workflow orchestration
- Building async background workers for long-running AI processes
- Implementing retry mechanisms and handling API failures from 3rd party services
- Managing the randomness and nondeterminism of LLM outputs
- Vector database integration (ChromaDB, Qdrant, Pinecone, etc.)
- Context injection and retrieval for improved AI accuracy and relevance

Technical Implementation Areas:
- Job orchestration for AI workflows
- Async background processing with queue systems
- Error handling for external API dependencies
- Safeguards against uncontrolled AI scenarios
- Performance optimization for AI-driven features
- Testing strategies for non-deterministic AI outputs
- Monitoring and observability for AI systems`,
			Type:     "ai_requirements",
			Category: "technical_skills",
		},
		{
			Text: `Overall Candidate Evaluation Framework

Evaluation Process:
1. CV Match Rate: Weighted Average (1–5) → Convert to 0-1 decimal (×0.2)
2. Project Score: Weighted Average (1–5) 
3. Overall Summary: Service should return 3–5 sentences (strengths, gaps, recommendations)

Expected Output Format:
{
  "cv_match_rate": 0.82,
  "cv_feedback": "Strong in backend and cloud, limited AI integration experience...",
  "project_score": 4.5,
  "project_feedback": "Meets prompt chaining requirements, lacks error handling robustness...",
  "overall_summary": "Good candidate fit, would benefit from deeper RAG knowledge..."
}

Evaluation Quality Factors:
- Consistency: Structured prompts with stable scoring
- Relevance: Job context injection improves accuracy
- Transparency: Detailed feedback explains scoring rationale
- Scalability: Background workers handle concurrent evaluations

Key Areas for Assessment:
- Backend engineering fundamentals
- AI/LLM integration capabilities
- System design and architecture skills
- Error handling and resilience
- Code quality and testing practices
- Documentation and communication skills`,
			Type:     "evaluation_framework",
			Category: "process",
		},
	}

	ctx := context.Background()
	allDocs := []vectordb.Document{}
	allTexts := []string{}

	fmt.Println("\n📝 Preparing documents...")
	for i, doc := range documents {
		docID := fmt.Sprintf("doc_%d", i+1)

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

	// Convert embeddings to float32 (Qdrant requirement)
	fmt.Println("\n🔄 Converting embeddings...")
	float32Embeddings := make([][]float32, len(embeddings))
	for i, embedding := range embeddings {
		float32Embedding := make([]float32, len(embedding))
		for j, val := range embedding {
			float32Embedding[j] = float32(val)
		}
		float32Embeddings[i] = float32Embedding
	}

	// Upload to Qdrant
	fmt.Println("\n📥 Uploading to Qdrant...")
	err = qdrant.AddDocuments(ctx, allDocs, float32Embeddings)
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
