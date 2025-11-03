import{S as Ce,i as Be,s as Te,ab as Le,O as J,h as d,d as ae,t as O,a as Q,I as G,ad as ye,ae as Se,C as De,af as Re,D as Ue,l as u,n as a,m as ne,u as c,A as $,v as k,c as ie,w as h,p as oe,J as je,k as N,o as qe,ac as Ee}from"./index-B4ZsHsKR.js";import{F as Fe}from"./FieldsQueryParam-K1y4zYh0.js";function $e(n,s,l){const o=n.slice();return o[8]=s[l],o}function Me(n,s,l){const o=n.slice();return o[8]=s[l],o}function Ae(n,s){let l,o=s[8].code+"",p,b,i,f;function m(){return s[6](s[8])}return{key:n,first:null,c(){l=c("button"),p=$(o),b=k(),h(l,"class","tab-item"),N(l,"active",s[1]===s[8].code),this.first=l},m(v,w){u(v,l,w),a(l,p),a(l,b),i||(f=qe(l,"click",m),i=!0)},p(v,w){s=v,w&4&&o!==(o=s[8].code+"")&&G(p,o),w&6&&N(l,"active",s[1]===s[8].code)},d(v){v&&d(l),i=!1,f()}}}function Pe(n,s){let l,o,p,b;return o=new Ee({props:{content:s[8].body}}),{key:n,first:null,c(){l=c("div"),ie(o.$$.fragment),p=k(),h(l,"class","tab-item"),N(l,"active",s[1]===s[8].code),this.first=l},m(i,f){u(i,l,f),ne(o,l,null),a(l,p),b=!0},p(i,f){s=i;const m={};f&4&&(m.content=s[8].body),o.$set(m),(!b||f&6)&&N(l,"active",s[1]===s[8].code)},i(i){b||(Q(o.$$.fragment,i),b=!0)},o(i){O(o.$$.fragment,i),b=!1},d(i){i&&d(l),ae(o)}}}function He(n){var ke,ge;let s,l,o=n[0].name+"",p,b,i,f,m,v,w,g=n[0].name+"",z,ce,K,M,V,L,W,A,E,re,F,S,de,X,H=n[0].name+"",Y,ue,Z,D,x,P,ee,fe,te,T,le,R,se,C,U,y=[],me=new Map,pe,j,_=[],be=new Map,B;M=new Le({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[3]}');

        ...

        const result = await pb.collection('${(ke=n[0])==null?void 0:ke.name}').listAuthMethods();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[3]}');

        ...

        final result = await pb.collection('${(ge=n[0])==null?void 0:ge.name}').listAuthMethods();
    `}}),T=new Fe({});let I=J(n[2]);const he=e=>e[8].code;for(let e=0;e<I.length;e+=1){let t=Me(n,I,e),r=he(t);me.set(r,y[e]=Ae(r,t))}let q=J(n[2]);const _e=e=>e[8].code;for(let e=0;e<q.length;e+=1){let t=$e(n,q,e),r=_e(t);be.set(r,_[e]=Pe(r,t))}return{c(){s=c("h3"),l=$("List auth methods ("),p=$(o),b=$(")"),i=k(),f=c("div"),m=c("p"),v=$("Returns a public list with all allowed "),w=c("strong"),z=$(g),ce=$(" authentication methods."),K=k(),ie(M.$$.fragment),V=k(),L=c("h6"),L.textContent="API details",W=k(),A=c("div"),E=c("strong"),E.textContent="GET",re=k(),F=c("div"),S=c("p"),de=$("/api/collections/"),X=c("strong"),Y=$(H),ue=$("/auth-methods"),Z=k(),D=c("div"),D.textContent="Query parameters",x=k(),P=c("table"),ee=c("thead"),ee.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',fe=k(),te=c("tbody"),ie(T.$$.fragment),le=k(),R=c("div"),R.textContent="Responses",se=k(),C=c("div"),U=c("div");for(let e=0;e<y.length;e+=1)y[e].c();pe=k(),j=c("div");for(let e=0;e<_.length;e+=1)_[e].c();h(s,"class","m-b-sm"),h(f,"class","content txt-lg m-b-sm"),h(L,"class","m-b-xs"),h(E,"class","label label-primary"),h(F,"class","content"),h(A,"class","alert alert-info"),h(D,"class","section-title"),h(P,"class","table-compact table-border m-b-base"),h(R,"class","section-title"),h(U,"class","tabs-header compact combined left"),h(j,"class","tabs-content"),h(C,"class","tabs")},m(e,t){u(e,s,t),a(s,l),a(s,p),a(s,b),u(e,i,t),u(e,f,t),a(f,m),a(m,v),a(m,w),a(w,z),a(m,ce),u(e,K,t),ne(M,e,t),u(e,V,t),u(e,L,t),u(e,W,t),u(e,A,t),a(A,E),a(A,re),a(A,F),a(F,S),a(S,de),a(S,X),a(X,Y),a(S,ue),u(e,Z,t),u(e,D,t),u(e,x,t),u(e,P,t),a(P,ee),a(P,fe),a(P,te),ne(T,te,null),u(e,le,t),u(e,R,t),u(e,se,t),u(e,C,t),a(C,U);for(let r=0;r<y.length;r+=1)y[r]&&y[r].m(U,null);a(C,pe),a(C,j);for(let r=0;r<_.length;r+=1)_[r]&&_[r].m(j,null);B=!0},p(e,[t]){var ve,we;(!B||t&1)&&o!==(o=e[0].name+"")&&G(p,o),(!B||t&1)&&g!==(g=e[0].name+"")&&G(z,g);const r={};t&9&&(r.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const result = await pb.collection('${(ve=e[0])==null?void 0:ve.name}').listAuthMethods();
    `),t&9&&(r.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final result = await pb.collection('${(we=e[0])==null?void 0:we.name}').listAuthMethods();
    `),M.$set(r),(!B||t&1)&&H!==(H=e[0].name+"")&&G(Y,H),t&6&&(I=J(e[2]),y=ye(y,t,he,1,e,I,me,U,Se,Ae,null,Me)),t&6&&(q=J(e[2]),De(),_=ye(_,t,_e,1,e,q,be,j,Re,Pe,null,$e),Ue())},i(e){if(!B){Q(M.$$.fragment,e),Q(T.$$.fragment,e);for(let t=0;t<q.length;t+=1)Q(_[t]);B=!0}},o(e){O(M.$$.fragment,e),O(T.$$.fragment,e);for(let t=0;t<_.length;t+=1)O(_[t]);B=!1},d(e){e&&(d(s),d(i),d(f),d(K),d(V),d(L),d(W),d(A),d(Z),d(D),d(x),d(P),d(le),d(R),d(se),d(C)),ae(M,e),ae(T);for(let t=0;t<y.length;t+=1)y[t].d();for(let t=0;t<_.length;t+=1)_[t].d()}}}function Ie(n,s,l){let o,{collection:p}=s,b=200,i=[],f={},m=!1;v();async function v(){l(5,m=!0);try{l(4,f=await oe.collection(p.name).listAuthMethods())}catch(g){oe.error(g)}l(5,m=!1)}const w=g=>l(1,b=g.code);return n.$$set=g=>{"collection"in g&&l(0,p=g.collection)},n.$$.update=()=>{n.$$.dirty&48&&l(2,i=[{code:200,body:m?"...":JSON.stringify(f,null,2)},{code:404,body:`
                {
                  "status": 404,
                  "message": "Missing collection context.",
                  "data": {}
                }
            `}])},l(3,o=je.getApiExampleUrl(oe.baseURL)),[p,b,i,o,f,m,w]}class Qe extends Ce{constructor(s){super(),Be(this,s,Ie,He,Te,{collection:0})}}export{Qe as default};
