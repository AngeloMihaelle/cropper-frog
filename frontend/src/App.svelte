<script>
    import { onMount } from "svelte";
    import { status, toast, showToast } from "./lib/Store.js";
    import ClipComponent from "./lib/Clip.svelte";

    // --- Wails Runtime Imports ---
    import { EventsOn } from "@wailsjs/runtime/runtime.js";
    
    // --- Wails Go Function Imports ---
    import {
        GetVideoDuration,
        ValidateClip,
        StartCropping,
        OpenFileDialog,
        OpenDirectoryDialog,
    } from "@wailsjs/go/main/App.js";

    // --- Local State ---
    let videoPath = "";
    let videoDuration = 0.0;
    let clipQueue = [];
    let isProcessing = false;

    // Form inputs
    let clipName = "";
    let startTime = "";
    let endTime = "";

    // --- Status Listeners  ---
    let statusMessage = "Ready to load video.";
    let progress = 0;
    let eta = "";

    // Subscribe to our global status store
    status.subscribe((s) => {
        statusMessage = s.message;
        progress = s.progress;
        eta = s.eta;
    });

    // Listen for all events from the Go backend
    onMount(() => {
        EventsOn("processing:status", (msg) => {
            status.set({ ...$status, message: msg });
            if (msg.includes("successfully") || msg.includes("Error!")) {
                isProcessing = false; // Re-enable buttons
                if (msg.includes("successfully")) {
                    clipQueue = []; // Clear queue on success
                }
            }
        });
        EventsOn("processing:progress", (p) =>
            status.set({ ...$status, progress: p }),
        );
        EventsOn("processing:eta", (e) => status.set({ ...$status, eta: e }));

        EventsOn("processing:error", (err) => {
            showToast(err, "error");
            isProcessing = false; // Re-enable buttons on error
            status.set({
                message: "Error! Ready for new job.",
                progress: 0,
                eta: "",
            });
        });
    });

    // --- UI Functions ---

    async function handleLoadVideo() {

        if (isProcessing) return;
        try {
            const path = await OpenFileDialog({
                title: "Select Video File",
                filters: [
                    {
                        displayName: "Video Files",
                        patterns: ["*.mp4;*.mov;*.mkv;*.avi;*.webm"],
                    },
                ],
            });
            if (!path) return; // User cancelled

            videoPath = path;
            status.set({
                message: "Reading video metadata...",
                progress: 0,
                eta: "",
            });

            const duration = await GetVideoDuration(path);
            videoDuration = duration;

            const friendlyDuration = new Date(duration * 1000)
                .toISOString()
                .substr(11, 8);
            status.set({
                message: `Loaded. Total time: ${friendlyDuration}`,
                progress: 0,
                eta: "",
            });
        } catch (err) {
            showToast(err.message || err, "error");
            status.set({
                message: "Failed to load video.",
                progress: 0,
                eta: "",
            });
        }
    }

    async function handleAddInterval() {
        if (isProcessing) return;

        const newClip = {
            Name: clipName,
            StartTime: startTime,
            EndTime: endTime,
        };

        try {
            // Call Go for Validation
            await ValidateClip(newClip, videoDuration);

            // Validation passed, add to UI queue
            clipQueue = [...clipQueue, newClip];

            // Reset form
            clipName = "";
            startTime = "";
            endTime = "";
        } catch (err) {
            // Validation failed, show error
            showToast(err.message || err, "error");
        }
    }

    function handleRemoveClip(index) {
        // (Req 3.3)
        if (isProcessing) return;
        clipQueue.splice(index, 1);
        clipQueue = clipQueue;
    }

    async function handleStartCropping() {
        if (isProcessing || clipQueue.length === 0) return;

        try {
            const outDir = await OpenDirectoryDialog({
                title: "Select Output Folder",
            });
            if (!outDir) return; // User cancelled

            isProcessing = true; // Disable buttons
            status.set({
                message: "Starting job...",
                progress: 0,
                eta: "Calculating...",
            });

            await StartCropping(clipQueue, videoPath, outDir);

        } catch (err) {
            showToast(err.message || err, "error");
            isProcessing = false;
            status.set({
                message: "Error starting job.",
                progress: 0,
                eta: "",
            });
        }
    }
</script>


