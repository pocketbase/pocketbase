import { derived } from 'svelte/store';
import { init, getLocaleFromNavigator, dictionary, locale, _ } from 'svelte-i18n';
 
const MESSAGE_FILE_URL_TEMPLATE = 'lang/{locale}.json';
const dir = derived(locale, $locale => $locale === 'ar' ? 'rtl' : 'ltr');
 
let cachedLocale;

function setupI18n({ withLocale: _locale } = { withLocale: 'en' }) {
  const messsagesFileUrl = MESSAGE_FILE_URL_TEMPLATE.replace('{locale}', _locale);
  return fetch(messsagesFileUrl)
      .then(response => response.json())
      .then((messages) => {
          dictionary.set({ [_locale]: messages });
          cachedLocale = _locale;
          locale.set(_locale);
          localStorage.setItem('locale', _locale);
      });
}

let current_locale = localStorage.getItem('locale')

if(current_locale == null){
  current_locale = "en"
}

setupI18n({ withLocale: current_locale })

init({
  fallbackLocale: 'en',
  initialLocale: getLocaleFromNavigator(),
})

export { _, locale, dir, setupI18n };