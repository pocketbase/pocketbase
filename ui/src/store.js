const notifyChannel = new BroadcastChannel("tabsSync");

const SETTINGS_STORAGE_KEY = "pbSettings";
const COLOR_SCHEME_STORAGE_KEY = "pbColorScheme";

window.app = window.app || {};
window.app.store = store({
    // flag used to track when the internal bootstrap process is done
    _ready: false,

    // the current authenticated superuser
    superuser: null,

    // used to force hiding the header even when authenticated
    showHeader: true,

    page: t.div({ className: "page" }, () => {
        if (!app.store._ready) {
            return t.span({ className: "loader lg m-auto", title: "Loading plugins..." });
        }
    }),

    mainLogo: import.meta.env.BASE_URL + "images/logo.svg",
    headerLogo: import.meta.env.BASE_URL + "images/logo_white.svg",
    favicon: "", // leave empty to fallback to the default one

    title: "",

    _mediaColorScheme: "",
    userColorScheme: window.localStorage.getItem(COLOR_SCHEME_STORAGE_KEY) || "",
    get activeColorScheme() {
        // explicitly set
        if (app.store.userColorScheme) {
            return app.store.userColorScheme;
        }

        // fallback to the loaded browser preference
        return app.store._mediaColorScheme || "light";
    },

    // api response errors
    errors: null,

    creditLinks: [
        {
            // optional: isActive
            href: import.meta.env.PB_DOCS_URL,
            icon: "ri-book-open-line",
            label: "Docs",
        },
        {
            href: import.meta.env.PB_RELEASES,
            icon: "ri-github-line",
            label: `PocketBase ${import.meta.env.PB_VERSION}`,
        },
    ],

    headerLinks: [
        {
            // optional: isActive
            href: "#/collections",
            icon: "ri-database-2-line",
            label: "Collections",
        },
        {
            href: "#/logs",
            icon: "ri-bar-chart-box-line",
            label: "Logs",
        },
        {
            href: "#/settings",
            icon: "ri-settings-3-line",
            label: "Settings",
        },
    ],

    settingsNavGroups: {
        System: [
            {
                // optional: isActive
                href: "#/settings",
                icon: "ri-home-gear-line",
                label: "Application",
            },
            {
                href: "#/settings/mail",
                icon: "ri-send-plane-2-line",
                label: "Mail settings",
            },
            {
                href: "#/settings/storage",
                icon: "ri-archive-drawer-line",
                label: "Files storage",
            },
            {
                href: "#/settings/backups",
                icon: "ri-archive-line",
                label: "Backups",
            },
            {
                href: "#/settings/crons",
                icon: "ri-time-line",
                label: "Crons",
            },
        ],
        Sync: [
            {
                href: "#/settings/export-collections",
                icon: "ri-uninstall-line",
                label: "Export collections",
            },
            {
                href: "#/settings/import-collections",
                icon: "ri-install-line",
                label: "Import collections",
            },
        ],
    },

    predefinedAccentColors: [
        "#1055c9",
        "#a3142a",
        "#096d5c",
        "#e6620a",
        "#007d9c",
        "#3f3da9",
    ],

    settings: app.utils.getLocalHistory(SETTINGS_STORAGE_KEY, {}),
    isLoadingSettings: false,
    async loadSettings() {
        app.store.isLoadingSettings = true;

        try {
            const settings = await app.pb.settings.getAll({ requestKey: "appStore.loadSettings" });

            app.store.settings = settings;
            app.store.isLoadingSettings = false;
        } catch (err) {
            if (!err.isAbort) {
                app.store.isLoadingSettings = false;
                app.checkApiError(err);
            }
        }
    },

    collections: [],
    collectionScaffolds: {},
    isLoadingCollections: false,
    _activeCollectionIdOrName: "",
    get activeCollection() {
        const idOrName = app.store._activeCollectionIdOrName;
        return app.store.collections.find((c) => c.id == idOrName || c.name == idOrName) || app.store.collections[0];
    },
    set activeCollection(collection) {
        if (typeof collection == "string") {
            app.store._activeCollectionIdOrName = collection;
        } else {
            app.store._activeCollectionIdOrName = collection?.id;
        }
    },
    async silentlyReloadCollections() {
        try {
            let newCollections = await app.pb.collections.getFullList({
                requestKey: "appStore.silentlyReloadCollections",
            });
            newCollections = app.utils.sortedCollectionsByType(newCollections);

            if (JSON.stringify(newCollections) != JSON.stringify(app.store.collections)) {
                app.store.collections = newCollections;
            }
        } catch (err) {
            if (!err.isAbort) {
                console.warn("failed to reload app store collections:", err);
            }
        }
    },
    async loadCollections(activeIdOrName = null) {
        app.store.isLoadingCollections = true;

        try {
            let [resultScaffolds, resultCollections] = await Promise.all([
                app.pb.collections.getScaffolds({ requestKey: "appStore.loadCollections.getScaffolds" }),
                app.pb.collections.getFullList({ requestKey: "appStore.loadCollections.getFullList" }),
            ]);

            resultCollections = app.utils.sortedCollectionsByType(resultCollections);

            // replace only if there are changes to minimize flickering
            if (JSON.stringify(app.store.collections) != JSON.stringify(resultCollections)) {
                app.store.collections = resultCollections;
            }

            app.store.collectionScaffolds = resultScaffolds;
            app.store._activeCollectionIdOrName = activeIdOrName || app.store._activeCollectionIdOrName
                || app.store.collections[0]?.id || "";
            app.store.isLoadingCollections = false;
        } catch (err) {
            if (!err.isAbort) {
                app.store.isLoadingCollections = false;
                app.checkApiError(err);
            }
        }
    },
    addOrUpdateCollection(collection) {
        const index = app.store.collections.findIndex((c) => c.id == collection.id);
        if (index >= 0) {
            if (app.store.activeCollection.id == collection.id) {
                app.store._activeCollectionIdOrName = collection.id;
            }

            app.store.collections[index] = collection;
        } else {
            app.store.collections.push(collection);
        }

        app.store.collections = app.utils.sortedCollectionsByType(app.store.collections);
    },

    oauth2Providers: [],
    isLoadingOAuth2Providers: false,
    async loadOAuth2Providers() {
        app.store.isLoadingOAuth2Providers = true;

        try {
            // @todo replace with SDK call
            app.store.oauth2Providers = await app.pb.send("/api/collections/meta/oauth2-providers");
            app.store.isLoadingOAuth2Providers = false;
        } catch (err) {
            if (!err.isAbort) {
                app.checkApiError(err);
                app.store.isLoadingOAuth2Providers = false;
            }
        }
    },
});

