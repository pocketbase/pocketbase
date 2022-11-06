import{S as Be,i as qe,s as Le,e as n,w as v,b as h,c as Te,f as b,g as r,h as s,m as Ie,x as U,P as ye,Q as Oe,k as Me,R as Re,n as Ve,t as te,a as le,o as d,d as Se,L as ze,C as De,p as He,r as j,u as Ue,O as je}from"./index.b110ca95.js";import{S as Ge}from"./SdkTabs.b01956c7.js";function Ee(a,l,o){const i=a.slice();return i[5]=l[o],i}function Ae(a,l,o){const i=a.slice();return i[5]=l[o],i}function Pe(a,l){let o,i=l[5].code+"",m,_,c,u;function f(){return l[4](l[5])}return{key:a,first:null,c(){o=n("button"),m=v(i),_=h(),b(o,"class","tab-item"),j(o,"active",l[1]===l[5].code),this.first=o},m(g,y){r(g,o,y),s(o,m),s(o,_),c||(u=Ue(o,"click",f),c=!0)},p(g,y){l=g,y&4&&i!==(i=l[5].code+"")&&U(m,i),y&6&&j(o,"active",l[1]===l[5].code)},d(g){g&&d(o),c=!1,u()}}}function Ce(a,l){let o,i,m,_;return i=new je({props:{content:l[5].body}}),{key:a,first:null,c(){o=n("div"),Te(i.$$.fragment),m=h(),b(o,"class","tab-item"),j(o,"active",l[1]===l[5].code),this.first=o},m(c,u){r(c,o,u),Ie(i,o,null),s(o,m),_=!0},p(c,u){l=c;const f={};u&4&&(f.content=l[5].body),i.$set(f),(!_||u&6)&&j(o,"active",l[1]===l[5].code)},i(c){_||(te(i.$$.fragment,c),_=!0)},o(c){le(i.$$.fragment,c),_=!1},d(c){c&&d(o),Se(i)}}}function Ke(a){var be,_e,he,ke;let l,o,i=a[0].name+"",m,_,c,u,f,g,y,M=a[0].name+"",G,oe,se,K,N,E,Q,T,F,w,R,ae,V,A,ie,J,z=a[0].name+"",W,ne,X,ce,re,D,Y,I,Z,S,x,B,ee,P,q,$=[],de=new Map,ue,L,k=[],pe=new Map,C;E=new Ge({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(be=a[0])==null?void 0:be.name}').authViaEmail('test@example.com', '123456');

        const result = await pb.collection('${(_e=a[0])==null?void 0:_e.name}').listExternalAuths(
            pb.authStore.model.id
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(he=a[0])==null?void 0:he.name}').authViaEmail('test@example.com', '123456');

        final result = await pb.collection('${(ke=a[0])==null?void 0:ke.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `}});let H=a[2];const me=e=>e[5].code;for(let e=0;e<H.length;e+=1){let t=Ae(a,H,e),p=me(t);de.set(p,$[e]=Pe(p,t))}let O=a[2];const fe=e=>e[5].code;for(let e=0;e<O.length;e+=1){let t=Ee(a,O,e),p=fe(t);pe.set(p,k[e]=Ce(p,t))}return{c(){l=n("h3"),o=v("List OAuth2 accounts ("),m=v(i),_=v(")"),c=h(),u=n("div"),f=n("p"),g=v("Returns a list with all OAuth2 providers linked to a single "),y=n("strong"),G=v(M),oe=v("."),se=h(),K=n("p"),K.textContent="Only admins and the account owner can access this action.",N=h(),Te(E.$$.fragment),Q=h(),T=n("h6"),T.textContent="API details",F=h(),w=n("div"),R=n("strong"),R.textContent="GET",ae=h(),V=n("div"),A=n("p"),ie=v("/api/collections/"),J=n("strong"),W=v(z),ne=v("/records/"),X=n("strong"),X.textContent=":id",ce=v("/external-auths"),re=h(),D=n("p"),D.innerHTML="Requires <code>Authorization:TOKEN</code> header",Y=h(),I=n("div"),I.textContent="Path Parameters",Z=h(),S=n("table"),S.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the auth record.</td></tr></tbody>`,x=h(),B=n("div"),B.textContent="Responses",ee=h(),P=n("div"),q=n("div");for(let e=0;e<$.length;e+=1)$[e].c();ue=h(),L=n("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(u,"class","content txt-lg m-b-sm"),b(T,"class","m-b-xs"),b(R,"class","label label-primary"),b(V,"class","content"),b(D,"class","txt-hint txt-sm txt-right"),b(w,"class","alert alert-info"),b(I,"class","section-title"),b(S,"class","table-compact table-border m-b-base"),b(B,"class","section-title"),b(q,"class","tabs-header compact left"),b(L,"class","tabs-content"),b(P,"class","tabs")},m(e,t){r(e,l,t),s(l,o),s(l,m),s(l,_),r(e,c,t),r(e,u,t),s(u,f),s(f,g),s(f,y),s(y,G),s(f,oe),s(u,se),s(u,K),r(e,N,t),Ie(E,e,t),r(e,Q,t),r(e,T,t),r(e,F,t),r(e,w,t),s(w,R),s(w,ae),s(w,V),s(V,A),s(A,ie),s(A,J),s(J,W),s(A,ne),s(A,X),s(A,ce),s(w,re),s(w,D),r(e,Y,t),r(e,I,t),r(e,Z,t),r(e,S,t),r(e,x,t),r(e,B,t),r(e,ee,t),r(e,P,t),s(P,q);for(let p=0;p<$.length;p+=1)$[p].m(q,null);s(P,ue),s(P,L);for(let p=0;p<k.length;p+=1)k[p].m(L,null);C=!0},p(e,[t]){var ve,we,$e,ge;(!C||t&1)&&i!==(i=e[0].name+"")&&U(m,i),(!C||t&1)&&M!==(M=e[0].name+"")&&U(G,M);const p={};t&9&&(p.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ve=e[0])==null?void 0:ve.name}').authViaEmail('test@example.com', '123456');

        const result = await pb.collection('${(we=e[0])==null?void 0:we.name}').listExternalAuths(
            pb.authStore.model.id
        );
    `),t&9&&(p.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${($e=e[0])==null?void 0:$e.name}').authViaEmail('test@example.com', '123456');

        final result = await pb.collection('${(ge=e[0])==null?void 0:ge.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `),E.$set(p),(!C||t&1)&&z!==(z=e[0].name+"")&&U(W,z),t&6&&(H=e[2],$=ye($,t,me,1,e,H,de,q,Oe,Pe,null,Ae)),t&6&&(O=e[2],Me(),k=ye(k,t,fe,1,e,O,pe,L,Re,Ce,null,Ee),Ve())},i(e){if(!C){te(E.$$.fragment,e);for(let t=0;t<O.length;t+=1)te(k[t]);C=!0}},o(e){le(E.$$.fragment,e);for(let t=0;t<k.length;t+=1)le(k[t]);C=!1},d(e){e&&d(l),e&&d(c),e&&d(u),e&&d(N),Se(E,e),e&&d(Q),e&&d(T),e&&d(F),e&&d(w),e&&d(Y),e&&d(I),e&&d(Z),e&&d(S),e&&d(x),e&&d(B),e&&d(ee),e&&d(P);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function Ne(a,l,o){let i,{collection:m=new ze}=l,_=200,c=[];const u=f=>o(1,_=f.code);return a.$$set=f=>{"collection"in f&&o(0,m=f.collection)},a.$$.update=()=>{a.$$.dirty&1&&o(2,c=[{code:200,body:`
                [
                    {
                      "id": "8171022dc95a4e8",
                      "created": "2022-09-01 10:24:18.434",
                      "updated": "2022-09-01 10:24:18.889",
                      "recordId": "e22581b6f1d44ea",
                      "collectionId": "${m.id}",
                      "provider": "google",
                      "providerId": "2da15468800514p",
                    },
                    {
                      "id": "171022dc895a4e8",
                      "created": "2022-09-01 10:24:18.434",
                      "updated": "2022-09-01 10:24:18.889",
                      "recordId": "e22581b6f1d44ea",
                      "collectionId": "${m.id}",
                      "provider": "twitter",
                      "providerId": "720688005140514",
                    }
                ]
            `},{code:401,body:`
                {
                  "code": 401,
                  "message": "The request requires valid record authorization token to be set.",
                  "data": {}
                }
            `},{code:403,body:`
                {
                  "code": 403,
                  "message": "The authorized record model is not allowed to perform this action.",
                  "data": {}
                }
            `},{code:404,body:`
                {
                  "code": 404,
                  "message": "The requested resource wasn't found.",
                  "data": {}
                }
            `}])},o(3,i=De.getApiExampleUrl(He.baseUrl)),[m,_,c,i,u]}class Je extends Be{constructor(l){super(),qe(this,l,Ne,Ke,Le,{collection:0})}}export{Je as default};
