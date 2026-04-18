import "./css/_main.css";

import "./utils";
import "./mimeTypes";
import "./store";
import "./pb";
import "./base/appHeader";
import "./base/autoAccordionOpenOnError";
import "./base/dropdownKeyboardNav";
import "./base/tooltip";
import "./base/confirm";
import "./base/dragline";
import "./base/slide";
import "./base/modal";
import "./base/toast";
import "./base/sortable";
import "./base/copyButton";
import "./base/codeBlock";
import "./base/codeEditor";
import "./base/codeBlockTabs";
import "./base/select";
import "./base/formattedDate";
import "./base/refreshButton";
import "./base/searchHistoryButton";
import "./base/s3Test";
import "./base/s3ConfigFields";
import "./base/leaflet";
import "./base/tinymce";
import "./base/colorPicker";
import "./base/ruleField";
import "./base/credits";
import "./base/searchbar";
import "./base/uploadedFileThumb";
import "./base/filePreviewModal";
import "./base/erd";
import "./base/pageSidebar";
import "./apiPreview/apiPreviewModal";
import "./settings/mail/mailTestModal";
import "./settings/sync/importCollectionsReviewModal";
import "./records/recordSummary";
import "./records/recordsSearchbar";
import "./records/recordFileThumb";
import "./records/recordFilePickerModal";
import "./records/recordsPickerModal";
import "./records/recordPreviewModal";
import "./records/recordImpersonateModal";
import "./records/recordUpsertModal";
import "./records/recordsList";
import "./base/fieldSettings";
import "./fields/text/init";
import "./fields/editor/init";
import "./fields/number/init";
import "./fields/bool/init";
import "./fields/email/init";
import "./fields/url/init";
import "./fields/date/init";
import "./fields/autodate/init";
import "./fields/file/init";
import "./fields/relation/init";
import "./fields/select/init";
import "./fields/json/init";
import "./fields/geoPoint/init";
import "./fields/password/init";
import "./collections/indexUpsertModal";
import "./collections/collectionUpsertModal";
import "./collections/addCollectionFieldButton";
import "./collections/autocomplete.utils";
import "./collections/providerPickerModal";
import "./collections/providerSettingsModal";
import "./collections/collectionChangesConfirmationModal";
import "./collections/collectionsOverviewModal";
import "./collections/oauth2/microsoftOptions";
import "./collections/oauth2/larkOptions";
import "./collections/oauth2/selfhostOptions";
import "./collections/oauth2/oidcOptions";
import "./collections/oauth2/appleOptions";
import "./logs/logsSettingsModal";
import "./logs/logPreviewModal";
import { appHeader } from "./base/appHeader";
import { initRouter } from "./router";

// tag proxy wrapper to register the global pbEvent mount:/unmount: events
// (events are intentionally not propagated to minimize performance issues)
const originalT = t;
t = new Proxy(
    {},
    {
        get(_, prop) {
            return function() {
                const config = arguments?.[0];
                if (config && config.pbEvent) {
                    const originalOnmount = config.onmount;
                    config.onmount = (el) => {
                        originalOnmount?.(el);
                        el.dataset.pb = config.pbEvent;
                        document.dispatchEvent(new CustomEvent("mount:" + config.pbEvent, { detail: el }));
                    };

                    const originalOnunmount = config.onunmount;
                    config.onunmount = (el) => {
                        document.dispatchEvent(new CustomEvent("unmount:" + config.pbEvent, { detail: el }));
                        originalOnunmount?.(el);
                    };
                }
                return originalT[prop](...arguments);
            };
        },
    },
);

document.body.prepend(t.main(
    {
        // @todo temp workaround to minimize onmount MutationObserver render flickering
        "html-class": "app",
        className: () => `app ${app.store.settings?.meta?.hideControls ? "hide-controls" : ""}`,
    },
    appHeader(),
    (el) => {
        if (typeof app.store.page == "function") {
            app.store.page(el);
        }

        return app.store.page;
    },
));

watch(
    () => app.store._ready,
    (isReady) => {
        if (isReady) {
            initRouter();
        }
    },
);

// load extensions
let extensionsScript = t.script({
    type: "module",
    src: app.pb.buildURL("/_/extensions.js"),
    onload: () => {
        app.store._ready = true;
    },
    onerror: (err) => {
        console.warn("Failed to load extensions:", err);
        app.store._ready = true;
    },
});
document.body.appendChild(extensionsScript);
