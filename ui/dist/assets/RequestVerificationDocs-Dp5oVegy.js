import{S as qe,i as we,s as Pe,O as F,e as r,v as g,b as h,c as ve,f as b,g as d,h as n,m as ge,w as I,P as pe,Q as ye,k as Ce,R as Be,n as Se,t as x,a as ee,o as f,d as $e,C as Te,A as Ae,q as L,r as Re,N as Ve}from"./index-D0DO79Dq.js";import{S as Me}from"./SdkTabs-DC6EUYpr.js";function be(o,l,s){const a=o.slice();return a[5]=l[s],a}function _e(o,l,s){const a=o.slice();return a[5]=l[s],a}function ke(o,l){let s,a=l[5].code+"",_,p,i,m;function u(){return l[4](l[5])}return{key:o,first:null,c(){s=r("button"),_=g(a),p=h(),b(s,"class","tab-item"),L(s,"active",l[1]===l[5].code),this.first=s},m($,q){d($,s,q),n(s,_),n(s,p),i||(m=Re(s,"click",u),i=!0)},p($,q){l=$,q&4&&a!==(a=l[5].code+"")&&I(_,a),q&6&&L(s,"active",l[1]===l[5].code)},d($){$&&f(s),i=!1,m()}}}function he(o,l){let s,a,_,p;return a=new Ve({props:{content:l[5].body}}),{key:o,first:null,c(){s=r("div"),ve(a.$$.fragment),_=h(),b(s,"class","tab-item"),L(s,"active",l[1]===l[5].code),this.first=s},m(i,m){d(i,s,m),ge(a,s,null),n(s,_),p=!0},p(i,m){l=i;const u={};m&4&&(u.content=l[5].body),a.$set(u),(!p||m&6)&&L(s,"active",l[1]===l[5].code)},i(i){p||(x(a.$$.fragment,i),p=!0)},o(i){ee(a.$$.fragment,i),p=!1},d(i){i&&f(s),$e(a)}}}function Ue(o){var de,fe;let l,s,a=o[0].name+"",_,p,i,m,u,$,q,j=o[0].name+"",N,te,Q,w,z,B,G,P,D,le,H,S,se,J,O=o[0].name+"",K,ae,W,T,X,A,Y,R,Z,y,V,v=[],oe=new Map,ne,M,k=[],ie=new Map,C;w=new Me({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(de=o[0])==null?void 0:de.name}').requestVerification('test@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(fe=o[0])==null?void 0:fe.name}').requestVerification('test@example.com');
    `}});let E=F(o[2]);const ce=e=>e[5].code;for(let e=0;e<E.length;e+=1){let t=_e(o,E,e),c=ce(t);oe.set(c,v[e]=ke(c,t))}let U=F(o[2]);const re=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=be(o,U,e),c=re(t);ie.set(c,k[e]=he(c,t))}return{c(){l=r("h3"),s=g("Request verification ("),_=g(a),p=g(")"),i=h(),m=r("div"),u=r("p"),$=g("Sends "),q=r("strong"),N=g(j),te=g(" verification email request."),Q=h(),ve(w.$$.fragment),z=h(),B=r("h6"),B.textContent="API details",G=h(),P=r("div"),D=r("strong"),D.textContent="POST",le=h(),H=r("div"),S=r("p"),se=g("/api/collections/"),J=r("strong"),K=g(O),ae=g("/request-verification"),W=h(),T=r("div"),T.textContent="Body Parameters",X=h(),A=r("table"),A.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address to send the verification request (if exists).</td></tr></tbody>',Y=h(),R=r("div"),R.textContent="Responses",Z=h(),y=r("div"),V=r("div");for(let e=0;e<v.length;e+=1)v[e].c();ne=h(),M=r("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(m,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(D,"class","label label-primary"),b(H,"class","content"),b(P,"class","alert alert-success"),b(T,"class","section-title"),b(A,"class","table-compact table-border m-b-base"),b(R,"class","section-title"),b(V,"class","tabs-header compact combined left"),b(M,"class","tabs-content"),b(y,"class","tabs")},m(e,t){d(e,l,t),n(l,s),n(l,_),n(l,p),d(e,i,t),d(e,m,t),n(m,u),n(u,$),n(u,q),n(q,N),n(u,te),d(e,Q,t),ge(w,e,t),d(e,z,t),d(e,B,t),d(e,G,t),d(e,P,t),n(P,D),n(P,le),n(P,H),n(H,S),n(S,se),n(S,J),n(J,K),n(S,ae),d(e,W,t),d(e,T,t),d(e,X,t),d(e,A,t),d(e,Y,t),d(e,R,t),d(e,Z,t),d(e,y,t),n(y,V);for(let c=0;c<v.length;c+=1)v[c]&&v[c].m(V,null);n(y,ne),n(y,M);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(M,null);C=!0},p(e,[t]){var me,ue;(!C||t&1)&&a!==(a=e[0].name+"")&&I(_,a),(!C||t&1)&&j!==(j=e[0].name+"")&&I(N,j);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').requestVerification('test@example.com');
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(ue=e[0])==null?void 0:ue.name}').requestVerification('test@example.com');
    `),w.$set(c),(!C||t&1)&&O!==(O=e[0].name+"")&&I(K,O),t&6&&(E=F(e[2]),v=pe(v,t,ce,1,e,E,oe,V,ye,ke,null,_e)),t&6&&(U=F(e[2]),Ce(),k=pe(k,t,re,1,e,U,ie,M,Be,he,null,be),Se())},i(e){if(!C){x(w.$$.fragment,e);for(let t=0;t<U.length;t+=1)x(k[t]);C=!0}},o(e){ee(w.$$.fragment,e);for(let t=0;t<k.length;t+=1)ee(k[t]);C=!1},d(e){e&&(f(l),f(i),f(m),f(Q),f(z),f(B),f(G),f(P),f(W),f(T),f(X),f(A),f(Y),f(R),f(Z),f(y)),$e(w,e);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function je(o,l,s){let a,{collection:_}=l,p=204,i=[];const m=u=>s(1,p=u.code);return o.$$set=u=>{"collection"in u&&s(0,_=u.collection)},s(3,a=Te.getApiExampleUrl(Ae.baseUrl)),s(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,p,i,a,m]}class Oe extends qe{constructor(l){super(),we(this,l,je,Ue,Pe,{collection:0})}}export{Oe as default};
