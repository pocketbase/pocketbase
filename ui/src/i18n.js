import { register, init, getLocaleFromNavigator } from 'svelte-i18n';

register('zh', () => import('./locales/zh.json'));
register('en', () => import('./locales/en.json'));

await init({
  fallbackLocale: 'en',
  initialLocale: getLocaleFromNavigator(),
});
