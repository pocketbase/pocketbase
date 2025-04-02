<script>
    import tooltip from "@/actions/tooltip";
    import Field from "@/components/base/Field.svelte";
    import FieldLabel from "@/components/records/fields/FieldLabel.svelte";
    import { slide } from "svelte/transition";

    export let original;
    export let field;
    export let value = undefined;

    let mapComponent;
    let isMapComponentLoading = false;
    let isMapVisible = false;

    $: if (typeof value === "undefined") {
        value = { lat: 0, lon: 0 };
    }

    $: if (value) {
        normalize();
    }

    function normalize() {
        if (value.lat > 90) {
            value.lat = 90;
        }

        if (value.lat < -90) {
            value.lat = -90;
        }

        if (value.lon > 180) {
            value.lon = 180;
        }

        if (value.lon < -180) {
            value.lon = -180;
        }
    }

    function toggleMapVisibility() {
        if (isMapVisible) {
            hideMap();
        } else {
            showMap();
        }
    }

    function showMap() {
        loadMapComponent();
        isMapVisible = true;
    }

    function hideMap() {
        isMapVisible = false;
    }

    async function loadMapComponent() {
        if (mapComponent || isMapComponentLoading) {
            return; // already loaded or in the process
        }

        isMapComponentLoading = true;

        mapComponent = (await import("@/components/base/Leaflet.svelte")).default;

        isMapComponentLoading = false;
    }
</script>

<Field class="form-field form-field-list {field.required ? 'required' : ''}" name={field.name} let:uniqueId>
    <FieldLabel {uniqueId} {field} />

    <div class="list">
        <div class="list-item">
            <Field class="form-field form-field-inline m-0" let:uniqueId>
                <label for={uniqueId}>Longitude:</label>
                <input
                    type="number"
                    id={uniqueId}
                    required={field.required}
                    placeholder="0"
                    step="any"
                    min="-180"
                    max="180"
                    bind:value={value.lon}
                />
            </Field>
            <span class="separator"></span>
            <Field class="form-field form-field-inline m-0" let:uniqueId>
                <label for={uniqueId}>Latitude:</label>
                <input
                    type="number"
                    id={uniqueId}
                    required={field.required}
                    placeholder="0"
                    step="any"
                    min="-90"
                    max="90"
                    bind:value={value.lat}
                />
            </Field>
            <span class="separator"></span>
            <button
                type="button"
                class="btn btn-circle btn-sm btn-circle {isMapVisible
                    ? 'btn-secondary'
                    : 'btn-hint btn-transparent'}"
                aria-label="Toggle map"
                use:tooltip={"Toggle map"}
                on:click={toggleMapVisibility}
            >
                <i class="ri-map-2-line"></i>
            </button>
        </div>

        {#if isMapVisible}
            <div class="block" style="height:200px" transition:slide={{ duration: 150 }}>
                {#if isMapComponentLoading}
                    <div class="block txt-center p-base">
                        <span class="loader loader-sm"></span>
                    </div>
                {:else}
                    <svelte:component this={mapComponent} height={200} bind:point={value} />
                {/if}
            </div>
        {/if}
    </div>
</Field>

<style lang="scss">
    .list-item {
        padding: 5px 10px;
        min-height: 0;
        gap: 10px;
    }
    .separator {
        align-self: stretch;
        background: var(--baseAlt2Color);
        width: 1px;
        margin: -5px 0;
    }
</style>