// reset title and errors on route change
window.addEventListener("hashchange", () => {
    app.store.title = "";
    app.store.errors = null;
});

// append the app settings name to document.title.
watch(() => {
    let titleParts = app.utils.toArray(app.store.title);

    const appName = app.store.settings?.meta?.appName || "";
    if (appName) {
        titleParts.push(appName);
    }

    document.title = titleParts.join(" - ");
});

// sync <meta="theme-color"> with the accent color
let metaThemeColor;
watch(() => app.store.settings?.meta?.accentColor, (newColor) => {
    if (!metaThemeColor) {
        metaThemeColor = t.meta({ name: "theme-color" });
        document.head.appendChild(metaThemeColor);
    }

    if (newColor) {
        metaThemeColor?.setAttribute("content", newColor);
        document.documentElement.style.setProperty("--accentColor", newColor);
    } else {
        metaThemeColor?.removeAttribute("content");
        document.documentElement.style.removeProperty("--accentColor");
    }
});

// sync favicon
let linkFavicon;
watch(() => app.store.favicon, (favicon) => {
    if (!linkFavicon) {
        linkFavicon = t.link({ rel: "icon" });
        document.head.appendChild(linkFavicon);
    }

    if (favicon) {
        linkFavicon.href = favicon;
    } else {
        linkFavicon.href = window.location.href.startsWith("https://")
            ? "./images/favicon_prod.png"
            : "./images/favicon.png";
    }
});

// sync color scheme
const colorSchemeMedia = window.matchMedia("(prefers-color-scheme: dark)");
app.store._mediaColorScheme = colorSchemeMedia.matches ? "dark" : "light";
colorSchemeMedia.addEventListener("change", ({ matches }) => {
    app.store._mediaColorScheme = matches ? "dark" : "light";
});
watch(() => app.store.userColorScheme, (colorScheme) => {
    if (!colorScheme) {
        window.localStorage.removeItem(COLOR_SCHEME_STORAGE_KEY);
    } else {
        window.localStorage.setItem(COLOR_SCHEME_STORAGE_KEY, colorScheme);
    }

    notifyChannel?.postMessage({ colorScheme });
});

