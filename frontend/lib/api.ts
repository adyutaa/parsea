export const API_URL = 'http://localhost:8080';

export interface UploadResponse {
    cv_id: number;
    report_id: number;
    message: string;
}

export interface StartEvaluationResponse {
    id: number;
    status: string;
    message?: string;
}

export interface EvaluationResult {
    scan_match: number;
    scan_feedback: string;
    project_score: number;
    project_feedback: string;
    summary: string;
}

export interface EvaluationJob {
    id: number;
    job_title?: string;
    status: 'queued' | 'processing' | 'completed' | 'failed';
    result?: EvaluationResult;
    error?: string;
    created_at?: string;
}

export const api = {
    /**
     * Uploads the CV and Project Report files.
     */
    async uploadDocuments(cv: File, report: File): Promise<UploadResponse> {
        const formData = new FormData();
        formData.append('cv', cv);
        formData.append('project_report', report);

        const response = await fetch(`${API_URL}/upload`, {
            method: 'POST',
            body: formData,
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to upload documents');
        }

        return response.json();
    },

    /**
     * Starts the evaluation process for the uploaded documents.
     */
    async startEvaluation(cvId: number, reportId: number, jobTitle: string): Promise<StartEvaluationResponse> {
        const response = await fetch(`${API_URL}/evaluate`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                cv_id: cvId,
                report_id: reportId,
                job_title: jobTitle,
            }),
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to start evaluation');
        }

        return response.json();
    },

    /**
     * Gets the status and result of an evaluation job.
     */
    async getResult(jobId: number): Promise<EvaluationJob> {
        const response = await fetch(`${API_URL}/result?id=${jobId}`, {
            method: 'GET',
        });

        if (!response.ok) {
            const errorText = await response.text();
            throw new Error(errorText || 'Failed to fetch results');
        }

        return response.json();
    },
};
