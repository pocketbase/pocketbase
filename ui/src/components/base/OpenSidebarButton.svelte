<script>
    import { showPageSidebar } from "@/stores/app";
    import tooltip from "@/actions/tooltip";

    let tooltipData = { text: "Open Sidebar", position: "right" };
</script>

<button
    type="button"
    aria-label="Open Sidebar"
    class="btn btn-transparent btn-circle btn-open-sidebar"
    use:tooltip={tooltipData}
    on:click={() => showPageSidebar.set(true)}
>
    <i class="ri-menu-unfold-line" />
</button>
