import{S as qe,i as we,s as Pe,e as r,w as $,b as v,c as ve,f as b,g as f,h as i,m as he,x as I,aa as me,ab as ge,k as ye,ac as Be,n as Ce,t as Z,a as x,o as u,d as $e,ae as Se,C as Te,p as Ve,r as L,u as Me,a9 as Re}from"./index-197def9c.js";import{S as Ae}from"./SdkTabs-e182a429.js";function pe(o,l,s){const a=o.slice();return a[5]=l[s],a}function be(o,l,s){const a=o.slice();return a[5]=l[s],a}function _e(o,l){let s,a=l[5].code+"",_,p,n,d;function m(){return l[4](l[5])}return{key:o,first:null,c(){s=r("button"),_=$(a),p=v(),b(s,"class","tab-item"),L(s,"active",l[1]===l[5].code),this.first=s},m(q,w){f(q,s,w),i(s,_),i(s,p),n||(d=Me(s,"click",m),n=!0)},p(q,w){l=q,w&4&&a!==(a=l[5].code+"")&&I(_,a),w&6&&L(s,"active",l[1]===l[5].code)},d(q){q&&u(s),n=!1,d()}}}function ke(o,l){let s,a,_,p;return a=new Re({props:{content:l[5].body}}),{key:o,first:null,c(){s=r("div"),ve(a.$$.fragment),_=v(),b(s,"class","tab-item"),L(s,"active",l[1]===l[5].code),this.first=s},m(n,d){f(n,s,d),he(a,s,null),i(s,_),p=!0},p(n,d){l=n;const m={};d&4&&(m.content=l[5].body),a.$set(m),(!p||d&6)&&L(s,"active",l[1]===l[5].code)},i(n){p||(Z(a.$$.fragment,n),p=!0)},o(n){x(a.$$.fragment,n),p=!1},d(n){n&&u(s),$e(a)}}}function Ue(o){var re,fe;let l,s,a=o[0].name+"",_,p,n,d,m,q,w,j=o[0].name+"",O,ee,z,P,G,C,J,g,D,te,H,S,le,K,E=o[0].name+"",N,se,Q,T,W,V,X,M,Y,y,R,h=[],ae=new Map,oe,A,k=[],ie=new Map,B;P=new Ae({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(re=o[0])==null?void 0:re.name}').requestVerification('test@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(fe=o[0])==null?void 0:fe.name}').requestVerification('test@example.com');
    `}});let F=o[2];const ne=e=>e[5].code;for(let e=0;e<F.length;e+=1){let t=be(o,F,e),c=ne(t);ae.set(c,h[e]=_e(c,t))}let U=o[2];const ce=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=pe(o,U,e),c=ce(t);ie.set(c,k[e]=ke(c,t))}return{c(){l=r("h3"),s=$("Request verification ("),_=$(a),p=$(")"),n=v(),d=r("div"),m=r("p"),q=$("Sends "),w=r("strong"),O=$(j),ee=$(" verification email request."),z=v(),ve(P.$$.fragment),G=v(),C=r("h6"),C.textContent="API details",J=v(),g=r("div"),D=r("strong"),D.textContent="POST",te=v(),H=r("div"),S=r("p"),le=$("/api/collections/"),K=r("strong"),N=$(E),se=$("/request-verification"),Q=v(),T=r("div"),T.textContent="Body Parameters",W=v(),V=r("table"),V.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>email</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The auth record email address to send the verification request (if exists).</td></tr></tbody>`,X=v(),M=r("div"),M.textContent="Responses",Y=v(),y=r("div"),R=r("div");for(let e=0;e<h.length;e+=1)h[e].c();oe=v(),A=r("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(d,"class","content txt-lg m-b-sm"),b(C,"class","m-b-xs"),b(D,"class","label label-primary"),b(H,"class","content"),b(g,"class","alert alert-success"),b(T,"class","section-title"),b(V,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(R,"class","tabs-header compact left"),b(A,"class","tabs-content"),b(y,"class","tabs")},m(e,t){f(e,l,t),i(l,s),i(l,_),i(l,p),f(e,n,t),f(e,d,t),i(d,m),i(m,q),i(m,w),i(w,O),i(m,ee),f(e,z,t),he(P,e,t),f(e,G,t),f(e,C,t),f(e,J,t),f(e,g,t),i(g,D),i(g,te),i(g,H),i(H,S),i(S,le),i(S,K),i(K,N),i(S,se),f(e,Q,t),f(e,T,t),f(e,W,t),f(e,V,t),f(e,X,t),f(e,M,t),f(e,Y,t),f(e,y,t),i(y,R);for(let c=0;c<h.length;c+=1)h[c]&&h[c].m(R,null);i(y,oe),i(y,A);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(A,null);B=!0},p(e,[t]){var ue,de;(!B||t&1)&&a!==(a=e[0].name+"")&&I(_,a),(!B||t&1)&&j!==(j=e[0].name+"")&&I(O,j);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ue=e[0])==null?void 0:ue.name}').requestVerification('test@example.com');
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(de=e[0])==null?void 0:de.name}').requestVerification('test@example.com');
    `),P.$set(c),(!B||t&1)&&E!==(E=e[0].name+"")&&I(N,E),t&6&&(F=e[2],h=me(h,t,ne,1,e,F,ae,R,ge,_e,null,be)),t&6&&(U=e[2],ye(),k=me(k,t,ce,1,e,U,ie,A,Be,ke,null,pe),Ce())},i(e){if(!B){Z(P.$$.fragment,e);for(let t=0;t<U.length;t+=1)Z(k[t]);B=!0}},o(e){x(P.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);B=!1},d(e){e&&u(l),e&&u(n),e&&u(d),e&&u(z),$e(P,e),e&&u(G),e&&u(C),e&&u(J),e&&u(g),e&&u(Q),e&&u(T),e&&u(W),e&&u(V),e&&u(X),e&&u(M),e&&u(Y),e&&u(y);for(let t=0;t<h.length;t+=1)h[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function je(o,l,s){let a,{collection:_=new Se}=l,p=204,n=[];const d=m=>s(1,p=m.code);return o.$$set=m=>{"collection"in m&&s(0,_=m.collection)},s(3,a=Te.getApiExampleUrl(Ve.baseUrl)),s(2,n=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,p,n,a,d]}class Ee extends qe{constructor(l){super(),we(this,l,je,Ue,Pe,{collection:0})}}export{Ee as default};
