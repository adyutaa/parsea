package worker

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/adyutaa/parsea/internal/domain"
	"github.com/adyutaa/parsea/internal/infrastructure/llm"
	"github.com/adyutaa/parsea/internal/repository"
	"github.com/adyutaa/parsea/internal/service"
	"github.com/adyutaa/parsea/pkg/pdf"
	"github.com/redis/go-redis/v9"
)

type EvaluationWorker struct {
	redis          *redis.Client
	evalRepo       *repository.EvaluationRepository
	docRepo        *repository.DocumentRepository
	llmClient      *llm.OpenAIClient
	contextService *service.ContextService
	pdfParser      *pdf.Parser
}

func NewEvaluationWorker(
	redis *redis.Client,
	evalRepo *repository.EvaluationRepository,
	docRepo *repository.DocumentRepository,
	llmClient *llm.OpenAIClient,
	contextService *service.ContextService,
) *EvaluationWorker {
	return &EvaluationWorker{
		redis:          redis,
		evalRepo:       evalRepo,
		docRepo:        docRepo,
		llmClient:      llmClient,
		contextService: contextService,
		pdfParser:      pdf.NewParser(),
	}
}

// Start begins processing jobs from the queue
func (w *EvaluationWorker) Start(ctx context.Context) {
	log.Println("🔄 Worker started, waiting for jobs...")

	for {
		select {
		case <-ctx.Done():
			log.Println("👋 Worker shutting down...")
			return
		default:
			w.processNextJob(ctx)
		}
	}
}

func (w *EvaluationWorker) processNextJob(ctx context.Context) {
	// Block and wait for job (timeout 1 second for faster shutdown)
	result, err := w.redis.BRPop(ctx, 1*time.Second, "evaluation_queue").Result()
	if err != nil {
		if err.Error() != "redis: nil" && err != context.Canceled {
			log.Printf("⚠️  Failed to pop from queue: %v\n", err)
		}
		return
	}

	if len(result) < 2 {
		return
	}

	jobID := result[1]
	log.Printf("\n" + strings.Repeat("=", 60))
	log.Printf("📋 Processing job: %s", jobID)
	log.Printf(strings.Repeat("=", 60) + "\n")

	// Process the job with timeout
	jobCtx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if err := w.processJob(jobCtx, jobID); err != nil {
		log.Printf("❌ Job %s failed: %v\n", jobID, err)
		w.evalRepo.UpdateError(jobID, err.Error())
	} else {
		log.Printf("\n✅ Job %s completed successfully!\n", jobID)
	}
}

