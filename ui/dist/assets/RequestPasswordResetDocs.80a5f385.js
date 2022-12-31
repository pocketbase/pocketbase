import{S as Pe,i as $e,s as qe,e as c,w,b as v,c as ve,f as b,g as r,h as n,m as we,x as F,O as ue,P as Re,k as ge,Q as ye,n as Be,t as Z,a as x,o as d,d as he,R as Ce,C as Se,p as Te,r as L,u as Me,N as Ae}from"./index.89a3f554.js";import{S as Ue}from"./SdkTabs.0a6ad1c9.js";function me(a,s,l){const o=a.slice();return o[5]=s[l],o}function be(a,s,l){const o=a.slice();return o[5]=s[l],o}function _e(a,s){let l,o=s[5].code+"",_,m,i,p;function u(){return s[4](s[5])}return{key:a,first:null,c(){l=c("button"),_=w(o),m=v(),b(l,"class","tab-item"),L(l,"active",s[1]===s[5].code),this.first=l},m(P,$){r(P,l,$),n(l,_),n(l,m),i||(p=Me(l,"click",u),i=!0)},p(P,$){s=P,$&4&&o!==(o=s[5].code+"")&&F(_,o),$&6&&L(l,"active",s[1]===s[5].code)},d(P){P&&d(l),i=!1,p()}}}function ke(a,s){let l,o,_,m;return o=new Ae({props:{content:s[5].body}}),{key:a,first:null,c(){l=c("div"),ve(o.$$.fragment),_=v(),b(l,"class","tab-item"),L(l,"active",s[1]===s[5].code),this.first=l},m(i,p){r(i,l,p),we(o,l,null),n(l,_),m=!0},p(i,p){s=i;const u={};p&4&&(u.content=s[5].body),o.$set(u),(!m||p&6)&&L(l,"active",s[1]===s[5].code)},i(i){m||(Z(o.$$.fragment,i),m=!0)},o(i){x(o.$$.fragment,i),m=!1},d(i){i&&d(l),he(o)}}}function je(a){var re,de;let s,l,o=a[0].name+"",_,m,i,p,u,P,$,D=a[0].name+"",N,ee,Q,q,z,B,G,R,H,te,I,C,se,J,O=a[0].name+"",K,le,V,S,W,T,X,M,Y,g,A,h=[],oe=new Map,ae,U,k=[],ne=new Map,y;q=new Ue({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(re=a[0])==null?void 0:re.name}').requestPasswordReset('test@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(de=a[0])==null?void 0:de.name}').requestPasswordReset('test@example.com');
    `}});let E=a[2];const ie=e=>e[5].code;for(let e=0;e<E.length;e+=1){let t=be(a,E,e),f=ie(t);oe.set(f,h[e]=_e(f,t))}let j=a[2];const ce=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=me(a,j,e),f=ce(t);ne.set(f,k[e]=ke(f,t))}return{c(){s=c("h3"),l=w("Request password reset ("),_=w(o),m=w(")"),i=v(),p=c("div"),u=c("p"),P=w("Sends "),$=c("strong"),N=w(D),ee=w(" password reset email request."),Q=v(),ve(q.$$.fragment),z=v(),B=c("h6"),B.textContent="API details",G=v(),R=c("div"),H=c("strong"),H.textContent="POST",te=v(),I=c("div"),C=c("p"),se=w("/api/collections/"),J=c("strong"),K=w(O),le=w("/request-password-reset"),V=v(),S=c("div"),S.textContent="Body Parameters",W=v(),T=c("table"),T.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>email</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The auth record email address to send the password reset request (if exists).</td></tr></tbody>`,X=v(),M=c("div"),M.textContent="Responses",Y=v(),g=c("div"),A=c("div");for(let e=0;e<h.length;e+=1)h[e].c();ae=v(),U=c("div");for(let e=0;e<k.length;e+=1)k[e].c();b(s,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(H,"class","label label-primary"),b(I,"class","content"),b(R,"class","alert alert-success"),b(S,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(A,"class","tabs-header compact left"),b(U,"class","tabs-content"),b(g,"class","tabs")},m(e,t){r(e,s,t),n(s,l),n(s,_),n(s,m),r(e,i,t),r(e,p,t),n(p,u),n(u,P),n(u,$),n($,N),n(u,ee),r(e,Q,t),we(q,e,t),r(e,z,t),r(e,B,t),r(e,G,t),r(e,R,t),n(R,H),n(R,te),n(R,I),n(I,C),n(C,se),n(C,J),n(J,K),n(C,le),r(e,V,t),r(e,S,t),r(e,W,t),r(e,T,t),r(e,X,t),r(e,M,t),r(e,Y,t),r(e,g,t),n(g,A);for(let f=0;f<h.length;f+=1)h[f].m(A,null);n(g,ae),n(g,U);for(let f=0;f<k.length;f+=1)k[f].m(U,null);y=!0},p(e,[t]){var fe,pe;(!y||t&1)&&o!==(o=e[0].name+"")&&F(_,o),(!y||t&1)&&D!==(D=e[0].name+"")&&F(N,D);const f={};t&9&&(f.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(fe=e[0])==null?void 0:fe.name}').requestPasswordReset('test@example.com');
    `),t&9&&(f.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(pe=e[0])==null?void 0:pe.name}').requestPasswordReset('test@example.com');
    `),q.$set(f),(!y||t&1)&&O!==(O=e[0].name+"")&&F(K,O),t&6&&(E=e[2],h=ue(h,t,ie,1,e,E,oe,A,Re,_e,null,be)),t&6&&(j=e[2],ge(),k=ue(k,t,ce,1,e,j,ne,U,ye,ke,null,me),Be())},i(e){if(!y){Z(q.$$.fragment,e);for(let t=0;t<j.length;t+=1)Z(k[t]);y=!0}},o(e){x(q.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);y=!1},d(e){e&&d(s),e&&d(i),e&&d(p),e&&d(Q),he(q,e),e&&d(z),e&&d(B),e&&d(G),e&&d(R),e&&d(V),e&&d(S),e&&d(W),e&&d(T),e&&d(X),e&&d(M),e&&d(Y),e&&d(g);for(let t=0;t<h.length;t+=1)h[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function De(a,s,l){let o,{collection:_=new Ce}=s,m=204,i=[];const p=u=>l(1,m=u.code);return a.$$set=u=>{"collection"in u&&l(0,_=u.collection)},l(3,o=Se.getApiExampleUrl(Te.baseUrl)),l(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,m,i,o,p]}class Oe extends Pe{constructor(s){super(),$e(this,s,De,je,qe,{collection:0})}}export{Oe as default};
