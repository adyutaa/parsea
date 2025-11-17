# Parsea - AI-Powered CV & Project Evaluation System

<div align="center">

![Go](https://img.shields.io/badge/Go-1.24-blue?style=flat-square&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue?style=flat-square&logo=postgresql)
![Redis](https://img.shields.io/badge/Redis-Cloud-red?style=flat-square&logo=redis)
![OpenAI](https://img.shields.io/badge/OpenAI-GPT--3.5-green?style=flat-square&logo=openai)
![Qdrant](https://img.shields.io/badge/Qdrant-Vector%20DB-purple?style=flat-square)

_Intelligent backend service for automated CV screening and project evaluation using advanced LLM technology_

</div>

## 📋 Table of Contents

- [Overview](#overview)
- [Features](#features)
- [Tech Stack](#tech-stack)
- [Architecture](#architecture)
- [API Documentation](#api-documentation)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Development](#development)

## 🎯 Overview

Parsea is a backend evaluation system that automates the initial screening of job applications. The service receives candidate CVs and project reports, evaluates them against specific job requirements using AI, and produces structured evaluation reports.

### Key Capabilities

- **Intelligent CV Analysis**: Evaluates technical skills, experience, and cultural fit
- **Project Report Assessment**: Analyzes code quality, implementation approach, and documentation
- **RAG-Enhanced Context**: Uses vector database for accurate job requirement matching
- **Asynchronous Processing**: Handles long-running AI evaluations with Redis queue
- **Structured Scoring**: Provides weighted scores based on predefined rubrics

### Design Philosophy

**Why GPT-3.5-turbo over GPT-4?**
While GPT-4 offers superior capabilities, GPT-3.5-turbo was chosen for practical reasons: it costs approximately 10x less and responds 3x faster. For structured evaluation tasks with clear prompts, the quality difference is minimal, but the cost difference becomes significant when processing hundreds of evaluations.

**Why Qdrant Cloud?**
Hosted vector database chosen to avoid infrastructure complexity while providing reliable semantic search for job requirement matching. The 1536-dimensional embeddings from OpenAI ensure compatibility and semantic accuracy.

**Why Redis Queue?**
AI evaluations take 30-60 seconds per document pair, requiring asynchronous processing. Redis provides reliable job persistence and enables horizontal worker scaling for concurrent evaluation processing.

## ✨ Features

### Core Functionality

- 📄 **PDF Document Processing** - Extract and analyze CV and project reports
- 🤖 **LLM-Powered Evaluation** - GPT-3.5 Turbo for intelligent assessment
- 📊 **Structured Scoring System** - Weighted evaluation based on job requirements
- 🔍 **Vector Search (RAG)** - Context-aware evaluation using Qdrant
- ⚡ **Async Job Processing** - Redis-backed queue for scalable processing

### Technical Features

- 🔒 **Input Validation** - Comprehensive file and data validation
- 🚦 **Error Handling** - Robust error management and recovery
- 📈 **Monitoring Ready** - Health checks and queue status endpoints
- 🐳 **Docker Support** - Containerized deployment ready
- 🔧 **Clean Architecture** - Separation of concerns with repository pattern

## 🛠 Tech Stack

### Backend

- **Language**: Go 1.24+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL (Supabase)
- **Cache/Queue**: Redis Cloud
- **ORM**: GORM

### AI & ML

- **LLM Provider**: OpenAI GPT-3.5 Turbo *(chosen for 10x lower cost and 3x faster response vs GPT-4)*
- **Vector Database**: Qdrant Cloud *(hosted solution to avoid infrastructure complexity)*
- **Embeddings**: OpenAI text-embedding-ada-002 *(1536-dimensional for semantic matching)*
- **PDF Processing**: Custom Go PDF parser *(lightweight text extraction)*

### Infrastructure

- **Deployment**: Docker + Docker Compose
- **Environment**: Supabase (PostgreSQL), Redis Cloud
- **File Storage**: Local filesystem with configurable upload paths

## 🏗 Architecture

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │    │   File Upload   │    │   Evaluation    │
│                 │    │                 │    │     Queue       │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          ▼                      ▼                      ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Gin Router    │───▶│   Handlers      │───▶│   Services      │
│                 │    │                 │    │                 │
└─────────────────┘    └─────────────────┘    └─────────┬───────┘
                                                        │
                                                        ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Worker Pool   │───▶│  Repositories   │───▶│   PostgreSQL    │
│                 │    │                 │    │   (Supabase)    │
└─────────┬───────┘    └─────────────────┘    └─────────────────┘
          │
          ▼
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   OpenAI API    │    │   Qdrant Cloud  │    │   Redis Cloud   │
│   (LLM + RAG)   │    │  (Vector Store) │    │     (Queue)     │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Data Flow

1. **Document Upload**: PDF files uploaded via multipart/form-data
2. **Job Creation**: Evaluation job created and queued in Redis
3. **Worker Processing**: Background worker processes job asynchronously
4. **AI Evaluation**: 3-step LLM chain (CV → Project → Summary)
5. **Result Storage**: Structured results saved to PostgreSQL

## 📚 API Documentation

### Base URL

```
http://localhost:8080
```

### Endpoints

#### 📤 Upload Documents

```http
POST /upload
Content-Type: multipart/form-data

Form Data:
- cv: PDF file (max 10MB)
- project_report: PDF file (max 10MB)
```

**Response:**

```json
{
  "cv_id": 1,
  "report_id": 2,
  "message": "Files uploaded successfully"
}
```

#### 🚀 Start Evaluation

```http
POST /evaluate
Content-Type: application/json

{
    "cv_id": 1,
    "report_id": 2,
    "job_title": "Backend Developer"
}
```

**Response:**

```json
{
  "id": 456,
  "status": "queued"
}
```

#### 📊 Get Results

```http
GET /result/456
```

**Response:**

```json
{
  "id": 456,
  "status": "completed",
  "result": {
    "cv_match_rate": 0.82,
    "cv_feedback": "Strong backend experience with good cloud knowledge",
    "project_score": 4.2,
    "project_feedback": "Well-structured code with good error handling",
    "overall_summary": "Good candidate fit, would benefit from deeper RAG knowledge"
  },
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:35:00Z"
}
```

#### ⚡ Health Check

```http
GET /health
```

#### 📈 Queue Status

```http
GET /queue/status
```

## 🚀 Installation

### Prerequisites

- **Go**: 1.24 or higher
- **PostgreSQL**: 16+ (Supabase account)
- **Redis**: Cloud instance or local Redis server
- **OpenAI API Key**: For LLM services
- **Qdrant Cloud**: For vector search (optional but recommended)

### Local Development Setup

1. **Clone the repository**

```bash
git clone https://github.com/adyutaa/parsea.git
cd parsea
```

2. **Install Go dependencies**

```bash
go mod download
```

3. **Create environment file**

```bash
cp .env.example .env
```

4. **Setup database**

```sql
-- Run this SQL in your Supabase SQL editor
-- Drop existing tables and recreate with auto-incrementing integers
DROP TABLE IF EXISTS public.evaluation_jobs CASCADE;
DROP TABLE IF EXISTS public.documents CASCADE;

CREATE TABLE public.documents (
  id SERIAL PRIMARY KEY,
  filename character varying NOT NULL,
  file_path text NOT NULL,
  doc_type character varying NOT NULL,
  file_size bigint,
  uploaded_at timestamp without time zone DEFAULT now()
);

CREATE TABLE public.evaluation_jobs (
  id SERIAL PRIMARY KEY,
  cv_id INTEGER NOT NULL,
  report_id INTEGER NOT NULL,
  job_title character varying NOT NULL,
  status character varying DEFAULT 'queued'::character varying,
  result jsonb,
  error_message text,
  created_at timestamp without time zone DEFAULT now(),
  updated_at timestamp without time zone DEFAULT now(),
  CONSTRAINT evaluation_jobs_cv_id_fkey FOREIGN KEY (cv_id) REFERENCES public.documents(id),
  CONSTRAINT evaluation_jobs_report_id_fkey FOREIGN KEY (report_id) REFERENCES public.documents(id)
);

CREATE INDEX idx_documents_type ON public.documents(doc_type);
CREATE INDEX idx_jobs_status ON public.evaluation_jobs(status);
CREATE INDEX idx_jobs_created ON public.evaluation_jobs(created_at);
```

5. **Seed vector database (optional)**

```bash
go run scripts/ingest.go
```

6. **Build and run**

```bash
# Development
go run cmd/server/main.go

# Production
go build -o bin/server cmd/server/main.go
./bin/server
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file in the project root:

```env
# Database Configuration (Supabase)
DATABASE_URL=postgresql://user:password@host:port/database

# Redis Configuration (Redis Cloud)
REDIS_HOST=your-redis-host
REDIS_PORT=6379
REDIS_USERNAME=default
REDIS_PASSWORD=your-redis-password

# OpenAI Configuration
OPENAI_API_KEY=sk-your-openai-api-key

# Qdrant Configuration (Optional)
QDRANT_HOST=your-qdrant-host
QDRANT_PORT=6333
QDRANT_API_KEY=your-qdrant-api-key

# Application Configuration
PORT=8080
UPLOAD_PATH=./uploads
```

### Docker Deployment

```bash
# Using Docker Compose
docker-compose up -d

# Manual Docker build
docker build -t parsea .
docker run -p 8080:8080 --env-file .env parsea
```

## 💻 Usage

### Quick Start Example

1. **Upload documents**

```bash
curl -X POST http://localhost:8080/upload \
  -F "cv=@candidate-cv.pdf" \
  -F "project_report=@project-report.pdf"
```

2. **Start evaluation**

```bash
curl -X POST http://localhost:8080/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "cv_id": 1,
    "report_id": 2,
    "job_title": "Backend Developer"
  }'
```

3. **Check results**

```bash
curl "http://localhost:8080/result/456"
```

### Evaluation Process

The system follows a **7-step evaluation pipeline** (approximately 45 seconds total):

1. **PDF Text Extraction** (5s) - Extract readable text from uploaded documents
2. **Job Requirements Retrieval** (3s) - RAG-based context gathering from vector database
3. **CV Evaluation** (15s) - LLM assessment with job-specific scoring rubric
4. **Project Text Extraction** (5s) - Extract content from project report
5. **Case Study Context** (3s) - Retrieve evaluation criteria from vector store
6. **Project Evaluation** (15s) - LLM assessment with technical rubric
7. **Final Summary** (5s) - Combine evaluations into comprehensive assessment

### AI Evaluation Strategy

**Consistency Through Structure**
- Temperature set to 0.3 to reduce response randomness
- Prompts designed to always return valid JSON structures
- Role-based prompting establishes evaluation perspective
- Context injection provides job-specific requirements

**RAG Implementation**
Vector search addresses evaluation consistency by providing relevant job requirements as context to the AI. Without proper context, a CV mentioning "Docker" might receive different scores depending on interpretation. The system stores job descriptions, case study briefs, and scoring rubrics as 1536-dimensional embeddings, using cosine similarity to find semantically related content.

**Fallback Mechanism**
When RAG fails, the system automatically switches to hardcoded job requirements and evaluation rubrics to ensure evaluations can proceed. This design ensures 100% system availability even when the vector database is unavailable.

## 🧪 Development

### Project Structure

```
parsea/
├── cmd/
│   └── server/           # Application entrypoint
├── internal/
│   ├── domain/          # Business entities and models
│   ├── handler/         # HTTP handlers (controllers)
│   ├── infrastructure/  # External services (DB, APIs)
│   ├── repository/      # Data access layer
│   ├── service/         # Business logic layer
│   ├── validation/      # Input validation
│   └── worker/          # Background job processors
├── pkg/
│   └── pdf/             # PDF processing utilities
├── scripts/
│   ├── database-schema.sql  # Database migrations
│   └── ingest.go           # Vector DB seeding
├── docs/                # Documentation
└── uploads/             # File upload directory
```

### Code Standards

- **Clean Architecture**: Clear separation between layers
- **Error Handling**: Comprehensive error management with context
- **Validation**: Input sanitization and security checks
- **Testing**: Unit tests for critical business logic
- **Documentation**: Inline comments and API documentation

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/service/...
```

### Building

```bash
# Development build
go build -o bin/parsea cmd/server/main.go

# Production build with optimizations
go build -ldflags="-s -w" -o bin/parsea cmd/server/main.go

# Cross-platform build
GOOS=linux GOARCH=amd64 go build -o bin/parsea-linux cmd/server/main.go
```

## 🙏 Acknowledgments

- [OpenAI](https://openai.com/) for providing powerful LLM capabilities
- [Qdrant](https://qdrant.tech/) for vector search infrastructure
- [Supabase](https://supabase.com/) for managed PostgreSQL services
- [Redis](https://redis.io/) for reliable queueing and caching

---

<div align="center">
Built with ❤️ by <a href="https://github.com/adyutaa">@adyutaa</a>
</div>
