<script>
    import { createEventDispatcher } from "svelte";

    const dispatch = createEventDispatcher();

    export let index;
    export let list = [];
    export let group = "default";
    export let disabled = false;

    let dragging = false;
    let dragover = false;

    function onDrag(event, i) {
        if (!event || disabled) {
            return;
        }

        dragging = true;

        event.dataTransfer.effectAllowed = "move";
        event.dataTransfer.dropEffect = "move";
        event.dataTransfer.setData(
            "text/plain",
            JSON.stringify({
                index: i,
                group: group,
            })
        );
    }

    function onDrop(event, target) {
        if (!event || disabled) {
            return;
        }

        dragover = false;
        dragging = false;

        event.dataTransfer.dropEffect = "move";

        let dragData = {};
        try {
            dragData = JSON.parse(event.dataTransfer.getData("text/plain"));
        } catch (_) {}

        if (dragData.group != group) {
            return; // different draggable group
        }

        const start = dragData.index << 0;

        if (start < target) {
            list.splice(target + 1, 0, list[start]);
            list.splice(start, 1);
        } else {
            list.splice(target, 0, list[start]);
            list.splice(start + 1, 1);
        }

        list = list;

        dispatch("sort", {
            oldIndex: start,
            newIndex: target,
            list: list,
        });
    }
</script>

<div
    draggable={!disabled}
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
