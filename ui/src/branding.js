export const productName = import.meta.env.PB_PRODUCT_NAME || "PocketBase";

function defaultLogoUrl() {
    const b = import.meta.env.BASE_URL ?? "";
    if (!b) {
        return "images/logo.svg";
    }
    return b.endsWith("/") ? `${b}images/logo.svg` : `${b}/images/logo.svg`;
}

const envLogo = import.meta.env.PB_LOGO_URL;
export const logoUrl =
    typeof envLogo === "string" && envLogo.trim() !== "" ? envLogo.trim() : defaultLogoUrl();
