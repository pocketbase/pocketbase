<script>
    import Select from "@/components/base/Select.svelte";
    import ApiClient from "@/utils/ApiClient";
    import { confirm } from "@/stores/confirmation";
    import { addSuccessToast, addErrorToast } from "@/stores/toasts";

    let sourceSelected = undefined;
    let destinationSelected = undefined;

    let items = ["beta", "staging", "production"];

    function migrateData() {
        if (sourceSelected == undefined || destinationSelected == undefined) {
            return;
        }

        // Destinations: beta, staging, production
        confirm(`Are you sure? This will overwrite all data in ${destinationSelected}`, () => {
            ApiClient.migrate(sourceSelected, destinationSelected)
                .then((resp) => {
                    if (resp.ok) {
                        addSuccessToast("Successfully migrated data");
                    } else {
                        addErrorToast("There was an error");
                    }
                })
                .catch((_) => addErrorToast("There was an error"));
        });
    }
</script>

<div class="panel top">
    <div class="content txt-xl m-b-base">
        <p>Migrations</p>
    </div>
    <div class="pathSelect">
        <div class="src">
            <p>Source</p>
            <Select bind:selected={sourceSelected} {items} />
        </div>
        <div class="dest">
            <p>Destination</p>
            <Select bind:selected={destinationSelected} {items} />
        </div>
    </div>

    <button
        class="btn btn-migrate"
        on:click={migrateData}
        disabled={sourceSelected === undefined ||
            destinationSelected === undefined ||
            sourceSelected === destinationSelected}>Migrate</button
    >
</div>

<style>
    .panel.top {
        margin-bottom: 20px;
    }

    .pathSelect {
        display: flex;
        justify-content: space-between;
        gap: 20px;
    }

    .pathSelect > div {
        width: 100%;
    }

    .btn-migrate {
        margin-top: 20px;
        width: 100%;
    }
</style>
