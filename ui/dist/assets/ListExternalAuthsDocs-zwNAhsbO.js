import{S as ze,i as Qe,s as Ue,O as F,e as i,w as v,b as m,c as pe,f as b,g as c,h as a,m as ue,x as N,P as Oe,Q as je,k as Fe,R as Ne,n as Ge,t as G,a as K,o as d,d as me,C as Ke,p as Je,r as J,u as Ve,N as Xe}from"./index-78piLIP3.js";import{S as Ye}from"./SdkTabs-c6VuPJvR.js";import{F as Ze}from"./FieldsQueryParam-oYzijp1d.js";function De(o,l,s){const n=o.slice();return n[5]=l[s],n}function He(o,l,s){const n=o.slice();return n[5]=l[s],n}function Re(o,l){let s,n=l[5].code+"",f,_,r,u;function h(){return l[4](l[5])}return{key:o,first:null,c(){s=i("button"),f=v(n),_=m(),b(s,"class","tab-item"),J(s,"active",l[1]===l[5].code),this.first=s},m(w,y){c(w,s,y),a(s,f),a(s,_),r||(u=Ve(s,"click",h),r=!0)},p(w,y){l=w,y&4&&n!==(n=l[5].code+"")&&N(f,n),y&6&&J(s,"active",l[1]===l[5].code)},d(w){w&&d(s),r=!1,u()}}}function We(o,l){let s,n,f,_;return n=new Xe({props:{content:l[5].body}}),{key:o,first:null,c(){s=i("div"),pe(n.$$.fragment),f=m(),b(s,"class","tab-item"),J(s,"active",l[1]===l[5].code),this.first=s},m(r,u){c(r,s,u),ue(n,s,null),a(s,f),_=!0},p(r,u){l=r;const h={};u&4&&(h.content=l[5].body),n.$set(h),(!_||u&6)&&J(s,"active",l[1]===l[5].code)},i(r){_||(G(n.$$.fragment,r),_=!0)},o(r){K(n.$$.fragment,r),_=!1},d(r){r&&d(s),me(n)}}}function xe(o){var Ce,Se,Ee,Ie;let l,s,n=o[0].name+"",f,_,r,u,h,w,y,R=o[0].name+"",V,be,fe,X,Y,P,Z,I,x,$,W,he,z,T,_e,ee,Q=o[0].name+"",te,ke,le,ve,ge,U,se,B,ae,q,oe,L,ne,A,ie,$e,ce,E,de,M,re,C,O,g=[],we=new Map,ye,D,k=[],Pe=new Map,S;P=new Ye({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(Ce=o[0])==null?void 0:Ce.name}').authWithPassword('test@example.com', '123456');

        const result = await pb.collection('${(Se=o[0])==null?void 0:Se.name}').listExternalAuths(
            pb.authStore.model.id
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(Ee=o[0])==null?void 0:Ee.name}').authWithPassword('test@example.com', '123456');

        final result = await pb.collection('${(Ie=o[0])==null?void 0:Ie.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `}}),E=new Ze({});let j=F(o[2]);const Te=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=He(o,j,e),p=Te(t);we.set(p,g[e]=Re(p,t))}let H=F(o[2]);const Ae=e=>e[5].code;for(let e=0;e<H.length;e+=1){let t=De(o,H,e),p=Ae(t);Pe.set(p,k[e]=We(p,t))}return{c(){l=i("h3"),s=v("List OAuth2 accounts ("),f=v(n),_=v(")"),r=m(),u=i("div"),h=i("p"),w=v("Returns a list with all OAuth2 providers linked to a single "),y=i("strong"),V=v(R),be=v("."),fe=m(),X=i("p"),X.textContent="Only admins and the account owner can access this action.",Y=m(),pe(P.$$.fragment),Z=m(),I=i("h6"),I.textContent="API details",x=m(),$=i("div"),W=i("strong"),W.textContent="GET",he=m(),z=i("div"),T=i("p"),_e=v("/api/collections/"),ee=i("strong"),te=v(Q),ke=v("/records/"),le=i("strong"),le.textContent=":id",ve=v("/external-auths"),ge=m(),U=i("p"),U.innerHTML="Requires <code>Authorization:TOKEN</code> header",se=m(),B=i("div"),B.textContent="Path Parameters",ae=m(),q=i("table"),q.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the auth record.</td></tr></tbody>',oe=m(),L=i("div"),L.textContent="Query parameters",ne=m(),A=i("table"),ie=i("thead"),ie.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',$e=m(),ce=i("tbody"),pe(E.$$.fragment),de=m(),M=i("div"),M.textContent="Responses",re=m(),C=i("div"),O=i("div");for(let e=0;e<g.length;e+=1)g[e].c();ye=m(),D=i("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(u,"class","content txt-lg m-b-sm"),b(I,"class","m-b-xs"),b(W,"class","label label-primary"),b(z,"class","content"),b(U,"class","txt-hint txt-sm txt-right"),b($,"class","alert alert-info"),b(B,"class","section-title"),b(q,"class","table-compact table-border m-b-base"),b(L,"class","section-title"),b(A,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(O,"class","tabs-header compact combined left"),b(D,"class","tabs-content"),b(C,"class","tabs")},m(e,t){c(e,l,t),a(l,s),a(l,f),a(l,_),c(e,r,t),c(e,u,t),a(u,h),a(h,w),a(h,y),a(y,V),a(h,be),a(u,fe),a(u,X),c(e,Y,t),ue(P,e,t),c(e,Z,t),c(e,I,t),c(e,x,t),c(e,$,t),a($,W),a($,he),a($,z),a(z,T),a(T,_e),a(T,ee),a(ee,te),a(T,ke),a(T,le),a(T,ve),a($,ge),a($,U),c(e,se,t),c(e,B,t),c(e,ae,t),c(e,q,t),c(e,oe,t),c(e,L,t),c(e,ne,t),c(e,A,t),a(A,ie),a(A,$e),a(A,ce),ue(E,ce,null),c(e,de,t),c(e,M,t),c(e,re,t),c(e,C,t),a(C,O);for(let p=0;p<g.length;p+=1)g[p]&&g[p].m(O,null);a(C,ye),a(C,D);for(let p=0;p<k.length;p+=1)k[p]&&k[p].m(D,null);S=!0},p(e,[t]){var Be,qe,Le,Me;(!S||t&1)&&n!==(n=e[0].name+"")&&N(f,n),(!S||t&1)&&R!==(R=e[0].name+"")&&N(V,R);const p={};t&9&&(p.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(Be=e[0])==null?void 0:Be.name}').authWithPassword('test@example.com', '123456');

        const result = await pb.collection('${(qe=e[0])==null?void 0:qe.name}').listExternalAuths(
            pb.authStore.model.id
        );
    `),t&9&&(p.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(Le=e[0])==null?void 0:Le.name}').authWithPassword('test@example.com', '123456');

        final result = await pb.collection('${(Me=e[0])==null?void 0:Me.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `),P.$set(p),(!S||t&1)&&Q!==(Q=e[0].name+"")&&N(te,Q),t&6&&(j=F(e[2]),g=Oe(g,t,Te,1,e,j,we,O,je,Re,null,He)),t&6&&(H=F(e[2]),Fe(),k=Oe(k,t,Ae,1,e,H,Pe,D,Ne,We,null,De),Ge())},i(e){if(!S){G(P.$$.fragment,e),G(E.$$.fragment,e);for(let t=0;t<H.length;t+=1)G(k[t]);S=!0}},o(e){K(P.$$.fragment,e),K(E.$$.fragment,e);for(let t=0;t<k.length;t+=1)K(k[t]);S=!1},d(e){e&&(d(l),d(r),d(u),d(Y),d(Z),d(I),d(x),d($),d(se),d(B),d(ae),d(q),d(oe),d(L),d(ne),d(A),d(de),d(M),d(re),d(C)),me(P,e),me(E);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function et(o,l,s){let n,{collection:f}=l,_=200,r=[];const u=h=>s(1,_=h.code);return o.$$set=h=>{"collection"in h&&s(0,f=h.collection)},o.$$.update=()=>{o.$$.dirty&1&&s(2,r=[{code:200,body:`
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
            `}])},s(3,n=Ke.getApiExampleUrl(Je.baseUrl)),[f,_,r,n,u]}class at extends ze{constructor(l){super(),Qe(this,l,et,xe,Ue,{collection:0})}}export{at as default};
