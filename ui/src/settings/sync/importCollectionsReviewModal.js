import { collectionsDiffTable } from "./collectionsDiffTable";

window.app = window.app || {};
window.app.modals = window.app.modals || {};

window.app.modals.openImportCollectionsReview = function(oldCollections, newCollections, settings = {
    deleteMissing: false,
    onsubmit: null,
}) {
    const modal = importCollectionsModal(oldCollections, newCollections, settings);
    if (!modal) {
        return;
    }

    document.body.appendChild(modal);

    app.modals.open(modal);
};

function importCollectionsModal(oldCollections, newCollections, settingsArg) {
    let modal;

    const settings = store({
        deleteMissing: false,
        onsubmit: function(newCollections) {},
    });

    const watchers = app.utils.extendStore(settings, settingsArg);

    const data = store({
        isImporting: false,
        pairs: [],
    });

    function loadPairs() {
        const pairs = [];

        // add modified and deleted (if deleteMissing is set)
        for (const oldCollection of oldCollections) {
            const newCollection = newCollections.find((c) => c.id == oldCollection.id);
            if (
                (settings.deleteMissing && !newCollection?.id)
                || (newCollection?.id
                    && app.utils.hasCollectionChanges(oldCollection, newCollection, settings.deleteMissing))
            ) {
                pairs.push({
                    old: oldCollection,
                    new: newCollection,
                });
            }
        }

        // add only new collections
        for (const newCollection of newCollections) {
            const oldCollection = oldCollections.find((c) => c.id == newCollection.id);
            if (!oldCollection?.id) {
                pairs.push({
                    old: oldCollection,
                    new: newCollection,
                });
            }
        }

        data.pairs = pairs;
    }

    function submitConfirm() {
        // find deleted fields
        const deletedFieldNames = [];
        if (settings.deleteMissing) {
            for (const old of oldCollections) {
                const imported = newCollections.find((c) => c.id == old.id);
                if (!imported) {
                    // add all fields
                    deletedFieldNames.push(old.name + ".*");
                } else {
                    // add only deleted fields
                    const fields = Array.isArray(old.fields) ? old.fields : [];
                    for (const field of fields) {
                        if (!imported.fields.find((f) => f.id == field.id)) {
                            deletedFieldNames.push(`${old.name}.${field.name} (${field.id})`);
                        }
                    }
                }
            }
        }

        if (deletedFieldNames.length) {
            app.modals.confirm(
                [
                    t.h6(
                        null,
                        "Do you really want to delete the following collection fields and their related records data:",
                    ),
                    t.ul(null, () => {
                        return deletedFieldNames.map((name) => {
                            return t.li(null, name);
                        });
                    }),
                ],
                () => submit(),
            );
        } else {
            submit();
        }
    }

    async function submit() {
        if (data.isImporting) {
            return;
        }

        data.isImporting = true;

        try {
            await app.pb.collections.import(newCollections, settings.deleteMissing);
            await app.store.loadCollections();
            settings.onsubmit?.(JSON.parse(JSON.stringify(app.store.collections)));
            app.toasts.success("Successfully imported collections configuration.");
        } catch (err) {
            app.checkApiError(err);
        }

        data.isImporting = false;

        app.modals.close(modal);
    }

    modal = t.div(
        {
            pbEvent: "importCollectionsReviewModal",
            className: "modal popup full import-collections-review-modal",
            onbeforeopen: () => {
                loadPairs();
            },
            onbeforeclose: () => {
                return !data.isImporting;
            },
            onafterclose: (el) => {
                el?.remove();
            },
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        t.header({ className: "modal-header" }, t.h5(null, "Side-by-side diff")),
        t.div({ className: "modal-content" }, () => {
            return data.pairs.map((pair) => {
                return collectionsDiffTable({
                    collectionA: pair.old,
                    collectionB: pair.new,
                    deleteMissing: settings.deleteMissing,
                });
            });
        }),
        t.footer(
            { className: "modal-footer" },
            t.button(
                {
                    type: "button",
                    className: "btn transparent m-r-auto",
                    disabled: () => data.isImporting,
                    onclick: () => app.modals.close(modal),
                },
                t.span({ className: "txt" }, "Close"),
            ),
            t.button(
                {
                    type: "button",
                    className: () => `btn expanded ${data.isImporting ? "loading" : ""}`,
                    disabled: () => data.isImporting,
                    onclick: () => submitConfirm(),
                },
                t.span({ className: "txt" }, "Confirm and import"),
            ),
        ),
    );

    return modal;
}
