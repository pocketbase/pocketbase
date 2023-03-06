import SelfHostedOptions from "@/components/settings/providers/SelfHostedOptions.svelte";
import MicrosoftOptions  from "@/components/settings/providers/MicrosoftOptions.svelte";
import OIDCOptions       from "@/components/settings/providers/OIDCOptions.svelte";
import AppleOptions      from "@/components/settings/providers/AppleOptions.svelte";

// Object list with all supported OAuth2 providers in the format:
// ```
// { settingsKey: { title, icon, hidden, optionsComponent? } }
// ```
//
// If `optionsComponent` is provided it will receive 2 parameters:
// - `key`    - the provider settings key (eg. "gitlabAuth")
// - `config` - the provider settings config that is currently being updated
export default {
    appleAuth: {
        title: "Apple",
        icon:  "ri-apple-fill",
        optionsComponent: AppleOptions,
    },
    googleAuth: {
        title: "Google",
        icon:  "ri-google-fill",
    },
    microsoftAuth: {
        title: "Microsoft",
        icon:  "ri-microsoft-fill",
        optionsComponent: MicrosoftOptions,
    },
    facebookAuth: {
        title: "Facebook",
        icon:  "ri-facebook-fill",
    },
    twitterAuth: {
        title: "Twitter",
        icon:  "ri-twitter-fill",
    },
    githubAuth: {
        title: "GitHub",
        icon:  "ri-github-fill",
    },
    gitlabAuth: {
        title: "GitLab",
        icon:  "ri-gitlab-fill",
        optionsComponent: SelfHostedOptions,
    },
    giteeAuth: {
        title: "Gitee",
        icon:  "ri-git-repository-fill",
    },
    giteaAuth: {
        title: "Gitea",
        icon:  "ri-cup-fill",
        optionsComponent: SelfHostedOptions,
    },
    discordAuth: {
        title: "Discord",
        icon:  "ri-discord-fill",
    },
    kakaoAuth: {
        title: "Kakao",
        icon:  "ri-kakao-talk-fill",
    },
    spotifyAuth: {
        title: "Spotify",
        icon:  "ri-spotify-fill",
    },
    twitchAuth: {
        title: "Twitch",
        icon:  "ri-twitch-fill",
    },
    stravaAuth: {
        title: "Strava",
        icon:  "ri-riding-fill",
    },
    livechatAuth: {
        title: "LiveChat",
        icon:  "ri-chat-1-fill",
    },
    oidcAuth: {
        title: "OpenID Connect - Authentik, Keycloak, Okta, etc.",
        icon:  "ri-lock-fill",
        optionsComponent: OIDCOptions,
    },
    oidc2Auth: {
        title: "(2) OpenID Connect - Authentik, Keycloak, Okta, etc.",
        icon:  "ri-lock-fill",
        hidden: true,
        optionsComponent: OIDCOptions,
    },
    oidc3Auth: {
        title: "(3) OpenID Connect - Authentik, Keycloak, Okta, etc.",
        icon:  "ri-lock-fill",
        hidden: true,
        optionsComponent: OIDCOptions,
    },
};
