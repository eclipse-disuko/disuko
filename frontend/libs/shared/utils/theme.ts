export function getDefaultTheme(): string {
  const theme = localStorage.getItem('disco-theme');
  if (theme === 'dark' || theme === 'light') {
    return theme;
  }
  return 'dark'; // Default theme
}
