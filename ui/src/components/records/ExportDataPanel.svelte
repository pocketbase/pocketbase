<script>
    import OverlayPanel from "@/components/base/OverlayPanel.svelte";
    import Field from "@/components/base/Field.svelte";
    import ObjectSelect from "@/components/base/ObjectSelect.svelte";
    import ApiClient from "@/utils/ApiClient";
    import CommonHelper from "@/utils/CommonHelper";
    import { addErrorToast, addSuccessToast } from "@/stores/toasts";

    export let collection;
    export let filter = "";

    let exportPanel;
    let recordScope = "all";
    let selectedFields = {}; // Start with empty selection
    let isExporting = false;
    let previewRecords = [];
    let isLoadingPreview = false;

    const scopeOptions = [
        { label: "All records", value: "all" },
        { label: "Current filtered results", value: "filtered" }
    ];

    $: fields = collection?.fields || [];
    
    // Create combined list of all available fields (system + collection fields)
    // Filter out any collection fields that might have the same names as system fields
    $: allFields = [
        'id', 
        'created', 
        'updated', 
        ...fields.map(f => f.name).filter(name => !['id', 'created', 'updated'].includes(name))
    ];
    
    $: totalSelectedFields = Object.keys(selectedFields).length;
    $: areAllSelected = allFields.length && totalSelectedFields === allFields.length;

    // Update preview data when scope or filter changes (but not when fields change)
    $: if (collection?.id && (recordScope || filter !== undefined)) {
        loadPreviewData();
    }

    function initializeFieldSelection() {
        selectAll();
    }

    function selectAll() {
        selectedFields = {};
        
        // Select all available fields
        for (const fieldName of allFields) {
            selectedFields[fieldName] = true;
        }
        selectedFields = selectedFields; // trigger reactivity
    }

    function deselectAll() {
        selectedFields = {};
    }

    async function loadPreviewData() {
        if (!collection?.id) return;

        isLoadingPreview = true;

        try {
            const queryParams = {
                limit: 3, // Only fetch first 3 records for preview
            };

            // Apply filter if "filtered" scope is selected
            if (recordScope === "filtered" && filter) {
                queryParams.filter = filter;
            }

            previewRecords = await ApiClient.collection(collection.id).getList(1, 3, queryParams);
        } catch (err) {
            console.error("Preview load error:", err);
            previewRecords = [];
        }

        isLoadingPreview = false;
    }

    // Generate preview JSON based on selected fields and current data
    $: previewData = generatePreviewData(previewRecords, selectedFields);

    function generatePreviewData(records, fields) {
        if (!records?.items?.length) {
            return [];
        }

        const fieldsToInclude = Object.keys(fields).filter(f => fields[f]);
        
        if (fieldsToInclude.length === 0) {
            return [];
        }
        
        // Filter records to only include selected fields
        const filteredRecords = records.items.map(record => {
            const filteredRecord = {};
            for (const field of fieldsToInclude) {
                if (record.hasOwnProperty(field)) {
                    filteredRecord[field] = record[field];
                }
            }
            return filteredRecord;
        });

        // Add ellipsis if there are more records
        if (records.totalItems > 3) {
            filteredRecords.push("... and " + (records.totalItems - 3) + " more records");
        }

        return filteredRecords;
    }

    function toggleSelectAll() {
        if (areAllSelected) {
            deselectAll();
        } else {
            selectAll();
        }
    }

    function toggleSelectField(fieldName) {
        if (selectedFields[fieldName]) {
            delete selectedFields[fieldName];
        } else {
            selectedFields[fieldName] = true;
        }
        selectedFields = selectedFields; // trigger reactivity
    }

    function getFieldTitle(fieldName) {
        // System fields
        if (['id', 'created', 'updated'].includes(fieldName)) {
            return `${fieldName} (system field)`;
        }
        
        // Collection fields
        const field = fields.find(f => f.name === fieldName);
        return field ? `${fieldName} (${field.type})` : fieldName;
    }

    export function show() {
        // Reset and initialize when showing panel
        selectedFields = {}; // Start fresh
        
        // Use setTimeout to ensure allFields is ready, then select all
        setTimeout(() => {
            if (allFields.length > 0) {
                selectAll();
                loadPreviewData();
            }
        }, 0);
        
        return exportPanel?.show();
    }

    async function exportData() {
        if (isExporting || !collection?.id || totalSelectedFields === 0) {
            return;
        }

        isExporting = true;

        try {
            // Build query parameters
            const queryParams = {
                batch: 500, // Export in batches to handle large datasets
            };

            // Apply filter if "filtered" scope is selected
            if (recordScope === "filtered" && filter) {
                queryParams.filter = filter;
            }

            // Build fields parameter - only include selected fields
            const fieldsToInclude = Object.keys(selectedFields).filter(f => selectedFields[f]);
            if (fieldsToInclude.length > 0) {
                queryParams.fields = fieldsToInclude.join(',');
            }

            // Fetch records using superuser authentication (already authenticated)
            const records = await ApiClient.collection(collection.id).getFullList(queryParams);

            // Generate filename
            const timestamp = new Date().toISOString().slice(0, 19).replace(/[:-]/g, '');
            const scopeSuffix = recordScope === "filtered" ? "_filtered" : "_all";
            const filename = `${collection.name}${scopeSuffix}_${timestamp}`;

            // Download JSON file
            CommonHelper.downloadJson(records, filename);
            
            addSuccessToast(`Successfully exported ${records.length} records from ${collection.name}`);
            
            // Close panel after successful export
            hide();

        } catch (err) {
            console.error("Export error:", err);
            addErrorToast(`Failed to export data: ${err.message || 'Unknown error'}`);
        }

        isExporting = false;
    }

    export function hide() {
        return exportPanel?.hide();
    }
