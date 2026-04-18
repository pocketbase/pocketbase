import{t as e}from"./expandInfo-DGS0CLSa.js";import{t as n}from"./fieldsInfo-Bz62125-.js";import{fullDummyPayload as r,primitivesDummyPayload as i,replaceDummyPayloadPlaceholder as a}from"./docsCreate-Be3S3y5K.js";function o(o){let s=app.utils.getApiExampleURL(),c=o.updateRule===null,l=o.type===`auth`?[`id`,`password`,`verified`,`email`,`emailVisibility`]:[`id`],u=o.fields?.filter(e=>!e.hidden&&e.type!=`autodate`&&!l.includes(e.name))||[],d={collectionId:o.id,collectionName:o.name},f=[{title:200,value:JSON.stringify(Object.assign(d,app.utils.getDummyFieldsData(o)),null,2)},{title:400,value:`
                {
                  "status": 400,
                  "message": "Failed to create record.",
                  "data": {
                    "${u.find(e=>!e.primaryKey)?.name||`someField`}": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}];return c&&f.push({title:403,value:`
                {
                  "status": 403,
                  "message": "Only superusers can perform this action.",
                  "data": {}
                }
            `}),f.push({title:404,value:`
            {
              "status": 404,
              "message": "The requested resource wasn't found.",
              "data": {}
            }
        `}),t.div({pbEvent:`apiPreviewUpdate`,className:`content`},t.p(null,`Updates an existing ${o.name} record.`),t.p(null,`Body parameters could be sent as `,t.code(null,`application/json`),` or `,t.code(null,`multipart/form-data`),`.`),t.p(null,`File upload is supported only via `,t.code(null,`multipart/form-data`),`. For more info and examples you could check the detailed `,t.a({href:`https://pocketbase.io/docs/files-handling`,target:`_blank`,rel:`noopener noreferrer`,textContent:`Files upload and handling docs`}),`.`),t.p(null,t.em(null,`Note that in case of a password change all previously issued tokens for the current record will be automatically invalidated and if you want your user to remain signed in you need to reauthenticate manually after the update call.`)),app.components.codeBlockTabs({className:`sdk-examples m-t-sm`,historyKey:`pbLastSDK`,tabs:[{title:`JS SDK`,language:`js`,value:`
import PocketBase from 'pocketbase';

const pb = new PocketBase('${s}');

...

// example update body
const body = ${a(JSON.stringify(r(o,!0),null,2))};

const record = await pb.collection('${o.name}').update('RECORD_ID', body);
`,footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/js-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`JS SDK docs`}))},{title:`Dart SDK`,language:`dart`,value:`
import 'package:pocketbase/pocketbase.dart';

final pb = PocketBase('${s}');

...

// example update body
final body = <String, dynamic>${JSON.stringify(i(o,!0),null,2)};

final record = await pb.collection('${o.name}').update(
  'RECORD_ID',
  body: body,
  files: [],
);
`,footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/dart-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`Dart SDK docs`}))},{title:`curl`,language:`bash`,value:`
                        curl -X PATCH \\
                          -H 'Authorization:TOKEN' \\
                          -H 'Content-Type:application/json' \\
                          -d '{ ... }' \\
                          '${s}/api/collections/${o.name}/records/RECORD_ID'
                    `}]}),t.div({className:`block m-t-base`},t.strong(null,`API details`)),t.div({className:`alert warning api-preview-alert`},t.span({className:`label method`},`PATCH`),t.span({className:`path`},`/api/collections/${o.name}/records/`,t.strong(null,`:id`)),()=>{if(c)return t.small({className:`extra`},`Requires superuser Authorization:TOKEN header`)}),t.table({className:`api-preview-table path-params`},t.thead(null,t.tr(null,t.th({className:`min-width txt-primary`},`Path params`),t.th({className:`min-width`},`Type`),t.th(null,`Description`))),t.tbody(null,t.tr(null,t.td({className:`min-width`},`id`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,`ID of the record to update.`)))),t.table({className:`api-preview-table query-params`},t.thead(null,t.tr(null,t.th({className:`min-width txt-primary`},`?query params`),t.th({className:`min-width`},`Type`),t.th(null,`Description`))),t.tbody(null,t.tr(null,t.td({className:`min-width`},`expand`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,e())),t.tr(null,t.td({className:`min-width`},`fields`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,n())))),t.div({className:`block m-t-base m-b-sm`},t.strong(null,`Example responses`)),app.components.codeBlockTabs({tabs:f}))}export{o as docsUpdate};