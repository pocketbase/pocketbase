const maxNestedLazyExpand = 10;
const lazyLoadBatchSize = 500;
const LAZY_EXPAND_EVENT_NAME = "pb:lazyExpandSummaryRels";

const recordsToLoad = {}; // collectionId: new Set(id1, id2, id3)
const queueTimeoutIds = {}; // collectionId: timeoutId

// {
//     record: undefined,
//     field:  undefined,
//     short:  false,
//     meta:   undefined,
// }
export function view(props) {
    let subsToRemove = new Set();

    return t.div(
        {
            className: "record-field-view field-type-relation",
            onunmount: () => {
                for (const sub of subsToRemove) {
                    document.removeEventListener(LAZY_EXPAND_EVENT_NAME, sub);
                }
                subsToRemove.clear();
                subsToRemove = null;
            },
        },
        () => {
            const ids = app.utils.toArray(props.record[props.field.name]);
            if (!ids.length) {
                return t.span({ className: "missing-value" });
            }

            // stop at cyclic references
            const meta = props.meta || {};
            let parents = app.utils.toArray(meta.parents);
            if (parents.includes(props.record.id)) {
                return t.span({ className: "marker recursive" }, "(recursive)");
            }

            const newMeta = JSON.parse(JSON.stringify(meta));
            newMeta.parents = parents.concat(props.record.id);

            const result = [];

            // truncate "full" view too to prevent freezing the browser tab
            const maxIndex = props.short ? 3 : 1000;

            const expanded = app.utils.toArray(props.record.expand?.[props.field.name]);

            for (let i = 0; i < ids.length; i++) {
                if (i >= maxIndex) {
                    result.push(t.span({ className: "marker more" }, "(", ids.length - maxIndex, " more)"));
                    break;
                }

                const id = ids[i];
                const rel = expanded.find((r) => r?.id == id);

                if (rel) {
                    result.push(app.components.recordSummary(rel, newMeta));
                } else {
                    result.push(
                        t.span(
                            { className: "label relation-id animate-delayed-fadeIn" },
                            app.components.copyButton(id),
                            id,
                        ),
                    );

                    // lazy expand
                    if (newMeta.parents.length < maxNestedLazyExpand) {
                        const preferredIndex = i;
                        recordsToLoad[props.field.collectionId] = recordsToLoad[props.field.collectionId] || new Set();
                        recordsToLoad[props.field.collectionId].add(id);
                        const sub = (e) => {
                            if (e.detail.id == id && e.detail.collectionId == props.field.collectionId) {
                                setExpand(props.record, props.field, structuredClone(e.detail), preferredIndex);

                                document.removeEventListener(LAZY_EXPAND_EVENT_NAME, sub);
                                subsToRemove.delete(sub);
                            }
                        };
                        subsToRemove.add(sub);
                        document.addEventListener(LAZY_EXPAND_EVENT_NAME, sub);
                    }
                }
            }

            fetchQueuedItems(props.field.collectionId);

            return result;
        },
    );
}

function setExpand(record, field, rel, preferredIndex = 0) {
    record.expand = record.expand || {};

    if (field.maxSelect > 1) {
        record.expand[field.name] = app.utils.toArray(record.expand[field.name]);

        const existingIndex = record.expand[field.name].findIndex((r) => r?.id == rel.id);
        if (existingIndex >= 0) {
            record.expand[field.name][existingIndex] = rel;
        } else if (!record.expand[field.name][preferredIndex]) {
            record.expand[field.name][preferredIndex] = rel;
        } else {
            record.expand[field.name].push(rel);
        }
    } else {
        record.expand[field.name] = rel;
    }
}

function fetchQueuedItems(collectionId) {
    if (!collectionId) {
        return;
    }

    clearTimeout(queueTimeoutIds[collectionId]);
    queueTimeoutIds[collectionId] = setTimeout(() => {
        const relIds = Array.from(recordsToLoad[collectionId] || []);
        if (!relIds.length) {
            return;
        }

        // clear without awaiting the fetch calls to allow other queue items to start loading
        //
        // it is OK if the same rel id is being requested multiple times,
        // since the event will be detached after the first call
        recordsToLoad[collectionId].clear();
        recordsToLoad[collectionId] = null;
        queueTimeoutIds[collectionId] = null;

        // split in multiple batches to minimize filter length errors
        while (relIds.length) {
            const ids = relIds.splice(0, lazyLoadBatchSize);

            // eagerly expand first level presentable relations (if any)
            let relExpands = [];
            const presentableRelationFields = app.store.collections
                .find((c) => c.id == collectionId)
                ?.fields?.filter((f) => !f.hidden && f.presentable && f.type == "relation") || [];
            for (let field of presentableRelationFields) {
                relExpands.push(field.name);
            }
            relExpands = relExpands.join(",") || undefined;

            const requestFields = fieldsWithExcerpt(collectionId, presentableRelationFields);

            let request;
            if (ids.length == 1) {
                request = app.pb.collection(collectionId).getOne(ids[0], {
                    requestKey: null,
                    expand: relExpands,
                    fields: requestFields,
                });
            } else {
                request = app.pb.collection(collectionId).getFullList({
                    requestKey: null,
                    expand: relExpands,
                    filter: ids.map((id) => app.pb.filter("id={:id}", { id })).join("||"),
                    fields: requestFields,
                });
            }

            request
                .then((expanded) => {
                    expanded = app.utils.toArray(expanded);
                    if (!expanded.length) {
                        return;
                    }

                    for (const item of expanded) {
                        document.dispatchEvent(
                            new CustomEvent(LAZY_EXPAND_EVENT_NAME, {
                                detail: item,
                            }),
                        );
                    }
                })
                .catch((err) => {
                    console.warn("failed to lazily expand presentable relation", err);
                });
        }
    }, 0);
}

// -------------------------------------------------------------------

export function fieldsWithExcerpt(collectionId, expandedRelFields = [], maxLength = 200) {
    const collection = app.store.collections.find((c) => c.id == collectionId);

    let requestFields = collection?.fields?.filter((f) => f.type == "editor")
        .map((f) => `${f.name}:excerpt(${maxLength},true)`) || [];

    for (const relField of expandedRelFields) {
        const excerptRelFields = app.store.collections?.find((c) => c.id == relField.collectionId)
            ?.fields?.filter((f) => f.type == "editor")?.map((f) =>
                `expand.${relField.name}.${f.name}:excerpt(${maxLength},true)`
            );
        if (excerptRelFields?.length) {
            requestFields.push(`expand.${relField.name}.*`);
            requestFields = requestFields.concat(excerptRelFields);
        }
    }

    if (requestFields.length > 0) {
        return ["*", "expand.*"].concat(requestFields).join(",");
    }

    return undefined;
}
