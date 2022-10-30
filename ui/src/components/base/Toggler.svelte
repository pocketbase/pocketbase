<script>
    import { onMount, createEventDispatcher } from "svelte";
    import { fly } from "svelte/transition";

    export let trigger = undefined;
    export let active = false;
    export let escClose = true;
    export let closableClass = "closable";
    let classes = "";
    export { classes as class }; // export reserved keyword

    let container = undefined;
    let activeTrigger = undefined;

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
    }

    export function show() {
        active = true;
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

    function handleOutsideClick(e) {
        if (active && !container?.contains(e.target) && !activeTrigger?.contains(e.target)) {
            hide();
        }
    }

    function handleEscPress(e) {
        if (active && escClose && e.code === "Escape") {
            e.preventDefault();
            hide();
        }
    }

    function handleFocusChange(e) {
        return handleOutsideClick(e);
    }

    function bindTrigger(newTrigger) {
        cleanup();

        activeTrigger = newTrigger || container?.parentNode;

        if (!activeTrigger) {
            return;
        }

        container?.addEventListener("click", handleClickToggle);
        activeTrigger.addEventListener("click", handleClickToggle);
        activeTrigger.addEventListener("keydown", handleKeydownToggle);
    }

    function cleanup() {
        if (!activeTrigger) {
            return;
        }

        container?.removeEventListener("click", handleClickToggle);
        activeTrigger.removeEventListener("click", handleClickToggle);
        activeTrigger.removeEventListener("keydown", handleKeydownToggle);
    }

    onMount(() => {
        bindTrigger();

        return () => cleanup();
    });
</script>

<svelte:window on:click={handleOutsideClick} on:keydown={handleEscPress} on:focusin={handleFocusChange} />

<div bind:this={container} class="toggler-container">
    {#if active}
        <div
            class={classes}
            class:active
            in:fly|local={{ duration: 150, y: -5 }}
            out:fly|local={{ duration: 150, y: 2 }}
        >
            <slot />
        </div>
    {/if}
</div>
