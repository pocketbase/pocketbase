const setUiTheme = (theme) => {
    switch (theme) {
        case "Light":
            localStorage.setItem("pb_ui_theme", "light");
            document.body.classList.remove("dark");
            break;
        case "Dark":
            localStorage.setItem("pb_ui_theme", "dark");
            document.body.classList.add("dark");
            break;
        default:
            localStorage.setItem("pb_ui_theme", "system");
            window.matchMedia("(prefers-color-scheme: light)").matches
                ? document.body.classList.remove("dark")
                : document.body.classList.add("dark");
            break;
    }
};

export default setUiTheme;
