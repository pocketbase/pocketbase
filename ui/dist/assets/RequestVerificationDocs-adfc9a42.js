import{S as qe,i as we,s as Pe,e as r,w as $,b as v,c as ve,f as b,g as f,h as i,m as he,x as I,N as me,P as ge,k as ye,Q as Be,n as Ce,t as Z,a as x,o as u,d as $e,T as Se,C as Te,p as Me,r as L,u as Ve,M as Re}from"./index-a65ca895.js";import{S as Ae}from"./SdkTabs-ad912c8f.js";function pe(a,l,s){const o=a.slice();return o[5]=l[s],o}function be(a,l,s){const o=a.slice();return o[5]=l[s],o}function _e(a,l){let s,o=l[5].code+"",_,p,n,d;function m(){return l[4](l[5])}return{key:a,first:null,c(){s=r("button"),_=$(o),p=v(),b(s,"class","tab-item"),L(s,"active",l[1]===l[5].code),this.first=s},m(q,w){f(q,s,w),i(s,_),i(s,p),n||(d=Ve(s,"click",m),n=!0)},p(q,w){l=q,w&4&&o!==(o=l[5].code+"")&&I(_,o),w&6&&L(s,"active",l[1]===l[5].code)},d(q){q&&u(s),n=!1,d()}}}function ke(a,l){let s,o,_,p;return o=new Re({props:{content:l[5].body}}),{key:a,first:null,c(){s=r("div"),ve(o.$$.fragment),_=v(),b(s,"class","tab-item"),L(s,"active",l[1]===l[5].code),this.first=s},m(n,d){f(n,s,d),he(o,s,null),i(s,_),p=!0},p(n,d){l=n;const m={};d&4&&(m.content=l[5].body),o.$set(m),(!p||d&6)&&L(s,"active",l[1]===l[5].code)},i(n){p||(Z(o.$$.fragment,n),p=!0)},o(n){x(o.$$.fragment,n),p=!1},d(n){n&&u(s),$e(o)}}}function Ue(a){var re,fe;let l,s,o=a[0].name+"",_,p,n,d,m,q,w,j=a[0].name+"",N,ee,O,P,Q,C,z,g,D,te,H,S,le,G,E=a[0].name+"",J,se,K,T,W,M,X,V,Y,y,R,h=[],oe=new Map,ae,A,k=[],ie=new Map,B;P=new Ae({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(re=a[0])==null?void 0:re.name}').requestVerification('test@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(fe=a[0])==null?void 0:fe.name}').requestVerification('test@example.com');
    `}});let F=a[2];const ne=e=>e[5].code;for(let e=0;e<F.length;e+=1){let t=be(a,F,e),c=ne(t);oe.set(c,h[e]=_e(c,t))}let U=a[2];const ce=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=pe(a,U,e),c=ce(t);ie.set(c,k[e]=ke(c,t))}return{c(){l=r("h3"),s=$("Request verification ("),_=$(o),p=$(")"),n=v(),d=r("div"),m=r("p"),q=$("Sends "),w=r("strong"),N=$(j),ee=$(" verification email request."),O=v(),ve(P.$$.fragment),Q=v(),C=r("h6"),C.textContent="API details",z=v(),g=r("div"),D=r("strong"),D.textContent="POST",te=v(),H=r("div"),S=r("p"),le=$("/api/collections/"),G=r("strong"),J=$(E),se=$("/request-verification"),K=v(),T=r("div"),T.textContent="Body Parameters",W=v(),M=r("table"),M.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>email</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The auth record email address to send the verification request (if exists).</td></tr></tbody>`,X=v(),V=r("div"),V.textContent="Responses",Y=v(),y=r("div"),R=r("div");for(let e=0;e<h.length;e+=1)h[e].c();ae=v(),A=r("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(d,"class","content txt-lg m-b-sm"),b(C,"class","m-b-xs"),b(D,"class","label label-primary"),b(H,"class","content"),b(g,"class","alert alert-success"),b(T,"class","section-title"),b(M,"class","table-compact table-border m-b-base"),b(V,"class","section-title"),b(R,"class","tabs-header compact left"),b(A,"class","tabs-content"),b(y,"class","tabs")},m(e,t){f(e,l,t),i(l,s),i(l,_),i(l,p),f(e,n,t),f(e,d,t),i(d,m),i(m,q),i(m,w),i(w,N),i(m,ee),f(e,O,t),he(P,e,t),f(e,Q,t),f(e,C,t),f(e,z,t),f(e,g,t),i(g,D),i(g,te),i(g,H),i(H,S),i(S,le),i(S,G),i(G,J),i(S,se),f(e,K,t),f(e,T,t),f(e,W,t),f(e,M,t),f(e,X,t),f(e,V,t),f(e,Y,t),f(e,y,t),i(y,R);for(let c=0;c<h.length;c+=1)h[c]&&h[c].m(R,null);i(y,ae),i(y,A);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(A,null);B=!0},p(e,[t]){var ue,de;(!B||t&1)&&o!==(o=e[0].name+"")&&I(_,o),(!B||t&1)&&j!==(j=e[0].name+"")&&I(N,j);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ue=e[0])==null?void 0:ue.name}').requestVerification('test@example.com');
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(de=e[0])==null?void 0:de.name}').requestVerification('test@example.com');
    `),P.$set(c),(!B||t&1)&&E!==(E=e[0].name+"")&&I(J,E),t&6&&(F=e[2],h=me(h,t,ne,1,e,F,oe,R,ge,_e,null,be)),t&6&&(U=e[2],ye(),k=me(k,t,ce,1,e,U,ie,A,Be,ke,null,pe),Ce())},i(e){if(!B){Z(P.$$.fragment,e);for(let t=0;t<U.length;t+=1)Z(k[t]);B=!0}},o(e){x(P.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);B=!1},d(e){e&&u(l),e&&u(n),e&&u(d),e&&u(O),$e(P,e),e&&u(Q),e&&u(C),e&&u(z),e&&u(g),e&&u(K),e&&u(T),e&&u(W),e&&u(M),e&&u(X),e&&u(V),e&&u(Y),e&&u(y);for(let t=0;t<h.length;t+=1)h[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function je(a,l,s){let o,{collection:_=new Se}=l,p=204,n=[];const d=m=>s(1,p=m.code);return a.$$set=m=>{"collection"in m&&s(0,_=m.collection)},s(3,o=Te.getApiExampleUrl(Me.baseUrl)),s(2,n=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,p,n,o,d]}class Ee extends qe{constructor(l){super(),we(this,l,je,Ue,Pe,{collection:0})}}export{Ee as default};
