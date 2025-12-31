'use client';

import React, { useState, useEffect, Suspense } from 'react';
import { useSearchParams } from 'next/navigation';
import { api, EvaluationJob } from '../../lib/api';

function ResultContent() {
    const searchParams = useSearchParams();
    const jobId = searchParams.get('id');
    const [job, setJob] = useState<EvaluationJob | null>(null);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        if (!jobId) return;
        api.getResult(parseInt(jobId))
            .then(data => {
                setJob(data);
                setLoading(false);
            })
            .catch(err => {
                console.error(err);
                setLoading(false);
            });
    }, [jobId]);

    if (loading) return <div className="flex h-screen items-center justify-center">Loading...</div>;
    if (!job || !job.result) return <div className="flex h-screen items-center justify-center">No results found</div>;

    const { result } = job;

    // Helper to determine score color
    const getScoreColor = (score: number, max: number = 5) => {
        const percentage = score / max;
        if (percentage >= 0.9) return 'score-high';
        if (percentage >= 0.7) return 'score-good';
        if (percentage >= 0.5) return 'score-avg';
        if (percentage >= 0.3) return 'score-low';
        return 'score-poor';
    };

    const getScoreColorHex = (score: number, max: number = 5) => {
        const percentage = score / max;
        if (percentage >= 0.9) return '#16a34a';
        if (percentage >= 0.7) return '#84cc16';
        if (percentage >= 0.5) return '#eab308';
        if (percentage >= 0.3) return '#f97316';
        return '#ef4444';
    };

    return (
        <div className="bg-background-light dark:bg-background-dark text-[#0e121b] dark:text-white font-display min-h-screen flex flex-col overflow-x-hidden transition-colors duration-300">
            {/* Top Navigation */}
            <header className="sticky top-0 z-50 w-full border-b border-[#e7ebf3] dark:border-gray-800 bg-white/80 dark:bg-[#111621]/90 backdrop-blur-md">
                <div className="px-4 md:px-10 py-3 flex items-center justify-between max-w-[1400px] mx-auto w-full">
                    <div className="flex items-center gap-4 text-[#0e121b] dark:text-white">
                        <div className="size-8 flex items-center justify-center rounded-lg bg-primary/10 text-primary">
                            <span className="material-symbols-outlined text-2xl">analytics</span>
                        </div>
                        <h2 className="text-lg font-bold leading-tight tracking-[-0.015em]">EvalAI</h2>
                    </div>
                    <div className="flex items-center gap-6">
                        <nav className="hidden md:flex gap-6 text-sm font-medium text-slate-600 dark:text-slate-300">
                            <a className="hover:text-primary transition-colors" href="/">Dashboard</a>
                            <a className="text-primary" href="#">Results</a>
                            <a className="hover:text-primary transition-colors" href="#">History</a>
                        </nav>
                        <div className="flex items-center gap-4">
                            <button className="p-2 text-slate-500 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-full">
                                <span className="material-symbols-outlined">settings</span>
                            </button>
                            <div className="bg-center bg-no-repeat aspect-square bg-cover rounded-full size-10 ring-2 ring-white dark:ring-gray-800 cursor-pointer"
                                style={{ backgroundImage: 'url("https://lh3.googleusercontent.com/aida-public/AB6AXuDUKenFExToJ9b6ncjNpVBm4PWwCxwXK7u9VytzADDWeoLvxtNcx8DSPyf11KxPtVw5fEwCnW0_tts1bpNJBpIuGTOtgV7YQsVtBo53LS5sYeL6D_7Ijs38Ng_pY5v3hZgMdnHhndyMTBlsj-db44R3xjrqanjbSQMGd7ihBMIF1NPX4T0Mxe8zwpa-oY8A87IUhZIPYNYkeLDTfqDM6-KBas6mw3wfOGQ31heehwEnU7HIsfxx71xWJ5_tUA_veRiFrpySEXCLOxOA")' }}>
                            </div>
                        </div>
                    </div>
                </div>
            </header>

            <main className="flex-1 w-full max-w-[1280px] mx-auto px-4 md:px-10 py-8">
                {/* Header & Actions */}
                <div className="flex flex-col md:flex-row justify-between items-start md:items-end gap-6 mb-8">
                    <div className="flex flex-col gap-2">
                        <div className="flex items-center gap-2 text-sm font-medium text-primary bg-primary/10 px-3 py-1 w-fit rounded-full">
                            <span className="material-symbols-outlined text-sm">check_circle</span>
                            Evaluation Complete
                        </div>
                        <h1 className="text-3xl md:text-4xl font-black leading-tight tracking-tight text-[#0e121b] dark:text-white">
                            Evaluation Results
                        </h1>
                        <p className="text-slate-500 dark:text-slate-400 text-base font-normal">
                            Job Role: {job.job_title || 'N/A'} <span className="mx-2">•</span> ID: {job.id}
                        </p>
                    </div>
                    <div className="flex gap-3">
                        <button className="flex items-center gap-2 px-4 py-2 bg-white dark:bg-gray-800 border border-slate-200 dark:border-gray-700 rounded-lg text-sm font-medium hover:bg-slate-50 dark:hover:bg-gray-700 transition-colors cursor-pointer">
                            <span className="material-symbols-outlined text-lg">share</span>
                            Share
                        </button>
                        <button className="flex items-center gap-2 px-4 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors shadow-lg shadow-blue-500/20 cursor-pointer">
                            <span className="material-symbols-outlined text-lg">download</span>
                            Download Report
                        </button>
                    </div>
                </div>

                {/* KPI Grid */}
                <div className="grid grid-cols-1 md:grid-cols-12 gap-6 mb-8">
                    {/* CV Match Rate Card */}
                    <div className="col-span-1 md:col-span-5 lg:col-span-4 bg-white dark:bg-[#1e2330] rounded-xl p-6 border border-slate-100 dark:border-gray-800 shadow-sm flex flex-col justify-between">
                        <div className="flex justify-between items-start mb-4">
                            <div>
                                <h3 className="text-lg font-bold text-[#0e121b] dark:text-white">CV Match Rate</h3>
                                <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">Relevance to Job Description</p>
                            </div>
                            <div className="group relative">
                                <span className="material-symbols-outlined text-slate-400 cursor-help text-xl">info</span>
                                <div className="absolute right-0 w-48 p-2 bg-gray-800 text-white text-xs rounded shadow-lg opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none z-10">
                                    Based on keyword density, skills overlap, and experience duration matching.
                                </div>
                            </div>
                        </div>
                        <div className="flex flex-col items-center justify-center py-4">
                            {/* Gauge Container */}
                            <div
                                className="relative size-40 conic-gauge flex items-center justify-center"
                                style={{
                                    '--gauge-value': `${result.scan_match * 100}%`,
                                    '--gauge-color': getScoreColorHex(result.scan_match, 1)
                                } as React.CSSProperties}
                            >
                                {/* Inner Circle */}
                                <div className="size-32 bg-white dark:bg-[#1e2330] rounded-full flex flex-col items-center justify-center shadow-inner">
                                    <span className="text-4xl font-black text-[#0e121b] dark:text-white">
                                        {Math.round(result.scan_match * 100)}%
                                    </span>
                                    <span className={`text-xs font-semibold uppercase tracking-wider text-${getScoreColor(result.scan_match, 1)}`}>
                                        {result.scan_match >= 0.8 ? 'Excellent' : result.scan_match >= 0.6 ? 'Good' : 'Average'}
                                    </span>
                                </div>
                            </div>
                        </div>
                        <div className="mt-4 pt-4 border-t border-slate-100 dark:border-gray-800">
                            <p className="text-sm text-slate-600 dark:text-slate-300 flex items-start gap-2 max-h-24 overflow-y-auto">
                                <span className={`material-symbols-outlined text-${getScoreColor(result.scan_match, 1)} shrink-0`}>check_circle</span>
                                {result.scan_feedback ? result.scan_feedback.substring(0, 100) + '...' : 'Analysis available below.'}
                            </p>
                        </div>
                    </div>

                    {/* Project Score Card */}
                    <div className="col-span-1 md:col-span-7 lg:col-span-8 bg-white dark:bg-[#1e2330] rounded-xl p-6 border border-slate-100 dark:border-gray-800 shadow-sm flex flex-col">
                        <div className="flex justify-between items-start mb-6">
                            <div>
                                <h3 className="text-lg font-bold text-[#0e121b] dark:text-white">Project Score</h3>
                                <p className="text-sm text-slate-500 dark:text-slate-400 mt-1">Technical Assessment Analysis</p>
                            </div>
                            <div className="flex items-center gap-1 bg-score-good/10 text-score-good px-3 py-1 rounded-full text-sm font-bold">
                                <span>{result.project_score}</span>
                                <span className="text-xs font-normal text-slate-500 dark:text-slate-400">/ 5.0</span>
                            </div>
                        </div>
                        <div className="flex flex-col lg:flex-row gap-8">
                            {/* Main Score Big */}
                            <div className="flex flex-col justify-center min-w-[180px]">
                                <div className="flex gap-1 mb-2 text-score-good" style={{ color: getScoreColorHex(result.project_score, 5) }}>
                                    {[1, 2, 3, 4, 5].map((star) => (
                                        <span key={star} className="material-symbols-outlined text-3xl fill-current">
                                            {result.project_score >= star ? 'star' : result.project_score >= star - 0.5 ? 'star_half' : 'star_outline'}
                                        </span>
                                    ))}
                                </div>
                                <p className="text-3xl font-black text-[#0e121b] dark:text-white">{result.project_score}</p>
                                <p className="text-sm text-slate-500 dark:text-slate-400">Weighted Average</p>
                            </div>
                            {/* Breakdown Bars (Mocked for visual parity as backend doesn't provide breakdown yet) */}
                            <div className="flex-1 grid gap-3">
                                {[
                                    { label: 'Code Quality', val: result.project_score + 0.3 > 5 ? 5 : result.project_score + 0.3 },
                                    { label: 'Functionality', val: result.project_score - 0.2 < 0 ? 0 : result.project_score - 0.2 },
                                    { label: 'Performance', val: result.project_score },
                                    { label: 'Documentation', val: result.project_score + 0.5 > 5 ? 5 : result.project_score + 0.5 }
                                ].map((item, i) => (
                                    <div key={i} className="grid grid-cols-[120px_1fr_40px] items-center gap-4">
                                        <span className="text-sm font-medium text-slate-700 dark:text-slate-200">{item.label}</span>
                                        <div className="h-2 w-full bg-slate-100 dark:bg-gray-700 rounded-full overflow-hidden">
                                            <div
                                                className="h-full rounded-full transition-all duration-1000"
                                                style={{
                                                    width: `${(item.val / 5) * 100}%`,
                                                    backgroundColor: getScoreColorHex(item.val, 5)
                                                }}
                                            ></div>
                                        </div>
                                        <span className="text-sm font-bold text-slate-700 dark:text-slate-200 text-right">{item.val.toFixed(1)}</span>
                                    </div>
                                ))}
                            </div>
                        </div>
                    </div>
                </div>

                {/* Overall Summary */}
                <div className="bg-white dark:bg-[#1e2330] rounded-xl p-6 border border-slate-100 dark:border-gray-800 shadow-sm mb-8">
                    <h2 className="text-xl font-bold text-[#0e121b] dark:text-white mb-4 flex items-center gap-2">
                        <span className="material-symbols-outlined text-primary">summarize</span>
                        Overall Summary
                    </h2>
                    <div className="bg-blue-50 dark:bg-blue-900/20 p-5 rounded-lg border-l-4 border-primary">
                        <p className="text-base leading-relaxed text-slate-700 dark:text-slate-200">
                            {result.summary}
                        </p>
                    </div>
                </div>

                {/* Detailed Feedback Grid */}
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* CV Feedback */}
                    <div className="bg-white dark:bg-[#1e2330] rounded-xl p-6 border border-slate-100 dark:border-gray-800 shadow-sm h-full">
                        <div className="flex items-center gap-3 mb-6 pb-4 border-b border-slate-100 dark:border-gray-800">
                            <div className="size-10 rounded-full bg-purple-100 dark:bg-purple-900/30 flex items-center justify-center text-purple-600 dark:text-purple-400">
                                <span className="material-symbols-outlined">description</span>
                            </div>
                            <div>
                                <h3 className="text-lg font-bold text-[#0e121b] dark:text-white">CV Analysis</h3>
                                <p className="text-sm text-slate-500 dark:text-slate-400">Strengths & Weaknesses</p>
                            </div>
                        </div>
                        <div className="space-y-6">
                            <div>
                                <h4 className="text-sm font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3">AI Feedback</h4>
                                <ul className="space-y-3">
                                    <li className="flex gap-3 text-sm text-slate-600 dark:text-slate-300">
                                        <span className="material-symbols-outlined text-score-high text-lg shrink-0">check</span>
                                        <span>{result.scan_feedback}</span>
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>

                    {/* Project Feedback */}
                    <div className="bg-white dark:bg-[#1e2330] rounded-xl p-6 border border-slate-100 dark:border-gray-800 shadow-sm h-full">
                        <div className="flex items-center gap-3 mb-6 pb-4 border-b border-slate-100 dark:border-gray-800">
                            <div className="size-10 rounded-full bg-orange-100 dark:bg-orange-900/30 flex items-center justify-center text-orange-600 dark:text-orange-400">
                                <span className="material-symbols-outlined">code</span>
                            </div>
                            <div>
                                <h3 className="text-lg font-bold text-[#0e121b] dark:text-white">Project Review</h3>
                                <p className="text-sm text-slate-500 dark:text-slate-400">Code Quality & Logic</p>
                            </div>
                        </div>
                        <div className="space-y-6">
                            <div>
                                <h4 className="text-sm font-bold text-slate-900 dark:text-white uppercase tracking-wider mb-3">Technical Analysis</h4>
                                <ul className="space-y-3">
                                    <li className="flex gap-3 text-sm text-slate-600 dark:text-slate-300">
                                        <span className="material-symbols-outlined text-score-high text-lg shrink-0">add_circle</span>
                                        <span>{result.project_feedback}</span>
                                    </li>
                                </ul>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Footer Actions */}
                <div className="flex justify-center mt-10">
                    <button className="text-slate-500 hover:text-primary text-sm font-medium flex items-center gap-1 transition-colors cursor-pointer" onClick={() => window.location.href = '/'}>
                        <span className="material-symbols-outlined text-lg">history</span>
                        Start New Evaluation
                    </button>
                </div>
            </main>
        </div>
    );
}

export default function ResultPage() {
    return (
        <Suspense fallback={<div className="flex h-screen items-center justify-center">Loading...</div>}>
            <ResultContent />
        </Suspense>
    );
}
