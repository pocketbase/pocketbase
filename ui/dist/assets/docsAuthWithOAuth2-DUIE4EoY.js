import{t as e}from"./expandInfo-DGS0CLSa.js";import{t as n}from"./fieldsInfo-Bz62125-.js";function r(r){let i=app.utils.getApiExampleURL(),a={collectionId:r.id,collectionName:r.name},o=[{title:200,value:JSON.stringify({token:`...JWT...`,record:Object.assign(a,app.utils.getDummyFieldsData(r)),meta:{id:`abc123`,name:`John Doe`,username:`john.doe`,email:`test@example.com`,avatarURL:`https://example.com/avatar.png`,accessToken:`...`,refreshToken:`...`,expiry:`2022-01-01 10:00:00.123Z`,isNew:!1,rawUser:{}}},null,2)},{title:400,value:`
                {
                  "status": 400,
                  "message": "An error occurred while submitting the form.",
                  "data": {
                    "provider": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}];return t.div({pbEvent:`apiPreviewAuthWithOAuth2`,className:`content`},t.p(null,`Authenticate with an OAuth2 provider and returns a new auth token and record data.`),t.p(null,`For more details please check the `,t.a({href:`https://pocketbase.io/docs/authentication#authenticate-with-oauth2`,target:`_blank`,rel:`noopener noreferrer`,textContent:`OAuth2 integration documentation`}),`.`),app.components.codeBlockTabs({className:`sdk-examples m-t-sm`,historyKey:`pbLastSDK`,tabs:[{title:`JS SDK`,language:`js`,value:`
                        import PocketBase from 'pocketbase';

                        const pb = new PocketBase('${i}');

                        ...

                        // OAuth2 authentication with a single realtime call.
                        //
                        // Make sure to register ${i}/api/oauth2-redirect
                        // as redirect url in the OAuth2 app configuration.
                        const authData = await pb.collection('${r.name}').authWithOAuth2({ provider: 'google' });

                        // OR authenticate with manual OAuth2 code exchange
                        // const authData = await pb.collection('${r.name}').authWithOAuth2Code(...);

                        // after the above you can also access the auth data from the authStore
                        console.log(pb.authStore.isValid);
                        console.log(pb.authStore.token);
                        console.log(pb.authStore.record.id);

                        // "logout"
                        pb.authStore.clear();
                    `,footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/js-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`JS SDK docs`}))},{title:`Dart SDK`,language:`dart`,value:`
                        import 'package:pocketbase/pocketbase.dart';
                        import 'package:url_launcher/url_launcher.dart';

                        final pb = PocketBase('${i}');

                        ...

                        // OAuth2 authentication with a single realtime call.
                        //
                        // Make sure to register ${i}/api/oauth2-redirect
                        // as redirect url in the OAuth2 app configuration.
                        final authData = await pb.collection('${r.name}').authWithOAuth2('google', (url) async {
                          await launchUrl(url);
                        });

                        // OR authenticate with manual OAuth2 code exchange
                        // final authData = await pb.collection('${r.name}').authWithOAuth2Code(...);

                        // after the above you can also access the auth data from the authStore
                        print(pb.authStore.isValid);
                        print(pb.authStore.token);
                        print(pb.authStore.record.id);

                        // "logout"
                        pb.authStore.clear();
                    `,footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/dart-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`Dart SDK docs`}))},{title:`curl`,language:`bash`,value:`
                        # authenticate with manual OAuth2 code exchange
                        curl -X POST \\
                          -H 'Content-Type:application/json' \\
                          -d '{ "provider":"google", "code":"OAUTH2_CODE", "codeVerifier":"...", "redirectURL":"..." }' \\
                          '${i}/api/collections/${r.name}/auth-with-oauth2'
                    `}]}),t.div({className:`m-t-base`},t.strong(null,`API details`)),t.div({className:`alert success api-preview-alert`},t.span({className:`label method`},`POST`),t.span({className:`path`},`/api/collections/${r.name}/auth-with-password`)),t.table({className:`api-preview-table body-params`},t.thead(null,t.tr(null,t.th({className:`min-width txt-primary`},`Body params`),t.th({className:`min-width`},`Type`),t.th(null,`Description`))),t.tbody(null,t.tr(null,t.td({className:`min-width`},`provider `,t.em(null,`(required)`)),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,`The name of the OAuth2 client provider (eg. "google").`)),t.tr(null,t.td({className:`min-width`},`code `,t.em(null,`(required)`)),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,`The authorization code returned from the initial request.`)),t.tr(null,t.td({className:`min-width`},`codeVerifier `,t.em(null,`(required)`)),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,`The code verifier sent with the initial request as part of the code_challenge.`)),t.tr(null,t.td({className:`min-width`},`redirectURL `,t.em(null,`(required)`)),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,`The redirect url sent with the initial request.`)),t.tr(null,t.td({className:`min-width`},`createData `,t.em(null,`(optional)`)),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,t.p(null,`Optional data that will be used when creating the auth record on OAuth2 sign-up.`),t.p(null,`The created auth record must comply with the same requirements and validations in the regular create action.`),t.p(null,`The data can only be in JSON, aka. user uploaded files currently are not supported during OAuth2 sign-ups.`))))),t.table({className:`api-preview-table query-params`},t.thead(null,t.tr(null,t.th({className:`min-width txt-primary`},`?query params`),t.th({className:`min-width`},`Type`),t.th(null,`Description`))),t.tbody(null,t.tr(null,t.td({className:`min-width`},`expand`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,e())),t.tr(null,t.td({className:`min-width`},`fields`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,n())))),t.div({className:`m-t-base m-b-sm`},t.strong(null,`Example responses`)),app.components.codeBlockTabs({tabs:o}))}export{r as docsAuthWithOAuth2};