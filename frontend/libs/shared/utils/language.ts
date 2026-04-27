export function isLanguage(locale: string | null): Lang | '' {
  if (locale === 'en' || locale === 'de') {
    return locale;
  }
  return '';
}

export function getStoreLanguage(): Lang | '' {
  const language = localStorage.getItem('appLanguage');
  return isLanguage(language);
}
