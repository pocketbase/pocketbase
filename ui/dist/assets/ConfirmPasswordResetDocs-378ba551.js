import{S as Pe,i as Se,s as Re,N as K,e as r,w,b as h,c as ve,f as _,g as d,h as o,m as we,x as U,P as ue,Q as Ne,k as Oe,R as Ce,n as We,t as x,a as ee,o as p,d as ge,U as $e,C as Ee,p as Te,r as j,u as Ae,M as ye}from"./index-a084d9d7.js";import{S as De}from"./SdkTabs-ba0ec979.js";function _e(n,s,l){const a=n.slice();return a[5]=s[l],a}function be(n,s,l){const a=n.slice();return a[5]=s[l],a}function ke(n,s){let l,a=s[5].code+"",b,u,i,f;function m(){return s[4](s[5])}return{key:n,first:null,c(){l=r("button"),b=w(a),u=h(),_(l,"class","tab-item"),j(l,"active",s[1]===s[5].code),this.first=l},m(g,P){d(g,l,P),o(l,b),o(l,u),i||(f=Ae(l,"click",m),i=!0)},p(g,P){s=g,P&4&&a!==(a=s[5].code+"")&&U(b,a),P&6&&j(l,"active",s[1]===s[5].code)},d(g){g&&p(l),i=!1,f()}}}function he(n,s){let l,a,b,u;return a=new ye({props:{content:s[5].body}}),{key:n,first:null,c(){l=r("div"),ve(a.$$.fragment),b=h(),_(l,"class","tab-item"),j(l,"active",s[1]===s[5].code),this.first=l},m(i,f){d(i,l,f),we(a,l,null),o(l,b),u=!0},p(i,f){s=i;const m={};f&4&&(m.content=s[5].body),a.$set(m),(!u||f&6)&&j(l,"active",s[1]===s[5].code)},i(i){u||(x(a.$$.fragment,i),u=!0)},o(i){ee(a.$$.fragment,i),u=!1},d(i){i&&p(l),ge(a)}}}function Me(n){var de,pe;let s,l,a=n[0].name+"",b,u,i,f,m,g,P,M=n[0].name+"",H,te,L,S,Q,C,z,R,q,se,B,W,le,G,F=n[0].name+"",J,ae,V,$,X,E,Y,T,Z,N,A,v=[],ne=new Map,oe,y,k=[],ie=new Map,O;S=new De({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[3]}');

        ...

        await pb.collection('${(de=n[0])==null?void 0:de.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[3]}');

        ...

        await pb.collection('${(pe=n[0])==null?void 0:pe.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );
    `}});let I=K(n[2]);const ce=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=be(n,I,e),c=ce(t);ne.set(c,v[e]=ke(c,t))}let D=K(n[2]);const re=e=>e[5].code;for(let e=0;e<D.length;e+=1){let t=_e(n,D,e),c=re(t);ie.set(c,k[e]=he(c,t))}return{c(){s=r("h3"),l=w("Confirm password reset ("),b=w(a),u=w(")"),i=h(),f=r("div"),m=r("p"),g=w("Confirms "),P=r("strong"),H=w(M),te=w(" password reset request and sets a new password."),L=h(),ve(S.$$.fragment),Q=h(),C=r("h6"),C.textContent="API details",z=h(),R=r("div"),q=r("strong"),q.textContent="POST",se=h(),B=r("div"),W=r("p"),le=w("/api/collections/"),G=r("strong"),J=w(F),ae=w("/confirm-password-reset"),V=h(),$=r("div"),$.textContent="Body Parameters",X=h(),E=r("table"),E.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the password reset request email.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The new password to set.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>The new password confirmation.</td></tr></tbody>',Y=h(),T=r("div"),T.textContent="Responses",Z=h(),N=r("div"),A=r("div");for(let e=0;e<v.length;e+=1)v[e].c();oe=h(),y=r("div");for(let e=0;e<k.length;e+=1)k[e].c();_(s,"class","m-b-sm"),_(f,"class","content txt-lg m-b-sm"),_(C,"class","m-b-xs"),_(q,"class","label label-primary"),_(B,"class","content"),_(R,"class","alert alert-success"),_($,"class","section-title"),_(E,"class","table-compact table-border m-b-base"),_(T,"class","section-title"),_(A,"class","tabs-header compact left"),_(y,"class","tabs-content"),_(N,"class","tabs")},m(e,t){d(e,s,t),o(s,l),o(s,b),o(s,u),d(e,i,t),d(e,f,t),o(f,m),o(m,g),o(m,P),o(P,H),o(m,te),d(e,L,t),we(S,e,t),d(e,Q,t),d(e,C,t),d(e,z,t),d(e,R,t),o(R,q),o(R,se),o(R,B),o(B,W),o(W,le),o(W,G),o(G,J),o(W,ae),d(e,V,t),d(e,$,t),d(e,X,t),d(e,E,t),d(e,Y,t),d(e,T,t),d(e,Z,t),d(e,N,t),o(N,A);for(let c=0;c<v.length;c+=1)v[c]&&v[c].m(A,null);o(N,oe),o(N,y);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(y,null);O=!0},p(e,[t]){var fe,me;(!O||t&1)&&a!==(a=e[0].name+"")&&U(b,a),(!O||t&1)&&M!==(M=e[0].name+"")&&U(H,M);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(fe=e[0])==null?void 0:fe.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );
    `),S.$set(c),(!O||t&1)&&F!==(F=e[0].name+"")&&U(J,F),t&6&&(I=K(e[2]),v=ue(v,t,ce,1,e,I,ne,A,Ne,ke,null,be)),t&6&&(D=K(e[2]),Oe(),k=ue(k,t,re,1,e,D,ie,y,Ce,he,null,_e),We())},i(e){if(!O){x(S.$$.fragment,e);for(let t=0;t<D.length;t+=1)x(k[t]);O=!0}},o(e){ee(S.$$.fragment,e);for(let t=0;t<k.length;t+=1)ee(k[t]);O=!1},d(e){e&&(p(s),p(i),p(f),p(L),p(Q),p(C),p(z),p(R),p(V),p($),p(X),p(E),p(Y),p(T),p(Z),p(N)),ge(S,e);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function qe(n,s,l){let a,{collection:b=new $e}=s,u=204,i=[];const f=m=>l(1,u=m.code);return n.$$set=m=>{"collection"in m&&l(0,b=m.collection)},l(3,a=Ee.getApiExampleUrl(Te.baseUrl)),l(2,i=[{code:204,body:"null"},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "token": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}]),[b,u,i,a,f]}class Ie extends Pe{constructor(s){super(),Se(this,s,qe,Me,Re,{collection:0})}}export{Ie as default};