</script>

<OverlayPanel bind:this={exportPanel} class="overlay-panel-lg export-data-panel">
    <svelte:fragment slot="header">
        <h4>Export Data</h4>
    </svelte:fragment>

    <div class="content">
        <Field class="form-field" let:uniqueId>
            <label for={uniqueId}>Records to Export</label>
            <ObjectSelect
                id={uniqueId}
                items={scopeOptions}
                bind:keyOfSelected={recordScope}
            />
        </Field>

        <div class="export-panel">
            <div class="export-list">
                <div class="list-item list-item-section">
                    <Field class="form-field" let:uniqueId>
                        <input
                            type="checkbox"
                            id={uniqueId}
                            disabled={!allFields.length}
                            checked={areAllSelected}
                            on:change={() => toggleSelectAll()}
                        />
                        <label for={uniqueId}>Select all</label>
                    </Field>
                </div>

                <!-- All Fields (System + Collection) -->
                {#each allFields as fieldName (fieldName)}
                    <div class="list-item list-item-collection">
                        <Field class="form-field" let:uniqueId>
                            <input
                                type="checkbox"
                                id={uniqueId}
                                checked={selectedFields[fieldName]}
                                on:change={() => toggleSelectField(fieldName)}
                            />
                            <label for={uniqueId} title={getFieldTitle(fieldName)}>{fieldName}</label>
                        </Field>
                    </div>
                {/each}
            </div>

            <div class="export-preview">
                {#if isLoadingPreview}
                    <div class="loader" />
                {:else}
                    <pre class="code-wrapper">{JSON.stringify(previewData, null, 4)}</pre>
                {/if}
            </div>
        </div>
    </div>

    <svelte:fragment slot="footer">
        <button 
            type="button" 
            class="btn btn-transparent" 
            disabled={isExporting}
            on:click={() => hide()}
        >
            Cancel
        </button>
        <button 
            type="button" 
            class="btn btn-expanded" 
            disabled={isExporting || totalSelectedFields === 0}
            on:click={exportData}
        >
            {#if isExporting}
                <span class="loader loader-sm"></span>
                <span class="txt">Exporting...</span>
            {:else}
                <i class="ri-download-2-line" />
                <span class="txt">Export Data</span>
            {/if}
        </button>
    </svelte:fragment>
</OverlayPanel>
