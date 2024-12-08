import { register, init, getLocaleFromNavigator } from 'svelte-i18n';

register('de', () => import('./locales/de.json'));
register('en', () => import('./locales/en.json'));
register('es', () => import('./locales/es.json'));
register('fr', () => import('./locales/fr.json'));
register('ja', () => import('./locales/ja.json'));
register('ko', () => import('./locales/ko.json'));
register('ru', () => import('./locales/ru.json'));
register('zh', () => import('./locales/zh.json'));

await init({
  fallbackLocale: 'en',
  initialLocale: getLocaleFromNavigator(),
});
