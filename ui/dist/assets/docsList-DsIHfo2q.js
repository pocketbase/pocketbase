import{t as e}from"./expandInfo-CJ9aeAaN.js";import{t as n}from"./fieldsInfo-CiLcXgNq.js";function r(){let e=store({show:!1});return t.div({className:`filter-details m-t-10`},t.button({type:`button`,className:`btn secondary sm`,onclick:()=>e.show=!e.show},()=>e.show?[t.span({className:`txt`},`Hide details`),t.i({className:`ri-arrow-up-s-line`,ariaHidden:!0})]:[t.span({className:`txt`},`Show details`),t.i({className:`ri-arrow-down-s-line`,ariaHidden:!0})]),app.components.slide(()=>e.show,t.div({className:`block p-t-5`},t.p(null,`The filter syntax follows the format `,t.code(null,t.span({className:`txt-success`},`OPERAND`),t.span({className:`txt-danger`},` OPERATOR `),t.span({className:`txt-success`},`OPERAND`)),`, where:`),t.ul(null,t.li(null,t.span({className:`txt-code txt-success`},`OPERAND`),` could be any of the above field literal, function, string (single or double quoted), number, null, true, false`),t.li(null,t.span({className:`txt-code txt-danger`},`OPERATOR`),` is one of:`,t.ul(null,t.li(null,t.code({className:`filter-op`},`=`),` Equal`),t.li(null,t.code({className:`filter-op`},`!=`),` Not equal`),t.li(null,t.code({className:`filter-op`},`>`),` Greater than`),t.li(null,t.code({className:`filter-op`},`>=`),` Greater than or equal`),t.li(null,t.code({className:`filter-op`},`<`),` Less than`),t.li(null,t.code({className:`filter-op`},`<=`),` Less than or equal`),t.li(null,t.code({className:`filter-op`},`~`),` Like/Contains`,t.div({className:`txt-sm txt-hint`},t.em(null,`(auto wraps the right string OPERAND in a "%" for wildcard match if not explicitly set)`))),t.li(null,t.code({className:`filter-op`},`!~`),` NOT Like/Contains`,t.div({className:`txt-sm txt-hint`},t.em(null,`(auto wraps the right string OPERAND in a "%" for wildcard match if not explicitly set)`))),t.li(null,t.code({className:`filter-op`},`?=`),t.span({className:`txt-hint`},` Any/At-least-one-of`),` Equal`),t.li(null,t.code({className:`filter-op`},`?!=`),t.span({className:`txt-hint`},` Any/At-least-one-of`),` Not equal`),t.li(null,t.code({className:`filter-op`},`?>`),t.span({className:`txt-hint`},` Any/At-least-one-of`),` Greater than`),t.li(null,t.code({className:`filter-op`},`?>=`),t.span({className:`txt-hint`},` Any/At-least-one-of`),` Greater than or equal`),t.li(null,t.code({className:`filter-op`},`?<`),t.span({className:`txt-hint`},` Any/At-least-one-of`),` Less than`),t.li(null,t.code({className:`filter-op`},`?<=`),t.span({className:`txt-hint`},` Any/At-least-one-of`),` Less than or equal`),t.li(null,t.code({className:`filter-op`},`?~`),t.span({className:`txt-hint`},` Any/At-least-one-of`),` Like/Contains`,t.div({className:`txt-sm txt-hint`},t.em(null,`(auto wraps the right string OPERAND in a "%" for wildcard match if not explicitly set)`))),t.li(null,t.code({className:`filter-op`},`?!~`),t.span({className:`txt-hint`},` Any/At-least-one-of`),` NOT Like/Contains`,t.div({className:`txt-sm txt-hint`},t.em(null,`(auto wraps the right string OPERAND in a "%" for wildcard match if not explicitly set)`)))),t.p(null,`To group and combine several expressions you could use brackets `,t.code(null,`(...)`),`, `,t.code(null,`&&`),`, (AND) and `,t.code(null,`||`),` (OR) tokens.`))))))}function i(i){let a=app.utils.getApiExampleURL(),o=i.listRule===null,s={collectionId:i.id,collectionName:i.name},c=[{title:200,value:JSON.stringify({page:1,perPage:30,totalPages:1,totalItems:2,items:[Object.assign(s,app.utils.getDummyFieldsData(i)),Object.assign(s,app.utils.getDummyFieldsData(i))]},null,2)},{title:400,value:`
                {
                  "status": 400,
                  "message": "Something went wrong while processing your request.",
                  "data": {}
                }
            `}];return o&&c.push({title:403,value:`
                {
                  "status": 403,
                  "message": "Only superusers can access this action.",
                  "data": {}
                }
            `}),t.div({pbEvent:`apiPreviewList`,className:`content`},t.p(null,`Fetch a paginated ${i.name} records list, supporting sorting and filtering.`),app.components.codeBlockTabs({className:`sdk-examples m-t-sm`,historyKey:`pbLastSDK`,tabs:[{title:`JS SDK`,language:`js`,value:`
                        import PocketBase from 'pocketbase';

                        const pb = new PocketBase('${a}');

                        ...

                        // fetch a paginated records list
                        const resultList = await pb.collection('${i.name}').getList(1, 50, {
                          filter: 'someField1 != someField2',
                        });

                        // you can also fetch all records at once via getFullList
                        const records = await pb.collection('${i.name}').getFullList({
                          sort: '-someField',
                        });

                        // or fetch only the first record that matches the specified filter
                        const record = await pb.collection('${i.name}').getFirstListItem(
                          'someField="test"',
                          { expand: 'relField1,relField2.subRelField' },
                        );
                    `,footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/js-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`JS SDK docs`}))},{title:`Dart SDK`,language:`dart`,value:`
                        import 'package:pocketbase/pocketbase.dart';

                        final pb = PocketBase('${a}');

                        ...

                        // fetch a paginated records list
                        final resultList = await pb.collection('${i.name}').getList(
                          page: 1,
                          perPage: 50,
                          filter: 'someField1 != someField2',
                        );

                        // you can also fetch all records at once via getFullList
                        final records = await pb.collection('${i.name}').getFullList(
                          sort: '-someField',
                        );

                        // or fetch only the first record that matches the specified filter
                        final record = await pb.collection('${i.name}').getFirstListItem(
                          'someField="test"',
                          expand: 'relField1,relField2.subRelField',
                        );
                    `,footnote:t.div({className:`txt-right`},t.a({href:`https://github.com/pocketbase/dart-sdk`,target:`_blank`,rel:`noopener noreferrer`,textContent:`Dart SDK docs`}))},{title:`curl`,language:`bash`,value:`
                        curl \\
                          -H 'Authorization:TOKEN' \\
                          '${a}/api/collections/${i.name}/records?perPage=50'
                    `}]}),t.div({className:`block m-t-base`},t.strong(null,`API details`)),t.div({className:`alert info api-preview-alert`},t.span({className:`label method`},`GET`),t.span({className:`path`},`/api/collections/${i.name}/records`),()=>{if(o)return t.small({className:`extra`},`Requires superuser Authorization:TOKEN header`)}),t.table({className:`api-preview-table query-params`},t.thead(null,t.tr(null,t.th({className:`min-width txt-primary`},`?query params`),t.th({className:`min-width`},`Type`),t.th(null,`Description`))),t.tbody(null,t.tr(null,t.td({className:`min-width`},`page`),t.td({className:`min-width`},t.span({className:`label`},`Number`)),t.td(null,`The page (aka. offset) of the paginated list (default to 1).`)),t.tr(null,t.td({className:`min-width`},`perPage`),t.td({className:`min-width`},t.span({className:`label`},`Number`)),t.td(null,`Specify the max returned records per page (default to 30).`)),t.tr(null,t.td({className:`min-width`},`sort`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,t.p(null,`Specify the records order attribute(s).`,t.br(),`Add -/+ (default) in front of the attribute for DESC / ASC order.`),t.p(null,`For example:`,app.components.codeBlock({value:`// DESC by created and ASC by id
?sort=-created,id`})),t.p(null,`In addition to the collection non-hidden fields, the following special sort fields could be also used: `,t.code(null,`@random`),` `,t.code({hidden:()=>i.type==`view`},`@rowid`),`.`))),t.tr(null,t.td({className:`min-width`},`filter`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,t.p(null,`Filter the returned records. For example:`),app.components.codeBlock({value:`?filter=(id='abc' && created>'2022-01-01')`,footnote:`All query params must be properly URL encoded (the SDKs do this automatically).`}),r())),t.tr(null,t.td({className:`min-width`},`expand`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,e())),t.tr(null,t.td({className:`min-width`},`fields`),t.td({className:`min-width`},t.span({className:`label`},`String`)),t.td(null,n())),t.tr(null,t.td({className:`min-width`},`skipTotal`),t.td({className:`min-width`},t.span({className:`label`},`Boolean`)),t.td(null,t.p(null,`If set to `,t.code(null,`1/true`),` the total counts query will be skipped and the response fields `,t.code(null,`totalItems`),` and `,t.code(null,`totalPages`),` will have -1 value.`),t.p(null,`This could drastically speed up the search queries when the total counters are not needed or cursor based pagination is used.`,` For optimization purposes, it is set by default in the `,t.code(null,`getFirstListItem()`),` and `,t.code(null,`getFullList()`),` SDKs methods.`))))),t.div({className:`block m-t-base m-b-sm`},t.strong(null,`Example responses`)),app.components.codeBlockTabs({tabs:c}))}export{i as docsList};