import{S as Pe,i as $e,s as qe,e as c,w,b as v,c as ve,f as b,g as r,h as n,m as we,x as I,O as me,P as Re,k as ge,Q as ye,n as Be,t as Z,a as x,o as d,d as he,R as Ce,C as Se,p as Te,r as L,u as Me,N as Ae}from"./index-8e98d27b.js";import{S as Ue}from"./SdkTabs-dc7da6d6.js";function ue(a,s,l){const o=a.slice();return o[5]=s[l],o}function be(a,s,l){const o=a.slice();return o[5]=s[l],o}function _e(a,s){let l,o=s[5].code+"",_,u,i,p;function m(){return s[4](s[5])}return{key:a,first:null,c(){l=c("button"),_=w(o),u=v(),b(l,"class","tab-item"),L(l,"active",s[1]===s[5].code),this.first=l},m(P,$){r(P,l,$),n(l,_),n(l,u),i||(p=Me(l,"click",m),i=!0)},p(P,$){s=P,$&4&&o!==(o=s[5].code+"")&&I(_,o),$&6&&L(l,"active",s[1]===s[5].code)},d(P){P&&d(l),i=!1,p()}}}function ke(a,s){let l,o,_,u;return o=new Ae({props:{content:s[5].body}}),{key:a,first:null,c(){l=c("div"),ve(o.$$.fragment),_=v(),b(l,"class","tab-item"),L(l,"active",s[1]===s[5].code),this.first=l},m(i,p){r(i,l,p),we(o,l,null),n(l,_),u=!0},p(i,p){s=i;const m={};p&4&&(m.content=s[5].body),o.$set(m),(!u||p&6)&&L(l,"active",s[1]===s[5].code)},i(i){u||(Z(o.$$.fragment,i),u=!0)},o(i){x(o.$$.fragment,i),u=!1},d(i){i&&d(l),he(o)}}}function je(a){var re,de;let s,l,o=a[0].name+"",_,u,i,p,m,P,$,D=a[0].name+"",N,ee,Q,q,z,B,G,R,H,te,O,C,se,J,E=a[0].name+"",K,le,V,S,W,T,X,M,Y,g,A,h=[],oe=new Map,ae,U,k=[],ne=new Map,y;q=new Ue({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(re=a[0])==null?void 0:re.name}').requestPasswordReset('test@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(de=a[0])==null?void 0:de.name}').requestPasswordReset('test@example.com');
    `}});let F=a[2];const ie=e=>e[5].code;for(let e=0;e<F.length;e+=1){let t=be(a,F,e),f=ie(t);oe.set(f,h[e]=_e(f,t))}let j=a[2];const ce=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=ue(a,j,e),f=ce(t);ne.set(f,k[e]=ke(f,t))}return{c(){s=c("h3"),l=w("Request password reset ("),_=w(o),u=w(")"),i=v(),p=c("div"),m=c("p"),P=w("Sends "),$=c("strong"),N=w(D),ee=w(" password reset email request."),Q=v(),ve(q.$$.fragment),z=v(),B=c("h6"),B.textContent="API details",G=v(),R=c("div"),H=c("strong"),H.textContent="POST",te=v(),O=c("div"),C=c("p"),se=w("/api/collections/"),J=c("strong"),K=w(E),le=w("/request-password-reset"),V=v(),S=c("div"),S.textContent="Body Parameters",W=v(),T=c("table"),T.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>email</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The auth record email address to send the password reset request (if exists).</td></tr></tbody>`,X=v(),M=c("div"),M.textContent="Responses",Y=v(),g=c("div"),A=c("div");for(let e=0;e<h.length;e+=1)h[e].c();ae=v(),U=c("div");for(let e=0;e<k.length;e+=1)k[e].c();b(s,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(H,"class","label label-primary"),b(O,"class","content"),b(R,"class","alert alert-success"),b(S,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(A,"class","tabs-header compact left"),b(U,"class","tabs-content"),b(g,"class","tabs")},m(e,t){r(e,s,t),n(s,l),n(s,_),n(s,u),r(e,i,t),r(e,p,t),n(p,m),n(m,P),n(m,$),n($,N),n(m,ee),r(e,Q,t),we(q,e,t),r(e,z,t),r(e,B,t),r(e,G,t),r(e,R,t),n(R,H),n(R,te),n(R,O),n(O,C),n(C,se),n(C,J),n(J,K),n(C,le),r(e,V,t),r(e,S,t),r(e,W,t),r(e,T,t),r(e,X,t),r(e,M,t),r(e,Y,t),r(e,g,t),n(g,A);for(let f=0;f<h.length;f+=1)h[f].m(A,null);n(g,ae),n(g,U);for(let f=0;f<k.length;f+=1)k[f].m(U,null);y=!0},p(e,[t]){var fe,pe;(!y||t&1)&&o!==(o=e[0].name+"")&&I(_,o),(!y||t&1)&&D!==(D=e[0].name+"")&&I(N,D);const f={};t&9&&(f.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(fe=e[0])==null?void 0:fe.name}').requestPasswordReset('test@example.com');
    `),t&9&&(f.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(pe=e[0])==null?void 0:pe.name}').requestPasswordReset('test@example.com');
    `),q.$set(f),(!y||t&1)&&E!==(E=e[0].name+"")&&I(K,E),t&6&&(F=e[2],h=me(h,t,ie,1,e,F,oe,A,Re,_e,null,be)),t&6&&(j=e[2],ge(),k=me(k,t,ce,1,e,j,ne,U,ye,ke,null,ue),Be())},i(e){if(!y){Z(q.$$.fragment,e);for(let t=0;t<j.length;t+=1)Z(k[t]);y=!0}},o(e){x(q.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);y=!1},d(e){e&&d(s),e&&d(i),e&&d(p),e&&d(Q),he(q,e),e&&d(z),e&&d(B),e&&d(G),e&&d(R),e&&d(V),e&&d(S),e&&d(W),e&&d(T),e&&d(X),e&&d(M),e&&d(Y),e&&d(g);for(let t=0;t<h.length;t+=1)h[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function De(a,s,l){let o,{collection:_=new Ce}=s,u=204,i=[];const p=m=>l(1,u=m.code);return a.$$set=m=>{"collection"in m&&l(0,_=m.collection)},l(3,o=Se.getApiExampleUrl(Te.baseUrl)),l(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,u,i,o,p]}class Ee extends Pe{constructor(s){super(),$e(this,s,De,je,qe,{collection:0})}}export{Ee as default};
