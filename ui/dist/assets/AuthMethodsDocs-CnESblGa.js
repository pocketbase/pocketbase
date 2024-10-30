import{S as Ce,i as Be,s as Te,Q as Le,T as G,e as c,w as $,b as k,c as ae,f as h,g as d,h as a,m as ne,x as I,U as ye,V as Re,k as Se,W as Ue,n as Qe,t as J,a as N,o as u,d as ie,p as oe,C as je,r as O,u as qe,R as De}from"./index-jQhhdClJ.js";import{F as Ee}from"./FieldsQueryParam-Cv9q4pej.js";function $e(n,s,l){const o=n.slice();return o[8]=s[l],o}function Me(n,s,l){const o=n.slice();return o[8]=s[l],o}function Ae(n,s){let l,o=s[8].code+"",p,b,i,f;function m(){return s[6](s[8])}return{key:n,first:null,c(){l=c("button"),p=$(o),b=k(),h(l,"class","tab-item"),O(l,"active",s[1]===s[8].code),this.first=l},m(v,w){d(v,l,w),a(l,p),a(l,b),i||(f=qe(l,"click",m),i=!0)},p(v,w){s=v,w&4&&o!==(o=s[8].code+"")&&I(p,o),w&6&&O(l,"active",s[1]===s[8].code)},d(v){v&&u(l),i=!1,f()}}}function Pe(n,s){let l,o,p,b;return o=new De({props:{content:s[8].body}}),{key:n,first:null,c(){l=c("div"),ae(o.$$.fragment),p=k(),h(l,"class","tab-item"),O(l,"active",s[1]===s[8].code),this.first=l},m(i,f){d(i,l,f),ne(o,l,null),a(l,p),b=!0},p(i,f){s=i;const m={};f&4&&(m.content=s[8].body),o.$set(m),(!b||f&6)&&O(l,"active",s[1]===s[8].code)},i(i){b||(J(o.$$.fragment,i),b=!0)},o(i){N(o.$$.fragment,i),b=!1},d(i){i&&u(l),ie(o)}}}function Fe(n){var ke,ge;let s,l,o=n[0].name+"",p,b,i,f,m,v,w,g=n[0].name+"",V,ce,W,M,z,L,K,A,D,re,E,R,de,X,F=n[0].name+"",Y,ue,Z,S,x,P,ee,fe,te,T,le,U,se,C,Q,y=[],me=new Map,pe,j,_=[],be=new Map,B;M=new Le({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[3]}');

        ...

        const result = await pb.collection('${(ke=n[0])==null?void 0:ke.name}').listAuthMethods();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[3]}');

        ...

        final result = await pb.collection('${(ge=n[0])==null?void 0:ge.name}').listAuthMethods();
    `}}),T=new Ee({});let H=G(n[2]);const he=e=>e[8].code;for(let e=0;e<H.length;e+=1){let t=Me(n,H,e),r=he(t);me.set(r,y[e]=Ae(r,t))}let q=G(n[2]);const _e=e=>e[8].code;for(let e=0;e<q.length;e+=1){let t=$e(n,q,e),r=_e(t);be.set(r,_[e]=Pe(r,t))}return{c(){s=c("h3"),l=$("List auth methods ("),p=$(o),b=$(")"),i=k(),f=c("div"),m=c("p"),v=$("Returns a public list with all allowed "),w=c("strong"),V=$(g),ce=$(" authentication methods."),W=k(),ae(M.$$.fragment),z=k(),L=c("h6"),L.textContent="API details",K=k(),A=c("div"),D=c("strong"),D.textContent="GET",re=k(),E=c("div"),R=c("p"),de=$("/api/collections/"),X=c("strong"),Y=$(F),ue=$("/auth-methods"),Z=k(),S=c("div"),S.textContent="Query parameters",x=k(),P=c("table"),ee=c("thead"),ee.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',fe=k(),te=c("tbody"),ae(T.$$.fragment),le=k(),U=c("div"),U.textContent="Responses",se=k(),C=c("div"),Q=c("div");for(let e=0;e<y.length;e+=1)y[e].c();pe=k(),j=c("div");for(let e=0;e<_.length;e+=1)_[e].c();h(s,"class","m-b-sm"),h(f,"class","content txt-lg m-b-sm"),h(L,"class","m-b-xs"),h(D,"class","label label-primary"),h(E,"class","content"),h(A,"class","alert alert-info"),h(S,"class","section-title"),h(P,"class","table-compact table-border m-b-base"),h(U,"class","section-title"),h(Q,"class","tabs-header compact combined left"),h(j,"class","tabs-content"),h(C,"class","tabs")},m(e,t){d(e,s,t),a(s,l),a(s,p),a(s,b),d(e,i,t),d(e,f,t),a(f,m),a(m,v),a(m,w),a(w,V),a(m,ce),d(e,W,t),ne(M,e,t),d(e,z,t),d(e,L,t),d(e,K,t),d(e,A,t),a(A,D),a(A,re),a(A,E),a(E,R),a(R,de),a(R,X),a(X,Y),a(R,ue),d(e,Z,t),d(e,S,t),d(e,x,t),d(e,P,t),a(P,ee),a(P,fe),a(P,te),ne(T,te,null),d(e,le,t),d(e,U,t),d(e,se,t),d(e,C,t),a(C,Q);for(let r=0;r<y.length;r+=1)y[r]&&y[r].m(Q,null);a(C,pe),a(C,j);for(let r=0;r<_.length;r+=1)_[r]&&_[r].m(j,null);B=!0},p(e,[t]){var ve,we;(!B||t&1)&&o!==(o=e[0].name+"")&&I(p,o),(!B||t&1)&&g!==(g=e[0].name+"")&&I(V,g);const r={};t&9&&(r.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const result = await pb.collection('${(ve=e[0])==null?void 0:ve.name}').listAuthMethods();
    `),t&9&&(r.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final result = await pb.collection('${(we=e[0])==null?void 0:we.name}').listAuthMethods();
    `),M.$set(r),(!B||t&1)&&F!==(F=e[0].name+"")&&I(Y,F),t&6&&(H=G(e[2]),y=ye(y,t,he,1,e,H,me,Q,Re,Ae,null,Me)),t&6&&(q=G(e[2]),Se(),_=ye(_,t,_e,1,e,q,be,j,Ue,Pe,null,$e),Qe())},i(e){if(!B){J(M.$$.fragment,e),J(T.$$.fragment,e);for(let t=0;t<q.length;t+=1)J(_[t]);B=!0}},o(e){N(M.$$.fragment,e),N(T.$$.fragment,e);for(let t=0;t<_.length;t+=1)N(_[t]);B=!1},d(e){e&&(u(s),u(i),u(f),u(W),u(z),u(L),u(K),u(A),u(Z),u(S),u(x),u(P),u(le),u(U),u(se),u(C)),ie(M,e),ie(T);for(let t=0;t<y.length;t+=1)y[t].d();for(let t=0;t<_.length;t+=1)_[t].d()}}}function He(n,s,l){let o,{collection:p}=s,b=200,i=[],f={},m=!1;v();async function v(){l(5,m=!0);try{l(4,f=await oe.collection(p.name).listAuthMethods())}catch(g){oe.error(g)}l(5,m=!1)}const w=g=>l(1,b=g.code);return n.$$set=g=>{"collection"in g&&l(0,p=g.collection)},n.$$.update=()=>{n.$$.dirty&48&&l(2,i=[{code:200,body:m?"...":JSON.stringify(f,null,2)},{code:404,body:`
                {
                  "code": 404,
                  "message": "Missing collection context.",
                  "data": {}
                }
            `}])},l(3,o=je.getApiExampleUrl(oe.baseURL)),[p,b,i,o,f,m,w]}class Je extends Ce{constructor(s){super(),Be(this,s,He,Fe,Te,{collection:0})}}export{Je as default};
