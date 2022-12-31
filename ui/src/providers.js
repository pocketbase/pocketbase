import SelfHostedOptions from "@/components/settings/providers/SelfHostedOptions.svelte";
import MicrosoftOptions  from "@/components/settings/providers/MicrosoftOptions.svelte";

// Object list with all supported OAuth2 providers in the format:
// ```
// { settingsKey: { title, icon, optionsComponent? } }
// ```
//
// If `optionsComponent` is provided it will receive 2 parameters:
// - `key`    - the provider settings key (eg. "gitlabAuth")
// - `config` - the provider settings config that is currently being updated
export default {
    googleAuth: {
        title: "Google",
        icon:  "ri-google-fill",
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
    discordAuth: {
        title: "Discord",
        icon:  "ri-discord-fill",
    },
    microsoftAuth: {
        title: "Microsoft",
        icon:  "ri-microsoft-fill",
        optionsComponent: MicrosoftOptions,
    },
    spotifyAuth: {
        title: "Spotify",
        icon:  "ri-spotify-fill",
    },
    kakaoAuth: {
        title: "Kakao",
        icon:  "ri-kakao-talk-fill",
    },
    twitchAuth: {
        title: "Twitch",
        icon:  "ri-twitch-fill",
    },
    stravaAuth: {
        title: "Strava",
        icon:  "ri-riding-fill",
    },
    giteeAuth: {
        title: "Gitee",
        icon:  "ri-git-repository-fill",
    },
};
