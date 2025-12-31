# Parsea - AI-Powered Recruitment Agent

<div align="center">

![Parsea Banner](public/results.png)

[![Go](https://img.shields.io/badge/Go-1.24-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
[![Next.js](https://img.shields.io/badge/Next.js-15-black?style=for-the-badge&logo=next.dot.js&logoColor=white)](https://nextjs.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)](https://www.postgresql.org/)
[![Redis](https://img.shields.io/badge/Redis-Stack-DC382D?style=for-the-badge&logo=redis&logoColor=white)](https://redis.io/)
[![OpenAI](https://img.shields.io/badge/OpenAI-GPT-412991?style=for-the-badge&logo=openai&logoColor=white)](https://openai.com/)
[![Docker](https://img.shields.io/badge/Docker-Compose-2496ED?style=for-the-badge&logo=docker&logoColor=white)](https://www.docker.com/)

**An intelligent recruitment automation platform that evaluates candidate CVs against project requirements using RAG (Retrieval-Augmented Generation) and LLMs.**

[Features](#-features) • [Architecture](#-architecture) • [Gallery](#-gallery) • [Getting Started](#-getting-started)

</div>

---

## 🎯 Project Overview

**Parsea** solves the bottleneck of manual technical screening. Instead of spending hours reading resumes and comparing them against inconsistent criteria, Parsea provides instant, data-driven evaluations.

The system ingests candidate **CVs (PDF)** and **Project Reports (PDF)**, analyzes them using a robust pipeline of AI agents, and produces a structured scorecard comparing the candidate's actual skills against the specific job context.

### Why Parsea?

- **🔍 Context-Aware Analysis**: Uses Vector DB (Qdrant) to understand job nuances, not just keyword matching.
- **⚡ Asynchronous Scale**: Handles heavy PDF processing and LLM inference background jobs via Redis queues.
- **📊 Structured Scoring**: Converts subjective text into objective metrics (Match Rate, Technical Score).
- **🎨 Premium UX**: sleek, modern dashboard for recruiters to view insights.

---

## 📸 Gallery

<div align="center">
  <table>
    <tr>
      <td align="center">
        <b>1. Upload Documents</b><br/>
        <img src="public/upload.png" width="400" alt="Upload Page" />
      </td>
      <td align="center">
        <b>2. Configure Evaluation</b><br/>
        <img src="public/evaluate.png" width="400" alt="Evaluate Page" />
      </td>
    </tr>
    <tr>
      <td align="center">
        <b>3. Real-time Analysis</b><br/>
        <img src="public/process.png" width="400" alt="Processing Status" />
      </td>
      <td align="center">
        <b>4. Detailed Results</b><br/>
        <img src="public/results.png" width="400" alt="Results Dashboard" />
      </td>
    </tr>
  </table>
</div>

---

## 🏗 Architecture

The system follows a **Clean Architecture** pattern to ensure scalability and maintainability.

```mermaid
graph TD
    Client[Next.js Frontend] -->|REST API| Router[Gin Router]
    
    subgraph "Backend Service (Go)"
        Router --> Handler
        Handler --> Service
        Service -->|Enqueue Job| Redis[(Redis Queue)]
        
        Worker[Background Worker] -->|Pop Job| Redis
        Worker -->|Store Result| DB[(PostgreSQL)]
        
        subgraph "AI Pipeline"
            Worker -->|Text Extraction| PDF[PDF Parser]
            Worker -->|Context Retrieval| VectorDB[(Qdrant Cloud)]
            Worker -->|Inference| LLM[OpenAI GPT-3.5]
        end
    end
```

### Technical Highlights

*   **Backend**: Built with **Go 1.24** and **Gin**, ensuring high performance and low latency.
*   **Database**: **PostgreSQL** (via Supabase) for transactional data, **Redis** for job queues.
*   **AI/LLM**: Integrates **OpenAI** for reasoning and **Qdrant** for semantic search (RAG).
*   **Frontend**: Built with **Next.js 15**, **Tailwind CSS**, and **Framer Motion** for a responsive, interactive UI.

---

## ✨ Features

### 🤖 Intelligent "Agentic" Workflow
The backend doesn't just "call an API". It orchestrates a multi-step agentic workflow:
1.  **Ingestion**: Extracts text from unstructured PDFs.
2.  **Context Loading (RAG)**: Fetches relevant job descriptions and scoring rubrics from the vector store.
3.  **Analysis**: Evaluates the CV against the job description.
4.  **Technical Review**: Deep-dives into the project report for code quality and architecture.
5.  **Synthesis**: Combines all data into a final weighted score and summary.

### ⚡ Performance-First
*   **Concurrent Processing**: Multiple workers process evaluations in parallel.
*   **Optimized Queries**: GORM with prepared statements and connection pooling.
*   **Real-time Updates**: Frontend polls status for immediate feedback.

---

## 🚀 Getting Started

### Prerequisites

- Go 1.24+
- Node.js 18+
- Docker (optional, for Redis/Postgres)
- OpenAI API Key

### Installation

1.  **Clone the repository**
    ```bash
    git clone https://github.com/adyutaa/parsea.git
    cd parsea
    ```

2.  **Start Backend**
    ```bash
    cd backend
    cp .env.example .env  # Configure your keys
    go run cmd/server/main.go
    ```

3.  **Start Frontend**
    ```bash
    cd frontend
    npm install
    npm run dev
    ```

4.  **Access App**
    Open `http://localhost:3000` to start evaluating.

---

## 📄 License

This project is open-sourced under the MIT License.

---

<div align="center">
  <sub>Built with ❤️ by Adyuta Indra Adyatma</sub>
</div>