<main class="flex flex-col h-screen p-4 gap-4 max-w-6xl mx-auto">

    <div class="flex-shrink-0">
        <h1 class="text-2xl font-bold text-white mb-2">Cropper Frog</h1>
        <div class="flex gap-4">
            <input
                type="text"
                readonly
                bind:value={videoPath}
                placeholder="Click 'Browse' to load a video file..."
                class="flex-1 p-2 bg-gray-800 border border-gray-700 rounded-md text-gray-300"
            />
            <button
                on:click={handleLoadVideo}
                disabled={isProcessing}
                class="px-4 py-2 bg-blue-600 text-white font-semibold rounded-md shadow-sm hover:bg-blue-500 disabled:bg-gray-600 disabled:opacity-50 transition-colors"
            >
                Browse
            </button>
        </div>
    </div>

    <!-- 2. Main Content (Queue & Form) -->
    <div class="flex-1 flex flex-col md:flex-row gap-4 overflow-hidden">
        <!-- 2a. Crop Intervals Queue -->
        <div class="flex flex-col flex-1 gap-4 overflow-hidden">
            <h2 class="text-xl font-semibold text-white">
                Crop Intervals ({clipQueue.length})
            </h2>
            <div
                class="flex-1 p-2 bg-gray-900 border border-gray-800 rounded-lg overflow-y-auto"
            >
                {#if clipQueue.length === 0}
                    <div
                        class="flex items-center justify-center h-full text-gray-500"
                    >
                        Add intervals using the form
                    </div>
                {:else}
                    <div class="grid grid-cols-1 xl:grid-cols-2 gap-3">
                        {#each clipQueue as clip, i}
                            <ClipComponent
                                {clip}
                                onRemove={() => handleRemoveClip(i)}
                            />
                        {/each}
                    </div>
                {/if}
            </div>
        </div>

        <!-- 2b. Add Interval Form -->
        <div
            class="w-full md:w-72 flex-shrink-0 flex flex-col gap-4 p-4 bg-gray-850 rounded-lg shadow-lg"
        >
            <h3 class="text-lg font-semibold text-white">Add New Clip</h3>
            <div class="flex flex-col gap-3">
                <label>
                    <span class="text-sm font-medium text-gray-300"
                        >New Clip Name</span
                    >
                    <input
                        type="text"
                        bind:value={clipName}
                        placeholder="e.g., Intro"
                        class="mt-1 w-full p-2 bg-gray-700 border border-gray-600 rounded-md text-white"
                    />
                </label>
                <label>
                    <span class="text-sm font-medium text-gray-300"
                        >Start Time (HH:MM:SS)</span
                    >
                    <input
                        type="text"
                        bind:value={startTime}
                        placeholder="00:01:30"
                        class="mt-1 w-full p-2 bg-gray-700 border border-gray-600 rounded-md text-white"
                    />
                </label>
                <label>
                    <span class="text-sm font-medium text-gray-300"
                        >End Time (HH:MM:SS)</span
                    >
                    <input
                        type="text"
                        bind:value={endTime}
                        placeholder="00:10:45"
                        class="mt-1 w-full p-2 bg-gray-700 border border-gray-600 rounded-md text-white"
                    />
                </label>
            </div>
            <button
                on:click={handleAddInterval}
                disabled={isProcessing || !videoPath}
                class="mt-2 w-full flex justify-center items-center gap-2 px-4 py-2 bg-indigo-600 text-white font-semibold rounded-md shadow-sm hover:bg-indigo-500 disabled:bg-gray-600 disabled:opacity-50 transition-colors"
            >
                <!-- Plus Icon -->
                <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-5 w-5"
                    viewBox="0 0 20 20"
                    fill="currentColor"
                >
                    <path
                        fill-rule="evenodd"
                        d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z"
                        clip-rule="evenodd"
                    />
                </svg>
                Add Interval
            </button>
        </div>
    </div>

    <!-- 3. Footer / Status & Controls-->
    <div
        class="flex-shrink-0 flex flex-col sm:flex-row items-center gap-4 p-4 bg-gray-850 rounded-lg"
    >
        <div class="flex-1 w-full">
            <!-- Status Label -->
            <div class="text-sm text-gray-300" id="status-label">
                {statusMessage}
            </div>
            <!-- ETA Label -->
            <div class="text-xs text-gray-400" id="eta-label">
                {eta ? `ETA: ${eta}` : ""}
            </div>

            <!-- Progress Bar -->
            <div class="w-full bg-gray-700 rounded-full h-2.5 mt-2">
                <div
                    class="bg-blue-600 h-2.5 rounded-full"
                    style="width: {progress * 100}%"
                ></div>
            </div>
        </div>
        <button
            on:click={handleStartCropping}
            disabled={isProcessing || clipQueue.length === 0}
            class="w-full sm:w-auto px-6 py-3 bg-green-600 text-white text-lg font-bold rounded-md shadow-lg hover:bg-green-500 disabled:bg-gray-600 disabled:opacity-50 transition-colors"
        >
            Start Cropping
        </button>
    </div>

    <!-- Global Toast/Error component  -->
    {#if $toast}
        <div
            class="fixed top-5 right-5 p-4 rounded-md shadow-lg z-50
                {$toast.type === 'error' ? 'bg-red-600' : 'bg-blue-600'}
                text-white"
            role="alert"
        >
            {$toast.message}
        </div>
    {/if}
</main>
