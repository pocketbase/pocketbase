import{S as qe,i as we,s as ye,O as F,b as r,d as ve,t as x,a as ee,r as I,Q as pe,R as Pe,g as Be,T as Ce,e as Se,f as d,h as n,m as ge,n as f,u as g,k as h,c as $e,o as b,C as Te,p as Re,w as L,x as Ve,N as Me}from"./index-lKVVd1Bs.js";import{S as Ae}from"./SdkTabs-CROZp_fs.js";function be(o,l,s){const a=o.slice();return a[5]=l[s],a}function _e(o,l,s){const a=o.slice();return a[5]=l[s],a}function ke(o,l){let s,a=l[5].code+"",_,p,i,u;function m(){return l[4](l[5])}return{key:o,first:null,c(){s=f("button"),_=g(a),p=h(),b(s,"class","tab-item"),L(s,"active",l[1]===l[5].code),this.first=s},m($,q){d($,s,q),n(s,_),n(s,p),i||(u=Ve(s,"click",m),i=!0)},p($,q){l=$,q&4&&a!==(a=l[5].code+"")&&I(_,a),q&6&&L(s,"active",l[1]===l[5].code)},d($){$&&r(s),i=!1,u()}}}function he(o,l){let s,a,_,p;return a=new Me({props:{content:l[5].body}}),{key:o,first:null,c(){s=f("div"),$e(a.$$.fragment),_=h(),b(s,"class","tab-item"),L(s,"active",l[1]===l[5].code),this.first=s},m(i,u){d(i,s,u),ge(a,s,null),n(s,_),p=!0},p(i,u){l=i;const m={};u&4&&(m.content=l[5].body),a.$set(m),(!p||u&6)&&L(s,"active",l[1]===l[5].code)},i(i){p||(ee(a.$$.fragment,i),p=!0)},o(i){x(a.$$.fragment,i),p=!1},d(i){i&&r(s),ve(a)}}}function Ue(o){var de,fe;let l,s,a=o[0].name+"",_,p,i,u,m,$,q,j=o[0].name+"",N,te,Q,w,z,C,G,y,D,le,H,S,se,J,O=o[0].name+"",K,ae,W,T,X,R,Y,V,Z,P,M,v=[],oe=new Map,ne,A,k=[],ie=new Map,B;w=new Ae({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(de=o[0])==null?void 0:de.name}').requestVerification('test@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(fe=o[0])==null?void 0:fe.name}').requestVerification('test@example.com');
    `}});let E=F(o[2]);const ce=e=>e[5].code;for(let e=0;e<E.length;e+=1){let t=_e(o,E,e),c=ce(t);oe.set(c,v[e]=ke(c,t))}let U=F(o[2]);const re=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=be(o,U,e),c=re(t);ie.set(c,k[e]=he(c,t))}return{c(){l=f("h3"),s=g("Request verification ("),_=g(a),p=g(")"),i=h(),u=f("div"),m=f("p"),$=g("Sends "),q=f("strong"),N=g(j),te=g(" verification email request."),Q=h(),$e(w.$$.fragment),z=h(),C=f("h6"),C.textContent="API details",G=h(),y=f("div"),D=f("strong"),D.textContent="POST",le=h(),H=f("div"),S=f("p"),se=g("/api/collections/"),J=f("strong"),K=g(O),ae=g("/request-verification"),W=h(),T=f("div"),T.textContent="Body Parameters",X=h(),R=f("table"),R.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address to send the verification request (if exists).</td></tr></tbody>',Y=h(),V=f("div"),V.textContent="Responses",Z=h(),P=f("div"),M=f("div");for(let e=0;e<v.length;e+=1)v[e].c();ne=h(),A=f("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(u,"class","content txt-lg m-b-sm"),b(C,"class","m-b-xs"),b(D,"class","label label-primary"),b(H,"class","content"),b(y,"class","alert alert-success"),b(T,"class","section-title"),b(R,"class","table-compact table-border m-b-base"),b(V,"class","section-title"),b(M,"class","tabs-header compact combined left"),b(A,"class","tabs-content"),b(P,"class","tabs")},m(e,t){d(e,l,t),n(l,s),n(l,_),n(l,p),d(e,i,t),d(e,u,t),n(u,m),n(m,$),n(m,q),n(q,N),n(m,te),d(e,Q,t),ge(w,e,t),d(e,z,t),d(e,C,t),d(e,G,t),d(e,y,t),n(y,D),n(y,le),n(y,H),n(H,S),n(S,se),n(S,J),n(J,K),n(S,ae),d(e,W,t),d(e,T,t),d(e,X,t),d(e,R,t),d(e,Y,t),d(e,V,t),d(e,Z,t),d(e,P,t),n(P,M);for(let c=0;c<v.length;c+=1)v[c]&&v[c].m(M,null);n(P,ne),n(P,A);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(A,null);B=!0},p(e,[t]){var ue,me;(!B||t&1)&&a!==(a=e[0].name+"")&&I(_,a),(!B||t&1)&&j!==(j=e[0].name+"")&&I(N,j);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ue=e[0])==null?void 0:ue.name}').requestVerification('test@example.com');
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').requestVerification('test@example.com');
    `),w.$set(c),(!B||t&1)&&O!==(O=e[0].name+"")&&I(K,O),t&6&&(E=F(e[2]),v=pe(v,t,ce,1,e,E,oe,M,Pe,ke,null,_e)),t&6&&(U=F(e[2]),Be(),k=pe(k,t,re,1,e,U,ie,A,Ce,he,null,be),Se())},i(e){if(!B){ee(w.$$.fragment,e);for(let t=0;t<U.length;t+=1)ee(k[t]);B=!0}},o(e){x(w.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);B=!1},d(e){e&&(r(l),r(i),r(u),r(Q),r(z),r(C),r(G),r(y),r(W),r(T),r(X),r(R),r(Y),r(V),r(Z),r(P)),ve(w,e);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function je(o,l,s){let a,{collection:_}=l,p=204,i=[];const u=m=>s(1,p=m.code);return o.$$set=m=>{"collection"in m&&s(0,_=m.collection)},s(3,a=Te.getApiExampleUrl(Re.baseUrl)),s(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,p,i,a,u]}class Oe extends qe{constructor(l){super(),we(this,l,je,Ue,ye,{collection:0})}}export{Oe as default};
