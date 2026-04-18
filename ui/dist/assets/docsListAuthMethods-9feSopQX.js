import{t as e}from"./fieldsInfo-Bz62125-.js";function n(n){let r=app.utils.getApiExampleURL(),i=store({isLoading:!1,authMethods:[],get responses(){return[{title:200,value:i.isLoading?`...`:JSON.stringify(i.authMethods,null,2)},{title:404,value:`
                        {
                          "status": 404,
                          "message": "Missing collection context.",
                          "data": {}
                        }
                    `}]}});async function a(){i.isLoading=!0;try{i.authMethods=await app.pb.collection(n.name).listAuthMethods()}catch(e){e.isAbort&&app.pb.checkApiError(e)}i.isLoading=!1}return t.div({pbEvent:`apiPreviewListAuthMethods`,className:`content`,onmount:()=>{a()}},t.p(null,`Returns a public list with all allowed ${n.name} authentication methods.`),app.components.codeBlockTabs({className:`sdk-examples m-t-sm`,historyKey:`pbLastSDK`,tabs:[{title:`JS SDK`,language:`js`,value:`
                        import PocketBase from 'pocketbase';

                        const pb = new PocketBase('${r}');

                        ...

                        const result = await pb.collection('${n.name}').listAuthMethods();
                    `,footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/js-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`JS SDK docs`}))},{title:`Dart SDK`,language:`dart`,value:`
                        import 'package:pocketbase/pocketbase.dart';

                        final pb = PocketBase('${r}');

                        ...

                        final result = await pb.collection('${n.name}').listAuthMethods();
                    `,footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/dart-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`Dart SDK docs`}))},{title:`curl`,language:`bash`,value:`
                        curl '${r}/api/collections/${n.name}/auth-methods'
                    `}]}),t.div({className:`block m-t-base`},t.strong(null,`API details`)),t.div({className:`alert info api-preview-alert`},t.span({className:`label method`},`GET`),t.span({className:`path`},`/api/collections/${n.name}/auth-methods`)),t.table({className:`api-preview-table query-params`},t.thead(null,t.tr(null,t.th({className:`min-width txt-primary`},`?query params`),t.th({className:`min-width`},`Type`),t.th(null,`Description`))),t.tbody(null,t.tr(null,t.td({className:`min-width`},`fields`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,e())))),t.div({className:`block m-t-base m-b-sm`},t.strong(null,`Example responses`)),app.components.codeBlockTabs({tabs:()=>i.responses}))}export{n as docsListAuthMethods};