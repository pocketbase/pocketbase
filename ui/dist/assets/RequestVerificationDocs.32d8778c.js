import{S as qe,i as we,s as Pe,e as c,w as h,b as v,c as ve,f as b,g as r,h as i,m as he,x as E,O as me,P as ge,k as ye,Q as Be,n as Ce,t as Z,a as x,o as f,d as $e,R as Se,C as Te,p as Re,r as F,u as Ve,N as Me}from"./index.89a3f554.js";import{S as Ae}from"./SdkTabs.0a6ad1c9.js";function pe(a,l,s){const o=a.slice();return o[5]=l[s],o}function be(a,l,s){const o=a.slice();return o[5]=l[s],o}function _e(a,l){let s,o=l[5].code+"",_,p,n,d;function m(){return l[4](l[5])}return{key:a,first:null,c(){s=c("button"),_=h(o),p=v(),b(s,"class","tab-item"),F(s,"active",l[1]===l[5].code),this.first=s},m(q,w){r(q,s,w),i(s,_),i(s,p),n||(d=Ve(s,"click",m),n=!0)},p(q,w){l=q,w&4&&o!==(o=l[5].code+"")&&E(_,o),w&6&&F(s,"active",l[1]===l[5].code)},d(q){q&&f(s),n=!1,d()}}}function ke(a,l){let s,o,_,p;return o=new Me({props:{content:l[5].body}}),{key:a,first:null,c(){s=c("div"),ve(o.$$.fragment),_=v(),b(s,"class","tab-item"),F(s,"active",l[1]===l[5].code),this.first=s},m(n,d){r(n,s,d),he(o,s,null),i(s,_),p=!0},p(n,d){l=n;const m={};d&4&&(m.content=l[5].body),o.$set(m),(!p||d&6)&&F(s,"active",l[1]===l[5].code)},i(n){p||(Z(o.$$.fragment,n),p=!0)},o(n){x(o.$$.fragment,n),p=!1},d(n){n&&f(s),$e(o)}}}function Ue(a){var re,fe;let l,s,o=a[0].name+"",_,p,n,d,m,q,w,j=a[0].name+"",L,ee,N,P,Q,C,z,g,D,te,H,S,le,G,I=a[0].name+"",J,se,K,T,W,R,X,V,Y,y,M,$=[],oe=new Map,ae,A,k=[],ie=new Map,B;P=new Ae({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(re=a[0])==null?void 0:re.name}').requestVerification('test@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(fe=a[0])==null?void 0:fe.name}').requestVerification('test@example.com');
    `}});let O=a[2];const ne=e=>e[5].code;for(let e=0;e<O.length;e+=1){let t=be(a,O,e),u=ne(t);oe.set(u,$[e]=_e(u,t))}let U=a[2];const ce=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=pe(a,U,e),u=ce(t);ie.set(u,k[e]=ke(u,t))}return{c(){l=c("h3"),s=h("Request verification ("),_=h(o),p=h(")"),n=v(),d=c("div"),m=c("p"),q=h("Sends "),w=c("strong"),L=h(j),ee=h(" verification email request."),N=v(),ve(P.$$.fragment),Q=v(),C=c("h6"),C.textContent="API details",z=v(),g=c("div"),D=c("strong"),D.textContent="POST",te=v(),H=c("div"),S=c("p"),le=h("/api/collections/"),G=c("strong"),J=h(I),se=h("/request-verification"),K=v(),T=c("div"),T.textContent="Body Parameters",W=v(),R=c("table"),R.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>email</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The auth record email address to send the verification request (if exists).</td></tr></tbody>`,X=v(),V=c("div"),V.textContent="Responses",Y=v(),y=c("div"),M=c("div");for(let e=0;e<$.length;e+=1)$[e].c();ae=v(),A=c("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(d,"class","content txt-lg m-b-sm"),b(C,"class","m-b-xs"),b(D,"class","label label-primary"),b(H,"class","content"),b(g,"class","alert alert-success"),b(T,"class","section-title"),b(R,"class","table-compact table-border m-b-base"),b(V,"class","section-title"),b(M,"class","tabs-header compact left"),b(A,"class","tabs-content"),b(y,"class","tabs")},m(e,t){r(e,l,t),i(l,s),i(l,_),i(l,p),r(e,n,t),r(e,d,t),i(d,m),i(m,q),i(m,w),i(w,L),i(m,ee),r(e,N,t),he(P,e,t),r(e,Q,t),r(e,C,t),r(e,z,t),r(e,g,t),i(g,D),i(g,te),i(g,H),i(H,S),i(S,le),i(S,G),i(G,J),i(S,se),r(e,K,t),r(e,T,t),r(e,W,t),r(e,R,t),r(e,X,t),r(e,V,t),r(e,Y,t),r(e,y,t),i(y,M);for(let u=0;u<$.length;u+=1)$[u].m(M,null);i(y,ae),i(y,A);for(let u=0;u<k.length;u+=1)k[u].m(A,null);B=!0},p(e,[t]){var ue,de;(!B||t&1)&&o!==(o=e[0].name+"")&&E(_,o),(!B||t&1)&&j!==(j=e[0].name+"")&&E(L,j);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ue=e[0])==null?void 0:ue.name}').requestVerification('test@example.com');
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(de=e[0])==null?void 0:de.name}').requestVerification('test@example.com');
    `),P.$set(u),(!B||t&1)&&I!==(I=e[0].name+"")&&E(J,I),t&6&&(O=e[2],$=me($,t,ne,1,e,O,oe,M,ge,_e,null,be)),t&6&&(U=e[2],ye(),k=me(k,t,ce,1,e,U,ie,A,Be,ke,null,pe),Ce())},i(e){if(!B){Z(P.$$.fragment,e);for(let t=0;t<U.length;t+=1)Z(k[t]);B=!0}},o(e){x(P.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);B=!1},d(e){e&&f(l),e&&f(n),e&&f(d),e&&f(N),$e(P,e),e&&f(Q),e&&f(C),e&&f(z),e&&f(g),e&&f(K),e&&f(T),e&&f(W),e&&f(R),e&&f(X),e&&f(V),e&&f(Y),e&&f(y);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function je(a,l,s){let o,{collection:_=new Se}=l,p=204,n=[];const d=m=>s(1,p=m.code);return a.$$set=m=>{"collection"in m&&s(0,_=m.collection)},s(3,o=Te.getApiExampleUrl(Re.baseUrl)),s(2,n=[{code:204,body:"null"},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "email": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}]),[_,p,n,o,d]}class Ie extends qe{constructor(l){super(),we(this,l,je,Ue,Pe,{collection:0})}}export{Ie as default};
