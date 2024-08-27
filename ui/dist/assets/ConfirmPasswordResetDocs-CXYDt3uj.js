import{S as Ne,i as $e,s as Ce,O as K,e as c,v as w,b as k,c as Ae,f as b,g as r,h as n,m as Re,w as U,P as we,Q as Ee,k as ye,R as De,n as Te,t as ee,a as te,o as p,d as Oe,C as qe,A as Be,q as j,r as Me,N as Fe}from"./index-D0DO79Dq.js";import{S as Ie}from"./SdkTabs-DC6EUYpr.js";function Se(o,l,s){const a=o.slice();return a[5]=l[s],a}function Pe(o,l,s){const a=o.slice();return a[5]=l[s],a}function We(o,l){let s,a=l[5].code+"",_,m,i,u;function f(){return l[4](l[5])}return{key:o,first:null,c(){s=c("button"),_=w(a),m=k(),b(s,"class","tab-item"),j(s,"active",l[1]===l[5].code),this.first=s},m(S,P){r(S,s,P),n(s,_),n(s,m),i||(u=Me(s,"click",f),i=!0)},p(S,P){l=S,P&4&&a!==(a=l[5].code+"")&&U(_,a),P&6&&j(s,"active",l[1]===l[5].code)},d(S){S&&p(s),i=!1,u()}}}function ge(o,l){let s,a,_,m;return a=new Fe({props:{content:l[5].body}}),{key:o,first:null,c(){s=c("div"),Ae(a.$$.fragment),_=k(),b(s,"class","tab-item"),j(s,"active",l[1]===l[5].code),this.first=s},m(i,u){r(i,s,u),Re(a,s,null),n(s,_),m=!0},p(i,u){l=i;const f={};u&4&&(f.content=l[5].body),a.$set(f),(!m||u&6)&&j(s,"active",l[1]===l[5].code)},i(i){m||(ee(a.$$.fragment,i),m=!0)},o(i){te(a.$$.fragment,i),m=!1},d(i){i&&p(s),Oe(a)}}}function Ke(o){var ue,fe,me,be;let l,s,a=o[0].name+"",_,m,i,u,f,S,P,q=o[0].name+"",H,le,se,L,Q,W,z,O,G,g,B,ae,M,N,oe,J,F=o[0].name+"",V,ne,X,$,Y,C,Z,E,x,A,y,v=[],ie=new Map,de,D,h=[],ce=new Map,R;W=new Ie({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        let oldAuth = pb.authStore.model;

        await pb.collection('${(ue=o[0])==null?void 0:ue.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );

        // reauthenticate if needed
        // (after the above call all previously issued tokens are invalidated)
        await pb.collection('${(fe=o[0])==null?void 0:fe.name}').authWithPassword(oldAuth.email, 'NEW_PASSWORD');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        final oldAuth = pb.authStore.model;

        await pb.collection('${(me=o[0])==null?void 0:me.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );

        // reauthenticate if needed
        // (after the above call all previously issued tokens are invalidated)
        await pb.collection('${(be=o[0])==null?void 0:be.name}').authWithPassword(oldAuth.email, 'NEW_PASSWORD');
    `}});let I=K(o[2]);const re=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=Pe(o,I,e),d=re(t);ie.set(d,v[e]=We(d,t))}let T=K(o[2]);const pe=e=>e[5].code;for(let e=0;e<T.length;e+=1){let t=Se(o,T,e),d=pe(t);ce.set(d,h[e]=ge(d,t))}return{c(){l=c("h3"),s=w("Confirm password reset ("),_=w(a),m=w(")"),i=k(),u=c("div"),f=c("p"),S=w("Confirms "),P=c("strong"),H=w(q),le=w(" password reset request and sets a new password."),se=k(),L=c("p"),L.textContent=`After this request all previously issued tokens for the specific record will be automatically
        invalidated.`,Q=k(),Ae(W.$$.fragment),z=k(),O=c("h6"),O.textContent="API details",G=k(),g=c("div"),B=c("strong"),B.textContent="POST",ae=k(),M=c("div"),N=c("p"),oe=w("/api/collections/"),J=c("strong"),V=w(F),ne=w("/confirm-password-reset"),X=k(),$=c("div"),$.textContent="Body Parameters",Y=k(),C=c("table"),C.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the password reset request email.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The new password to set.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>passwordConfirm</span></div></td> <td><span class="label">String</span></td> <td>The new password confirmation.</td></tr></tbody>',Z=k(),E=c("div"),E.textContent="Responses",x=k(),A=c("div"),y=c("div");for(let e=0;e<v.length;e+=1)v[e].c();de=k(),D=c("div");for(let e=0;e<h.length;e+=1)h[e].c();b(l,"class","m-b-sm"),b(u,"class","content txt-lg m-b-sm"),b(O,"class","m-b-xs"),b(B,"class","label label-primary"),b(M,"class","content"),b(g,"class","alert alert-success"),b($,"class","section-title"),b(C,"class","table-compact table-border m-b-base"),b(E,"class","section-title"),b(y,"class","tabs-header compact combined left"),b(D,"class","tabs-content"),b(A,"class","tabs")},m(e,t){r(e,l,t),n(l,s),n(l,_),n(l,m),r(e,i,t),r(e,u,t),n(u,f),n(f,S),n(f,P),n(P,H),n(f,le),n(u,se),n(u,L),r(e,Q,t),Re(W,e,t),r(e,z,t),r(e,O,t),r(e,G,t),r(e,g,t),n(g,B),n(g,ae),n(g,M),n(M,N),n(N,oe),n(N,J),n(J,V),n(N,ne),r(e,X,t),r(e,$,t),r(e,Y,t),r(e,C,t),r(e,Z,t),r(e,E,t),r(e,x,t),r(e,A,t),n(A,y);for(let d=0;d<v.length;d+=1)v[d]&&v[d].m(y,null);n(A,de),n(A,D);for(let d=0;d<h.length;d+=1)h[d]&&h[d].m(D,null);R=!0},p(e,[t]){var _e,he,ke,ve;(!R||t&1)&&a!==(a=e[0].name+"")&&U(_,a),(!R||t&1)&&q!==(q=e[0].name+"")&&U(H,q);const d={};t&9&&(d.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        let oldAuth = pb.authStore.model;

        await pb.collection('${(_e=e[0])==null?void 0:_e.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );

        // reauthenticate if needed
        // (after the above call all previously issued tokens are invalidated)
        await pb.collection('${(he=e[0])==null?void 0:he.name}').authWithPassword(oldAuth.email, 'NEW_PASSWORD');
    `),t&9&&(d.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final oldAuth = pb.authStore.model;

        await pb.collection('${(ke=e[0])==null?void 0:ke.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );

        // reauthenticate if needed
        // (after the above call all previously issued tokens are invalidated)
        await pb.collection('${(ve=e[0])==null?void 0:ve.name}').authWithPassword(oldAuth.email, 'NEW_PASSWORD');
    `),W.$set(d),(!R||t&1)&&F!==(F=e[0].name+"")&&U(V,F),t&6&&(I=K(e[2]),v=we(v,t,re,1,e,I,ie,y,Ee,We,null,Pe)),t&6&&(T=K(e[2]),ye(),h=we(h,t,pe,1,e,T,ce,D,De,ge,null,Se),Te())},i(e){if(!R){ee(W.$$.fragment,e);for(let t=0;t<T.length;t+=1)ee(h[t]);R=!0}},o(e){te(W.$$.fragment,e);for(let t=0;t<h.length;t+=1)te(h[t]);R=!1},d(e){e&&(p(l),p(i),p(u),p(Q),p(z),p(O),p(G),p(g),p(X),p($),p(Y),p(C),p(Z),p(E),p(x),p(A)),Oe(W,e);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<h.length;t+=1)h[t].d()}}}function Ue(o,l,s){let a,{collection:_}=l,m=204,i=[];const u=f=>s(1,m=f.code);return o.$$set=f=>{"collection"in f&&s(0,_=f.collection)},s(3,a=qe.getApiExampleUrl(Be.baseUrl)),s(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,m,i,a,u]}class Le extends Ne{constructor(l){super(),$e(this,l,Ue,Ke,Ce,{collection:0})}}export{Le as default};
