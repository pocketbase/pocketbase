<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import PageWrapper from "@/components/base/PageWrapper.svelte";
    import SettingsSidebar from "@/components/settings/SettingsSidebar.svelte";
    import { pageTitle } from "@/stores/app";
    import { setErrors } from "@/stores/errors";
    import { addSuccessToast } from "@/stores/toasts";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";

    $pageTitle = "Theme";

    let currentTheme = 'Default';
    let currentThemeLink;
    let isLoading = false;

    // Function to load the theme from localStorage
    function loadTheme() {
        const savedTheme = localStorage.getItem('pbTheme');
        currentTheme = savedTheme ?? "Default"
        if (savedTheme) {
            toggle(savedTheme);
        }
    }

    async function toggle(name) {
        isLoading = true; // Start loading

        // Create a new link element for the new theme
        const newThemeLink = document.createElement('link');
        newThemeLink.rel = 'stylesheet';

        // Set the href based on the selected theme
        if (name) {
            newThemeLink.href = import.meta.env.BASE_URL + 'themecss/' + name.toLowerCase() + '.css';
            alert(newThemeLink.href)
            currentTheme = name;
        }

        // Append the new theme link to the document head
        document.head.appendChild(newThemeLink);

        // Wait for the new theme to load
        await new Promise((resolve) => {
            newThemeLink.onload = () => {
                // Remove the old theme link after the new one has loaded
                if (currentThemeLink) {
                    document.head.removeChild(currentThemeLink);
                }
                currentThemeLink = newThemeLink; // Update the current theme link reference
                isLoading = false; // End loading

                // Store the selected theme in localStorage
                localStorage.setItem('pbTheme', currentTheme);
                resolve();
            };
        });
    }

    // Load the theme when the component is initialized
    loadTheme();
</script>

<SettingsSidebar />

<PageWrapper>
    <header class="page-header">
        <nav class="breadcrumbs">
            <div class="breadcrumb-item">Settings</div>
            <div class="breadcrumb-item">{$pageTitle}</div>
        </nav>
    </header>
    <div class="wrapper">
    <h1>Interface Theme</h1>
    <h6>choose your workspace theme</h6>
    <div class='separator'/>
    <div class="theme_switch_container">
            <button on:click={() => toggle('default')} class={`apperacance_btn ${currentTheme == 'default' && "active_btn"}`}>
        <img src="{import.meta.env.BASE_URL}images/highlight/default.png"/>
         <div class='apperanace_innner_div'>
            <p>Default</p>
                       <i class={`${currentTheme == 'default' ? "ri-checkbox-circle-line" : 'ri-checkbox-blank-circle-line'}`}></i>

        </div>
    </button>
    <button on:click={() => toggle('oz')} class={`apperacance_btn ${currentTheme == 'oz' && "active_btn"}`}>
        <img src="{import.meta.env.BASE_URL}images/highlight/oz.png"/>
        <div class='apperanace_innner_div'>
            <p>OZ</p>
            <i class={`${currentTheme == 'oz' ? "ri-checkbox-circle-line" : 'ri-checkbox-blank-circle-line'}`}></i>

        </div>
    </button>
    <button on:click={() => toggle('hardoz')} class={`apperacance_btn ${currentTheme == 'hardoz' && "active_btn"}`}>
       <img src="{import.meta.env.BASE_URL}images/highlight/hardoz.png"/>
        <div class='apperanace_innner_div'>
            <p>HardOZ</p>
                        <i class={`${currentTheme == 'hardoz' ? "ri-checkbox-circle-line" : 'ri-checkbox-blank-circle-line'}`}></i>

        </div>
    </button>
    <button on:click={() => toggle('blueshi')} class={`apperacance_btn ${currentTheme == 'blueshi' && "active_btn"}`}>
         <img src="{import.meta.env.BASE_URL}images/highlight/blueshi.png"/>
          <div class='apperanace_innner_div'>
            <p>Blueshi</p>
                       <i class={`${currentTheme == 'blueshi' ? "ri-checkbox-circle-line" : 'ri-checkbox-blank-circle-line'}`}></i>

        </div>
    </button>

    </div>
    
    </div>
</PageWrapper>

<style>
    .theme_switch_container{
        display: flex;
        overflow-x: scroll;
        gap: 10px;
        margin: 20px 0px;
/*        flex-wrap: wrap;*/
    }
    .apperacance_btn{
        background: var(--baseAlt1Color);
        overflow: hidden;
        padding: 17px 0px 0px 17px;
        border-radius: 7px;
        cursor: pointer;
        width: 26em;
        height: 9em;
        border: var(--baseAlt4Color) solid 2px;
        position: relative;
        color: var(--txtPrimaryColor);
    }
    .apperanace_innner_div{
        position: absolute;
        bottom: 0px;
        width: 100%;
        height: 2.2em;
        background: var(--baseAlt1Color);
        left: 0px;
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 0px 10px;
        border-top: var(--baseAlt4Color) solid 2px;
    }
    .active_btn{
        border: var(--primaryColor) solid 2px;
        color: var(--primaryColor);
    }
    .separator{
        width: 100%;
        height: 1px;
        background: var(--baseAlt4Color);
        margin: 2em 0px;
    }
</style>