<script>
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    export let index;
    export let list = [];
    export let disabled = false;

    let dragging = false;
    let dragover = false;

    function onDrag(event, i) {
        if (!event && !disabled) {
            return;
        }

        dragging = true;

        event.dataTransfer.effectAllowed = "move";
        event.dataTransfer.dropEffect = "move";
        event.dataTransfer.setData("text/plain", i);
    }

    function onDrop(event, target) {
        if (!event && !disabled) {
            return;
        }

        dragover = false;
        dragging = false;

        event.dataTransfer.dropEffect = "move";

        const start = parseInt(event.dataTransfer.getData("text/plain"));

        if (start < target) {
            list.splice(target + 1, 0, list[start]);
            list.splice(start, 1);
        } else {
            list.splice(target, 0, list[start]);
            list.splice(start + 1, 1);
        }

        list = list;

        dispatch("sort", list);
    }
</script>

<div
    draggable={true}
    class="draggable"
    class:dragging
    class:dragover
    on:dragover|preventDefault={() => {
        dragover = true;
    }}
    on:dragleave|preventDefault={() => {
        dragover = false;
    }}
    on:dragend={() => {
        dragover = false;
        dragging = false;
    }}
    on:dragstart={(e) => onDrag(e, index)}
    on:drop={(e) => onDrop(e, index)}
>
    <slot {dragging} {dragover} />
</div>

<style>
    .draggable {
        user-select: none;
        outline: 0;
        min-width: 0;
    }
</style>
