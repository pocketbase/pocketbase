import{S as Re,i as Oe,s as Ce,O as K,e as r,w,b as h,c as ge,f as b,g as d,h as n,m as Pe,x as U,P as _e,Q as Ne,k as We,R as $e,n as Ee,t as ee,a as te,o as p,d as Se,C as ye,p as Ae,r as j,u as Te,N as De}from"./index-8354bde7.js";import{S as qe}from"./SdkTabs-86785e52.js";function ke(o,s,l){const a=o.slice();return a[5]=s[l],a}function he(o,s,l){const a=o.slice();return a[5]=s[l],a}function ve(o,s){let l,a=s[5].code+"",_,u,i,f;function m(){return s[4](s[5])}return{key:o,first:null,c(){l=r("button"),_=w(a),u=h(),b(l,"class","tab-item"),j(l,"active",s[1]===s[5].code),this.first=l},m(g,P){d(g,l,P),n(l,_),n(l,u),i||(f=Te(l,"click",m),i=!0)},p(g,P){s=g,P&4&&a!==(a=s[5].code+"")&&U(_,a),P&6&&j(l,"active",s[1]===s[5].code)},d(g){g&&p(l),i=!1,f()}}}function we(o,s){let l,a,_,u;return a=new De({props:{content:s[5].body}}),{key:o,first:null,c(){l=r("div"),ge(a.$$.fragment),_=h(),b(l,"class","tab-item"),j(l,"active",s[1]===s[5].code),this.first=l},m(i,f){d(i,l,f),Pe(a,l,null),n(l,_),u=!0},p(i,f){s=i;const m={};f&4&&(m.content=s[5].body),a.$set(m),(!u||f&6)&&j(l,"active",s[1]===s[5].code)},i(i){u||(ee(a.$$.fragment,i),u=!0)},o(i){te(a.$$.fragment,i),u=!1},d(i){i&&p(l),Se(a)}}}function Be(o){var fe,me;let s,l,a=o[0].name+"",_,u,i,f,m,g,P,q=o[0].name+"",H,se,le,L,Q,S,z,N,G,R,B,ae,M,W,ne,J,F=o[0].name+"",V,oe,X,$,Y,E,Z,y,x,O,A,v=[],ie=new Map,ce,T,k=[],re=new Map,C;S=new qe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(fe=o[0])==null?void 0:fe.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(me=o[0])==null?void 0:me.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );
    `}});let I=K(o[2]);const de=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=he(o,I,e),c=de(t);ie.set(c,v[e]=ve(c,t))}let D=K(o[2]);const pe=e=>e[5].code;for(let e=0;e<D.length;e+=1){let t=ke(o,D,e),c=pe(t);re.set(c,k[e]=we(c,t))}return{c(){s=r("h3"),l=w("Confirm password reset ("),_=w(a),u=w(")"),i=h(),f=r("div"),m=r("p"),g=w("Confirms "),P=r("strong"),H=w(q),se=w(" password reset request and sets a new password."),le=h(),L=r("p"),L.textContent=`After this request all previously issued tokens for the specific record will be automatically
        invalidated.`,Q=h(),ge(S.$$.fragment),z=h(),N=r("h6"),N.textContent="API details",G=h(),R=r("div"),B=r("strong"),B.textContent="POST",ae=h(),M=r("div"),W=r("p"),ne=w("/api/collections/"),J=r("strong"),V=w(F),oe=w("/confirm-password-reset"),X=h(),$=r("div"),$.textContent="Body Parameters",Y=h(),E=r("table"),E.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the password reset request email.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The new password to set.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>The new password confirmation.</td></tr></tbody>',Z=h(),y=r("div"),y.textContent="Responses",x=h(),O=r("div"),A=r("div");for(let e=0;e<v.length;e+=1)v[e].c();ce=h(),T=r("div");for(let e=0;e<k.length;e+=1)k[e].c();b(s,"class","m-b-sm"),b(f,"class","content txt-lg m-b-sm"),b(N,"class","m-b-xs"),b(B,"class","label label-primary"),b(M,"class","content"),b(R,"class","alert alert-success"),b($,"class","section-title"),b(E,"class","table-compact table-border m-b-base"),b(y,"class","section-title"),b(A,"class","tabs-header compact combined left"),b(T,"class","tabs-content"),b(O,"class","tabs")},m(e,t){d(e,s,t),n(s,l),n(s,_),n(s,u),d(e,i,t),d(e,f,t),n(f,m),n(m,g),n(m,P),n(P,H),n(m,se),n(f,le),n(f,L),d(e,Q,t),Pe(S,e,t),d(e,z,t),d(e,N,t),d(e,G,t),d(e,R,t),n(R,B),n(R,ae),n(R,M),n(M,W),n(W,ne),n(W,J),n(J,V),n(W,oe),d(e,X,t),d(e,$,t),d(e,Y,t),d(e,E,t),d(e,Z,t),d(e,y,t),d(e,x,t),d(e,O,t),n(O,A);for(let c=0;c<v.length;c+=1)v[c]&&v[c].m(A,null);n(O,ce),n(O,T);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(T,null);C=!0},p(e,[t]){var ue,be;(!C||t&1)&&a!==(a=e[0].name+"")&&U(_,a),(!C||t&1)&&q!==(q=e[0].name+"")&&U(H,q);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ue=e[0])==null?void 0:ue.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );
    `),S.$set(c),(!C||t&1)&&F!==(F=e[0].name+"")&&U(V,F),t&6&&(I=K(e[2]),v=_e(v,t,de,1,e,I,ie,A,Ne,ve,null,he)),t&6&&(D=K(e[2]),We(),k=_e(k,t,pe,1,e,D,re,T,$e,we,null,ke),Ee())},i(e){if(!C){ee(S.$$.fragment,e);for(let t=0;t<D.length;t+=1)ee(k[t]);C=!0}},o(e){te(S.$$.fragment,e);for(let t=0;t<k.length;t+=1)te(k[t]);C=!1},d(e){e&&(p(s),p(i),p(f),p(Q),p(z),p(N),p(G),p(R),p(X),p($),p(Y),p(E),p(Z),p(y),p(x),p(O)),Se(S,e);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function Me(o,s,l){let a,{collection:_}=s,u=204,i=[];const f=m=>l(1,u=m.code);return o.$$set=m=>{"collection"in m&&l(0,_=m.collection)},l(3,a=ye.getApiExampleUrl(Ae.baseUrl)),l(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,u,i,a,f]}class Ke extends Re{constructor(s){super(),Oe(this,s,Me,Be,Ce,{collection:0})}}export{Ke as default};
