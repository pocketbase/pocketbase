<script>
    export let item = {}; // {ext, mimeType}
</script>

<span class="txt">{item.ext || "N/A"}</span>
<small class="txt-hint">{item.mimeType}</small>
