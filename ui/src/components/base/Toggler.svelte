<script>
    import { onMount, createEventDispatcher } from "svelte";
    import { fly } from "svelte/transition";

    export let trigger = undefined;
    export let active = false;
    export let escClose = true;
    export let closableClass = "closable";
    let classes = "";
    export { classes as class }; // export reserved keyword

    let container;

    const dispatch = createEventDispatcher();

    $: if (active) {
        trigger?.classList?.add("active");
        dispatch("show");
    } else {
        trigger?.classList?.remove("active");
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
            (trigger?.contains(elem) && !container.contains(elem)) ||
            // is closable toggler child
            (container.contains(elem) && elem.closest && elem.closest("." + closableClass))
        );
    }

    function handleClickToggle(e) {
        if (!active || isClosable(e.target)) {
            e.preventDefault();
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
        if (active && !container?.contains(e.target) && !trigger?.contains(e.target)) {
            hide();
        }
    }

    function handleEscPress(e) {
        if (active && escClose && e.code == "Escape") {
            e.preventDefault();
            hide();
        }
    }

    function handleFocusChange(e) {
        return handleOutsideClick(e);
    }

    onMount(() => {
        trigger = trigger || container.parentNode;

        trigger.addEventListener("click", handleClickToggle);
        trigger.addEventListener("keydown", handleKeydownToggle);

        return () => {
            trigger.removeEventListener("click", handleClickToggle);
            trigger.removeEventListener("keydown", handleKeydownToggle);
        };
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
