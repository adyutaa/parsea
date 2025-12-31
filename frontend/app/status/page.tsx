'use client';

import React, { useState, useEffect, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { api, EvaluationJob } from '../../lib/api';

function StatusContent() {
    const router = useRouter();
    const searchParams = useSearchParams();
    const jobId = searchParams.get('id');

    const [job, setJob] = useState<EvaluationJob | null>(null);
    const [loading, setLoading] = useState(true); // Initial load
    const [fakeProgress, setFakeProgress] = useState(0);
    const [currentStepIndex, setCurrentStepIndex] = useState(1);

    const steps = [
        'Uploading Documents',
        'Extracting Text Data',
        'Identifying Key Skills',
        'Analyzing Project History',
        'Calculating Fit Score',
        'Generating Summary',
        'Finalizing Report',
    ];

    // Polling Effect
    useEffect(() => {
        if (!jobId) return;

        const fetchStatus = async () => {
            try {
                const data = await api.getResult(parseInt(jobId));
                setJob(data);
                setLoading(false);
            } catch (err) {
                console.error('Failed to poll status:', err);
            }
        };

        // Poll every 3 seconds
        fetchStatus();
        const interval = setInterval(fetchStatus, 3000);

        return () => clearInterval(interval);
    }, [jobId]);

    // Simulation Effect for Progress Bar
    useEffect(() => {
        if (job?.status === 'completed') {
            setFakeProgress(100);
            setCurrentStepIndex(steps.length);
            return;
        }

        if (job?.status === 'failed') {
            return;
        }

        // Determine target step based on time or just loop for now
        // Since backend doesn't give granular steps, we simulate a slow progress 
        // over ~45 seconds (approx 2% per sec)
        const progressInterval = setInterval(() => {
            setFakeProgress((prev) => {
                if (prev >= 95) return prev; // stall at 95 until completes
                return prev + 0.5;
            });
        }, 200);

        return () => clearInterval(progressInterval);
    }, [job?.status]);

    // Derive current step index from fake progress
    useEffect(() => {
        if (job?.status !== 'completed' && job?.status !== 'failed') {
            // Map 0-95% to steps 1-6
            const step = Math.floor((fakeProgress / 95) * (steps.length - 2)) + 1;
            setCurrentStepIndex(step);
        }
    }, [fakeProgress, job?.status]);


    if (!jobId) {
        return <div className="flex justify-center p-10 text-red-500">Invalid Job ID</div>;
    }

    const isComplete = job?.status === 'completed';
    const isFailed = job?.status === 'failed';

    return (
        <div className="bg-background-light dark:bg-background-dark text-[#0e121b] dark:text-white font-display min-h-screen flex flex-col">
            {/* Top Navigation */}
            <header className="flex items-center justify-between whitespace-nowrap border-b border-solid border-slate-200 dark:border-slate-800 bg-surface-light dark:bg-surface-dark px-10 py-3 sticky top-0 z-50 shadow-sm">
                <div className="flex items-center gap-4">
                    <div className="size-8 text-primary">
                        <svg
                            className="h-full w-full"
                            fill="none"
                            viewBox="0 0 48 48"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <path
                                clipRule="evenodd"
                                d="M24 18.4228L42 11.475V34.3663C42 34.7796 41.7457 35.1504 41.3601 35.2992L24 42V18.4228Z"
                                fill="currentColor"
                                fillRule="evenodd"
                            ></path>
                            <path
                                clipRule="evenodd"
                                d="M24 8.18819L33.4123 11.574L24 15.2071L14.5877 11.574L24 8.18819ZM9 15.8487L21 20.4805V37.6263L9 32.9945V15.8487ZM27 37.6263V20.4805L39 15.8487V32.9945L27 37.6263ZM25.354 2.29885C24.4788 1.98402 23.5212 1.98402 22.646 2.29885L4.98454 8.65208C3.7939 9.08038 3 10.2097 3 11.475V34.3663C3 36.0196 4.01719 37.5026 5.55962 38.098L22.9197 44.7987C23.6149 45.0671 24.3851 45.0671 25.0803 44.7987L42.4404 38.098C43.9828 37.5026 45 36.0196 45 34.3663V11.475C45 10.2097 44.2061 9.08038 43.0155 8.65208L25.354 2.29885Z"
                                fill="currentColor"
                                fillRule="evenodd"
                            ></path>
                        </svg>
                    </div>
                    <h2 className="text-lg font-bold leading-tight tracking-[-0.015em]">AI Evaluator</h2>
                </div>
                <div className="flex flex-1 justify-end gap-8 items-center">
                    {/* Nav Links */}
                    <div className="flex items-center gap-9 hidden md:flex">
                        <a className="text-sm font-medium leading-normal hover:text-primary transition-colors" href="#">Dashboard</a>
                        <a className="text-sm font-medium leading-normal text-primary" href="#">Evaluations</a>
                        <a className="text-sm font-medium leading-normal hover:text-primary transition-colors" href="#">Settings</a>
                    </div>
                    {/* Avatar */}
                    <div className="bg-center bg-no-repeat aspect-square bg-cover rounded-full size-10 ring-2 ring-slate-200 dark:ring-slate-700" style={{ backgroundImage: 'url("https://lh3.googleusercontent.com/aida-public/AB6AXuAqNmyBVDN8b4lled9HFhO1vXPvaeP62iWcA-eKt3WrdlD7C10Zqqkr__DPOH5AYNPg2vaKhhIpb5RcI4BlGS_7sUIoE8SUXguqjyO8iHSvCvKTu7bLkgxUmoI1WZnBoQApRC_-ad30XwUX2zBBFtCowZbpekxfQKAsgzwzMVd3DbvpOl21jku33oL7OHlsN-uvTbhmIGFIgqqFaqj5pMMeeVQ8iD6meLigUwvemExZ7gk7sCEBL3Bozrv5uRjusWyoOuZOaHAF7YKm")' }}></div>
                </div>
            </header>

            {/* Main Content Area */}
            <main className="flex-1 flex flex-col items-center justify-start py-10 px-4 md:px-10 overflow-y-auto w-full">
                <div className="w-full max-w-4xl flex flex-col gap-6">
                    {/* Page Header */}
                    <div className="flex flex-col gap-2 text-center md:text-left">
                        <h1 className="text-3xl font-bold tracking-tight text-slate-900 dark:text-white">
                            {isComplete ? 'Evaluation Completed' : isFailed ? 'Evaluation Failed' : 'Evaluation in Progress'}
                        </h1>
                        <p className="text-slate-500 dark:text-slate-400 text-base">
                            {isComplete
                                ? 'Your results are ready. View the detailed report below.'
                                : isFailed
                                    ? 'Something went wrong during the analysis.'
                                    : 'Please wait while our AI analyzes the candidate profile against project requirements.'}
                        </p>
                    </div>

                    <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
                        {/* Left Column: Status Visualization */}
                        <div className="lg:col-span-2 flex flex-col gap-6">
                            {/* Progress Card */}
                            <div className="bg-surface-light dark:bg-surface-dark rounded-xl shadow-sm border border-slate-200 dark:border-slate-700 p-6">
                                <div className="flex flex-col gap-4">
                                    <div className="flex justify-between items-end">
                                        <div className="flex flex-col gap-1">
                                            <span className="text-sm font-medium text-primary uppercase tracking-wider">
                                                Current Action
                                            </span>
                                            <h3 className="text-xl font-bold text-slate-900 dark:text-white flex items-center gap-2">
                                                {isFailed ? 'Error Occurred' : isComplete ? 'Finalized' : steps[currentStepIndex]}
                                                {!isComplete && !isFailed && (
                                                    <span className="relative flex h-3 w-3">
                                                        <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary opacity-75"></span>
                                                        <span className="relative inline-flex rounded-full h-3 w-3 bg-primary"></span>
                                                    </span>
                                                )}
                                            </h3>
                                        </div>
                                        <div className="text-right">
                                            <span className="text-3xl font-bold text-primary">{Math.round(fakeProgress)}%</span>
                                        </div>
                                    </div>
                                    {/* Master Progress Bar */}
                                    <div className="w-full bg-slate-100 dark:bg-slate-800 rounded-full h-3 overflow-hidden">
                                        <div
                                            className={`h-3 rounded-full transition-all duration-500 ease-out relative overflow-hidden ${isFailed ? 'bg-red-500' : 'bg-primary'}`}
                                            style={{ width: `${fakeProgress}%` }}
                                        >
                                            {!isComplete && !isFailed && (
                                                <div
                                                    className="absolute inset-0 bg-white/20 w-full h-full animate-shimmer -skew-x-12 origin-top-left"
                                                    style={{
                                                        backgroundImage:
                                                            'linear-gradient(90deg, transparent, rgba(255,255,255,0.2), transparent)',
                                                    }}
                                                ></div>
                                            )}
                                        </div>
                                    </div>
                                    <div className="flex justify-between items-center text-sm text-slate-500 dark:text-slate-400 bg-slate-50 dark:bg-slate-800/50 p-3 rounded-lg border border-slate-100 dark:border-slate-800">
                                        <span className="flex items-center gap-2">
                                            <span className="material-symbols-outlined text-[18px]">timer</span>
                                            Estimated time remaining
                                        </span>
                                        <span className="font-mono font-medium text-slate-700 dark:text-slate-300">
                                            {isComplete ? '0s' : '~45s'}
                                        </span>
                                    </div>
                                </div>
                            </div>

                            {/* Visual Context Card */}
                            <div className="bg-surface-light dark:bg-surface-dark rounded-xl shadow-sm border border-slate-200 dark:border-slate-700 p-0 overflow-hidden relative group">
                                <div className="absolute inset-0 bg-gradient-to-t from-slate-900/80 to-transparent z-10 flex items-end p-6">
                                    <div className="text-white">
                                        <p className="font-bold text-lg mb-1">{isComplete ? 'Analysis Complete' : 'Scanning Document Structure'}</p>
                                        <p className="text-sm text-slate-200 opacity-90">
                                            AI model v4.2 is parsing layout and semantic blocks.
                                        </p>
                                    </div>
                                </div>
                                <div className="h-64 w-full bg-slate-200 dark:bg-slate-800 relative">
                                    {/* Abstract representation of scanning */}
                                    <div
                                        className="w-full h-full bg-cover bg-center opacity-50 dark:opacity-30"
                                        style={{
                                            backgroundImage:
                                                'url("https://lh3.googleusercontent.com/aida-public/AB6AXuAEbfC3_aN_ylg9oSADimrr6vB0dnTL74l4-uHob3Et4yUg0j7xr_XOOcWXYtsnrb4vIPT7tVm96Q-SlRvEcoQJMbfAyncWS-M43bNdu6FdhQNFhSg1kXGAkLm_5K_U4S21QAFjvGeTljFonBJODWYfSibMRYbt4cbA6wpsZykbPrRUY-3Wi2nf4gQPgYd5SnTXlKmnX7Fsz0yzDd5mAYf1FmllYYYQNDxJVUBxGU4WySizjsFLSvMJZ8_BJ95PU7l3heVkmgeF2JNo")',
                                        }}
                                    ></div>
                                    <div className="absolute inset-0 bg-primary/10 backdrop-blur-[1px]"></div>
                                    {/* Scanning Line Animation - Hide when complete */}
                                    {!isComplete && !isFailed && (
                                        <div className="absolute top-0 left-0 w-full h-1 bg-primary shadow-[0_0_15px_rgba(25,93,230,0.8)] animate-scan"></div>
                                    )}
                                </div>
                            </div>
                        </div>

                        {/* Right Column: Vertical Stepper */}
                        <div className="bg-surface-light dark:bg-surface-dark rounded-xl shadow-sm border border-slate-200 dark:border-slate-700 p-6 h-fit">
                            <h3 className="text-lg font-bold mb-6 text-slate-900 dark:text-white">
                                Process Checklist
                            </h3>
                            <div className="relative flex flex-col gap-0">
                                {/* Connecting Line */}
                                <div className="absolute left-[15px] top-4 bottom-4 w-[2px] bg-slate-100 dark:bg-slate-800 z-0"></div>

                                {/* Steps */}
                                {steps.map((step, index) => {
                                    const isActive = index === currentStepIndex;
                                    const isPassed = index < currentStepIndex;

                                    return (
                                        <div key={index} className={`flex gap-4 relative z-10 pb-6 group ${!isActive && !isPassed ? 'opacity-50' : ''}`}>
                                            <div className="flex-none">
                                                {isPassed ? (
                                                    <div className="size-8 rounded-full bg-primary text-white flex items-center justify-center ring-4 ring-surface-light dark:ring-surface-dark">
                                                        <span className="material-symbols-outlined text-[18px]">check</span>
                                                    </div>
                                                ) : isActive ? (
                                                    <div className="size-8 rounded-full bg-white dark:bg-slate-800 border-2 border-primary text-primary flex items-center justify-center ring-4 ring-surface-light dark:ring-surface-dark shadow-[0_0_0_4px_rgba(25,93,230,0.1)]">
                                                        <span className="material-symbols-outlined text-[18px] animate-spin">refresh</span>
                                                    </div>
                                                ) : (
                                                    <div className="size-8 rounded-full bg-slate-100 dark:bg-slate-800 border-2 border-slate-200 dark:border-slate-700 text-slate-400 flex items-center justify-center ring-4 ring-surface-light dark:ring-surface-dark">
                                                        <span className="text-xs font-bold">{index + 1}</span>
                                                    </div>
                                                )}
                                            </div>
                                            <div className="flex-1 pt-1">
                                                <p className={`text-sm font-bold ${isActive ? 'text-primary' : 'text-slate-900 dark:text-white'}`}>
                                                    {step}
                                                </p>
                                                {isActive && <p className="text-xs text-primary/80 mt-0.5 font-medium">Processing...</p>}
                                            </div>
                                        </div>
                                    );
                                })}
                            </div>

                            <div className="mt-8 pt-6 border-t border-slate-100 dark:border-slate-800">
                                <button
                                    className={`w-full py-2.5 px-4 font-bold rounded-lg transition-colors flex items-center justify-center gap-2 ${isComplete
                                        ? 'bg-primary text-white hover:bg-blue-600 shadow-lg cursor-pointer'
                                        : 'bg-slate-100 dark:bg-slate-800 text-slate-400 dark:text-slate-500 cursor-not-allowed'
                                        }`}
                                    disabled={!isComplete}
                                    onClick={() => router.push(`/result?id=${jobId}`)}
                                >
                                    <span>View Results</span>
                                    {isComplete ? (
                                        <span className="material-symbols-outlined text-[18px]">visibility</span>
                                    ) : (
                                        <span className="material-symbols-outlined text-[18px]">lock</span>
                                    )}
                                </button>
                                {!isComplete && (
                                    <p className="text-xs text-center mt-3 text-slate-400">
                                        Action will be available upon completion
                                    </p>
                                )}
                            </div>
                        </div>
                    </div>
                </div>
            </main>
        </div>
    );
}

export default function StatusPage() {
    return (
        <Suspense fallback={<div className="flex h-screen items-center justify-center">Loading...</div>}>
            <StatusContent />
        </Suspense>
    );
}
