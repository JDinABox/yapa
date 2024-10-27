import { browser } from '$app/environment';
import { writable, type Writable } from 'svelte/store';

export const dark: Writable<boolean> = writable();
const userToggledDark: Writable<boolean> = writable();

export const toggleDark = () => {
	userToggledDark.set(true);
	dark.update((value) => !value);
};

function getInitDarkMode(): boolean {
	// If the user has previously set a preference, use it
	// Otherwise, use the system preference
	const setDark = localStorage.dark ?? '';
	if (setDark.length > 0) {
		return setDark === 'true';
	}
	return window.matchMedia('(prefers-color-scheme: dark)').matches;
}

function localUserToggledDark(): boolean {
	return localStorage.userToggledDark === 'true';
}

if (browser) {
	// Update the localStorage value whenever the store changes
	dark.subscribe((value) => {
		if (!localUserToggledDark()) {
			return;
		}
		try {
			if (typeof value === 'boolean') {
				localStorage.dark = value.toString();
			}
		} catch (e) {
			console.error('Failed to save dark mode preference to localStorage:', e);
		}
	});
	userToggledDark.subscribe((value) => {
		try {
			if (typeof value === 'boolean') {
				localStorage.userToggledDark = value.toString();
			}
		} catch (e) {
			console.error('Failed to save user toggled dark mode preference to localStorage:', e);
		}
	});

	// Initialize the store
	dark.set(getInitDarkMode());
	userToggledDark.set(localUserToggledDark());
} else {
	dark.set(true);
	userToggledDark.set(false);
}
