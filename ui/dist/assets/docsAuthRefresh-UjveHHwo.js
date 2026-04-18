import{t as e}from"./expandInfo-DGS0CLSa.js";import{t as n}from"./fieldsInfo-Bz62125-.js";function r(r){let i=app.utils.getApiExampleURL(),a={collectionId:r.id,collectionName:r.name},o=[{title:200,value:JSON.stringify({token:`...JWT...`,record:Object.assign(a,app.utils.getDummyFieldsData(r))},null,2)},{title:401,value:`
                {
                  "status": 401,
                  "message": "The request requires valid record authorization token to be set.",
                  "data": {}
                }
            `},{title:403,value:`
                {
                  "status": 403,
                  "message": "The authorized record model is not allowed to perform this action.",
                  "data": {}
                }
            `},{title:404,value:`
                {
                  "status": 404,
                  "message": "Missing auth record context.",
                  "data": {}
                }
            `}];return t.div({pbEvent:`apiPreviewAuthRefresh`,className:`content`},t.p(null,`Returns a new auth response (token and record data) for an already authenticated record.`),t.p(null,`This method is usually called by users on page/screen reload to ensure that the previously stored data in `,t.code(null,`pb.authStore`),` is still valid and up-to-date.`),app.components.codeBlockTabs({className:`sdk-examples m-t-sm`,historyKey:`pbLastSDK`,tabs:[{title:`JS SDK`,language:`js`,value:`
                        import PocketBase from 'pocketbase';

                        const pb = new PocketBase('${i}');

                        ...

                        const authData = await pb.collection('${r.name}').authRefresh();

                        // after the above you can also access the refreshed auth data from the authStore
                        console.log(pb.authStore.isValid);
                        console.log(pb.authStore.token);
                        console.log(pb.authStore.record.id);
                    `,footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/js-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`JS SDK docs`}))},{title:`Dart SDK`,language:`dart`,value:`
                        import 'package:pocketbase/pocketbase.dart';

                        final pb = PocketBase('${i}');

                        ...

                        final authData = await pb.collection('${r.name}').authRefresh();

                        // after the above you can also access the refreshed auth data from the authStore
                        print(pb.authStore.isValid);
                        print(pb.authStore.token);
                        print(pb.authStore.record.id);
                    `,footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/dart-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`Dart SDK docs`}))},{title:`curl`,language:`bash`,value:`
                        curl -X POST \\
                          -H 'Authorization:TOKEN' \\
                          '${i}/api/collections/${r.name}/auth-refresh'
                    `}]}),t.div({className:`m-t-base`},t.strong(null,`API details`)),t.div({className:`alert success api-preview-alert`},t.span({className:`label method`},`POST`),t.span({className:`path`},`/api/collections/${r.name}/auth-refresh`),t.small({className:`extra`},`Requires`,t.br(),`Authorization:TOKEN header`)),t.table({className:`api-preview-table query-params`},t.thead(null,t.tr(null,t.th({className:`min-width txt-primary`},`?query params`),t.th({className:`min-width`},`Type`),t.th(null,`Description`))),t.tbody(null,t.tr(null,t.td({className:`min-width`},`expand`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,e())),t.tr(null,t.td({className:`min-width`},`fields`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,n())))),t.div({className:`m-t-base m-b-sm`},t.strong(null,`Example responses`)),app.components.codeBlockTabs({tabs:o}))}export{r as docsAuthRefresh};