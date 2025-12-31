'use client';

import React, { useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { api } from '../../lib/api';

function EvaluateContent() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const cvId = searchParams.get('cv_id');
    const reportId = searchParams.get('report_id');

    const [jobTitle, setJobTitle] = useState('Senior Product Manager');
    const [isStarting, setIsStarting] = useState(false);
    const [error, setError] = useState<string | null>(null);
    const [showSuggestions, setShowSuggestions] = useState(false);

    const suggestions = [
        { title: 'Senior Product Manager', dept: 'Product Team' },
        { title: 'Senior Project Lead', dept: 'Operations' },
        { title: 'Senior Programmer Analyst', dept: 'Engineering' },
    ];

    const handleStartEvaluation = async () => {
        if (!cvId || !reportId) {
            setError('Missing document IDs. Please upload files first.');
            return;
        }

        if (!jobTitle.trim()) {
            setError('Please enter a target job role.');
            return;
        }

        setIsStarting(true);
        setError(null);

        try {
            const response = await api.startEvaluation(
                parseInt(cvId),
                parseInt(reportId),
                jobTitle
            );
            router.push(`/status?id=${response.id}`);
        } catch (err) {
            setError(err instanceof Error ? err.message : 'Failed to start evaluation');
            setIsStarting(false);
        }
    };

    return (
        <div className="bg-background-light dark:bg-background-dark font-display text-slate-900 dark:text-slate-100 min-h-screen flex flex-col transition-colors duration-200">
            {/* Top Navigation Bar */}
            <header className="sticky top-0 z-50 flex items-center justify-between border-b border-solid border-slate-200 dark:border-slate-800 bg-surface-light dark:bg-surface-dark px-6 py-3 lg:px-10">
                <div className="flex items-center gap-4">
                    <div className="flex items-center justify-center size-10 rounded-xl bg-primary/10 text-primary">
                        <span className="material-symbols-outlined text-3xl">analytics</span>
                    </div>
                    <h2 className="text-xl font-bold leading-tight tracking-tight">AI Evaluator</h2>
                </div>
                <div className="hidden md:flex flex-1 justify-center gap-8">
                    <nav className="flex items-center gap-6">
                        <a
                            className="text-slate-600 dark:text-slate-400 hover:text-primary dark:hover:text-primary text-sm font-medium transition-colors"
                            href="#"
                        >
                            Dashboard
                        </a>
                        <a className="text-primary font-bold text-sm leading-normal" href="#">
                            Evaluations
                        </a>
                        <a
                            className="text-slate-600 dark:text-slate-400 hover:text-primary dark:hover:text-primary text-sm font-medium transition-colors"
                            href="#"
                        >
                            Settings
                        </a>
                    </nav>
                </div>
                <div className="flex items-center gap-4">
                    <button className="p-2 rounded-full hover:bg-slate-100 dark:hover:bg-slate-800 text-slate-600 dark:text-slate-400 cursor-pointer">
                        <span className="material-symbols-outlined">notifications</span>
                    </button>
                    <div
                        className="bg-center bg-no-repeat bg-cover rounded-full size-10 ring-2 ring-slate-100 dark:ring-slate-700 cursor-pointer"
                        style={{
                            backgroundImage:
                                'url("https://lh3.googleusercontent.com/aida-public/AB6AXuAVYyol2AfaxfjUwoVyQ4EQs-OH0eWCG8RJLhfCfJWFlMNJJMKwoTsTnZ9YRq2s_gOdkoRjemhYP_4brHljKsVMboSx_6_Es23xNBaSK5OAZWoTEt5dsiE39qK1ZJu9VFk5WBsM_Oj5ycgF68Yxxft0ha8r8Z3_MSzfIkioFMbtXrlabFGUEl8uPxeOhuj9QNA_YkLJ8bP-M6E330PnEfWEXNJ_bAm_mrvklp0DfUAv81WyfZxFot6Th_d_oTLJmUzabakAO-LgO-nN")',
                        }}
                    ></div>
                </div>
            </header>

            {/* Main Content */}
            <main className="flex-grow flex flex-col items-center py-10 px-4 sm:px-6 lg:px-8">
                <div className="w-full max-w-4xl">
                    {/* Breadcrumbs */}
                    <nav aria-label="Breadcrumb" className="flex mb-6">
                        <ol className="inline-flex items-center space-x-1 md:space-x-3">
                            <li className="inline-flex items-center">
                                <a
                                    className="inline-flex items-center text-sm font-medium text-slate-500 hover:text-primary dark:text-slate-400"
                                    href="/"
                                >
                                    <span className="material-symbols-outlined text-lg mr-2">upload_file</span>
                                    Upload
                                </a>
                            </li>
                            <li>
                                <div className="flex items-center">
                                    <span className="material-symbols-outlined text-slate-400 mx-2">chevron_right</span>
                                    <span className="text-sm font-medium text-slate-900 dark:text-white">
                                        Configure Evaluation
                                    </span>
                                </div>
                            </li>
                        </ol>
                    </nav>

                    {/* Card Container */}
                    <div className="bg-surface-light dark:bg-surface-dark rounded-xl shadow-sm border border-slate-200 dark:border-slate-800 overflow-hidden">
                        {/* Page Heading Section */}
                        <div className="p-8 pb-4">
                            <div className="flex flex-col gap-2">
                                <h1 className="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
                                    Set Evaluation Context
                                </h1>
                                <p className="text-slate-500 dark:text-slate-400 text-base">
                                    Review your uploads and specify the target role for the AI analysis.
                                </p>
                            </div>
                        </div>

                        {/* Divider */}
                        <hr className="border-slate-100 dark:border-slate-800 mx-8" />

                        {/* Uploaded Documents Section */}
                        <div className="p-8 py-6">
                            <div className="flex items-center justify-between mb-4">
                                <h3 className="text-lg font-bold text-slate-900 dark:text-white">
                                    Uploaded Documents
                                </h3>
                                <span className="text-xs font-medium px-2.5 py-0.5 rounded-full bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400 flex items-center gap-1">
                                    <span className="material-symbols-outlined text-sm">check_circle</span>
                                    Ready
                                </span>
                            </div>

                            {!cvId || !reportId ? (
                                <div className="text-center p-4 text-red-500 bg-red-50 dark:bg-red-900/10 rounded-lg">
                                    ⚠️ Documents not found. Please upload them first.
                                </div>
                            ) : (
                                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                                    {/* File Item 1: Resume */}
                                    <div className="group relative flex items-center p-4 rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/50 hover:border-primary/50 transition-colors">
                                        <div className="flex-shrink-0 size-12 rounded-lg bg-white dark:bg-slate-700 flex items-center justify-center text-primary shadow-sm">
                                            <span className="material-symbols-outlined">description</span>
                                        </div>
                                        <div className="ml-4 flex-1 min-w-0">
                                            <p className="text-sm font-semibold text-slate-900 dark:text-white truncate">
                                                Uploaded CV (ID: {cvId})
                                            </p>
                                            <p className="text-xs text-slate-500 dark:text-slate-400">PDF Document</p>
                                        </div>
                                        <div className="ml-2 flex-shrink-0 text-green-500">
                                            <span className="material-symbols-outlined">check_circle</span>
                                        </div>
                                    </div>

                                    {/* File Item 2: Project Report */}
                                    <div className="group relative flex items-center p-4 rounded-xl border border-slate-200 dark:border-slate-700 bg-slate-50 dark:bg-slate-800/50 hover:border-primary/50 transition-colors">
                                        <div className="flex-shrink-0 size-12 rounded-lg bg-white dark:bg-slate-700 flex items-center justify-center text-purple-600 shadow-sm">
                                            <span className="material-symbols-outlined">folder_zip</span>
                                        </div>
                                        <div className="ml-4 flex-1 min-w-0">
                                            <p className="text-sm font-semibold text-slate-900 dark:text-white truncate">
                                                Uploaded Report (ID: {reportId})
                                            </p>
                                            <p className="text-xs text-slate-500 dark:text-slate-400">
                                                PDF Document
                                            </p>
                                        </div>
                                        <div className="ml-2 flex-shrink-0 text-green-500">
                                            <span className="material-symbols-outlined">check_circle</span>
                                        </div>
                                    </div>
                                </div>
                            )}
                        </div>

                        {/* Job Title Input Section */}
                        <div className="px-8 pb-8">
                            <div className="relative max-w-2xl">
                                <label
                                    className="block text-sm font-semibold text-slate-900 dark:text-white mb-2"
                                    htmlFor="job-role"
                                >
                                    Target Job Role
                                    <span className="text-red-500 ml-1">*</span>
                                </label>
                                <div className="relative group">
                                    <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none">
                                        <span className="material-symbols-outlined text-slate-400">work</span>
                                    </div>
                                    <input
                                        autoComplete="off"
                                        className="block w-full pl-11 pr-4 py-3.5 bg-white dark:bg-slate-800 border border-slate-300 dark:border-slate-600 rounded-lg text-slate-900 dark:text-white placeholder-slate-400 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-shadow shadow-sm"
                                        id="job-role"
                                        placeholder="e.g. Senior Product Manager"
                                        type="text"
                                        value={jobTitle}
                                        onChange={(e) => setJobTitle(e.target.value)}
                                        onFocus={() => setShowSuggestions(true)}
                                        onBlur={() => setTimeout(() => setShowSuggestions(false), 200)}
                                    />
                                    {/* Autocomplete Dropdown */}
                                    {showSuggestions && (
                                        <div className="absolute z-10 mt-2 w-full bg-white dark:bg-slate-800 rounded-xl shadow-xl border border-slate-200 dark:border-slate-700 overflow-hidden animate-fade-in-down">
                                            <div className="px-4 py-2 text-xs font-semibold text-slate-500 dark:text-slate-400 bg-slate-50 dark:bg-slate-800/80 border-b border-slate-100 dark:border-slate-700 uppercase tracking-wider">
                                                Suggestions
                                            </div>
                                            <ul className="max-h-60 overflow-auto py-1">
                                                {suggestions.map((s, i) => (
                                                    <li key={i}>
                                                        <button
                                                            className="w-full text-left px-4 py-3 hover:bg-primary/5 dark:hover:bg-primary/20 flex items-center justify-between group/item transition-colors"
                                                            onClick={() => {
                                                                setJobTitle(s.title);
                                                                setShowSuggestions(false);
                                                            }}
                                                        >
                                                            <div className="flex flex-col">
                                                                <span className="text-sm font-medium text-slate-900 dark:text-white group-hover/item:text-primary">
                                                                    {s.title}
                                                                </span>
                                                                <span className="text-xs text-slate-500">{s.dept}</span>
                                                            </div>
                                                            <span className="material-symbols-outlined text-slate-300 group-hover/item:text-primary text-sm transform -rotate-45">
                                                                arrow_forward
                                                            </span>
                                                        </button>
                                                    </li>
                                                ))}
                                            </ul>
                                        </div>
                                    )}
                                </div>
                                {/* Validation Message Helper */}
                                <p className="mt-2 text-sm text-slate-500 dark:text-slate-400 flex items-center gap-1">
                                    <span className="material-symbols-outlined text-base">info</span>
                                    Select a standardized role for better AI accuracy.
                                </p>
                                {error && <p className="mt-2 text-sm text-red-500 font-bold">{error}</p>}
                            </div>
                        </div>

                        {/* Action Area */}
                        <div className="bg-slate-50 dark:bg-slate-900/50 px-8 py-6 border-t border-slate-200 dark:border-slate-800 flex flex-col sm:flex-row items-center justify-between gap-4">
                            <a
                                className="text-slate-600 dark:text-slate-400 hover:text-slate-900 dark:hover:text-white font-medium text-sm transition-colors flex items-center gap-2 group"
                                href="/"
                            >
                                <span className="material-symbols-outlined text-lg group-hover:-translate-x-1 transition-transform">
                                    arrow_back
                                </span>
                                Back to Upload
                            </a>
                            <button
                                className={`w-full sm:w-auto inline-flex items-center justify-center gap-2 rounded-lg shadow-lg shadow-blue-500/20 transition-all transform active:scale-95 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary dark:focus:ring-offset-slate-900 font-semibold py-3 px-8 ${isStarting || !cvId || !reportId
                                        ? 'bg-gray-400 cursor-not-allowed'
                                        : 'bg-primary hover:bg-blue-700 text-white'
                                    }`}
                                disabled={isStarting || !cvId || !reportId}
                                onClick={handleStartEvaluation}
                            >
                                {isStarting ? (
                                    <>
                                        <span className="material-symbols-outlined text-lg animate-spin">refresh</span>
                                        Starting...
                                    </>
                                ) : (
                                    <>
                                        <span className="material-symbols-outlined">smart_toy</span>
                                        Start AI Evaluation
                                    </>
                                )}
                            </button>
                        </div>
                    </div>
                    {/* Additional Info / Footer Text */}
                    <div className="mt-8 text-center">
                        <p className="text-sm text-slate-400 dark:text-slate-500">
                            By starting the evaluation, you agree to our{' '}
                            <a className="underline hover:text-primary" href="#">
                                Data Processing Agreement
                            </a>
                            .
                        </p>
                    </div>
                </div>
            </main>
        </div>
    );
}

export default function EvaluatePage() {
    return (
        <Suspense fallback={<div className="flex h-screen items-center justify-center">Loading...</div>}>
            <EvaluateContent />
        </Suspense>
    );
}
