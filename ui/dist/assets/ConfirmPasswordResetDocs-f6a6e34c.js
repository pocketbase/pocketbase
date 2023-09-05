import{S as Pe,i as Se,s as Re,O as K,e as r,w,b as h,c as ve,f as b,g as d,h as o,m as we,x as U,P as ue,Q as Oe,k as Ne,R as Ce,n as We,t as x,a as ee,o as p,d as ge,C as $e,p as Ee,r as j,u as Te,N as Ae}from"./index-93a5b881.js";import{S as ye}from"./SdkTabs-366cf0cf.js";function be(n,s,l){const a=n.slice();return a[5]=s[l],a}function _e(n,s,l){const a=n.slice();return a[5]=s[l],a}function ke(n,s){let l,a=s[5].code+"",_,u,i,f;function m(){return s[4](s[5])}return{key:n,first:null,c(){l=r("button"),_=w(a),u=h(),b(l,"class","tab-item"),j(l,"active",s[1]===s[5].code),this.first=l},m(g,P){d(g,l,P),o(l,_),o(l,u),i||(f=Te(l,"click",m),i=!0)},p(g,P){s=g,P&4&&a!==(a=s[5].code+"")&&U(_,a),P&6&&j(l,"active",s[1]===s[5].code)},d(g){g&&p(l),i=!1,f()}}}function he(n,s){let l,a,_,u;return a=new Ae({props:{content:s[5].body}}),{key:n,first:null,c(){l=r("div"),ve(a.$$.fragment),_=h(),b(l,"class","tab-item"),j(l,"active",s[1]===s[5].code),this.first=l},m(i,f){d(i,l,f),we(a,l,null),o(l,_),u=!0},p(i,f){s=i;const m={};f&4&&(m.content=s[5].body),a.$set(m),(!u||f&6)&&j(l,"active",s[1]===s[5].code)},i(i){u||(x(a.$$.fragment,i),u=!0)},o(i){ee(a.$$.fragment,i),u=!1},d(i){i&&p(l),ge(a)}}}function De(n){var de,pe;let s,l,a=n[0].name+"",_,u,i,f,m,g,P,q=n[0].name+"",H,te,L,S,Q,C,z,R,B,se,M,W,le,G,F=n[0].name+"",J,ae,V,$,X,E,Y,T,Z,O,A,v=[],ne=new Map,oe,y,k=[],ie=new Map,N;S=new ye({props:{js:`
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
    `}});let I=K(n[2]);const ce=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=_e(n,I,e),c=ce(t);ne.set(c,v[e]=ke(c,t))}let D=K(n[2]);const re=e=>e[5].code;for(let e=0;e<D.length;e+=1){let t=be(n,D,e),c=re(t);ie.set(c,k[e]=he(c,t))}return{c(){s=r("h3"),l=w("Confirm password reset ("),_=w(a),u=w(")"),i=h(),f=r("div"),m=r("p"),g=w("Confirms "),P=r("strong"),H=w(q),te=w(" password reset request and sets a new password."),L=h(),ve(S.$$.fragment),Q=h(),C=r("h6"),C.textContent="API details",z=h(),R=r("div"),B=r("strong"),B.textContent="POST",se=h(),M=r("div"),W=r("p"),le=w("/api/collections/"),G=r("strong"),J=w(F),ae=w("/confirm-password-reset"),V=h(),$=r("div"),$.textContent="Body Parameters",X=h(),E=r("table"),E.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the password reset request email.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The new password to set.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>The new password confirmation.</td></tr></tbody>',Y=h(),T=r("div"),T.textContent="Responses",Z=h(),O=r("div"),A=r("div");for(let e=0;e<v.length;e+=1)v[e].c();oe=h(),y=r("div");for(let e=0;e<k.length;e+=1)k[e].c();b(s,"class","m-b-sm"),b(f,"class","content txt-lg m-b-sm"),b(C,"class","m-b-xs"),b(B,"class","label label-primary"),b(M,"class","content"),b(R,"class","alert alert-success"),b($,"class","section-title"),b(E,"class","table-compact table-border m-b-base"),b(T,"class","section-title"),b(A,"class","tabs-header compact combined left"),b(y,"class","tabs-content"),b(O,"class","tabs")},m(e,t){d(e,s,t),o(s,l),o(s,_),o(s,u),d(e,i,t),d(e,f,t),o(f,m),o(m,g),o(m,P),o(P,H),o(m,te),d(e,L,t),we(S,e,t),d(e,Q,t),d(e,C,t),d(e,z,t),d(e,R,t),o(R,B),o(R,se),o(R,M),o(M,W),o(W,le),o(W,G),o(G,J),o(W,ae),d(e,V,t),d(e,$,t),d(e,X,t),d(e,E,t),d(e,Y,t),d(e,T,t),d(e,Z,t),d(e,O,t),o(O,A);for(let c=0;c<v.length;c+=1)v[c]&&v[c].m(A,null);o(O,oe),o(O,y);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(y,null);N=!0},p(e,[t]){var fe,me;(!N||t&1)&&a!==(a=e[0].name+"")&&U(_,a),(!N||t&1)&&q!==(q=e[0].name+"")&&U(H,q);const c={};t&9&&(c.js=`
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
    `),S.$set(c),(!N||t&1)&&F!==(F=e[0].name+"")&&U(J,F),t&6&&(I=K(e[2]),v=ue(v,t,ce,1,e,I,ne,A,Oe,ke,null,_e)),t&6&&(D=K(e[2]),Ne(),k=ue(k,t,re,1,e,D,ie,y,Ce,he,null,be),We())},i(e){if(!N){x(S.$$.fragment,e);for(let t=0;t<D.length;t+=1)x(k[t]);N=!0}},o(e){ee(S.$$.fragment,e);for(let t=0;t<k.length;t+=1)ee(k[t]);N=!1},d(e){e&&(p(s),p(i),p(f),p(L),p(Q),p(C),p(z),p(R),p(V),p($),p(X),p(E),p(Y),p(T),p(Z),p(O)),ge(S,e);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function qe(n,s,l){let a,{collection:_}=s,u=204,i=[];const f=m=>l(1,u=m.code);return n.$$set=m=>{"collection"in m&&l(0,_=m.collection)},l(3,a=$e.getApiExampleUrl(Ee.baseUrl)),l(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,u,i,a,f]}class Fe extends Pe{constructor(s){super(),Se(this,s,qe,De,Re,{collection:0})}}export{Fe as default};
