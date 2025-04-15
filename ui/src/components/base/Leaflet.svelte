<script>
    import { onMount } from "svelte";
    import tooltip from "@/actions/tooltip";
    import L from "leaflet";
    import "leaflet/dist/leaflet.css";

    // manually load the markers so that they can be embedded in the prod bundle
    import markerIconUrl from "leaflet/dist/images/marker-icon.png";
    import markerIconRetinaUrl from "leaflet/dist/images/marker-icon-2x.png";
    import markerShadowUrl from "leaflet/dist/images/marker-shadow.png";

    export let height = 225;
    export let point = { lat: 0, lon: 0 };

    let map;
    let mapEl;
    let marker;
    let isSearching = false;
    let searchTerm = "";
    let searchResults = [];
    let searchTimeoutId;
    let searchAbortController;
    let panTimeoutId;

    const defaultZoomLevel = 8;

    $: search(searchTerm);

    $: if (point.lat && point.lon) {
        panInside();
    }

    function normalizeCoordinate(coord) {
        return +(+coord).toFixed(6);
    }

    function panInside(debounce = 200) {
        clearTimeout(panTimeoutId);
        panTimeoutId = setTimeout(() => {
            marker?.setLatLng([point.lat, point.lon]);
            map?.panInside([point.lat, point.lon], { padding: [20, 40] });
        }, debounce);
    }

    function initMap() {
        const latlon = [normalizeCoordinate(point.lat), normalizeCoordinate(point.lon)];

        map = L.map(mapEl, { zoomControl: false }).setView(latlon, defaultZoomLevel);

        L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
            attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a>',
        }).addTo(map);

        // reassign the default marker images with the loaded ones
        // (https://leafletjs.com/reference.html#icon-default-option)
        L.Icon.Default.prototype.options.iconUrl = markerIconUrl;
        L.Icon.Default.prototype.options.iconRetinaUrl = markerIconRetinaUrl;
        L.Icon.Default.prototype.options.shadowUrl = markerShadowUrl;
        L.Icon.Default.imagePath = "";

        marker = L.marker(latlon, {
            draggable: true,
            autoPan: true,
        }).addTo(map);

        marker.bindTooltip("drag or right click anywhere on the map to move");

        marker.on("moveend", (e) => {
            if (e.sourceTarget?._latlng) {
                select(e.sourceTarget._latlng.lat, e.sourceTarget._latlng.lng, false);
            }
        });

        map.on("contextmenu", (e) => {
            select(e.latlng.lat, e.latlng.lng, false);
        });
    }

    function destroyMap() {
        resetSearch();
        marker?.remove();
        map?.remove();
    }

    function resetSearch() {
        searchAbortController?.abort();
        clearTimeout(searchTimeoutId);
        isSearching = false;
        searchResults = [];
        searchTerm = "";
    }

    // note: using debounce > 1s to minimize hitting the API rate limits
    // (see also https://operations.osmfoundation.org/policies/nominatim/)
    function search(q, debounce = 1100) {
        isSearching = true;
        searchResults = [];
        clearTimeout(searchTimeoutId);
        searchAbortController?.abort();

        if (!q) {
            isSearching = false;
            return;
        }

        searchTimeoutId = setTimeout(async () => {
            searchAbortController = new AbortController();

            try {
                const response = await fetch(
                    "https://nominatim.openstreetmap.org/search.php?format=jsonv2&q=" + encodeURIComponent(q),
                    { signal: searchAbortController.signal },
                );
                if (response.status != 200) {
                    throw new Error("OpenStreetMap API error " + response.status);
                }

                const addresses = await response.json();
                for (const item of addresses) {
                    searchResults.push({
                        lat: item.lat,
                        lon: item.lon,
                        name: item.display_name,
                    });
                }
            } catch (err) {
                console.warn("[address search failed]", err);
            }

            searchResults = searchResults;
            isSearching = false;
        }, debounce);
    }

    function select(lat, lon, centerMap = true) {
        point.lat = normalizeCoordinate(lat);
        point.lon = normalizeCoordinate(lon);

        // center the map
        if (centerMap) {
            marker?.setLatLng([point.lat, point.lon]); // optimistic marker update
            map?.panTo([point.lat, point.lon], { animate: false });
        }

        resetSearch();
    }

    onMount(() => {
        initMap();

        return () => {
            destroyMap();
        };
    });
</script>

<div class="map-wrapper" style="{height ? `height:${height}px` : null};">
    <div class="map-search">
        <div class="form-field m-0">
            {#if isSearching}
                <div class="form-field-addon">
                    <span class="loader loader-xs"></span>
                </div>
            {:else if searchTerm.length}
                <div class="form-field-addon">
                    <button
                        type="button"
                        class="btn btn-circle btn-xs btn-transparent"
                        on:click={resetSearch}
                    >
                        <i class="ri-close-line"></i>
                    </button>
                </div>
            {/if}
            <input type="text" placeholder="Search address..." bind:value={searchTerm} />
        </div>
        {#if searchTerm.length && searchResults.length}
            <div class="dropdown dropdown-sm dropdown-block">
                {#each searchResults as result, i}
                    <button
                        type="button"
                        class="dropdown-item"
                        use:tooltip={"Select address coordinates"}
                        on:click={() => select(result.lat, result.lon)}
                    >
                        {result.name}
                    </button>
                {/each}
            </div>
        {/if}
    </div>
    <div bind:this={mapEl} class="map-box"></div>
</div>

<style lang="scss">
    .map-wrapper {
        position: relative;
        display: block;
        height: 100%;
        width: 100%;
    }
    .map-box {
        z-index: 1;
        height: 100%;
        width: 100%;
    }
    .map-search {
        position: absolute;
        z-index: 2;
        top: 10px;
        width: 70%;
        max-width: 400px;
        margin-left: 15%;
        height: auto;
        input {
            opacity: 0.7;
            background: var(--baseColor);
            border: 0;
            box-shadow: 0 0 3px 0 var(--shadowColor);
            transition: opacity var(--baseAnimationSpeed);
        }
        .dropdown {
            max-height: 150px;
            border: 0;
            box-shadow: 0 0 3px 0 var(--shadowColor);
        }
        &:focus-within {
            input {
                opacity: 1;
            }
        }
    }
</style>
