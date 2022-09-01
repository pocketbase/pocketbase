// Object list with all supported OAuth2 providers in the format:
// ```
// { settingsKey: { title, icon, selfHosted } }
// ```
export default {
    googleAuth: {
        title: "Google",
        icon:  "ri-google-line",
    },
    facebookAuth: {
        title: "Facebook",
        icon:  "ri-facebook-line",
    },
    twitterAuth: {
        title: "Twitter",
        icon:  "ri-twitter-line",
    },
    githubAuth: {
        title: "GitHub",
        icon:  "ri-github-line",
    },
    gitlabAuth: {
        title:    "GitLab",
        icon:     "ri-gitlab-line",
        selfHosted: true,
    },
    discordAuth: {
        title: "Discord",
        icon:  "ri-discord-line",
    },
};
