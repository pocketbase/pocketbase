<script>
    import { onMount, createEventDispatcher } from "svelte";
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
    let isOutsideMouseDown = false;

    const dispatch = createEventDispatcher();

    $: if (container) {
        bindTrigger(trigger);
    }

    $: if (active) {
        activeTrigger?.classList?.add("active");
        dispatch("show");
    } else {
        activeTrigger?.classList?.remove("active");
        dispatch("hide");
    }

    export function hide() {
        active = false;
        isOutsideMouseDown = false;
        clearTimeout(scrollTimeoutId);
    }

    export function show() {
        active = true;

        clearTimeout(scrollTimeoutId);
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
            // is the trigger itself (or a direct child)
            (activeTrigger?.contains(elem) && !container.contains(elem)) ||
            // is closable toggler child
            (container.contains(elem) && elem.closest && elem.closest("." + closableClass))
        );
    }

    function handleClickToggle(e) {
        if (!active || isClosable(e.target)) {
            e.preventDefault();
            e.stopPropagation();

            toggle();
        }
    }

    function handleKeydownToggle(e) {
        if (
            (e.code === "Enter" || e.code === "Space") && // enter or spacebar
            (!active || isClosable(e.target))
        ) {
            e.preventDefault();
            e.stopPropagation();
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
        if (active && !container?.contains(e.target)) {
            isOutsideMouseDown = true;
        } else if (isOutsideMouseDown) {
            isOutsideMouseDown = false;
        }
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

    function handleFocusChange(e) {
        handleOutsideMousedown(e);
        handleOutsideClick(e);
    }

    function bindTrigger(newTrigger) {
        cleanup();

        container?.addEventListener("click", handleClickToggle);

        activeTrigger = newTrigger || container?.parentNode;
        activeTrigger?.addEventListener("click", handleClickToggle);
        activeTrigger?.addEventListener("keydown", handleKeydownToggle);
    }

    function cleanup() {
        clearTimeout(scrollTimeoutId);

        container?.removeEventListener("click", handleClickToggle);
        activeTrigger?.removeEventListener("click", handleClickToggle);
        activeTrigger?.removeEventListener("keydown", handleKeydownToggle);
    }

    onMount(() => {
        bindTrigger();

        return () => cleanup();
    });
</script>

<svelte:window
    on:mousedown={handleOutsideMousedown}
    on:click={handleOutsideClick}
    on:keydown={handleEscPress}
    on:focusin={handleFocusChange}
/>

<div bind:this={container} class="toggler-container" tabindex="-1">
    {#if active}
        <div
            bind:this={containerChild}
            class={classes}
            class:active
            transition:fly|local={{ duration: 150, y: 3 }}
        >
            <slot />
        </div>
    {/if}
</div>
