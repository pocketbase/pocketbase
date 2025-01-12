import{S as Be,i as Ce,s as Te,V as Le,X as J,h as c,z as w,j as k,c as ae,k as h,n as d,o as a,m as ne,H as Q,Y as $e,Z as Se,E as je,_ as qe,G as Ee,t as z,a as I,v as u,d as ie,p as oe,J as He,l as N,q as Re,W as Ue}from"./index-SKn09NMF.js";import{F as De}from"./FieldsQueryParam-NXRpAlxi.js";function we(n,s,l){const o=n.slice();return o[8]=s[l],o}function Me(n,s,l){const o=n.slice();return o[8]=s[l],o}function Ae(n,s){let l,o=s[8].code+"",p,b,i,f;function m(){return s[6](s[8])}return{key:n,first:null,c(){l=c("button"),p=w(o),b=k(),h(l,"class","tab-item"),N(l,"active",s[1]===s[8].code),this.first=l},m(v,y){d(v,l,y),a(l,p),a(l,b),i||(f=Re(l,"click",m),i=!0)},p(v,y){s=v,y&4&&o!==(o=s[8].code+"")&&Q(p,o),y&6&&N(l,"active",s[1]===s[8].code)},d(v){v&&u(l),i=!1,f()}}}function Pe(n,s){let l,o,p,b;return o=new Ue({props:{content:s[8].body}}),{key:n,first:null,c(){l=c("div"),ae(o.$$.fragment),p=k(),h(l,"class","tab-item"),N(l,"active",s[1]===s[8].code),this.first=l},m(i,f){d(i,l,f),ne(o,l,null),a(l,p),b=!0},p(i,f){s=i;const m={};f&4&&(m.content=s[8].body),o.$set(m),(!b||f&6)&&N(l,"active",s[1]===s[8].code)},i(i){b||(z(o.$$.fragment,i),b=!0)},o(i){I(o.$$.fragment,i),b=!1},d(i){i&&u(l),ie(o)}}}function Fe(n){var ke,ge;let s,l,o=n[0].name+"",p,b,i,f,m,v,y,g=n[0].name+"",O,ce,V,M,W,L,X,A,U,re,D,S,de,Y,F=n[0].name+"",Z,ue,K,j,x,P,ee,fe,te,T,le,q,se,B,E,$=[],me=new Map,pe,H,_=[],be=new Map,C;M=new Le({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[3]}');

        ...

        const result = await pb.collection('${(ke=n[0])==null?void 0:ke.name}').listAuthMethods();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[3]}');

        ...

        final result = await pb.collection('${(ge=n[0])==null?void 0:ge.name}').listAuthMethods();
    `}}),T=new De({});let G=J(n[2]);const he=e=>e[8].code;for(let e=0;e<G.length;e+=1){let t=Me(n,G,e),r=he(t);me.set(r,$[e]=Ae(r,t))}let R=J(n[2]);const _e=e=>e[8].code;for(let e=0;e<R.length;e+=1){let t=we(n,R,e),r=_e(t);be.set(r,_[e]=Pe(r,t))}return{c(){s=c("h3"),l=w("List auth methods ("),p=w(o),b=w(")"),i=k(),f=c("div"),m=c("p"),v=w("Returns a public list with all allowed "),y=c("strong"),O=w(g),ce=w(" authentication methods."),V=k(),ae(M.$$.fragment),W=k(),L=c("h6"),L.textContent="API details",X=k(),A=c("div"),U=c("strong"),U.textContent="GET",re=k(),D=c("div"),S=c("p"),de=w("/api/collections/"),Y=c("strong"),Z=w(F),ue=w("/auth-methods"),K=k(),j=c("div"),j.textContent="Query parameters",x=k(),P=c("table"),ee=c("thead"),ee.innerHTML='<tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr>',fe=k(),te=c("tbody"),ae(T.$$.fragment),le=k(),q=c("div"),q.textContent="Responses",se=k(),B=c("div"),E=c("div");for(let e=0;e<$.length;e+=1)$[e].c();pe=k(),H=c("div");for(let e=0;e<_.length;e+=1)_[e].c();h(s,"class","m-b-sm"),h(f,"class","content txt-lg m-b-sm"),h(L,"class","m-b-xs"),h(U,"class","label label-primary"),h(D,"class","content"),h(A,"class","alert alert-info"),h(j,"class","section-title"),h(P,"class","table-compact table-border m-b-base"),h(q,"class","section-title"),h(E,"class","tabs-header compact combined left"),h(H,"class","tabs-content"),h(B,"class","tabs")},m(e,t){d(e,s,t),a(s,l),a(s,p),a(s,b),d(e,i,t),d(e,f,t),a(f,m),a(m,v),a(m,y),a(y,O),a(m,ce),d(e,V,t),ne(M,e,t),d(e,W,t),d(e,L,t),d(e,X,t),d(e,A,t),a(A,U),a(A,re),a(A,D),a(D,S),a(S,de),a(S,Y),a(Y,Z),a(S,ue),d(e,K,t),d(e,j,t),d(e,x,t),d(e,P,t),a(P,ee),a(P,fe),a(P,te),ne(T,te,null),d(e,le,t),d(e,q,t),d(e,se,t),d(e,B,t),a(B,E);for(let r=0;r<$.length;r+=1)$[r]&&$[r].m(E,null);a(B,pe),a(B,H);for(let r=0;r<_.length;r+=1)_[r]&&_[r].m(H,null);C=!0},p(e,[t]){var ve,ye;(!C||t&1)&&o!==(o=e[0].name+"")&&Q(p,o),(!C||t&1)&&g!==(g=e[0].name+"")&&Q(O,g);const r={};t&9&&(r.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const result = await pb.collection('${(ve=e[0])==null?void 0:ve.name}').listAuthMethods();
    `),t&9&&(r.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final result = await pb.collection('${(ye=e[0])==null?void 0:ye.name}').listAuthMethods();
    `),M.$set(r),(!C||t&1)&&F!==(F=e[0].name+"")&&Q(Z,F),t&6&&(G=J(e[2]),$=$e($,t,he,1,e,G,me,E,Se,Ae,null,Me)),t&6&&(R=J(e[2]),je(),_=$e(_,t,_e,1,e,R,be,H,qe,Pe,null,we),Ee())},i(e){if(!C){z(M.$$.fragment,e),z(T.$$.fragment,e);for(let t=0;t<R.length;t+=1)z(_[t]);C=!0}},o(e){I(M.$$.fragment,e),I(T.$$.fragment,e);for(let t=0;t<_.length;t+=1)I(_[t]);C=!1},d(e){e&&(u(s),u(i),u(f),u(V),u(W),u(L),u(X),u(A),u(K),u(j),u(x),u(P),u(le),u(q),u(se),u(B)),ie(M,e),ie(T);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<_.length;t+=1)_[t].d()}}}function Ge(n,s,l){let o,{collection:p}=s,b=200,i=[],f={},m=!1;v();async function v(){l(5,m=!0);try{l(4,f=await oe.collection(p.name).listAuthMethods())}catch(g){oe.error(g)}l(5,m=!1)}const y=g=>l(1,b=g.code);return n.$$set=g=>{"collection"in g&&l(0,p=g.collection)},n.$$.update=()=>{n.$$.dirty&48&&l(2,i=[{code:200,body:m?"...":JSON.stringify(f,null,2)},{code:404,body:`
                {
                  "code": 404,
                  "message": "Missing collection context.",
                  "data": {}
                }
            `}])},l(3,o=He.getApiExampleUrl(oe.baseURL)),[p,b,i,o,f,m,y]}class ze extends Be{constructor(s){super(),Ce(this,s,Ge,Fe,Te,{collection:0})}}export{ze as default};
