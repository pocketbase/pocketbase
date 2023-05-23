import{S as ze,i as Qe,s as Re,e as n,w as v,b as f,c as de,f as m,g as r,h as o,m as pe,x as F,N as Le,P as Ue,k as je,Q as Fe,n as Ne,t as N,a as G,o as c,d as ue,T as Ge,C as Ke,p as Je,r as K,u as Ve,M as Xe}from"./index-a65ca895.js";import{S as Ye}from"./SdkTabs-ad912c8f.js";import{F as Ze}from"./FieldsQueryParam-ba250473.js";function De(a,l,s){const i=a.slice();return i[5]=l[s],i}function He(a,l,s){const i=a.slice();return i[5]=l[s],i}function Oe(a,l){let s,i=l[5].code+"",b,_,d,u;function h(){return l[4](l[5])}return{key:a,first:null,c(){s=n("button"),b=v(i),_=f(),m(s,"class","tab-item"),K(s,"active",l[1]===l[5].code),this.first=s},m(y,P){r(y,s,P),o(s,b),o(s,_),d||(u=Ve(s,"click",h),d=!0)},p(y,P){l=y,P&4&&i!==(i=l[5].code+"")&&F(b,i),P&6&&K(s,"active",l[1]===l[5].code)},d(y){y&&c(s),d=!1,u()}}}function We(a,l){let s,i,b,_;return i=new Xe({props:{content:l[5].body}}),{key:a,first:null,c(){s=n("div"),de(i.$$.fragment),b=f(),m(s,"class","tab-item"),K(s,"active",l[1]===l[5].code),this.first=s},m(d,u){r(d,s,u),pe(i,s,null),o(s,b),_=!0},p(d,u){l=d;const h={};u&4&&(h.content=l[5].body),i.$set(h),(!_||u&6)&&K(s,"active",l[1]===l[5].code)},i(d){_||(N(i.$$.fragment,d),_=!0)},o(d){G(i.$$.fragment,d),_=!1},d(d){d&&c(s),ue(i)}}}function xe(a){var Ce,ge,Se,Ee;let l,s,i=a[0].name+"",b,_,d,u,h,y,P,W=a[0].name+"",J,fe,me,V,X,T,Y,I,Z,$,z,be,Q,A,he,x,R=a[0].name+"",ee,_e,te,ke,ve,U,le,B,se,q,oe,M,ae,C,ie,we,ne,E,re,L,ce,g,D,w=[],$e=new Map,ye,H,k=[],Pe=new Map,S;T=new Ye({props:{js:`
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
    `}}),E=new Ze({});let j=a[2];const Te=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=He(a,j,e),p=Te(t);$e.set(p,w[e]=Oe(p,t))}let O=a[2];const Ae=e=>e[5].code;for(let e=0;e<O.length;e+=1){let t=De(a,O,e),p=Ae(t);Pe.set(p,k[e]=We(p,t))}return{c(){l=n("h3"),s=v("List OAuth2 accounts ("),b=v(i),_=v(")"),d=f(),u=n("div"),h=n("p"),y=v("Returns a list with all OAuth2 providers linked to a single "),P=n("strong"),J=v(W),fe=v("."),me=f(),V=n("p"),V.textContent="Only admins and the account owner can access this action.",X=f(),de(T.$$.fragment),Y=f(),I=n("h6"),I.textContent="API details",Z=f(),$=n("div"),z=n("strong"),z.textContent="GET",be=f(),Q=n("div"),A=n("p"),he=v("/api/collections/"),x=n("strong"),ee=v(R),_e=v("/records/"),te=n("strong"),te.textContent=":id",ke=v("/external-auths"),ve=f(),U=n("p"),U.innerHTML="Requires <code>Authorization:TOKEN</code> header",le=f(),B=n("div"),B.textContent="Path Parameters",se=f(),q=n("table"),q.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the auth record.</td></tr></tbody>`,oe=f(),M=n("div"),M.textContent="Query parameters",ae=f(),C=n("table"),ie=n("thead"),ie.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr>`,we=f(),ne=n("tbody"),de(E.$$.fragment),re=f(),L=n("div"),L.textContent="Responses",ce=f(),g=n("div"),D=n("div");for(let e=0;e<w.length;e+=1)w[e].c();ye=f(),H=n("div");for(let e=0;e<k.length;e+=1)k[e].c();m(l,"class","m-b-sm"),m(u,"class","content txt-lg m-b-sm"),m(I,"class","m-b-xs"),m(z,"class","label label-primary"),m(Q,"class","content"),m(U,"class","txt-hint txt-sm txt-right"),m($,"class","alert alert-info"),m(B,"class","section-title"),m(q,"class","table-compact table-border m-b-base"),m(M,"class","section-title"),m(C,"class","table-compact table-border m-b-base"),m(L,"class","section-title"),m(D,"class","tabs-header compact left"),m(H,"class","tabs-content"),m(g,"class","tabs")},m(e,t){r(e,l,t),o(l,s),o(l,b),o(l,_),r(e,d,t),r(e,u,t),o(u,h),o(h,y),o(h,P),o(P,J),o(h,fe),o(u,me),o(u,V),r(e,X,t),pe(T,e,t),r(e,Y,t),r(e,I,t),r(e,Z,t),r(e,$,t),o($,z),o($,be),o($,Q),o(Q,A),o(A,he),o(A,x),o(x,ee),o(A,_e),o(A,te),o(A,ke),o($,ve),o($,U),r(e,le,t),r(e,B,t),r(e,se,t),r(e,q,t),r(e,oe,t),r(e,M,t),r(e,ae,t),r(e,C,t),o(C,ie),o(C,we),o(C,ne),pe(E,ne,null),r(e,re,t),r(e,L,t),r(e,ce,t),r(e,g,t),o(g,D);for(let p=0;p<w.length;p+=1)w[p]&&w[p].m(D,null);o(g,ye),o(g,H);for(let p=0;p<k.length;p+=1)k[p]&&k[p].m(H,null);S=!0},p(e,[t]){var Ie,Be,qe,Me;(!S||t&1)&&i!==(i=e[0].name+"")&&F(b,i),(!S||t&1)&&W!==(W=e[0].name+"")&&F(J,W);const p={};t&9&&(p.js=`
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

        final result = await pb.collection('${(Me=e[0])==null?void 0:Me.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `),T.$set(p),(!S||t&1)&&R!==(R=e[0].name+"")&&F(ee,R),t&6&&(j=e[2],w=Le(w,t,Te,1,e,j,$e,D,Ue,Oe,null,He)),t&6&&(O=e[2],je(),k=Le(k,t,Ae,1,e,O,Pe,H,Fe,We,null,De),Ne())},i(e){if(!S){N(T.$$.fragment,e),N(E.$$.fragment,e);for(let t=0;t<O.length;t+=1)N(k[t]);S=!0}},o(e){G(T.$$.fragment,e),G(E.$$.fragment,e);for(let t=0;t<k.length;t+=1)G(k[t]);S=!1},d(e){e&&c(l),e&&c(d),e&&c(u),e&&c(X),ue(T,e),e&&c(Y),e&&c(I),e&&c(Z),e&&c($),e&&c(le),e&&c(B),e&&c(se),e&&c(q),e&&c(oe),e&&c(M),e&&c(ae),e&&c(C),ue(E),e&&c(re),e&&c(L),e&&c(ce),e&&c(g);for(let t=0;t<w.length;t+=1)w[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function et(a,l,s){let i,{collection:b=new Ge}=l,_=200,d=[];const u=h=>s(1,_=h.code);return a.$$set=h=>{"collection"in h&&s(0,b=h.collection)},a.$$.update=()=>{a.$$.dirty&1&&s(2,d=[{code:200,body:`
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
            `}])},s(3,i=Ke.getApiExampleUrl(Je.baseUrl)),[b,_,d,i,u]}class ot extends ze{constructor(l){super(),Qe(this,l,et,xe,Re,{collection:0})}}export{ot as default};
