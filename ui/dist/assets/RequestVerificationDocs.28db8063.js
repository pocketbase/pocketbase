import{S as we,i as qe,s as Pe,e as c,w as h,b as v,c as ve,f as b,g as r,h as i,m as he,x as O,P as ue,Q as ge,k as ye,R as Be,n as Ce,t as Z,a as x,o as f,d as $e,L as Se,C as Te,p as Re,r as E,u as Ve,O as Me}from"./index.be8ffbe5.js";import{S as Ae}from"./SdkTabs.8f55857f.js";function me(a,l,s){const o=a.slice();return o[5]=l[s],o}function be(a,l,s){const o=a.slice();return o[5]=l[s],o}function _e(a,l){let s,o=l[5].code+"",_,m,n,p;function u(){return l[4](l[5])}return{key:a,first:null,c(){s=c("button"),_=h(o),m=v(),b(s,"class","tab-item"),E(s,"active",l[1]===l[5].code),this.first=s},m(w,q){r(w,s,q),i(s,_),i(s,m),n||(p=Ve(s,"click",u),n=!0)},p(w,q){l=w,q&4&&o!==(o=l[5].code+"")&&O(_,o),q&6&&E(s,"active",l[1]===l[5].code)},d(w){w&&f(s),n=!1,p()}}}function ke(a,l){let s,o,_,m;return o=new Me({props:{content:l[5].body}}),{key:a,first:null,c(){s=c("div"),ve(o.$$.fragment),_=v(),b(s,"class","tab-item"),E(s,"active",l[1]===l[5].code),this.first=s},m(n,p){r(n,s,p),he(o,s,null),i(s,_),m=!0},p(n,p){l=n;const u={};p&4&&(u.content=l[5].body),o.$set(u),(!m||p&6)&&E(s,"active",l[1]===l[5].code)},i(n){m||(Z(o.$$.fragment,n),m=!0)},o(n){x(o.$$.fragment,n),m=!1},d(n){n&&f(s),$e(o)}}}function Ue(a){var re,fe;let l,s,o=a[0].name+"",_,m,n,p,u,w,q,j=a[0].name+"",F,ee,Q,P,z,C,G,g,D,te,H,S,le,J,I=a[0].name+"",K,se,N,T,W,R,X,V,Y,y,M,$=[],oe=new Map,ae,A,k=[],ie=new Map,B;P=new Ae({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(re=a[0])==null?void 0:re.name}').requestVerification('test@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(fe=a[0])==null?void 0:fe.name}').requestVerification('test@example.com');
    `}});let L=a[2];const ne=e=>e[5].code;for(let e=0;e<L.length;e+=1){let t=be(a,L,e),d=ne(t);oe.set(d,$[e]=_e(d,t))}let U=a[2];const ce=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=me(a,U,e),d=ce(t);ie.set(d,k[e]=ke(d,t))}return{c(){l=c("h3"),s=h("Request verification ("),_=h(o),m=h(")"),n=v(),p=c("div"),u=c("p"),w=h("Sends "),q=c("strong"),F=h(j),ee=h(" verification email request."),Q=v(),ve(P.$$.fragment),z=v(),C=c("h6"),C.textContent="API details",G=v(),g=c("div"),D=c("strong"),D.textContent="POST",te=v(),H=c("div"),S=c("p"),le=h("/api/collections/"),J=c("strong"),K=h(I),se=h("/request-password-reset"),N=v(),T=c("div"),T.textContent="Body Parameters",W=v(),R=c("table"),R.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>email</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The auth record email address to send the verification request (if exists).</td></tr></tbody>`,X=v(),V=c("div"),V.textContent="Responses",Y=v(),y=c("div"),M=c("div");for(let e=0;e<$.length;e+=1)$[e].c();ae=v(),A=c("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(C,"class","m-b-xs"),b(D,"class","label label-primary"),b(H,"class","content"),b(g,"class","alert alert-success"),b(T,"class","section-title"),b(R,"class","table-compact table-border m-b-base"),b(V,"class","section-title"),b(M,"class","tabs-header compact left"),b(A,"class","tabs-content"),b(y,"class","tabs")},m(e,t){r(e,l,t),i(l,s),i(l,_),i(l,m),r(e,n,t),r(e,p,t),i(p,u),i(u,w),i(u,q),i(q,F),i(u,ee),r(e,Q,t),he(P,e,t),r(e,z,t),r(e,C,t),r(e,G,t),r(e,g,t),i(g,D),i(g,te),i(g,H),i(H,S),i(S,le),i(S,J),i(J,K),i(S,se),r(e,N,t),r(e,T,t),r(e,W,t),r(e,R,t),r(e,X,t),r(e,V,t),r(e,Y,t),r(e,y,t),i(y,M);for(let d=0;d<$.length;d+=1)$[d].m(M,null);i(y,ae),i(y,A);for(let d=0;d<k.length;d+=1)k[d].m(A,null);B=!0},p(e,[t]){var de,pe;(!B||t&1)&&o!==(o=e[0].name+"")&&O(_,o),(!B||t&1)&&j!==(j=e[0].name+"")&&O(F,j);const d={};t&9&&(d.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(de=e[0])==null?void 0:de.name}').requestVerification('test@example.com');
    `),t&9&&(d.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(pe=e[0])==null?void 0:pe.name}').requestVerification('test@example.com');
    `),P.$set(d),(!B||t&1)&&I!==(I=e[0].name+"")&&O(K,I),t&6&&(L=e[2],$=ue($,t,ne,1,e,L,oe,M,ge,_e,null,be)),t&6&&(U=e[2],ye(),k=ue(k,t,ce,1,e,U,ie,A,Be,ke,null,me),Ce())},i(e){if(!B){Z(P.$$.fragment,e);for(let t=0;t<U.length;t+=1)Z(k[t]);B=!0}},o(e){x(P.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);B=!1},d(e){e&&f(l),e&&f(n),e&&f(p),e&&f(Q),$e(P,e),e&&f(z),e&&f(C),e&&f(G),e&&f(g),e&&f(N),e&&f(T),e&&f(W),e&&f(R),e&&f(X),e&&f(V),e&&f(Y),e&&f(y);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function je(a,l,s){let o,{collection:_=new Se}=l,m=204,n=[];const p=u=>s(1,m=u.code);return a.$$set=u=>{"collection"in u&&s(0,_=u.collection)},s(3,o=Te.getApiExampleUrl(Re.baseUrl)),s(2,n=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,m,n,o,p]}class Ie extends we{constructor(l){super(),qe(this,l,je,Ue,Pe,{collection:0})}}export{Ie as default};
