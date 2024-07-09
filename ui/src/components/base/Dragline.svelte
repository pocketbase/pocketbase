<script>
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    export let tolerance = 0;

    let elem;
    let startX = 0;
    let startY = 0;
    let shiftX = 0;
    let shiftY = 0;
    let dragStarted = false;

    function dragInit(e) {
        e.stopPropagation();

        startX = e.clientX;
        startY = e.clientY;
        shiftX = e.clientX - elem.offsetLeft;
        shiftY = e.clientY - elem.offsetTop;

        document.addEventListener("touchmove", onMove);
        document.addEventListener("mousemove", onMove);
        document.addEventListener("touchend", onStop);
        document.addEventListener("mouseup", onStop);
    }

    function onStop(e) {
        if (dragStarted) {
            e.preventDefault();

            dragStarted = false;

            elem.classList.remove("no-pointer-events");

            dispatch("dragstop", { event: e, elem: elem });
        }

        document.removeEventListener("touchmove", onMove);
        document.removeEventListener("mousemove", onMove);
        document.removeEventListener("touchend", onStop);
        document.removeEventListener("mouseup", onStop);
    }

    function onMove(e) {
        let diffX = e.clientX - startX;
        let diffY = e.clientY - startY;
        let left = e.clientX - shiftX;
        let top = e.clientY - shiftY;

        if (
            !dragStarted &&
            Math.abs(left - elem.offsetLeft) < tolerance &&
            Math.abs(top - elem.offsetTop) < tolerance
        ) {
            return;
        }

        e.preventDefault();

        if (!dragStarted) {
            dragStarted = true;
            elem.classList.add("no-pointer-events");
            dispatch("dragstart", { event: e, elem: elem });
        }

        dispatch("dragging", { event: e, elem: elem, diffX: diffX, diffY: diffY });
    }
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<span
    bind:this={elem}
    class="dragline"
    class:dragging={dragStarted}
    on:mousedown={(e) => {
        if (e.button == 0) {
            dragInit(e);
        }
    }}
    on:touchstart={dragInit}
/>

<style lang="scss">
    .dragline {
        position: relative;
        z-index: 101;
        left: 0;
        top: 0;
        height: 100%;
        width: 5px;
        padding: 0;
        margin: 0 -3px 0 -1px;
        background: none;
        cursor: ew-resize;
        box-sizing: content-box;
        user-select: none;
        transition: box-shadow var(--activeAnimationSpeed);
        box-shadow: inset 1px 0 0 0 var(--baseAlt2Color);
        &:hover,
        &.dragging {
            box-shadow: inset 3px 0 0 0 var(--baseAlt2Color);
        }
    }
</style>
