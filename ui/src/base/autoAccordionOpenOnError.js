// auto open accordions if there is an invalid item
document.addEventListener(
    "invalid",
    (e) => {
        const details = e.target.closest("details");
        if (details && !details.open && !e.target.closest("summary")) {
            details.open = true;

            // revalidate and show the error message
            e.target.reportValidity && e.target.reportValidity();
        }
    },
    true,
);
