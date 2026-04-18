import L from "leaflet";
import "leaflet/dist/leaflet.css";

// manually load the markers so that they can be embedded in the prod bundle
import markerIconRetinaUrl from "leaflet/dist/images/marker-icon-2x.png";
import markerIconUrl from "leaflet/dist/images/marker-icon.png";
import markerShadowUrl from "leaflet/dist/images/marker-shadow.png";

const defaultZoomLevel = 8;

window.app = window.app || {};
window.app.components = window.app.components || {};

/**
 * Leaflet component for showing and adjust a single geo point value.
 *
 * @param  {Object} [propsArg]
 * @return {Element}
 */
window.app.components.leaflet = function(propsArg = {}) {
    const props = store({
        rid: undefined,
        id: undefined,
        hidden: undefined,
        inert: undefined,
        className: "",
        point: { lat: 0, lon: 0 },
        onchange: function(point) {},
    });

    const watchers = app.utils.extendStore(props, propsArg);

    let map;
    let marker;
    let panTimeoutId;

    watchers.push(
        watch(
            () => {
                if (props.point.lat > 90) {
                    props.point.lat = 90;
                }
                if (props.point.lat < -90) {
                    props.point.lat = -90;
                }

                if (props.point.lon > 180) {
                    props.point.lon = 180;
                }
                if (props.point.lon < -180) {
                    props.point.lon = -180;
                }
            },
            () => {
                panInside();
            },
        ),
    );

    function panInside(debounce = 200) {
        if (!map) {
            return;
        }

        clearTimeout(panTimeoutId);
        panTimeoutId = setTimeout(() => {
            marker?.setLatLng([props.point.lat, props.point.lon]);
            map?.panInside([props.point.lat, props.point.lon], { padding: [20, 40] });
        }, debounce);
    }

    function initMap(mapEl) {
        const latlon = [toFixedCoord(props.point.lat), toFixedCoord(props.point.lon)];

        map = L.map(mapEl, { zoomControl: false }).setView(latlon, defaultZoomLevel);

        L.tileLayer("https://tile.openstreetmap.org/{z}/{x}/{y}.png", {
            attribution: "&copy; <a href=\"https://www.openstreetmap.org/copyright\">OpenStreetMap</a>",
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
        clearTimeout(panTimeoutId);
        marker?.remove();
        map?.remove();
    }

    function select(lat, lon, centerMap = true) {
        const point = {
            lat: toFixedCoord(lat),
            lon: toFixedCoord(lon),
        };

        if (props.onchange && props.onchange(point) === false) {
            return;
        }

        props.point = point;

        // center the map
        if (centerMap) {
            marker?.setLatLng([props.point.lat, props.point.lon]); // optimistic marker update
            map?.panTo([props.point.lat, props.point.lon], { animate: false });
        }

        map.getContainer()?.dispatchEvent(new CustomEvent("change", { detail: point }));

        resetSearch?.();
    }

    const [searchEl, resetSearch] = initSearch(select);

    return t.div(
        {
            rid: props.rid,
            id: () => props.id,
            hidden: () => props.hidden,
            inert: () => props.inert,
            className: "map-container",
            onunmount: () => {
                watchers.forEach((w) => w?.unwatch());
            },
        },
        searchEl,
        t.div({
            className: "map-box",
            onmount: (el) => {
                initMap(el);
            },
            onunmount: () => {
                destroyMap();
            },
        }),
    );
};

function toFixedCoord(coord) {
    return +(+coord).toFixed(6) || 0;
}

function initSearch(selectFunc = null) {
    const data = store({
        searchTerm: "",
        isSearching: false,
        searchResults: [],
    });

    let searchTimeoutId;
    let searchAbortController;

    function reset() {
        searchAbortController?.abort("reset");
        clearTimeout(searchTimeoutId);

        data.isSearching = false;
        data.searchResults = [];
        data.searchTerm = "";
    }

    // note: using debounce > 1s to minimize hitting the API rate limits
    // (see also https://operations.osmfoundation.org/policies/nominatim/)
    function search(debounce = 1100) {
        clearTimeout(searchTimeoutId);
        searchAbortController?.abort("search debounce");

        data.isSearching = true;
        data.searchResults = [];

        if (!data.searchTerm) {
            data.isSearching = false;
            return;
        }

        searchTimeoutId = setTimeout(async () => {
            try {
                searchAbortController = new AbortController();

                const response = await fetch(
                    "https://nominatim.openstreetmap.org/search.php?format=jsonv2&q="
                        + encodeURIComponent(data.searchTerm),
                    { signal: searchAbortController.signal },
                );
                if (response.status != 200) {
                    throw new Error("OpenStreetMap API error " + response.status);
                }

                const results = [];

                const addresses = await response.json();
                for (const item of addresses) {
                    results.push({
                        lat: item.lat,
                        lon: item.lon,
                        name: item.display_name,
                    });
                }

                data.searchResults = results;
            } catch (err) {
                console.warn("[address search failed]", err);
            }

            data.isSearching = false;
        }, debounce);
    }

    const searchInput = t.div(
        { className: "fields" },
        t.div(
            { className: "field" },
            t.input({
                type: "text",
                placeholder: "Search address...",
                value: () => data.searchTerm,
                oninput: (e) => (data.searchTerm = e.target.value),
            }),
        ),
        t.div({ className: "field addon p-l-10 p-r-10" }, () => {
            if (data.isSearching) {
                return t.span({ className: "loader sm" });
            }

            if (data.searchTerm.length) {
                return t.button(
                    {
                        className: "link-hint",
                        title: "Clear search",
                        onclick: () => reset(),
                    },
                    t.i({ className: "ri-close-line", ariaHidden: true }),
                );
            }
        }),
    );

    const searchDropdown = t.div({ className: "dropdown", popover: "manual" }, () => {
        return data.searchResults.map((item) => {
            return t.button(
                {
                    type: "button",
                    className: "dropdown-item",
                    title: "Select address coordinates",
                    onclick: () => selectFunc?.(item.lat, item.lon),
                },
                item.name,
            );
        });
    });

    const watchers = [];

    return [
        t.div(
            {
                className: "map-search",
                onmount: () => {
                    watchers.push(
                        watch(
                            () => data.searchTerm,
                            (searchTerm) => {
                                search(searchTerm);
                            },
                        ),
                    );

                    watchers.push(
                        watch(
                            () => data.searchResults,
                            (results) => {
                                if (results.length) {
                                    searchDropdown.showPopover({ source: searchInput });
                                } else {
                                    searchDropdown.hidePopover();
                                }
                            },
                        ),
                    );
                },
                onunmount: () => {
                    watchers.forEach((w) => w?.unwatch());
                    reset();
                },
            },
            searchInput,
            searchDropdown,
        ),
        reset,
    ];
}
