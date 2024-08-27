import{S as Pe,i as $e,s as qe,O as I,e as r,v as g,b as h,c as ve,f as b,g as d,h as n,m as ge,w as L,P as fe,Q as ye,k as Re,R as Ce,n as Be,t as x,a as ee,o as p,d as we,C as Se,A as Te,q as N,r as Ae,N as Me}from"./index-D0DO79Dq.js";import{S as Ue}from"./SdkTabs-DC6EUYpr.js";function be(o,s,l){const a=o.slice();return a[5]=s[l],a}function _e(o,s,l){const a=o.slice();return a[5]=s[l],a}function ke(o,s){let l,a=s[5].code+"",_,f,i,m;function u(){return s[4](s[5])}return{key:o,first:null,c(){l=r("button"),_=g(a),f=h(),b(l,"class","tab-item"),N(l,"active",s[1]===s[5].code),this.first=l},m(w,P){d(w,l,P),n(l,_),n(l,f),i||(m=Ae(l,"click",u),i=!0)},p(w,P){s=w,P&4&&a!==(a=s[5].code+"")&&L(_,a),P&6&&N(l,"active",s[1]===s[5].code)},d(w){w&&p(l),i=!1,m()}}}function he(o,s){let l,a,_,f;return a=new Me({props:{content:s[5].body}}),{key:o,first:null,c(){l=r("div"),ve(a.$$.fragment),_=h(),b(l,"class","tab-item"),N(l,"active",s[1]===s[5].code),this.first=l},m(i,m){d(i,l,m),ge(a,l,null),n(l,_),f=!0},p(i,m){s=i;const u={};m&4&&(u.content=s[5].body),a.$set(u),(!f||m&6)&&N(l,"active",s[1]===s[5].code)},i(i){f||(x(a.$$.fragment,i),f=!0)},o(i){ee(a.$$.fragment,i),f=!1},d(i){i&&p(l),we(a)}}}function je(o){var de,pe;let s,l,a=o[0].name+"",_,f,i,m,u,w,P,D=o[0].name+"",Q,te,z,$,G,C,J,q,H,se,O,B,le,K,E=o[0].name+"",V,ae,W,S,X,T,Y,A,Z,y,M,v=[],oe=new Map,ne,U,k=[],ie=new Map,R;$=new Ue({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(de=o[0])==null?void 0:de.name}').requestPasswordReset('test@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(pe=o[0])==null?void 0:pe.name}').requestPasswordReset('test@example.com');
    `}});let F=I(o[2]);const ce=e=>e[5].code;for(let e=0;e<F.length;e+=1){let t=_e(o,F,e),c=ce(t);oe.set(c,v[e]=ke(c,t))}let j=I(o[2]);const re=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=be(o,j,e),c=re(t);ie.set(c,k[e]=he(c,t))}return{c(){s=r("h3"),l=g("Request password reset ("),_=g(a),f=g(")"),i=h(),m=r("div"),u=r("p"),w=g("Sends "),P=r("strong"),Q=g(D),te=g(" password reset email request."),z=h(),ve($.$$.fragment),G=h(),C=r("h6"),C.textContent="API details",J=h(),q=r("div"),H=r("strong"),H.textContent="POST",se=h(),O=r("div"),B=r("p"),le=g("/api/collections/"),K=r("strong"),V=g(E),ae=g("/request-password-reset"),W=h(),S=r("div"),S.textContent="Body Parameters",X=h(),T=r("table"),T.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address to send the password reset request (if exists).</td></tr></tbody>',Y=h(),A=r("div"),A.textContent="Responses",Z=h(),y=r("div"),M=r("div");for(let e=0;e<v.length;e+=1)v[e].c();ne=h(),U=r("div");for(let e=0;e<k.length;e+=1)k[e].c();b(s,"class","m-b-sm"),b(m,"class","content txt-lg m-b-sm"),b(C,"class","m-b-xs"),b(H,"class","label label-primary"),b(O,"class","content"),b(q,"class","alert alert-success"),b(S,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(A,"class","section-title"),b(M,"class","tabs-header compact combined left"),b(U,"class","tabs-content"),b(y,"class","tabs")},m(e,t){d(e,s,t),n(s,l),n(s,_),n(s,f),d(e,i,t),d(e,m,t),n(m,u),n(u,w),n(u,P),n(P,Q),n(u,te),d(e,z,t),ge($,e,t),d(e,G,t),d(e,C,t),d(e,J,t),d(e,q,t),n(q,H),n(q,se),n(q,O),n(O,B),n(B,le),n(B,K),n(K,V),n(B,ae),d(e,W,t),d(e,S,t),d(e,X,t),d(e,T,t),d(e,Y,t),d(e,A,t),d(e,Z,t),d(e,y,t),n(y,M);for(let c=0;c<v.length;c+=1)v[c]&&v[c].m(M,null);n(y,ne),n(y,U);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(U,null);R=!0},p(e,[t]){var me,ue;(!R||t&1)&&a!==(a=e[0].name+"")&&L(_,a),(!R||t&1)&&D!==(D=e[0].name+"")&&L(Q,D);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').requestPasswordReset('test@example.com');
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(ue=e[0])==null?void 0:ue.name}').requestPasswordReset('test@example.com');
    `),$.$set(c),(!R||t&1)&&E!==(E=e[0].name+"")&&L(V,E),t&6&&(F=I(e[2]),v=fe(v,t,ce,1,e,F,oe,M,ye,ke,null,_e)),t&6&&(j=I(e[2]),Re(),k=fe(k,t,re,1,e,j,ie,U,Ce,he,null,be),Be())},i(e){if(!R){x($.$$.fragment,e);for(let t=0;t<j.length;t+=1)x(k[t]);R=!0}},o(e){ee($.$$.fragment,e);for(let t=0;t<k.length;t+=1)ee(k[t]);R=!1},d(e){e&&(p(s),p(i),p(m),p(z),p(G),p(C),p(J),p(q),p(W),p(S),p(X),p(T),p(Y),p(A),p(Z),p(y)),we($,e);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function De(o,s,l){let a,{collection:_}=s,f=204,i=[];const m=u=>l(1,f=u.code);return o.$$set=u=>{"collection"in u&&l(0,_=u.collection)},l(3,a=Se.getApiExampleUrl(Te.baseUrl)),l(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,f,i,a,m]}class Ee extends Pe{constructor(s){super(),$e(this,s,De,je,qe,{collection:0})}}export{Ee as default};
