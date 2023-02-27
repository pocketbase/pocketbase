<script>
    let classes = "";
    export { classes as class }; // export reserved keyword

    export let name;
    export let sort = "";
    export let disable = false;

    function toggleSort() {
        if (disable) {
            return;
        }

        if ("-" + name === sort) {
            sort = "+" + name;
        } else {
            sort = "-" + name;
        }
    }
</script>

<th
    tabindex="0"
    title={name}
    class="col-sort {classes}"
    class:col-sort-disabled={disable}
    class:sort-active={sort === "-" + name || sort === "+" + name}
    class:sort-desc={sort === "-" + name}
    class:sort-asc={sort === "+" + name}
    on:click={() => toggleSort()}
    on:keydown={(e) => {
        if (e.code === "Enter" || e.code === "Space") {
            e.preventDefault();
            toggleSort();
        }
    }}
>
    <slot />
</th>
