<script>
    import { onMount } from "svelte";

    let classes = "";
    export { classes as class }; // export reserved keyword

    let wrapper = null;
    let scrollClasses = "";
    let scrollTimeoutId = null;
    let observer;

    export function refresh() {
        if (!wrapper) {
            return;
        }

        clearTimeout(scrollTimeoutId);

        scrollTimeoutId = setTimeout(() => {
            const offsetWidth = wrapper.offsetWidth;
            const scrollWidth = wrapper.scrollWidth;

            if (scrollWidth - offsetWidth) {
                scrollClasses = "scrollable";
                if (wrapper.scrollLeft === 0) {
                    scrollClasses += " scroll-start";
                } else if (wrapper.scrollLeft + offsetWidth == scrollWidth) {
                    scrollClasses += " scroll-end";
                }
            } else {
                scrollClasses = "";
            }
        }, 100);
    }

    onMount(() => {
        refresh();

        observer = new MutationObserver(() => {
            refresh();
        });

        observer.observe(wrapper, {
            attributeFilter: ["width"],
            childList: true,
            subtree: true,
        });

        return () => {
            observer?.disconnect();
            clearTimeout(scrollTimeoutId);
        };
    });
</script>

<svelte:window on:resize={refresh} />

<div class="horizontal-scroller-wrapper">
    <slot name="before" />

    <div bind:this={wrapper} class="horizontal-scroller {classes} {scrollClasses}" on:scroll={refresh}>
        <slot />
    </div>

    <slot name="after" />
</div>

<style>
    .horizontal-scroller {
        width: auto;
        overflow-x: auto;
    }
    .horizontal-scroller-wrapper {
        position: relative;
    }
    :global(.horizontal-scroller-wrapper .columns-dropdown) {
        top: 40px;
        z-index: 100;
        max-height: 340px;
    }
</style>
