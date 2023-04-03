import{S as Be,i as qe,s as Oe,e as i,w as v,b as _,c as Se,f as b,g as d,h as s,m as Ee,x as U,O as Pe,P as Le,k as Me,Q as We,n as ze,t as te,a as le,o as u,d as Ie,T as De,C as He,p as Re,r as j,u as Ue,N as je}from"./index-4eea3e34.js";import{S as Ne}from"./SdkTabs-5d6cc1d4.js";function ye(a,l,o){const n=a.slice();return n[5]=l[o],n}function Ae(a,l,o){const n=a.slice();return n[5]=l[o],n}function Te(a,l){let o,n=l[5].code+"",f,h,c,p;function m(){return l[4](l[5])}return{key:a,first:null,c(){o=i("button"),f=v(n),h=_(),b(o,"class","tab-item"),j(o,"active",l[1]===l[5].code),this.first=o},m(g,P){d(g,o,P),s(o,f),s(o,h),c||(p=Ue(o,"click",m),c=!0)},p(g,P){l=g,P&4&&n!==(n=l[5].code+"")&&U(f,n),P&6&&j(o,"active",l[1]===l[5].code)},d(g){g&&u(o),c=!1,p()}}}function Ce(a,l){let o,n,f,h;return n=new je({props:{content:l[5].body}}),{key:a,first:null,c(){o=i("div"),Se(n.$$.fragment),f=_(),b(o,"class","tab-item"),j(o,"active",l[1]===l[5].code),this.first=o},m(c,p){d(c,o,p),Ee(n,o,null),s(o,f),h=!0},p(c,p){l=c;const m={};p&4&&(m.content=l[5].body),n.$set(m),(!h||p&6)&&j(o,"active",l[1]===l[5].code)},i(c){h||(te(n.$$.fragment,c),h=!0)},o(c){le(n.$$.fragment,c),h=!1},d(c){c&&u(o),Ie(n)}}}function Ge(a){var be,he,_e,ke;let l,o,n=a[0].name+"",f,h,c,p,m,g,P,M=a[0].name+"",N,oe,se,G,K,y,Q,S,F,$,W,ae,z,A,ne,J,D=a[0].name+"",V,ie,X,ce,re,H,Y,E,Z,I,x,B,ee,T,q,w=[],de=new Map,ue,O,k=[],pe=new Map,C;y=new Ne({props:{js:`
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
    `}});let R=a[2];const fe=e=>e[5].code;for(let e=0;e<R.length;e+=1){let t=Ae(a,R,e),r=fe(t);de.set(r,w[e]=Te(r,t))}let L=a[2];const me=e=>e[5].code;for(let e=0;e<L.length;e+=1){let t=ye(a,L,e),r=me(t);pe.set(r,k[e]=Ce(r,t))}return{c(){l=i("h3"),o=v("List OAuth2 accounts ("),f=v(n),h=v(")"),c=_(),p=i("div"),m=i("p"),g=v("Returns a list with all OAuth2 providers linked to a single "),P=i("strong"),N=v(M),oe=v("."),se=_(),G=i("p"),G.textContent="Only admins and the account owner can access this action.",K=_(),Se(y.$$.fragment),Q=_(),S=i("h6"),S.textContent="API details",F=_(),$=i("div"),W=i("strong"),W.textContent="GET",ae=_(),z=i("div"),A=i("p"),ne=v("/api/collections/"),J=i("strong"),V=v(D),ie=v("/records/"),X=i("strong"),X.textContent=":id",ce=v("/external-auths"),re=_(),H=i("p"),H.innerHTML="Requires <code>Authorization:TOKEN</code> header",Y=_(),E=i("div"),E.textContent="Path Parameters",Z=_(),I=i("table"),I.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr></thead> 
    <tbody><tr><td>id</td> 
            <td><span class="label">String</span></td> 
            <td>ID of the auth record.</td></tr></tbody>`,x=_(),B=i("div"),B.textContent="Responses",ee=_(),T=i("div"),q=i("div");for(let e=0;e<w.length;e+=1)w[e].c();ue=_(),O=i("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(S,"class","m-b-xs"),b(W,"class","label label-primary"),b(z,"class","content"),b(H,"class","txt-hint txt-sm txt-right"),b($,"class","alert alert-info"),b(E,"class","section-title"),b(I,"class","table-compact table-border m-b-base"),b(B,"class","section-title"),b(q,"class","tabs-header compact left"),b(O,"class","tabs-content"),b(T,"class","tabs")},m(e,t){d(e,l,t),s(l,o),s(l,f),s(l,h),d(e,c,t),d(e,p,t),s(p,m),s(m,g),s(m,P),s(P,N),s(m,oe),s(p,se),s(p,G),d(e,K,t),Ee(y,e,t),d(e,Q,t),d(e,S,t),d(e,F,t),d(e,$,t),s($,W),s($,ae),s($,z),s(z,A),s(A,ne),s(A,J),s(J,V),s(A,ie),s(A,X),s(A,ce),s($,re),s($,H),d(e,Y,t),d(e,E,t),d(e,Z,t),d(e,I,t),d(e,x,t),d(e,B,t),d(e,ee,t),d(e,T,t),s(T,q);for(let r=0;r<w.length;r+=1)w[r]&&w[r].m(q,null);s(T,ue),s(T,O);for(let r=0;r<k.length;r+=1)k[r]&&k[r].m(O,null);C=!0},p(e,[t]){var ve,we,$e,ge;(!C||t&1)&&n!==(n=e[0].name+"")&&U(f,n),(!C||t&1)&&M!==(M=e[0].name+"")&&U(N,M);const r={};t&9&&(r.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ve=e[0])==null?void 0:ve.name}').authWithPassword('test@example.com', '123456');

        const result = await pb.collection('${(we=e[0])==null?void 0:we.name}').listExternalAuths(
            pb.authStore.model.id
        );
    `),t&9&&(r.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${($e=e[0])==null?void 0:$e.name}').authWithPassword('test@example.com', '123456');

        final result = await pb.collection('${(ge=e[0])==null?void 0:ge.name}').listExternalAuths(
          pb.authStore.model.id,
        );
    `),y.$set(r),(!C||t&1)&&D!==(D=e[0].name+"")&&U(V,D),t&6&&(R=e[2],w=Pe(w,t,fe,1,e,R,de,q,Le,Te,null,Ae)),t&6&&(L=e[2],Me(),k=Pe(k,t,me,1,e,L,pe,O,We,Ce,null,ye),ze())},i(e){if(!C){te(y.$$.fragment,e);for(let t=0;t<L.length;t+=1)te(k[t]);C=!0}},o(e){le(y.$$.fragment,e);for(let t=0;t<k.length;t+=1)le(k[t]);C=!1},d(e){e&&u(l),e&&u(c),e&&u(p),e&&u(K),Ie(y,e),e&&u(Q),e&&u(S),e&&u(F),e&&u($),e&&u(Y),e&&u(E),e&&u(Z),e&&u(I),e&&u(x),e&&u(B),e&&u(ee),e&&u(T);for(let t=0;t<w.length;t+=1)w[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function Ke(a,l,o){let n,{collection:f=new De}=l,h=200,c=[];const p=m=>o(1,h=m.code);return a.$$set=m=>{"collection"in m&&o(0,f=m.collection)},a.$$.update=()=>{a.$$.dirty&1&&o(2,c=[{code:200,body:`
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
            `}])},o(3,n=He.getApiExampleUrl(Re.baseUrl)),[f,h,c,n,p]}class Je extends Be{constructor(l){super(),qe(this,l,Ke,Ge,Oe,{collection:0})}}export{Je as default};
