import{S as We,i as ze,s as Qe,N as F,e as i,w as v,b as m,c as pe,f as b,g as c,h as a,m as ue,x as N,P as De,Q as je,k as Fe,R as Ne,n as Ge,t as G,a as K,o as d,d as me,U as Ke,C as Je,p as Ve,r as J,u as Xe,M as Ye}from"./index-a084d9d7.js";import{S as Ze}from"./SdkTabs-ba0ec979.js";import{F as xe}from"./FieldsQueryParam-71e01e64.js";function He(o,l,s){const n=o.slice();return n[5]=l[s],n}function Oe(o,l,s){const n=o.slice();return n[5]=l[s],n}function Re(o,l){let s,n=l[5].code+"",f,_,r,u;function h(){return l[4](l[5])}return{key:o,first:null,c(){s=i("button"),f=v(n),_=m(),b(s,"class","tab-item"),J(s,"active",l[1]===l[5].code),this.first=s},m($,y){c($,s,y),a(s,f),a(s,_),r||(u=Xe(s,"click",h),r=!0)},p($,y){l=$,y&4&&n!==(n=l[5].code+"")&&N(f,n),y&6&&J(s,"active",l[1]===l[5].code)},d($){$&&d(s),r=!1,u()}}}function Ue(o,l){let s,n,f,_;return n=new Ye({props:{content:l[5].body}}),{key:o,first:null,c(){s=i("div"),pe(n.$$.fragment),f=m(),b(s,"class","tab-item"),J(s,"active",l[1]===l[5].code),this.first=s},m(r,u){c(r,s,u),ue(n,s,null),a(s,f),_=!0},p(r,u){l=r;const h={};u&4&&(h.content=l[5].body),n.$set(h),(!_||u&6)&&J(s,"active",l[1]===l[5].code)},i(r){_||(G(n.$$.fragment,r),_=!0)},o(r){K(n.$$.fragment,r),_=!1},d(r){r&&d(s),me(n)}}}function et(o){var Ce,Se,Ee,Ie;let l,s,n=o[0].name+"",f,_,r,u,h,$,y,R=o[0].name+"",V,be,fe,X,Y,P,Z,I,x,w,U,he,W,T,_e,ee,z=o[0].name+"",te,ke,le,ve,ge,Q,se,B,ae,q,oe,M,ne,A,ie,we,ce,E,de,L,re,C,D,g=[],$e=new Map,ye,H,k=[],Pe=new Map,S;P=new Ze({props:{js:`
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
    `}}),E=new xe({});let j=F(o[2]);const Te=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=Oe(o,j,e),p=Te(t);$e.set(p,g[e]=Re(p,t))}let O=F(o[2]);const Ae=e=>e[5].code;for(let e=0;e<O.length;e+=1){let t=He(o,O,e),p=Ae(t);Pe.set(p,k[e]=Ue(p,t))}return{c(){l=i("h3"),s=v("List OAuth2 accounts ("),f=v(n),_=v(")"),r=m(),u=i("div"),h=i("p"),$=v("Returns a list with all OAuth2 providers linked to a single "),y=i("strong"),V=v(R),be=v("."),fe=m(),X=i("p"),X.textContent="Only admins and the account owner can access this action.",Y=m(),pe(P.$$.fragment),Z=m(),I=i("h6"),I.textContent="API details",x=m(),w=i("div"),U=i("strong"),U.textContent="GET",he=m(),W=i("div"),T=i("p"),_e=v("/api/collections/"),ee=i("strong"),te=v(z),ke=v("/records/"),le=i("strong"),le.textContent=":id",ve=v("/external-auths"),ge=m(),Q=i("p"),Q.innerHTML="Requires <code>Authorization:TOKEN</code> header",se=m(),B=i("div"),B.textContent="Path Parameters",ae=m(),q=i("table"),q.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="60%">Description</th></tr></thead> <tbody><tr><td>id</td> <td><span class="label">String</span></td> <td>ID of the auth record.</td></tr></tbody>',oe=m(),M=i("div"),M.textContent="Query parameters",ne=m(),A=i("table"),ie=i("thead"),ie.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',we=m(),ce=i("tbody"),pe(E.$$.fragment),de=m(),L=i("div"),L.textContent="Responses",re=m(),C=i("div"),D=i("div");for(let e=0;e<g.length;e+=1)g[e].c();ye=m(),H=i("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(u,"class","content txt-lg m-b-sm"),b(I,"class","m-b-xs"),b(U,"class","label label-primary"),b(W,"class","content"),b(Q,"class","txt-hint txt-sm txt-right"),b(w,"class","alert alert-info"),b(B,"class","section-title"),b(q,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(A,"class","table-compact table-border m-b-base"),b(L,"class","section-title"),b(D,"class","tabs-header compact left"),b(H,"class","tabs-content"),b(C,"class","tabs")},m(e,t){c(e,l,t),a(l,s),a(l,f),a(l,_),c(e,r,t),c(e,u,t),a(u,h),a(h,$),a(h,y),a(y,V),a(h,be),a(u,fe),a(u,X),c(e,Y,t),ue(P,e,t),c(e,Z,t),c(e,I,t),c(e,x,t),c(e,w,t),a(w,U),a(w,he),a(w,W),a(W,T),a(T,_e),a(T,ee),a(ee,te),a(T,ke),a(T,le),a(T,ve),a(w,ge),a(w,Q),c(e,se,t),c(e,B,t),c(e,ae,t),c(e,q,t),c(e,oe,t),c(e,M,t),c(e,ne,t),c(e,A,t),a(A,ie),a(A,we),a(A,ce),ue(E,ce,null),c(e,de,t),c(e,L,t),c(e,re,t),c(e,C,t),a(C,D);for(let p=0;p<g.length;p+=1)g[p]&&g[p].m(D,null);a(C,ye),a(C,H);for(let p=0;p<k.length;p+=1)k[p]&&k[p].m(H,null);S=!0},p(e,[t]){var Be,qe,Me,Le;(!S||t&1)&&n!==(n=e[0].name+"")&&N(f,n),(!S||t&1)&&R!==(R=e[0].name+"")&&N(V,R);const p={};t&9&&(p.js=`
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

        await pb.collection('${(Me=e[0])==null?void 0:Me.name}').authWithPassword('test@example.com', '123456');

        final result = await pb.collection('${(Le=e[0])==null?void 0:Le.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `),P.$set(p),(!S||t&1)&&z!==(z=e[0].name+"")&&N(te,z),t&6&&(j=F(e[2]),g=De(g,t,Te,1,e,j,$e,D,je,Re,null,Oe)),t&6&&(O=F(e[2]),Fe(),k=De(k,t,Ae,1,e,O,Pe,H,Ne,Ue,null,He),Ge())},i(e){if(!S){G(P.$$.fragment,e),G(E.$$.fragment,e);for(let t=0;t<O.length;t+=1)G(k[t]);S=!0}},o(e){K(P.$$.fragment,e),K(E.$$.fragment,e);for(let t=0;t<k.length;t+=1)K(k[t]);S=!1},d(e){e&&(d(l),d(r),d(u),d(Y),d(Z),d(I),d(x),d(w),d(se),d(B),d(ae),d(q),d(oe),d(M),d(ne),d(A),d(de),d(L),d(re),d(C)),me(P,e),me(E);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function tt(o,l,s){let n,{collection:f=new Ke}=l,_=200,r=[];const u=h=>s(1,_=h.code);return o.$$set=h=>{"collection"in h&&s(0,f=h.collection)},o.$$.update=()=>{o.$$.dirty&1&&s(2,r=[{code:200,body:`
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
            `}])},s(3,n=Je.getApiExampleUrl(Ve.baseUrl)),[f,_,r,n,u]}class ot extends We{constructor(l){super(),ze(this,l,tt,et,Qe,{collection:0})}}export{ot as default};
