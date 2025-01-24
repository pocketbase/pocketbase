import { writable } from 'svelte/store';

// Initialize darkMode store with the value from localStorage (if it exists)
export const darkMode = writable(localStorage.getItem('darkMode') === 'true');

// Subscribe to changes and update localStorage and the theme
darkMode.subscribe(value => {
    localStorage.setItem('darkMode', value);
    document.documentElement.setAttribute('data-theme', value ? 'dark' : 'light');
});

// Initialize the theme on page load
document.documentElement.setAttribute('data-theme', localStorage.getItem('darkMode') === 'true' ? 'dark' : 'light');
