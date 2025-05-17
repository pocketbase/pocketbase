import AppleOptions from "@/components/collections/providers/AppleOptions.svelte";
import MicrosoftOptions from "@/components/collections/providers/MicrosoftOptions.svelte";
import OIDCOptions from "@/components/collections/providers/OIDCOptions.svelte";
import SelfHostedOptions from "@/components/collections/providers/SelfHostedOptions.svelte";
import ApiClient from "@/utils/ApiClient";
import CommonHelper from "@/utils/CommonHelper";

// Object list with all supported OAuth2 providers in the format:
// ```
// [ { key, title, logo, optionsComponent?, optionComponentProps? }, ... ]
// ```
//
// If `optionsComponent` is provided it will receive 2 parameters:
// - `key`    - the provider settings key (eg. "gitlabAuth")
// - `config` - the provider settings config that is currently being updated
// - any other prop from optionComponentProps
export async function loadOAuth2Providers() {
    const result = await ApiClient.send("/api/settings/list-oauth2-providers", {
        method: "GET",
        requestKey: CommonHelper.randomString(10),
    });
    const providers = result.providers;
    for (const entry of providers) {
        if (entry.key in optionsOverride) {
            Object.assign(entry, optionsOverride[entry.key]);
        }
    }
    return providers;
}

const optionsOverride = {
    apple: {
        optionsComponent: AppleOptions,
    },
    microsoft: {
        optionsComponent: MicrosoftOptions,
    },
    gitlab: {
        optionsComponent: SelfHostedOptions,
        optionsComponentProps: { title: "Self-hosted endpoints (optional)" },
    },
    gitea: {
        optionsComponent: SelfHostedOptions,
        optionsComponentProps: { title: "Self-hosted endpoints (optional)" },
    },
    mailcow: {
        optionsComponent: SelfHostedOptions,
        optionsComponentProps: { required: true },
    },
    oidc: {
        optionsComponent: OIDCOptions,
    },
    oidc2: {
        optionsComponent: OIDCOptions,
    },
    oidc3: {
        optionsComponent: OIDCOptions,
    },
};