// temporary disable animations on color scheme change to minimize flickering
let tempNoAnimationTimeoutId;
watch(() => app.store.activeColorScheme, (colorScheme) => {
    clearTimeout(tempNoAnimationTimeoutId);
    document.documentElement.style.setProperty("--animationSpeed", "0");

    document.documentElement.setAttribute("data-color-scheme", colorScheme);

    // restore animation
    tempNoAnimationTimeoutId = setTimeout(() => {
        document.documentElement.style.removeProperty("--animationSpeed");
    }, 100);
});

// Errors handler
// -------------------------------------------------------------------

function removeErrorState(input, container) {
    if (input.__errListener) {
        input.removeEventListener("input", input.__errListener);
        input.removeEventListener("change", input.__errListener);
        input.__errListener = null;
    }

    if (input.setCustomValidity) {
        input.setCustomValidity("");
        if (input._oldTitle) {
            input.setAttribute("title", input._oldTitle);
        } else {
            input.removeAttribute("title");
        }
    }

    input.removeAttribute("data-error");

    const helpElem = container.nextSibling;
    if (
        helpElem
        && helpElem.classList?.contains("generated-error")
        // remove only the error help text related to the input
        && helpElem.getAttribute("data-input-name") == input.getAttribute("name")
    ) {
        helpElem.remove();
    }

    // no other error inputs
    if (!container.querySelector("[data-error]")) {
        container.classList.remove("error");
    }
}

watch(
    () => JSON.stringify(app.store.errors) && app.store.errors,
    (errs) => {
        // search for input or other elements with "name" attribute
        const inputs = document.querySelectorAll(`[name]`);

        for (let input of inputs) {
            if (input.classList.contains("no-error")) {
                continue;
            }

            // find the top-most wrapper field element
            const container = input.closest(".field-list") || input.closest(".fields") || input.closest(".field");
            if (!container) {
                continue;
            }

            const name = input.getAttribute("name");

            removeErrorState(input, container);

            const errMsg = app.utils.getByPath(errs, name)?.message;
            if (!errMsg) {
                continue;
            }

            container.classList.add("error");

            input.__errListener = function() {
                removeErrorState(input, container);
                app.utils.deleteByPath(app.store.errors, name);
            };
            input.addEventListener("input", input.__errListener);
            input.addEventListener("change", input.__errListener);
            input.setAttribute("data-error", true);

            if (input.setCustomValidity && input.reportValidity && input.classList.contains("inline-error")) {
                input.setCustomValidity(errMsg);
                input.reportValidity();

                input._oldTitle = input.title;
                input.title = errMsg;
            } else {
                container.after(
                    t.div({
                        "html-data-input-name": name,
                        "className": "field-help error generated-error",
                        "textContent": errMsg,
                    }),
                );
            }
        }
    },
);

// Tabs sync
// -------------------------------------------------------------------

notifyChannel.onmessage = (e) => {
    if (
        e.data?.collections
        // replace only if there are changes to minimize flickering
        && JSON.stringify(app.store.collections) != JSON.stringify(e.data.collections)
    ) {
        app.store.collections = e.data.collections;
    }

    if (
        e.data?.settings
        // replace only if there are changes to minimize flickering
        && JSON.stringify(app.store.settings) != JSON.stringify(e.data.settings)
    ) {
        app.store.settings = e.data.settings;
    }

    if (e.data?.colorScheme) {
        app.store.userColorScheme = e.data.colorScheme;
    }
};

watch(
    () => JSON.stringify(app.store.collections),
    (newHash, oldHash) => {
        if (newHash && newHash != "[]" && oldHash && oldHash != "[]" && newHash != oldHash) {
            notifyChannel?.postMessage({
                collections: JSON.parse(newHash),
            });
        }
    },
);

watch(
    () => JSON.stringify(app.store.settings),
    (newHash, oldHash) => {
        if (newHash && newHash != "{}" && oldHash && oldHash != "{}" && newHash != oldHash) {
            notifyChannel?.postMessage({
                settings: JSON.parse(newHash),
            });
        }

        window.localStorage.setItem(SETTINGS_STORAGE_KEY, newHash);
    },
);
