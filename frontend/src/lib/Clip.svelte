<script>
    export let clip;
    export let onRemove;

    function formatTime(time) {
        // Simple helper to make 00:01:30 look nice
        if (!time) return "N/A";
        return time;
    }
</script>


<div
    class="bg-gray-800 p-4 rounded-lg shadow-md flex justify-between items-center relative transition-all hover:bg-gray-700"
>
    <div>
        <h3 class="font-bold text-white text-lg">{clip.Name}</h3>
        <div class="text-sm text-gray-400">
            <span>Start: {formatTime(clip.StartTime)}</span>
            <span class="mx-2">|</span>
            <span>End: {formatTime(clip.EndTime)}</span>
        </div>
    </div>
    <button
        on:click={onRemove}
        class="absolute top-2 right-2 p-1 text-gray-500 hover:text-red-400 transition-colors rounded-full"
        aria-label="Remove clip"
    >
        <svg
            xmlns="http://www.w3.org/2000/svg"
            class="h-5 w-5"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
        >
            <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M6 18L18 6M6 6l12 12"
            />
        </svg>
    </button>
</div>
