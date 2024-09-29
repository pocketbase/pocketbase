<script>
    import { createEventDispatcher, onMount } from "svelte";
    import { slide } from "svelte/transition";

    const dispatch = createEventDispatcher();

    let accordionElem;
    let expandTimeoutId;

    let classes = "";
    export { classes as class }; // export reserved keyword

    export let draggable = false;
    export let active = false;
    export let interactive = true;
    export let single = false; // ensures that only one accordion is expanded in its given parent container

    let isDragOver = false;

    $: if (active) {
        clearTimeout(expandTimeoutId);
        expandTimeoutId = setTimeout(() => {
            if (accordionElem?.scrollIntoViewIfNeeded) {
                accordionElem.scrollIntoViewIfNeeded();
            } else if (accordionElem?.scrollIntoView) {
                accordionElem.scrollIntoView({
                    behavior: "smooth",
                    block: "nearest",
                });
            }
        }, 200);
    }

    export function isExpanded() {
        return !!active;
    }

    export function expand() {
        collapseSiblings();
        active = true;
        dispatch("expand");
    }

    export function collapse() {
        active = false;
        clearTimeout(expandTimeoutId);
        dispatch("collapse");
    }

    export function toggle() {
        dispatch("toggle");

        if (active) {
            collapse();
        } else {
            expand();
        }
    }

    export function collapseSiblings() {
        if (single && accordionElem.closest(".accordions")) {
            const handlers = accordionElem
                .closest(".accordions")
                .querySelectorAll(".accordion.active .accordion-header.interactive");
            for (const handler of handlers) {
                handler.click(); // @todo consider using store or other more reliable approach
            }
        }
    }

    onMount(() => {
        return () => clearTimeout(expandTimeoutId);
    });
</script>

<div bind:this={accordionElem} class="accordion {isDragOver ? 'drag-over' : ''} {classes}" class:active>
    <button
        type="button"
        class="accordion-header"
        {draggable}
        class:interactive
        aria-expanded={active}
        on:click|preventDefault={() => interactive && toggle()}
        on:drop|preventDefault={(e) => {
            if (draggable) {
                isDragOver = false;
                collapseSiblings();
                dispatch("drop", e);
            }
        }}
        on:dragstart={(e) => draggable && dispatch("dragstart", e)}
        on:dragenter={(e) => {
            if (draggable) {
                isDragOver = true;
                dispatch("dragenter", e);
            }
        }}
        on:dragleave={(e) => {
            if (draggable) {
                isDragOver = false;
                dispatch("dragleave", e);
            }
        }}
        on:dragover|preventDefault
    >
        <slot name="header" {active} />
    </button>

    {#if active}
        <div class="accordion-content" transition:slide={{ delay: 10, duration: 150 }}>
            <slot />
        </div>
    {/if}
</div>
