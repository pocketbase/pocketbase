import{S as ze,i as Qe,s as Re,e as n,w as v,b as f,c as de,f as m,g as r,h as o,m as pe,x as F,O as Me,P as Ue,k as je,Q as Fe,n as Ne,t as N,a as G,o as c,d as ue,C as Ge,p as Ke,r as K,u as Je,N as Ve}from"./index-24aa4d9d.js";import{S as Xe}from"./SdkTabs-3de29185.js";import{F as Ye}from"./FieldsQueryParam-8f27aa35.js";function Oe(a,l,s){const i=a.slice();return i[5]=l[s],i}function De(a,l,s){const i=a.slice();return i[5]=l[s],i}function He(a,l){let s,i=l[5].code+"",b,_,d,u;function h(){return l[4](l[5])}return{key:a,first:null,c(){s=n("button"),b=v(i),_=f(),m(s,"class","tab-item"),K(s,"active",l[1]===l[5].code),this.first=s},m(y,P){r(y,s,P),o(s,b),o(s,_),d||(u=Je(s,"click",h),d=!0)},p(y,P){l=y,P&4&&i!==(i=l[5].code+"")&&F(b,i),P&6&&K(s,"active",l[1]===l[5].code)},d(y){y&&c(s),d=!1,u()}}}function We(a,l){let s,i,b,_;return i=new Ve({props:{content:l[5].body}}),{key:a,first:null,c(){s=n("div"),de(i.$$.fragment),b=f(),m(s,"class","tab-item"),K(s,"active",l[1]===l[5].code),this.first=s},m(d,u){r(d,s,u),pe(i,s,null),o(s,b),_=!0},p(d,u){l=d;const h={};u&4&&(h.content=l[5].body),i.$set(h),(!_||u&6)&&K(s,"active",l[1]===l[5].code)},i(d){_||(N(i.$$.fragment,d),_=!0)},o(d){G(i.$$.fragment,d),_=!1},d(d){d&&c(s),ue(i)}}}function Ze(a){var Ce,ge,Se,Ee;let l,s,i=a[0].name+"",b,_,d,u,h,y,P,W=a[0].name+"",J,fe,me,V,X,T,Y,I,Z,w,z,be,Q,A,he,x,R=a[0].name+"",ee,_e,te,ke,ve,U,le,B,se,q,oe,L,ae,C,ie,$e,ne,E,re,M,ce,g,O,$=[],we=new Map,ye,D,k=[],Pe=new Map,S;T=new Xe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(Ce=a[0])==null?void 0:Ce.name}').authWithPassword('test@example.com', '123456');

        const result = await pb.collection('${(ge=a[0])==null?void 0:ge.name}').listExternalAuths(
            pb.authStore.model.id
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(Se=a[0])==null?void 0:Se.name}').authWithPassword('test@example.com', '123456');

        final result = await pb.collection('${(Ee=a[0])==null?void 0:Ee.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `}}),E=new Ye({});let j=a[2];const Te=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=De(a,j,e),p=Te(t);we.set(p,$[e]=He(p,t))}let H=a[2];const Ae=e=>e[5].code;for(let e=0;e<H.length;e+=1){let t=Oe(a,H,e),p=Ae(t);Pe.set(p,k[e]=We(p,t))}return{c(){l=n("h3"),s=v("List OAuth2 accounts ("),b=v(i),_=v(")"),d=f(),u=n("div"),h=n("p"),y=v("Returns a list with all OAuth2 providers linked to a single "),P=n("strong"),J=v(W),fe=v("."),me=f(),V=n("p"),V.textContent="Only admins and the account owner can access this action.",X=f(),de(T.$$.fragment),Y=f(),I=n("h6"),I.textContent="API details",Z=f(),w=n("div"),z=n("strong"),z.textContent="GET",be=f(),Q=n("div"),A=n("p"),he=v("/api/collections/"),x=n("strong"),ee=v(R),_e=v("/records/"),te=n("strong"),te.textContent=":id",ke=v("/external-auths"),ve=f(),U=n("p"),U.innerHTML="Requires <code>Authorization:TOKEN</code> header",le=f(),B=n("div"),B.textContent="Path Parameters",se=f(),q=n("table"),q.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the auth record.</td></tr></tbody>`,oe=f(),L=n("div"),L.textContent="Query parameters",ae=f(),C=n("table"),ie=n("thead"),ie.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,$e=f(),ne=n("tbody"),de(E.$$.fragment),re=f(),M=n("div"),M.textContent="Responses",ce=f(),g=n("div"),O=n("div");for(let e=0;e<$.length;e+=1)$[e].c();ye=f(),D=n("div");for(let e=0;e<k.length;e+=1)k[e].c();m(l,"class","m-b-sm"),m(u,"class","content txt-lg m-b-sm"),m(I,"class","m-b-xs"),m(z,"class","label label-primary"),m(Q,"class","content"),m(U,"class","txt-hint txt-sm txt-right"),m(w,"class","alert alert-info"),m(B,"class","section-title"),m(q,"class","table-compact table-border m-b-base"),m(L,"class","section-title"),m(C,"class","table-compact table-border m-b-base"),m(M,"class","section-title"),m(O,"class","tabs-header compact left"),m(D,"class","tabs-content"),m(g,"class","tabs")},m(e,t){r(e,l,t),o(l,s),o(l,b),o(l,_),r(e,d,t),r(e,u,t),o(u,h),o(h,y),o(h,P),o(P,J),o(h,fe),o(u,me),o(u,V),r(e,X,t),pe(T,e,t),r(e,Y,t),r(e,I,t),r(e,Z,t),r(e,w,t),o(w,z),o(w,be),o(w,Q),o(Q,A),o(A,he),o(A,x),o(x,ee),o(A,_e),o(A,te),o(A,ke),o(w,ve),o(w,U),r(e,le,t),r(e,B,t),r(e,se,t),r(e,q,t),r(e,oe,t),r(e,L,t),r(e,ae,t),r(e,C,t),o(C,ie),o(C,$e),o(C,ne),pe(E,ne,null),r(e,re,t),r(e,M,t),r(e,ce,t),r(e,g,t),o(g,O);for(let p=0;p<$.length;p+=1)$[p]&&$[p].m(O,null);o(g,ye),o(g,D);for(let p=0;p<k.length;p+=1)k[p]&&k[p].m(D,null);S=!0},p(e,[t]){var Ie,Be,qe,Le;(!S||t&1)&&i!==(i=e[0].name+"")&&F(b,i),(!S||t&1)&&W!==(W=e[0].name+"")&&F(J,W);const p={};t&9&&(p.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(Ie=e[0])==null?void 0:Ie.name}').authWithPassword('test@example.com', '123456');

        const result = await pb.collection('${(Be=e[0])==null?void 0:Be.name}').listExternalAuths(
            pb.authStore.model.id
        );
    `),t&9&&(p.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(qe=e[0])==null?void 0:qe.name}').authWithPassword('test@example.com', '123456');

        final result = await pb.collection('${(Le=e[0])==null?void 0:Le.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `),T.$set(p),(!S||t&1)&&R!==(R=e[0].name+"")&&F(ee,R),t&6&&(j=e[2],$=Me($,t,Te,1,e,j,we,O,Ue,He,null,De)),t&6&&(H=e[2],je(),k=Me(k,t,Ae,1,e,H,Pe,D,Fe,We,null,Oe),Ne())},i(e){if(!S){N(T.$$.fragment,e),N(E.$$.fragment,e);for(let t=0;t<H.length;t+=1)N(k[t]);S=!0}},o(e){G(T.$$.fragment,e),G(E.$$.fragment,e);for(let t=0;t<k.length;t+=1)G(k[t]);S=!1},d(e){e&&c(l),e&&c(d),e&&c(u),e&&c(X),ue(T,e),e&&c(Y),e&&c(I),e&&c(Z),e&&c(w),e&&c(le),e&&c(B),e&&c(se),e&&c(q),e&&c(oe),e&&c(L),e&&c(ae),e&&c(C),ue(E),e&&c(re),e&&c(M),e&&c(ce),e&&c(g);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function xe(a,l,s){let i,{collection:b}=l,_=200,d=[];const u=h=>s(1,_=h.code);return a.$$set=h=>{"collection"in h&&s(0,b=h.collection)},a.$$.update=()=>{a.$$.dirty&1&&s(2,d=[{code:200,body:`
                [
                    {
                      "id": "8171022dc95a4e8",
                      "created": "2022-09-01 10:24:18.434",
                      "updated": "2022-09-01 10:24:18.889",
                      "recordId": "e22581b6f1d44ea",
                      "collectionId": "${b.id}",
                      "provider": "google",
                      "providerId": "2da15468800514p",
                    },
                    {
                      "id": "171022dc895a4e8",
                      "created": "2022-09-01 10:24:18.434",
                      "updated": "2022-09-01 10:24:18.889",
                      "recordId": "e22581b6f1d44ea",
                      "collectionId": "${b.id}",
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
            `}])},s(3,i=Ge.getApiExampleUrl(Ke.baseUrl)),[b,_,d,i,u]}class st extends ze{constructor(l){super(),Qe(this,l,xe,Ze,Re,{collection:0})}}export{st as default};
