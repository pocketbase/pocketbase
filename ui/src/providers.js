import SelfHostedOptions from "@/components/settings/providers/SelfHostedOptions.svelte";
import MicrosoftOptions  from "@/components/settings/providers/MicrosoftOptions.svelte";
import OIDCOptions       from "@/components/settings/providers/OIDCOptions.svelte";
import AppleOptions      from "@/components/settings/providers/AppleOptions.svelte";

// Object list with all supported OAuth2 providers in the format:
// ```
// [ { key, title, logo, optionsComponent? }, ... ]
// ```
//
// If `optionsComponent` is provided it will receive 2 parameters:
// - `key`    - the provider settings key (eg. "gitlabAuth")
// - `config` - the provider settings config that is currently being updated
export default [
    {
        key:   "appleAuth",
        title: "Apple",
        logo:  "/images/oauth2/apple.svg",
        optionsComponent: AppleOptions,
    },
    {
        key:   "googleAuth",
        title: "Google",
        logo:  "/images/oauth2/google.svg",
    },
    {
        key:   "facebookAuth",
        title: "Facebook",
        logo:  "/images/oauth2/facebook.svg",
    },
    {
        key:   "microsoftAuth",
        title: "Microsoft",
        logo:  "/images/oauth2/microsoft.svg",
        optionsComponent: MicrosoftOptions,
    },
    {
        key:   "githubAuth",
        title: "GitHub",
        logo:  "/images/oauth2/github.svg",
    },
    {
        key:   "gitlabAuth",
        title: "GitLab",
        logo:  "/images/oauth2/gitlab.svg",
        optionsComponent: SelfHostedOptions,
    },
    {
        key:   "giteeAuth",
        title: "Gitee",
        logo:  "/images/oauth2/gitee.svg",
    },
    {
        key:   "giteaAuth",
        title: "Gitea",
        logo:  "/images/oauth2/gitea.svg",
        optionsComponent: SelfHostedOptions,
    },
    {
        key:   "discordAuth",
        title: "Discord",
        logo:  "/images/oauth2/discord.svg",
    },
    {
        key:   "twitterAuth",
        title: "Twitter",
        logo:  "/images/oauth2/twitter.svg",
    },
    {
        key:   "kakaoAuth",
        title: "Kakao",
        logo:  "/images/oauth2/kakao.svg",
    },
    {
        key:   "spotifyAuth",
        title: "Spotify",
        logo:  "/images/oauth2/spotify.svg",
    },
    {
        key:   "twitchAuth",
        title: "Twitch",
        logo:  "/images/oauth2/twitch.svg",
    },
    {
        key:   "stravaAuth",
        title: "Strava",
        logo:  "/images/oauth2/strava.svg",
    },
    {
        key:   "livechatAuth",
        title: "LiveChat",
        logo:  "/images/oauth2/livechat.svg",
    },
    {
        key:   "oidcAuth",
        title: "OpenID Connect",
        logo:  "/images/oauth2/oidc.svg",
        optionsComponent: OIDCOptions,
    },
    {
        key:   "oidc2Auth",
        title: "(2) OpenID Connect",
        logo:  "/images/oauth2/oidc.svg",
        optionsComponent: OIDCOptions,
    },
    {
        key:   "oidc3Auth",
        title: "(3) OpenID Connect",
        logo:  "/images/oauth2/oidc.svg",
        optionsComponent: OIDCOptions,
    },
];
