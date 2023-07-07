import{S as $e,i as Pe,s as qe,e as r,w as h,b as v,c as ve,f as b,g as d,h as n,m as we,x as L,aa as ue,ab as ge,k as ye,ac as Re,n as Be,t as Z,a as x,o as f,d as he,ae as Ce,C as Se,p as Te,r as O,u as Me,a9 as Ae}from"./index-c45c880c.js";import{S as Ue}from"./SdkTabs-04dd5574.js";function me(o,s,l){const a=o.slice();return a[5]=s[l],a}function be(o,s,l){const a=o.slice();return a[5]=s[l],a}function _e(o,s){let l,a=s[5].code+"",_,m,i,p;function u(){return s[4](s[5])}return{key:o,first:null,c(){l=r("button"),_=h(a),m=v(),b(l,"class","tab-item"),O(l,"active",s[1]===s[5].code),this.first=l},m($,P){d($,l,P),n(l,_),n(l,m),i||(p=Me(l,"click",u),i=!0)},p($,P){s=$,P&4&&a!==(a=s[5].code+"")&&L(_,a),P&6&&O(l,"active",s[1]===s[5].code)},d($){$&&f(l),i=!1,p()}}}function ke(o,s){let l,a,_,m;return a=new Ae({props:{content:s[5].body}}),{key:o,first:null,c(){l=r("div"),ve(a.$$.fragment),_=v(),b(l,"class","tab-item"),O(l,"active",s[1]===s[5].code),this.first=l},m(i,p){d(i,l,p),we(a,l,null),n(l,_),m=!0},p(i,p){s=i;const u={};p&4&&(u.content=s[5].body),a.$set(u),(!m||p&6)&&O(l,"active",s[1]===s[5].code)},i(i){m||(Z(a.$$.fragment,i),m=!0)},o(i){x(a.$$.fragment,i),m=!1},d(i){i&&f(l),he(a)}}}function je(o){var re,de;let s,l,a=o[0].name+"",_,m,i,p,u,$,P,D=o[0].name+"",z,ee,G,q,J,B,K,g,H,te,E,C,se,N,F=o[0].name+"",Q,le,V,S,W,T,X,M,Y,y,A,w=[],ae=new Map,oe,U,k=[],ne=new Map,R;q=new Ue({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(re=o[0])==null?void 0:re.name}').requestPasswordReset('test@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(de=o[0])==null?void 0:de.name}').requestPasswordReset('test@example.com');
    `}});let I=o[2];const ie=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=be(o,I,e),c=ie(t);ae.set(c,w[e]=_e(c,t))}let j=o[2];const ce=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=me(o,j,e),c=ce(t);ne.set(c,k[e]=ke(c,t))}return{c(){s=r("h3"),l=h("Request password reset ("),_=h(a),m=h(")"),i=v(),p=r("div"),u=r("p"),$=h("Sends "),P=r("strong"),z=h(D),ee=h(" password reset email request."),G=v(),ve(q.$$.fragment),J=v(),B=r("h6"),B.textContent="API details",K=v(),g=r("div"),H=r("strong"),H.textContent="POST",te=v(),E=r("div"),C=r("p"),se=h("/api/collections/"),N=r("strong"),Q=h(F),le=h("/request-password-reset"),V=v(),S=r("div"),S.textContent="Body Parameters",W=v(),T=r("table"),T.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>email</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The auth record email address to send the password reset request (if exists).</td></tr></tbody>`,X=v(),M=r("div"),M.textContent="Responses",Y=v(),y=r("div"),A=r("div");for(let e=0;e<w.length;e+=1)w[e].c();oe=v(),U=r("div");for(let e=0;e<k.length;e+=1)k[e].c();b(s,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(H,"class","label label-primary"),b(E,"class","content"),b(g,"class","alert alert-success"),b(S,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(A,"class","tabs-header compact left"),b(U,"class","tabs-content"),b(y,"class","tabs")},m(e,t){d(e,s,t),n(s,l),n(s,_),n(s,m),d(e,i,t),d(e,p,t),n(p,u),n(u,$),n(u,P),n(P,z),n(u,ee),d(e,G,t),we(q,e,t),d(e,J,t),d(e,B,t),d(e,K,t),d(e,g,t),n(g,H),n(g,te),n(g,E),n(E,C),n(C,se),n(C,N),n(N,Q),n(C,le),d(e,V,t),d(e,S,t),d(e,W,t),d(e,T,t),d(e,X,t),d(e,M,t),d(e,Y,t),d(e,y,t),n(y,A);for(let c=0;c<w.length;c+=1)w[c]&&w[c].m(A,null);n(y,oe),n(y,U);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(U,null);R=!0},p(e,[t]){var fe,pe;(!R||t&1)&&a!==(a=e[0].name+"")&&L(_,a),(!R||t&1)&&D!==(D=e[0].name+"")&&L(z,D);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(fe=e[0])==null?void 0:fe.name}').requestPasswordReset('test@example.com');
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(pe=e[0])==null?void 0:pe.name}').requestPasswordReset('test@example.com');
    `),q.$set(c),(!R||t&1)&&F!==(F=e[0].name+"")&&L(Q,F),t&6&&(I=e[2],w=ue(w,t,ie,1,e,I,ae,A,ge,_e,null,be)),t&6&&(j=e[2],ye(),k=ue(k,t,ce,1,e,j,ne,U,Re,ke,null,me),Be())},i(e){if(!R){Z(q.$$.fragment,e);for(let t=0;t<j.length;t+=1)Z(k[t]);R=!0}},o(e){x(q.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);R=!1},d(e){e&&f(s),e&&f(i),e&&f(p),e&&f(G),he(q,e),e&&f(J),e&&f(B),e&&f(K),e&&f(g),e&&f(V),e&&f(S),e&&f(W),e&&f(T),e&&f(X),e&&f(M),e&&f(Y),e&&f(y);for(let t=0;t<w.length;t+=1)w[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function De(o,s,l){let a,{collection:_=new Ce}=s,m=204,i=[];const p=u=>l(1,m=u.code);return o.$$set=u=>{"collection"in u&&l(0,_=u.collection)},l(3,a=Se.getApiExampleUrl(Te.baseUrl)),l(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,m,i,a,p]}class Fe extends $e{constructor(s){super(),Pe(this,s,De,je,qe,{collection:0})}}export{Fe as default};
