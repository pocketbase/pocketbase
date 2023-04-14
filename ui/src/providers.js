import SelfHostedOptions from "@/components/settings/providers/SelfHostedOptions.svelte";
import MicrosoftOptions  from "@/components/settings/providers/MicrosoftOptions.svelte";
import OIDCOptions       from "@/components/settings/providers/OIDCOptions.svelte";
import AppleOptions      from "@/components/settings/providers/AppleOptions.svelte";

// Object list with all supported OAuth2 providers in the format:
// ```
// [ { key, title, logo, optionsComponent? }, ... ]
// ```
//
// The logo images must be placed inside the /public/images/oauth2 directory.
//
// If `optionsComponent` is provided it will receive 2 parameters:
// - `key`    - the provider settings key (eg. "gitlabAuth")
// - `config` - the provider settings config that is currently being updated
export default [
    {
        key:   "appleAuth",
        title: "Apple",
        logo:  "apple.svg",
        optionsComponent: AppleOptions,
    },
    {
        key:   "googleAuth",
        title: "Google",
        logo:  "google.svg",
    },
    {
        key:   "facebookAuth",
        title: "Facebook",
        logo:  "facebook.svg",
    },
    {
        key:   "microsoftAuth",
        title: "Microsoft",
        logo:  "microsoft.svg",
        optionsComponent: MicrosoftOptions,
    },
    {
        key:   "githubAuth",
        title: "GitHub",
        logo:  "github.svg",
    },
    {
        key:   "gitlabAuth",
        title: "GitLab",
        logo:  "gitlab.svg",
        optionsComponent: SelfHostedOptions,
    },
    {
        key:   "giteeAuth",
        title: "Gitee",
        logo:  "gitee.svg",
    },
    {
        key:   "giteaAuth",
        title: "Gitea",
        logo:  "gitea.svg",
        optionsComponent: SelfHostedOptions,
    },
    {
        key:   "discordAuth",
        title: "Discord",
        logo:  "discord.svg",
    },
    {
        key:   "twitterAuth",
        title: "Twitter",
        logo:  "twitter.svg",
    },
    {
        key:   "kakaoAuth",
        title: "Kakao",
        logo:  "kakao.svg",
    },
    {
        key:   "spotifyAuth",
        title: "Spotify",
        logo:  "spotify.svg",
    },
    {
        key:   "twitchAuth",
        title: "Twitch",
        logo:  "twitch.svg",
    },
    {
        key:   "stravaAuth",
        title: "Strava",
        logo:  "strava.svg",
    },
    {
        key:   "livechatAuth",
        title: "LiveChat",
        logo:  "livechat.svg",
    },
    {
        key:   "oidcAuth",
        title: "OpenID Connect",
        logo:  "oidc.svg",
        optionsComponent: OIDCOptions,
    },
    {
        key:   "oidc2Auth",
        title: "(2) OpenID Connect",
        logo:  "oidc.svg",
        optionsComponent: OIDCOptions,
    },
    {
        key:   "oidc3Auth",
        title: "(3) OpenID Connect",
        logo:  "oidc.svg",
        optionsComponent: OIDCOptions,
    },
];
