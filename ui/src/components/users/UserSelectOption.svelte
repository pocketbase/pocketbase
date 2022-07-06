<script>
    import tooltip from "@/actions/tooltip";

    export let item = {}; // model
</script>

<i
    class="ri-information-line link-hint"
    use:tooltip={{ text: JSON.stringify(item, null, 2), position: "left", class: "code" }}
/>

<div class="content">
    <div class="block txt-ellipsis">{item.id}</div>
    <small class="block txt-hint txt-ellipsis">{item.email}</small>
</div>
