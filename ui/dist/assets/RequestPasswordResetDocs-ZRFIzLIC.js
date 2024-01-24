import{S as Pe,i as $e,s as qe,O as I,e as r,w as g,b as h,c as ve,f as b,g as d,h as n,m as ge,x as L,P as fe,Q as ye,k as Re,R as Be,n as Ce,t as x,a as ee,o as p,d as we,C as Se,p as Te,r as N,u as Me,N as Ae}from"./index-78piLIP3.js";import{S as Ue}from"./SdkTabs-c6VuPJvR.js";function be(o,s,l){const a=o.slice();return a[5]=s[l],a}function _e(o,s,l){const a=o.slice();return a[5]=s[l],a}function ke(o,s){let l,a=s[5].code+"",_,f,i,u;function m(){return s[4](s[5])}return{key:o,first:null,c(){l=r("button"),_=g(a),f=h(),b(l,"class","tab-item"),N(l,"active",s[1]===s[5].code),this.first=l},m(w,P){d(w,l,P),n(l,_),n(l,f),i||(u=Me(l,"click",m),i=!0)},p(w,P){s=w,P&4&&a!==(a=s[5].code+"")&&L(_,a),P&6&&N(l,"active",s[1]===s[5].code)},d(w){w&&p(l),i=!1,u()}}}function he(o,s){let l,a,_,f;return a=new Ae({props:{content:s[5].body}}),{key:o,first:null,c(){l=r("div"),ve(a.$$.fragment),_=h(),b(l,"class","tab-item"),N(l,"active",s[1]===s[5].code),this.first=l},m(i,u){d(i,l,u),ge(a,l,null),n(l,_),f=!0},p(i,u){s=i;const m={};u&4&&(m.content=s[5].body),a.$set(m),(!f||u&6)&&N(l,"active",s[1]===s[5].code)},i(i){f||(x(a.$$.fragment,i),f=!0)},o(i){ee(a.$$.fragment,i),f=!1},d(i){i&&p(l),we(a)}}}function je(o){var de,pe;let s,l,a=o[0].name+"",_,f,i,u,m,w,P,D=o[0].name+"",Q,te,z,$,G,B,J,q,H,se,O,C,le,K,E=o[0].name+"",V,ae,W,S,X,T,Y,M,Z,y,A,v=[],oe=new Map,ne,U,k=[],ie=new Map,R;$=new Ue({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(de=o[0])==null?void 0:de.name}').requestPasswordReset('test@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(pe=o[0])==null?void 0:pe.name}').requestPasswordReset('test@example.com');
    `}});let F=I(o[2]);const ce=e=>e[5].code;for(let e=0;e<F.length;e+=1){let t=_e(o,F,e),c=ce(t);oe.set(c,v[e]=ke(c,t))}let j=I(o[2]);const re=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=be(o,j,e),c=re(t);ie.set(c,k[e]=he(c,t))}return{c(){s=r("h3"),l=g("Request password reset ("),_=g(a),f=g(")"),i=h(),u=r("div"),m=r("p"),w=g("Sends "),P=r("strong"),Q=g(D),te=g(" password reset email request."),z=h(),ve($.$$.fragment),G=h(),B=r("h6"),B.textContent="API details",J=h(),q=r("div"),H=r("strong"),H.textContent="POST",se=h(),O=r("div"),C=r("p"),le=g("/api/collections/"),K=r("strong"),V=g(E),ae=g("/request-password-reset"),W=h(),S=r("div"),S.textContent="Body Parameters",X=h(),T=r("table"),T.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>email</span></div></td> <td><span class="label">String</span></td> <td>The auth record email address to send the password reset request (if exists).</td></tr></tbody>',Y=h(),M=r("div"),M.textContent="Responses",Z=h(),y=r("div"),A=r("div");for(let e=0;e<v.length;e+=1)v[e].c();ne=h(),U=r("div");for(let e=0;e<k.length;e+=1)k[e].c();b(s,"class","m-b-sm"),b(u,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(H,"class","label label-primary"),b(O,"class","content"),b(q,"class","alert alert-success"),b(S,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(A,"class","tabs-header compact combined left"),b(U,"class","tabs-content"),b(y,"class","tabs")},m(e,t){d(e,s,t),n(s,l),n(s,_),n(s,f),d(e,i,t),d(e,u,t),n(u,m),n(m,w),n(m,P),n(P,Q),n(m,te),d(e,z,t),ge($,e,t),d(e,G,t),d(e,B,t),d(e,J,t),d(e,q,t),n(q,H),n(q,se),n(q,O),n(O,C),n(C,le),n(C,K),n(K,V),n(C,ae),d(e,W,t),d(e,S,t),d(e,X,t),d(e,T,t),d(e,Y,t),d(e,M,t),d(e,Z,t),d(e,y,t),n(y,A);for(let c=0;c<v.length;c+=1)v[c]&&v[c].m(A,null);n(y,ne),n(y,U);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(U,null);R=!0},p(e,[t]){var ue,me;(!R||t&1)&&a!==(a=e[0].name+"")&&L(_,a),(!R||t&1)&&D!==(D=e[0].name+"")&&L(Q,D);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ue=e[0])==null?void 0:ue.name}').requestPasswordReset('test@example.com');
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').requestPasswordReset('test@example.com');
    `),$.$set(c),(!R||t&1)&&E!==(E=e[0].name+"")&&L(V,E),t&6&&(F=I(e[2]),v=fe(v,t,ce,1,e,F,oe,A,ye,ke,null,_e)),t&6&&(j=I(e[2]),Re(),k=fe(k,t,re,1,e,j,ie,U,Be,he,null,be),Ce())},i(e){if(!R){x($.$$.fragment,e);for(let t=0;t<j.length;t+=1)x(k[t]);R=!0}},o(e){ee($.$$.fragment,e);for(let t=0;t<k.length;t+=1)ee(k[t]);R=!1},d(e){e&&(p(s),p(i),p(u),p(z),p(G),p(B),p(J),p(q),p(W),p(S),p(X),p(T),p(Y),p(M),p(Z),p(y)),we($,e);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function De(o,s,l){let a,{collection:_}=s,f=204,i=[];const u=m=>l(1,f=m.code);return o.$$set=m=>{"collection"in m&&l(0,_=m.collection)},l(3,a=Se.getApiExampleUrl(Te.baseUrl)),l(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,f,i,a,u]}class Ee extends Pe{constructor(s){super(),$e(this,s,De,je,qe,{collection:0})}}export{Ee as default};
