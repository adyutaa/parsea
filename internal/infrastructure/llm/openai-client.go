package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/openai/openai-go"
)

type OpenAIClient struct {
	client *openai.Client
}

func NewOpenAIClient(apiKey string) *OpenAIClient {
	client := openai.NewClient()
	return &OpenAIClient{
		client: &client,
	}
}

// CVEvaluationResult represents the structured output for CV evaluation
type CVEvaluationResult struct {
	MatchRate float64 `json:"match_rate"`
	Feedback  string  `json:"feedback"`
}

// ProjectEvaluationResult represents the structured output for project evaluation
type ProjectEvaluationResult struct {
	Score    float64 `json:"score"`
	Feedback string  `json:"feedback"`
}

// EvaluateCV evaluates a candidate's CV against job requirements
func (c *OpenAIClient) EvaluateCV(cvText, jobContext string) (*CVEvaluationResult, error) {
	prompt := fmt.Sprintf(`You are an expert technical recruiter evaluating a candidate's CV for a Backend Engineer position.

Job Requirements and Context:
%s

Candidate's CV:
%s

Evaluate this CV and provide:
1. match_rate: A score from 0.0 to 1.0 representing how well the candidate matches the job requirements
2. feedback: Detailed feedback (3-5 sentences) covering:
   - Technical skills match
   - Experience level
   - Relevant achievements
   - Areas for improvement

Return ONLY a valid JSON object with this exact structure:
{
  "match_rate": 0.85,
  "feedback": "Your detailed feedback here..."
}`, jobContext, cvText)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT3_5Turbo,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a technical recruiter. Always respond with valid JSON only, no markdown or extra text."),
			openai.UserMessage(prompt),
		},
		Temperature: openai.Float(0.3),
		MaxTokens:   openai.Int(1000),
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content

	var result CVEvaluationResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w, content: %s", err, content)
	}

	// Validate match_rate range
	if result.MatchRate < 0 {
		result.MatchRate = 0
	}
	if result.MatchRate > 1 {
		result.MatchRate = 1
	}

	return &result, nil
}

// EvaluateProject evaluates a candidate's project report
func (c *OpenAIClient) EvaluateProject(reportText, caseStudyContext string) (*ProjectEvaluationResult, error) {
	prompt := fmt.Sprintf(`You are an expert technical evaluator assessing a candidate's project submission.

Case Study Requirements and Rubric:
%s

Candidate's Project Report:
%s

Evaluate this project and provide:
1. score: A score from 1.0 to 5.0 based on the rubric
2. feedback: Detailed feedback (4-6 sentences) covering:
   - Correctness (prompt design, LLM chaining, RAG)
   - Code quality and structure
   - Resilience and error handling
   - Documentation quality

Return ONLY a valid JSON object with this exact structure:
{
  "score": 4.5,
  "feedback": "Your detailed feedback here..."
}`, caseStudyContext, reportText)

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT3_5Turbo,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a technical evaluator. Always respond with valid JSON only, no markdown or extra text."),
			openai.UserMessage(prompt),
		},
		Temperature: openai.Float(0.3),
		MaxTokens:   openai.Int(1200),
	})

	if err != nil {
		return nil, fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content

	var result ProjectEvaluationResult
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("failed to parse OpenAI response: %w, content: %s", err, content)
	}

	// Validate score range
	if result.Score < 1 {
		result.Score = 1
	}
	if result.Score > 5 {
		result.Score = 5
	}

	return &result, nil
}

// GenerateSummary creates a final summary from CV and project evaluations
func (c *OpenAIClient) GenerateSummary(cvFeedback, projectFeedback string, cvMatchRate, projectScore float64) (string, error) {
	prompt := fmt.Sprintf(`You are an expert hiring manager making a final recommendation.

			CV Evaluation:
			- Match Rate: %.2f
			- Feedback: %s

			Project Evaluation:
			- Score: %.1f/5.0
			- Feedback: %s

			Provide a concise overall summary (3-5 sentences) that:
			1. Highlights the candidate's strengths
			2. Notes any gaps or concerns
			3. Makes a clear hiring recommendation (e.g., "Recommended for interview", "Strong candidate", "Needs more experience")

			Return ONLY the summary text, no JSON.`, cvMatchRate, cvFeedback, projectScore, projectFeedback)

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	resp, err := c.client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: openai.ChatModelGPT3_5Turbo,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage("You are a hiring manager providing concise, actionable recommendations."),
			openai.UserMessage(prompt),
		},
		Temperature: openai.Float(0.4),
		MaxTokens:   openai.Int(500),
	})

	if err != nil {
		return "", fmt.Errorf("OpenAI API call failed: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateEmbeddings generates embeddings for a batch of texts
func (c *OpenAIClient) GenerateEmbeddings(ctx context.Context, texts []string) ([][]float64, error) {
	embeddings := make([][]float64, len(texts))

	// Process one text at a time to avoid complex union type issues
	for i, text := range texts {
		resp, err := c.client.Embeddings.New(ctx, openai.EmbeddingNewParams{
			Model: openai.EmbeddingModelTextEmbeddingAda002,
			Input: openai.EmbeddingNewParamsInputUnion{
				OfString: openai.String(text),
			},
		})

		if err != nil {
			return nil, fmt.Errorf("failed to generate embedding for text %d: %w", i, err)
		}

		if len(resp.Data) == 0 {
			return nil, fmt.Errorf("no embedding data returned for text %d", i)
		}

		embeddings[i] = resp.Data[0].Embedding

		// Add small delay to avoid rate limits
		time.Sleep(200 * time.Millisecond)
	}

	return embeddings, nil
}
