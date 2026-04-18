import{t as e}from"./expandInfo-DGS0CLSa.js";import{t as n}from"./fieldsInfo-Bz62125-.js";function r(r){let s=app.utils.getApiExampleURL(),c=r.createRule===null,l=r.type===`auth`,u=l?[`password`,`verified`,`email`,`emailVisibility`]:[],d=r.fields?.filter(e=>!e.hidden&&e.type!=`autodate`&&!u.includes(e.name))||[],f={collectionId:r.id,collectionName:r.name},p=[{title:200,value:JSON.stringify(Object.assign(f,app.utils.getDummyFieldsData(r)),null,2)},{title:400,value:`
                {
                  "status": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${l?`email`:d.find(e=>!e.primaryKey)?.name||`someField`}": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}];return c&&p.push({title:403,value:`
                {
                  "status": 403,
                  "message": "Only superusers can perform this action.",
                  "data": {}
                }
            `}),t.div({pbEvent:`apiPreviewCreate`,className:`content`},t.p(null,`Creates a new ${r.name} record.`),t.p(null,`Body parameters could be sent as `,t.code(null,`application/json`),` or `,t.code(null,`multipart/form-data`),`.`),t.p(null,`File upload is supported only via `,t.code(null,`multipart/form-data`),`. For more info and examples you could check the detailed `,t.a({href:`https://pocketbase.io/docs/files-handling`,target:`_blank`,rel:`noopener noreferrer`,textContent:`Files upload and handling docs`}),`.`),app.components.codeBlockTabs({className:`sdk-examples m-t-sm`,historyKey:`pbLastSDK`,tabs:[{title:`JS SDK`,language:`js`,value:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${s}');

...

// example create body
const body = ${i(JSON.stringify(a(r),null,2))};

const record = await pb.collection('${r.name}').create(body);
`+(l?`
// (optional) send an email verification request
await pb.collection('${r?.name}').requestVerification('test@example.com');
`:``),footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/js-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`JS SDK docs`}))},{title:`Dart SDK`,language:`dart`,value:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${s}');

...

// example create body
final body = <String, dynamic>${JSON.stringify(o(r),null,2)};

final record = await pb.collection('${r.name}').create(body: body, files: []);
`+(l?`
// (optional) send an email verification request
await pb.collection('${r?.name}').requestVerification('test@example.com');
`:``),footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/dart-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`Dart SDK docs`}))},{title:`curl`,language:`bash`,value:`
                        curl -X POST \\
                          -H 'Authorization:TOKEN' \\
                          -H 'Content-Type:application/json' \\
                          -d '{ ... }' \\
                          '${s}/api/collections/${r.name}/records/RECORD_ID'
                    `}]}),t.div({className:`block m-t-base`},t.strong(null,`API details`)),t.div({className:`alert success api-preview-alert`},t.span({className:`label method`},`POST`),t.span({className:`path`},`/api/collections/${r.name}/records`),()=>{if(c)return t.small({className:`extra`},`Requires superuser Authorization:TOKEN header`)}),t.table({className:`api-preview-table body-params`},t.thead(null,t.tr(null,t.th({className:`min-width txt-primary`},`Body params`),t.th({className:`min-width`},`Type`),t.th(null,`Description`))),t.tbody(null,()=>{if(l)return[t.tr(null,t.th({colSpan:99},`Auth specific fields`)),t.tr(null,t.td({className:`min-width`},`email `,()=>r.fields?.find(e=>e.name==`email`)?.required?t.em(null,`(required)`):t.em(null,`(optional)`)),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,`Auth record email address.`)),t.tr(null,t.td({className:`min-width`},`emailVisibility `,()=>r.fields?.find(e=>e.name==`emailVisibility`)?.required?t.em(null,`(required)`):t.em(null,`(optional)`)),t.td({className:`min-width`},t.span({className:`label`},`Boolean`)),t.td(null,`Whether to show/hide the auth record email when fetching the record data.`,t.br(),`Superusers and the owner of the record always have access to the email address.`)),t.tr(null,t.td({className:`min-width`},`password `,t.em(null,`(required)`)),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,`Auth record password.`)),t.tr(null,t.td({className:`min-width`},`passwordConfirm `,t.em(null,`(required)`)),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,`Auth record password confirmation.`)),t.tr(null,t.td({className:`min-width`},`verified `,t.em(null,`(optional)`)),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,t.p(null,`Indicates whether the auth record is verified or not.`),t.p(null,`This field can be set only by superusers or auth records with "Manage" access.`))),t.tr(null,t.th({colSpan:99},`Other fields`))]},()=>d.map(e=>t.tr(null,t.td({className:`min-width`},e.name,t.em(null,e.required&&!e.autogeneratePattern?` (required)`:` (optional)`)),t.td({className:`min-width`},t.span({className:`label`},()=>{let n=app.fieldTypes[e.type]?.dummyData(e,!0),r=typeof n;return e.type==`file`?`File`:r===`string`?`String`:r==`number`?`Number`:r==`bool`?`Boolean`:Array.isArray(n)?`Array`:app.utils.isObject(n)?`Object`:`Mixed`})),t.td(null,t.code(null,e.type),` field type value.`,t.br(),t.small({className:`txt-hint`},`For more details you could check the `,t.a({href:`https://pocketbase.io/docs/collections/#fields`,target:`_blank`,rel:`noopener noreferrer`,textContent:`Fields docs`}),`.`)))))),t.table({className:`api-preview-table query-params`},t.thead(null,t.tr(null,t.th({className:`min-width txt-primary`},`?query params`),t.th({className:`min-width`},`Type`),t.th(null,`Description`))),t.tbody(null,t.tr(null,t.td({className:`min-width`},`expand`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,e())),t.tr(null,t.td({className:`min-width`},`fields`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,n())))),t.div({className:`block m-t-base m-b-sm`},t.strong(null,`Example responses`)),app.components.codeBlockTabs({tabs:p}))}function i(e){return e.replaceAll(`"[[`,``).replaceAll(`]]"`,``)}function a(e,n=!1){let r=app.utils.getDummyFieldsData(e,!0);return delete r.id,e.type==`auth`&&(n&&(r.oldPassword=`987654321`,delete r.email),r.password=`123456789`,r.passwordConfirm=`123456789`,delete r.verified),r}function o(e,n=!1){let r=a(e,n);for(let e in r){let n=typeof r[e];(r[e]?.startsWith?.(`[[`)||![`number`,`string`,`boolean`].includes(n)&&!Array.isArray(r[e]))&&delete r[e]}return r}export{r as docsCreate,a as fullDummyPayload,o as primitivesDummyPayload,i as replaceDummyPayloadPlaceholder};