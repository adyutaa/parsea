CREATE TABLE documents (
    id UUID PRIMARY KEY,
    filename VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    doc_type VARCHAR(50) NOT NULL,
    file_size BIGINT,
    uploaded_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE evaluation_jobs (
    id UUID PRIMARY KEY,
    cv_id UUID REFERENCES documents(id),
    report_id UUID REFERENCES documents(id),
    job_title VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'queued',
    result JSONB,
    error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_documents_type ON documents(doc_type);
CREATE INDEX idx_jobs_status ON evaluation_jobs(status);
CREATE INDEX idx_jobs_created ON evaluation_jobs(created_at);