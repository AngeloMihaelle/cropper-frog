import { writable } from 'svelte/store';

// Store for status messages, progress, and ETA
export const status = writable({
    message: "Ready to load video.", // Idle status
    progress: 0,
    eta: ""
});

// Store for error toasts
export const toast = writable(null);

export function showToast(message, type = 'error', duration = 4000) {
    toast.set({ message, type });
    setTimeout(() => {
        toast.set(null);
    }, duration);
}