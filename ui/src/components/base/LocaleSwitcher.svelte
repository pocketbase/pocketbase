<script>
  import { createEventDispatcher } from 'svelte';  
  import Toggler from "@/components/base/Toggler.svelte";

  export let value = 'en';
  let langs = [
    { id: 'en', name: 'English' },
    { id: 'zh', name: '中文' },
    { id: 'fr', name: 'Français' },
  ];

  const dispatch = createEventDispatcher();
  function switchLocale(event) {
    event.preventDefault();
    dispatch('locale-changed', event.target.getAttribute('lang'));
  }
  let lscss = `<style>
  .localeSwitcher {
      background-color: transparent;
  }
  .localeSwitchPanel {
    right: auto;
    left: 44px;
    top: 0;
  }
</style>`;
</script>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/gh/syone/locale-icons/css/icons.css">
{@html lscss}

<figure class="thumb link-hint closable localeSwitcher">
  <i class="lc lc-{value}"></i>
  <Toggler class="dropdown dropdown-nowrap localeSwitchPanel">
    {#each langs as lang}
      <div tabindex="0" class="dropdown-item closable" lang="{lang.id}" on:click={switchLocale}>
        <i class="lc lc-{lang.id}" lang="{lang.id}"></i> {lang.name}
      </div>
    {/each}
  </Toggler>
</figure>