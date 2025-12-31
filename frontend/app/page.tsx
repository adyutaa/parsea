'use client';

import React, { useState, useRef } from 'react';
import { useRouter } from 'next/navigation';
import { api } from '../lib/api';

export default function Home() {
  const router = useRouter();
  const [cvFile, setCvFile] = useState<File | null>(null);
  const [reportFile, setReportFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const cvInputRef = useRef<HTMLInputElement>(null);
  const reportInputRef = useRef<HTMLInputElement>(null);

  const handleFileChange = (
    e: React.ChangeEvent<HTMLInputElement>,
    setFile: React.Dispatch<React.SetStateAction<File | null>>
  ) => {
    if (e.target.files && e.target.files[0]) {
      setFile(e.target.files[0]);
      setError(null);
    }
  };

  const handleUpload = async () => {
    if (!cvFile || !reportFile) {
      setError('Please select both your CV and Project Report.');
      return;
    }

    setIsUploading(true);
    setError(null);

    try {
      const response = await api.uploadDocuments(cvFile, reportFile);
      // Redirect to evaluate page with document IDs
      router.push(`/evaluate?cv_id=${response.cv_id}&report_id=${response.report_id}`);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Upload failed');
      setIsUploading(false);
    }
  };

  return (
    <div className="relative flex min-h-screen w-full flex-col overflow-x-hidden">
      {/* Top Navigation */}
      <header className="sticky top-0 z-50 w-full border-b border-gray-200 dark:border-gray-800 bg-surface-light/80 dark:bg-background-dark/80 backdrop-blur-md">
        <div className="flex items-center justify-between px-6 py-4 lg:px-10 max-w-7xl mx-auto">
          <div className="flex items-center gap-4">
            <div className="flex items-center justify-center text-primary bg-primary/10 rounded-lg p-2">
              <span className="material-symbols-outlined text-2xl">analytics</span>
            </div>
            <h2 className="text-xl font-extrabold leading-tight tracking-tight text-gray-900 dark:text-white">
              AI Evaluator
            </h2>
          </div>
          {/* ... existing header content ... */}
          <div className="hidden md:flex items-center gap-8">
            <nav className="flex items-center gap-6">
              <a
                className="text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-primary dark:hover:text-primary transition-colors"
                href="#"
              >
                Dashboard
              </a>
              <a className="text-sm font-medium text-primary dark:text-primary" href="#">
                Evaluations
              </a>
              <a
                className="text-sm font-medium text-gray-600 dark:text-gray-400 hover:text-primary dark:hover:text-primary transition-colors"
                href="#"
              >
                Settings
              </a>
            </nav>
            <div className="h-6 w-px bg-gray-200 dark:bg-gray-700"></div>
            <div className="flex items-center gap-3">
              <div className="relative">
                <div
                  className="bg-center bg-no-repeat bg-cover rounded-full size-10 ring-2 ring-gray-100 dark:ring-gray-800 cursor-pointer"
                  style={{
                    backgroundImage:
                      'url("https://lh3.googleusercontent.com/aida-public/AB6AXuAIVo84K8JNjPo6dE_KqjIMoB4cXiOvAzjbeDOtVHIlm9lsfC_CXjKR3YgEl-IjxWBpe7kXh6-rHZsNO7-qHXYVTmEejWgE0-48GegwFTNqydKQE2WX4qmhiInH235aXF-q3wYF7MzwYckIchF5LqGfgO8impiCsiStZO-t6bPIDmMUB0o-xwUDpBk3UGiIQpR83bS153wfi4iNN8mbFtmr-m5pOj9vFYRdL2kmY4u6_SEKT9Dpnc6TjcLi_GaXgIwt1mRWyyWKT88z")',
                  }}
                ></div>
                <div className="absolute bottom-0 right-0 size-3 rounded-full bg-green-500 border-2 border-white dark:border-background-dark"></div>
              </div>
            </div>
          </div>
          <button className="md:hidden p-2 text-gray-600 dark:text-gray-400 hover:bg-gray-100 dark:hover:bg-gray-800 rounded-lg">
            <span className="material-symbols-outlined">menu</span>
          </button>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-grow flex flex-col items-center justify-center py-12 px-4 sm:px-6 lg:px-8">
        <div className="w-full max-w-4xl space-y-8">
          {/* Page Heading */}
          <div className="text-center space-y-4 max-w-2xl mx-auto">
            <h1 className="text-3xl sm:text-4xl font-black tracking-tight text-gray-900 dark:text-white">
              Document Evaluation Setup
            </h1>
            <p className="text-lg text-slate-500 dark:text-slate-400 font-normal leading-relaxed">
              Please upload your Curriculum Vitae and Project Report in PDF format. Our AI will
              analyze both to generate your comprehensive evaluation.
            </p>
          </div>

          {/* Upload Section */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mt-8">
            {/* Zone A: CV / Resume */}
            <div className="group relative flex flex-col rounded-2xl bg-surface-light dark:bg-surface-dark shadow-sm border border-gray-200 dark:border-gray-800 overflow-hidden transition-all hover:shadow-md">
              <input
                type="file"
                accept=".pdf"
                className="hidden"
                ref={cvInputRef}
                onChange={(e) => handleFileChange(e, setCvFile)}
              />
              <div
                className="absolute inset-0 bg-primary/0 group-hover:bg-primary/5 transition-colors cursor-pointer"
                onClick={() => cvInputRef.current?.click()}
              ></div>
              <div className="flex-1 flex flex-col items-center justify-center p-8 text-center border-2 border-dashed border-gray-200 dark:border-gray-700 m-2 rounded-xl group-hover:border-primary/50 transition-colors pointer-events-none">
                <div className="size-16 rounded-full bg-blue-50 dark:bg-blue-900/20 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                  <span className="material-symbols-outlined text-3xl text-primary">
                    {cvFile ? 'check_circle' : 'description'}
                  </span>
                </div>
                <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">
                  {cvFile ? cvFile.name : 'CV / Resume'}
                </h3>
                <p className="text-sm text-slate-500 dark:text-slate-400 mb-8 max-w-[240px]">
                  {cvFile
                    ? `${(cvFile.size / 1024 / 1024).toFixed(2)} MB`
                    : 'Drag & drop your CV here, or click to browse.'}
                </p>
                <button className="relative overflow-hidden rounded-lg bg-primary text-white px-6 py-2.5 text-sm font-bold shadow-lg shadow-primary/25 hover:bg-blue-600 transition-all focus:ring-4 focus:ring-primary/20 active:scale-95">
                  <span className="relative z-10 flex items-center gap-2">
                    <span className="material-symbols-outlined text-lg">upload_file</span>
                    {cvFile ? 'Change File' : 'Select PDF File'}
                  </span>
                </button>
              </div>
            </div>

            {/* Zone B: Project Report */}
            <div className="group relative flex flex-col rounded-2xl bg-surface-light dark:bg-surface-dark shadow-sm border border-gray-200 dark:border-gray-800 overflow-hidden transition-all hover:shadow-md">
              <input
                type="file"
                accept=".pdf"
                className="hidden"
                ref={reportInputRef}
                onChange={(e) => handleFileChange(e, setReportFile)}
              />
              <div
                className="absolute inset-0 bg-primary/0 group-hover:bg-primary/5 transition-colors cursor-pointer"
                onClick={() => reportInputRef.current?.click()}
              ></div>
              <div className="flex-1 flex flex-col items-center justify-center p-8 text-center border-2 border-dashed border-gray-200 dark:border-gray-700 m-2 rounded-xl group-hover:border-primary/50 transition-colors pointer-events-none">
                <div className="size-16 rounded-full bg-indigo-50 dark:bg-indigo-900/20 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300">
                  <span className="material-symbols-outlined text-3xl text-indigo-600 dark:text-indigo-400">
                    {reportFile ? 'check_circle' : 'folder_open'}
                  </span>
                </div>
                <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-2">
                  {reportFile ? reportFile.name : 'Project Report'}
                </h3>
                <p className="text-sm text-slate-500 dark:text-slate-400 mb-8 max-w-[240px]">
                  {reportFile
                    ? `${(reportFile.size / 1024 / 1024).toFixed(2)} MB`
                    : 'Drag & drop your Project Report here, or click to browse.'}
                </p>
                <button className="relative overflow-hidden rounded-lg bg-surface-light dark:bg-gray-800 text-gray-900 dark:text-white border border-gray-200 dark:border-gray-600 px-6 py-2.5 text-sm font-bold shadow-sm hover:bg-gray-50 dark:hover:bg-gray-700 transition-all focus:ring-4 focus:ring-gray-200 dark:focus:ring-gray-700 active:scale-95">
                  <span className="relative z-10 flex items-center gap-2">
                    <span className="material-symbols-outlined text-lg">upload_file</span>
                    {reportFile ? 'Change File' : 'Select PDF File'}
                  </span>
                </button>
              </div>
            </div>
          </div>

          {/* Meta Info & Validation */}
          <div className="flex flex-col items-center gap-2 text-center">
            {error && <div className="text-red-500 text-sm font-bold mb-2">{error}</div>}
            <div className="flex items-center gap-2 text-xs font-medium text-slate-500 dark:text-slate-400 bg-gray-100 dark:bg-gray-800/50 px-3 py-1.5 rounded-full">
              <span className="material-symbols-outlined text-base">info</span>
              <span>Max file size: 10MB per file</span>
              <span className="w-1 h-1 bg-slate-400 rounded-full mx-1"></span>
              <span>Format: PDF only</span>
            </div>
          </div>

          {/* Action Bar */}
          <div className="pt-6 border-t border-gray-200 dark:border-gray-800 flex flex-col sm:flex-row items-center justify-center sm:justify-end gap-4">
            <button
              className="w-full sm:w-auto px-6 py-3 text-sm font-bold text-slate-600 dark:text-slate-400 hover:text-gray-900 dark:hover:text-white transition-colors cursor-pointer"
              onClick={() => {
                setCvFile(null);
                setReportFile(null);
                setError(null);
              }}
            >
              Cancel
            </button>
            <button
              className={`w-full sm:w-auto flex items-center justify-center gap-2 rounded-lg px-8 py-3 text-sm font-bold transition-all ${cvFile && reportFile && !isUploading
                  ? 'bg-primary text-white hover:bg-blue-600 shadow-lg shadow-primary/25 cursor-pointer'
                  : 'bg-gray-200 dark:bg-gray-700 text-gray-400 dark:text-gray-500 cursor-not-allowed'
                }`}
              disabled={!cvFile || !reportFile || isUploading}
              onClick={handleUpload}
            >
              {isUploading ? (
                <>
                  <span className="material-symbols-outlined text-lg animate-spin">refresh</span>
                  Uploading...
                </>
              ) : (
                <>
                  <span className="material-symbols-outlined text-lg">auto_awesome</span>
                  Proceed to Setup
                </>
              )}
            </button>
          </div>
        </div>
      </main>

      {/* Decorative Background Elements */}
      <div className="fixed top-0 left-0 w-full h-full overflow-hidden -z-10 pointer-events-none">
        <div className="absolute top-[-10%] right-[-5%] w-[500px] h-[500px] bg-primary/5 rounded-full blur-3xl opacity-50 dark:opacity-20"></div>
        <div className="absolute bottom-[-10%] left-[-10%] w-[600px] h-[600px] bg-indigo-500/5 rounded-full blur-3xl opacity-50 dark:opacity-20"></div>
      </div>
    </div>
  );
}
