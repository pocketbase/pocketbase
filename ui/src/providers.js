// Object list with all supported OAuth2 providers in the format:
// ```
// { settingsKey: { title, icon, selfHosted } }
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
    },
    discordAuth: {
        title: "Discord",
        icon:  "ri-discord-fill",
    },
};
