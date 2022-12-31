import{S as Be,i as qe,s as Oe,e as i,w as v,b as _,c as Ie,f as b,g as r,h as s,m as Se,x as U,O as Pe,P as Le,k as Me,Q as Re,n as We,t as te,a as le,o as d,d as Ee,R as ze,C as De,p as He,r as j,u as Ue,N as je}from"./index.89a3f554.js";import{S as Ne}from"./SdkTabs.0a6ad1c9.js";function ye(a,l,o){const n=a.slice();return n[5]=l[o],n}function Ae(a,l,o){const n=a.slice();return n[5]=l[o],n}function Ce(a,l){let o,n=l[5].code+"",f,h,c,u;function m(){return l[4](l[5])}return{key:a,first:null,c(){o=i("button"),f=v(n),h=_(),b(o,"class","tab-item"),j(o,"active",l[1]===l[5].code),this.first=o},m(g,P){r(g,o,P),s(o,f),s(o,h),c||(u=Ue(o,"click",m),c=!0)},p(g,P){l=g,P&4&&n!==(n=l[5].code+"")&&U(f,n),P&6&&j(o,"active",l[1]===l[5].code)},d(g){g&&d(o),c=!1,u()}}}function Te(a,l){let o,n,f,h;return n=new je({props:{content:l[5].body}}),{key:a,first:null,c(){o=i("div"),Ie(n.$$.fragment),f=_(),b(o,"class","tab-item"),j(o,"active",l[1]===l[5].code),this.first=o},m(c,u){r(c,o,u),Se(n,o,null),s(o,f),h=!0},p(c,u){l=c;const m={};u&4&&(m.content=l[5].body),n.$set(m),(!h||u&6)&&j(o,"active",l[1]===l[5].code)},i(c){h||(te(n.$$.fragment,c),h=!0)},o(c){le(n.$$.fragment,c),h=!1},d(c){c&&d(o),Ee(n)}}}function Ge(a){var be,he,_e,ke;let l,o,n=a[0].name+"",f,h,c,u,m,g,P,M=a[0].name+"",N,oe,se,G,K,y,Q,I,F,w,R,ae,W,A,ne,J,z=a[0].name+"",V,ie,X,ce,re,D,Y,S,Z,E,x,B,ee,C,q,$=[],de=new Map,ue,O,k=[],pe=new Map,T;y=new Ne({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(be=a[0])==null?void 0:be.name}').authWithPassword('test@example.com', '123456');

        const result = await pb.collection('${(he=a[0])==null?void 0:he.name}').listExternalAuths(
            pb.authStore.model.id
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(_e=a[0])==null?void 0:_e.name}').authWithPassword('test@example.com', '123456');

        final result = await pb.collection('${(ke=a[0])==null?void 0:ke.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `}});let H=a[2];const fe=e=>e[5].code;for(let e=0;e<H.length;e+=1){let t=Ae(a,H,e),p=fe(t);de.set(p,$[e]=Ce(p,t))}let L=a[2];const me=e=>e[5].code;for(let e=0;e<L.length;e+=1){let t=ye(a,L,e),p=me(t);pe.set(p,k[e]=Te(p,t))}return{c(){l=i("h3"),o=v("List OAuth2 accounts ("),f=v(n),h=v(")"),c=_(),u=i("div"),m=i("p"),g=v("Returns a list with all OAuth2 providers linked to a single "),P=i("strong"),N=v(M),oe=v("."),se=_(),G=i("p"),G.textContent="Only admins and the account owner can access this action.",K=_(),Ie(y.$$.fragment),Q=_(),I=i("h6"),I.textContent="API details",F=_(),w=i("div"),R=i("strong"),R.textContent="GET",ae=_(),W=i("div"),A=i("p"),ne=v("/api/collections/"),J=i("strong"),V=v(z),ie=v("/records/"),X=i("strong"),X.textContent=":id",ce=v("/external-auths"),re=_(),D=i("p"),D.innerHTML="Requires <code>Authorization:TOKEN</code> header",Y=_(),S=i("div"),S.textContent="Path Parameters",Z=_(),E=i("table"),E.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the auth record.</td></tr></tbody>`,x=_(),B=i("div"),B.textContent="Responses",ee=_(),C=i("div"),q=i("div");for(let e=0;e<$.length;e+=1)$[e].c();ue=_(),O=i("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(u,"class","content txt-lg m-b-sm"),b(I,"class","m-b-xs"),b(R,"class","label label-primary"),b(W,"class","content"),b(D,"class","txt-hint txt-sm txt-right"),b(w,"class","alert alert-info"),b(S,"class","section-title"),b(E,"class","table-compact table-border m-b-base"),b(B,"class","section-title"),b(q,"class","tabs-header compact left"),b(O,"class","tabs-content"),b(C,"class","tabs")},m(e,t){r(e,l,t),s(l,o),s(l,f),s(l,h),r(e,c,t),r(e,u,t),s(u,m),s(m,g),s(m,P),s(P,N),s(m,oe),s(u,se),s(u,G),r(e,K,t),Se(y,e,t),r(e,Q,t),r(e,I,t),r(e,F,t),r(e,w,t),s(w,R),s(w,ae),s(w,W),s(W,A),s(A,ne),s(A,J),s(J,V),s(A,ie),s(A,X),s(A,ce),s(w,re),s(w,D),r(e,Y,t),r(e,S,t),r(e,Z,t),r(e,E,t),r(e,x,t),r(e,B,t),r(e,ee,t),r(e,C,t),s(C,q);for(let p=0;p<$.length;p+=1)$[p].m(q,null);s(C,ue),s(C,O);for(let p=0;p<k.length;p+=1)k[p].m(O,null);T=!0},p(e,[t]){var ve,we,$e,ge;(!T||t&1)&&n!==(n=e[0].name+"")&&U(f,n),(!T||t&1)&&M!==(M=e[0].name+"")&&U(N,M);const p={};t&9&&(p.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ve=e[0])==null?void 0:ve.name}').authWithPassword('test@example.com', '123456');

        const result = await pb.collection('${(we=e[0])==null?void 0:we.name}').listExternalAuths(
            pb.authStore.model.id
        );
    `),t&9&&(p.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${($e=e[0])==null?void 0:$e.name}').authWithPassword('test@example.com', '123456');

        final result = await pb.collection('${(ge=e[0])==null?void 0:ge.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `),y.$set(p),(!T||t&1)&&z!==(z=e[0].name+"")&&U(V,z),t&6&&(H=e[2],$=Pe($,t,fe,1,e,H,de,q,Le,Ce,null,Ae)),t&6&&(L=e[2],Me(),k=Pe(k,t,me,1,e,L,pe,O,Re,Te,null,ye),We())},i(e){if(!T){te(y.$$.fragment,e);for(let t=0;t<L.length;t+=1)te(k[t]);T=!0}},o(e){le(y.$$.fragment,e);for(let t=0;t<k.length;t+=1)le(k[t]);T=!1},d(e){e&&d(l),e&&d(c),e&&d(u),e&&d(K),Ee(y,e),e&&d(Q),e&&d(I),e&&d(F),e&&d(w),e&&d(Y),e&&d(S),e&&d(Z),e&&d(E),e&&d(x),e&&d(B),e&&d(ee),e&&d(C);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function Ke(a,l,o){let n,{collection:f=new ze}=l,h=200,c=[];const u=m=>o(1,h=m.code);return a.$$set=m=>{"collection"in m&&o(0,f=m.collection)},a.$$.update=()=>{a.$$.dirty&1&&o(2,c=[{code:200,body:`
                [
                    {
                      "id": "8171022dc95a4e8",
                      "created": "2022-09-01 10:24:18.434",
                      "updated": "2022-09-01 10:24:18.889",
                      "recordId": "e22581b6f1d44ea",
                      "collectionId": "${f.id}",
                      "provider": "google",
                      "providerId": "2da15468800514p",
                    },
                    {
                      "id": "171022dc895a4e8",
                      "created": "2022-09-01 10:24:18.434",
                      "updated": "2022-09-01 10:24:18.889",
                      "recordId": "e22581b6f1d44ea",
                      "collectionId": "${f.id}",
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
            `}])},o(3,n=De.getApiExampleUrl(He.baseUrl)),[f,h,c,n,u]}class Je extends Be{constructor(l){super(),qe(this,l,Ke,Ge,Oe,{collection:0})}}export{Je as default};
