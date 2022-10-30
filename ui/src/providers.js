// Object list with all supported OAuth2 providers in the format:
// ```
// { settingsKey: { title, icon, selfHosted, selfHostedRequired, selfHostedDescription} }
// ```
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
        title:    "GitLab",
        icon:     "ri-gitlab-fill",
        selfHosted: true,
        selfHostedDescription: "Optional endpoints (if you self host the OAUTH2 service)",
    },
    discordAuth: {
        title: "Discord",
        icon:  "ri-discord-fill",
    },
    microsoftAdAuth: {
        title: "MicrosoftAd",
        icon:  "ri-microsoft-fill",
        selfHosted: true,
        selfHostedRequired: "required",
        selfHostedDescription: "Provide the endpoints for your tenant."
    },
};
