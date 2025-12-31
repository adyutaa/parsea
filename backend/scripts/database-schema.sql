

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