func (w *EvaluationWorker) processJob(ctx context.Context, jobID string) error {
	// Update status to processing
	if err := w.evalRepo.UpdateStatus(jobID, "processing"); err != nil {
		return fmt.Errorf("failed to update status: %w", err)
	}

	// Get job details
	job, err := w.evalRepo.GetByID(jobID)
	if err != nil {
		return fmt.Errorf("failed to get job: %w", err)
	}

	// Get CV document
	cvDoc, err := w.docRepo.GetByID(job.CVID)
	if err != nil {
		return fmt.Errorf("failed to get CV document: %w", err)
	}

	// Get Project Report document
	reportDoc, err := w.docRepo.GetByID(job.ReportID)
	if err != nil {
		return fmt.Errorf("failed to get report document: %w", err)
	}

	// ========================================
	// STEP 1: Extract text from CV
	// ========================================
	log.Println("📄 [1/7] Extracting text from CV...")
	cvText, err := w.pdfParser.ExtractText(cvDoc.FilePath)
	if err != nil {
		return fmt.Errorf("failed to extract CV text: %w", err)
	}
	cvText = w.pdfParser.CleanText(cvText)
	log.Printf("   ✅ Extracted %d characters from CV\n", len(cvText))

	// ========================================
	// STEP 2: Get job requirements context (RAG!)
	// ========================================
	log.Println("\n🔍 [2/7] Retrieving job requirements context (RAG)...")
	var jobContext string
	if w.contextService != nil {
		jobContext, err = w.contextService.GetJobRequirementsContext(ctx, job.JobTitle)
		if err != nil {
			log.Printf("   ⚠️  RAG failed, using fallback: %v\n", err)
			jobContext = service.GetHardcodedJobContext()
		} else {
			log.Println("   ✅ Retrieved context from vector database")
		}
	} else {
		log.Println("   ⚠️  No context service, using fallback")
		jobContext = service.GetHardcodedJobContext()
	}

	// Add CV scoring rubric
	if w.contextService != nil {
		cvRubric, _ := w.contextService.GetCVScoringContext(ctx)
		jobContext += "\n\n" + cvRubric
	} else {
		jobContext += "\n\n" + service.GetCVScoringRubric()
	}

	// ========================================
	// STEP 3: Evaluate CV with LLM
	// ========================================
	log.Println("\n🤖 [3/7] Evaluating CV with LLM...")
	cvResult, err := w.llmClient.EvaluateCV(cvText, jobContext)
	if err != nil {
		return fmt.Errorf("failed to evaluate CV: %w", err)
	}
	log.Printf("   ✅ CV Match Rate: %.2f (%.0f%%)\n", cvResult.MatchRate, cvResult.MatchRate*100)

	// ========================================
	// STEP 4: Extract text from Project Report
	// ========================================
	log.Println("\n📄 [4/7] Extracting text from Project Report...")
	reportText, err := w.pdfParser.ExtractText(reportDoc.FilePath)
	if err != nil {
		return fmt.Errorf("failed to extract report text: %w", err)
	}
	reportText = w.pdfParser.CleanText(reportText)
	log.Printf("   ✅ Extracted %d characters from report\n", len(reportText))

	// ========================================
	// STEP 5: Get case study context (RAG!)
	// ========================================
	log.Println("\n🔍 [5/7] Retrieving case study context (RAG)...")
	var caseContext string
	if w.contextService != nil {
		caseContext, err = w.contextService.GetCaseStudyContext(ctx)
		if err != nil {
			log.Printf("   ⚠️  RAG failed, using fallback: %v\n", err)
			caseContext = service.GetHardcodedCaseStudyContext()
		} else {
			log.Println("   ✅ Retrieved context from vector database")
		}
	} else {
		log.Println("   ⚠️  No context service, using fallback")
		caseContext = service.GetHardcodedCaseStudyContext()
	}

	// Add project scoring rubric
	if w.contextService != nil {
		projectRubric, _ := w.contextService.GetProjectScoringContext(ctx)
		caseContext += "\n\n" + projectRubric
	} else {
		caseContext += "\n\n" + service.GetProjectScoringRubric()
	}

	// ========================================
	// STEP 6: Evaluate Project with LLM
	// ========================================
	log.Println("\n🤖 [6/7] Evaluating Project with LLM...")
	projectResult, err := w.llmClient.EvaluateProject(reportText, caseContext)
	if err != nil {
		return fmt.Errorf("failed to evaluate project: %w", err)
	}
	log.Printf("   ✅ Project Score: %.1f/5.0\n", projectResult.Score)

	// ========================================
	// STEP 7: Generate final summary
	// ========================================
	log.Println("\n🤖 [7/7] Generating final summary...")
	summary, err := w.llmClient.GenerateSummary(
		cvResult.Feedback,
		projectResult.Feedback,
		cvResult.MatchRate,
		projectResult.Score,
	)
	if err != nil {
		return fmt.Errorf("failed to generate summary: %w", err)
	}
	log.Println("   ✅ Summary generated")

	// ========================================
	// Save results
	// ========================================
	log.Println("\n💾 Saving results to database...")
	result := &domain.EvaluationResult{
		CVMatchRate:     cvResult.MatchRate,
		CVFeedback:      cvResult.Feedback,
		ProjectScore:    projectResult.Score,
		ProjectFeedback: projectResult.Feedback,
		OverallSummary:  summary,
	}

	if err := w.evalRepo.UpdateResult(jobID, result); err != nil {
		return fmt.Errorf("failed to save results: %w", err)
	}

	log.Println("   ✅ Results saved")

	return nil
}