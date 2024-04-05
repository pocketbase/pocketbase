<script>
    import { createEventDispatcher, onMount } from "svelte";
    import { fly } from "svelte/transition";

    export let trigger = undefined;
    export let active = false;
    export let escClose = true;
    export let autoScroll = true;
    export let closableClass = "closable";
    let classes = "";
    export { classes as class }; // export reserved keyword

    let container;
    let containerChild;
    let activeTrigger;
    let scrollTimeoutId;
    let hideTimeoutId;
    let isOutsideMouseDown = false;

    const dispatch = createEventDispatcher();

    $: if (container) {
        bindTrigger(trigger);
    }

    $: if (active) {
        activeTrigger?.classList?.add("active");
        activeTrigger?.setAttribute("aria-expanded", true);
        dispatch("show");
    } else {
        activeTrigger?.classList?.remove("active");
        activeTrigger?.setAttribute("aria-expanded", false);
        dispatch("hide");
    }

    export function hideWithDelay(delay = 0) {
        if (!active) {
            return;
        }

        clearTimeout(hideTimeoutId);
        hideTimeoutId = setTimeout(hide, delay);
    }

    export function hide() {
        if (!active) {
            return; // already hidden
        }

        active = false;
        isOutsideMouseDown = false;
        clearTimeout(scrollTimeoutId);
        clearTimeout(hideTimeoutId);
    }

    export function show() {
        clearTimeout(hideTimeoutId);
        clearTimeout(scrollTimeoutId);

        if (active) {
            return; // already active
        }

        active = true;

        // focus toggler container not nested into the trigger
        if (!activeTrigger?.contains(container)) {
            container?.focus();
        }

        scrollTimeoutId = setTimeout(() => {
            if (!autoScroll) {
                return;
            }

            if (containerChild?.scrollIntoViewIfNeeded) {
                containerChild?.scrollIntoViewIfNeeded();
            } else if (containerChild?.scrollIntoView) {
                containerChild?.scrollIntoView({
                    behavior: "smooth",
                    block: "nearest",
                });
            }
        }, 180);
    }

    export function toggle() {
        if (active) {
            hide();
        } else {
            show();
        }
    }

    function isClosable(elem) {
        return (
            !container ||
            elem.classList.contains(closableClass) ||
            (container.contains(elem) && elem.closest && elem.closest("." + closableClass))
        );
    }

    function bindTrigger(newTrigger) {
        cleanup();

        container?.addEventListener("click", handleContainerClick);
        container?.addEventListener("keydown", handleContainerKeydown);

        activeTrigger = newTrigger || container?.parentNode;
        activeTrigger?.addEventListener("click", handleTriggerClick);
        activeTrigger?.addEventListener("keydown", handleTriggerKeydown);
    }

    function cleanup() {
        clearTimeout(scrollTimeoutId);
        clearTimeout(hideTimeoutId);

        container?.removeEventListener("click", handleContainerClick);
        container?.removeEventListener("keydown", handleContainerKeydown);

        activeTrigger?.removeEventListener("click", handleTriggerClick);
        activeTrigger?.removeEventListener("keydown", handleTriggerKeydown);
    }

    // toggler container handlers
    // ---------------------------------------------------------------
    function handleContainerClick(e) {
        e.stopPropagation(); // prevents firing the trigger click event in case it is nested

        if (isClosable(e.target)) {
            hide();
        }
    }

    function handleContainerKeydown(e) {
        if (e.code === "Enter" || e.code === "Space") {
            e.stopPropagation(); // prevents firing the trigger keydown event in case it is nested

            if (isClosable(e.target)) {
                // hide with a short delay since the button on:click events
                // doesn't fire if the element is not visible
                hideWithDelay(150);
            }
        }
    }

    // trigger handlers
    // ---------------------------------------------------------------
    function handleTriggerClick(e) {
        e.preventDefault();
        e.stopPropagation();
        toggle();
    }

    function handleTriggerKeydown(e) {
        if (e.code === "Enter" || e.code === "Space") {
            e.preventDefault();
            e.stopPropagation();
            toggle();
        }
    }

    function handleFocusChange(e) {
        if (active && !activeTrigger?.contains(e.target) && !container?.contains(e.target)) {
            toggle();
        }
    }

    function handleEscPress(e) {
        if (active && escClose && e.code === "Escape") {
            e.preventDefault();
            hide();
        }
    }

    function handleOutsideMousedown(e) {
        if (!active) {
            return;
        }

        isOutsideMouseDown = !container?.contains(e.target);
    }

    function handleOutsideClick(e) {
        if (
            active &&
            isOutsideMouseDown &&
            !container?.contains(e.target) &&
            !activeTrigger?.contains(e.target) &&
            !e.target?.closest(".flatpickr-calendar")
        ) {
            hide();
        }
    }

    onMount(() => {
        bindTrigger();

        return () => cleanup();
    });
</script>

<svelte:window
    on:click={handleOutsideClick}
    on:mousedown={handleOutsideMousedown}
    on:keydown={handleEscPress}
    on:focusin={handleFocusChange}
/>

<div bind:this={container} class="toggler-container" tabindex="-1" role="menu">
    {#if active}
        <div bind:this={containerChild} class={classes} class:active transition:fly={{ duration: 150, y: 3 }}>
            <slot />
        </div>
    {/if}
</div>
