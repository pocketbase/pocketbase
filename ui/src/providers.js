import AppleOptions from "@/components/collections/providers/AppleOptions.svelte";
import LarkOptions from "@/components/collections/providers/LarkOptions.svelte";
import MicrosoftOptions from "@/components/collections/providers/MicrosoftOptions.svelte";
import OIDCOptions from "@/components/collections/providers/OIDCOptions.svelte";
import SelfHostedOptions from "@/components/collections/providers/SelfHostedOptions.svelte";

// @todo remove after allowing custom OAuth2 UI extendability
//
// Object list with all supported OAuth2 providers in the format:
// ```
// [ { key, title, logo, optionsComponent?, optionComponentProps? }, ... ]
// ```
//
// The logo images must be placed inside the /public/images/oauth2 directory.
//
// If `optionsComponent` is provided it will receive 2 parameters:
// - `key`    - the provider settings key (eg. "gitlabAuth")
// - `config` - the provider settings config that is currently being updated
// - any other prop from optionComponentProps
export default [
    {
        key: "apple",
        title: "Apple",
        logo: "apple.svg",
        optionsComponent: AppleOptions,
    },
    {
        key: "google",
        title: "Google",
        logo: "google.svg",
    },
    {
        key: "microsoft",
        title: "Microsoft",
        logo: "microsoft.svg",
        optionsComponent: MicrosoftOptions,
    },
    {
        key: "yandex",
        title: "Yandex",
        logo: "yandex.svg",
    },
    {
        key: "facebook",
        title: "Facebook",
        logo: "facebook.svg",
    },
    {
        key: "instagram2",
        title: "Instagram",
        logo: "instagram.svg",
    },
    {
        key: "github",
        title: "GitHub",
        logo: "github.svg",
    },
    {
        key: "gitlab",
        title: "GitLab",
        logo: "gitlab.svg",
        optionsComponent: SelfHostedOptions,
        optionsComponentProps: { title: "Self-hosted endpoints (optional)" },
    },
    {
        key: "bitbucket",
        title: "Bitbucket",
        logo: "bitbucket.svg",
    },
    {
        key: "gitee",
        title: "Gitee",
        logo: "gitee.svg",
    },
    {
        key: "gitea",
        title: "Gitea",
        logo: "gitea.svg",
        optionsComponent: SelfHostedOptions,
        optionsComponentProps: { title: "Self-hosted endpoints (optional)" },
    },
    {
        key: "discord",
        title: "Discord",
        logo: "discord.svg",
    },
    {
        key: "twitter",
        title: "X/Twitter",
        logo: "twitter.svg",
    },
    {
        key: "kakao",
        title: "Kakao",
        logo: "kakao.svg",
    },
    {
        key: "vk",
        title: "VK",
        logo: "vk.svg"
    },
    {
        key: "linear",
        title: "Linear",
        logo: "linear.svg",
    },
    {
        key:   "notion",
        title: "Notion",
        logo:  "notion.svg",
    },
    {
        key:   "monday",
        title: "monday.com",
        logo:  "monday.svg",
    },
    {
        key: "lark",
        title: "Lark",
        logo: "lark.svg",
        optionsComponent: LarkOptions,
    },
    {
        key:   "box",
        title: "Box",
        logo:  "box.svg",
    },
    {
        key: "spotify",
        title: "Spotify",
        logo: "spotify.svg",
    },
    {
        key: "trakt",
        title: "Trakt",
        logo: "trakt.svg",
    },
    {
        key: "twitch",
        title: "Twitch",
        logo: "twitch.svg",
    },
    {
        key: "patreon",
        title: "Patreon (v2)",
        logo: "patreon.svg"
    },
    {
        key: "strava",
        title: "Strava",
        logo: "strava.svg",
    },
    {
        key: "wakatime",
        title: "WakaTime",
        logo: "wakatime.svg",
    },
    {
        key: "livechat",
        title: "LiveChat",
        logo: "livechat.svg",
    },
    {
        key: "mailcow",
        title: "mailcow",
        logo: "mailcow.svg",
        optionsComponent: SelfHostedOptions,
        optionsComponentProps: { required: true },
    },
    {
        key: "planningcenter",
        title: "Planning Center",
        logo: "planningcenter.svg",
    },
    {
        key: "oidc",
        title: "OpenID Connect",
        logo: "oidc.svg",
        optionsComponent: OIDCOptions,
    },
    {
        key: "oidc2",
        title: "(2) OpenID Connect",
        logo: "oidc.svg",
        optionsComponent: OIDCOptions,
    },
    {
        key: "oidc3",
        title: "(3) OpenID Connect",
        logo: "oidc.svg",
        optionsComponent: OIDCOptions,
    },
];
