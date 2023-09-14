import{S as $e,i as Pe,s as Se,O as Y,e as r,w as v,b as k,c as ge,f as b,g as d,h as o,m as ve,x as j,P as ue,Q as we,k as Oe,R as Re,n as Te,t as x,a as ee,o as m,d as Ce,C as ye,p as Ee,r as H,u as Be,N as qe}from"./index-8354bde7.js";import{S as Ae}from"./SdkTabs-86785e52.js";function be(n,l,s){const a=n.slice();return a[5]=l[s],a}function _e(n,l,s){const a=n.slice();return a[5]=l[s],a}function he(n,l){let s,a=l[5].code+"",_,u,i,p;function f(){return l[4](l[5])}return{key:n,first:null,c(){s=r("button"),_=v(a),u=k(),b(s,"class","tab-item"),H(s,"active",l[1]===l[5].code),this.first=s},m(C,$){d(C,s,$),o(s,_),o(s,u),i||(p=Be(s,"click",f),i=!0)},p(C,$){l=C,$&4&&a!==(a=l[5].code+"")&&j(_,a),$&6&&H(s,"active",l[1]===l[5].code)},d(C){C&&m(s),i=!1,p()}}}function ke(n,l){let s,a,_,u;return a=new qe({props:{content:l[5].body}}),{key:n,first:null,c(){s=r("div"),ge(a.$$.fragment),_=k(),b(s,"class","tab-item"),H(s,"active",l[1]===l[5].code),this.first=s},m(i,p){d(i,s,p),ve(a,s,null),o(s,_),u=!0},p(i,p){l=i;const f={};p&4&&(f.content=l[5].body),a.$set(f),(!u||p&6)&&H(s,"active",l[1]===l[5].code)},i(i){u||(x(a.$$.fragment,i),u=!0)},o(i){ee(a.$$.fragment,i),u=!1},d(i){i&&m(s),Ce(a)}}}function Ue(n){var de,me;let l,s,a=n[0].name+"",_,u,i,p,f,C,$,D=n[0].name+"",F,te,I,P,L,R,Q,S,N,le,K,T,se,z,M=n[0].name+"",G,ae,J,y,V,E,X,B,Z,w,q,g=[],ne=new Map,oe,A,h=[],ie=new Map,O;P=new Ae({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${n[3]}');

        ...

        await pb.collection('${(de=n[0])==null?void 0:de.name}').confirmEmailChange(
            'TOKEN',
            'YOUR_PASSWORD',
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${n[3]}');

        ...

        await pb.collection('${(me=n[0])==null?void 0:me.name}').confirmEmailChange(
          'TOKEN',
          'YOUR_PASSWORD',
        );
    `}});let W=Y(n[2]);const ce=e=>e[5].code;for(let e=0;e<W.length;e+=1){let t=_e(n,W,e),c=ce(t);ne.set(c,g[e]=he(c,t))}let U=Y(n[2]);const re=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=be(n,U,e),c=re(t);ie.set(c,h[e]=ke(c,t))}return{c(){l=r("h3"),s=v("Confirm email change ("),_=v(a),u=v(")"),i=k(),p=r("div"),f=r("p"),C=v("Confirms "),$=r("strong"),F=v(D),te=v(" email change request."),I=k(),ge(P.$$.fragment),L=k(),R=r("h6"),R.textContent="API details",Q=k(),S=r("div"),N=r("strong"),N.textContent="POST",le=k(),K=r("div"),T=r("p"),se=v("/api/collections/"),z=r("strong"),G=v(M),ae=v("/confirm-email-change"),J=k(),y=r("div"),y.textContent="Body Parameters",V=k(),E=r("table"),E.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the change email request email.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The account password to confirm the email change.</td></tr></tbody>',X=k(),B=r("div"),B.textContent="Responses",Z=k(),w=r("div"),q=r("div");for(let e=0;e<g.length;e+=1)g[e].c();oe=k(),A=r("div");for(let e=0;e<h.length;e+=1)h[e].c();b(l,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(R,"class","m-b-xs"),b(N,"class","label label-primary"),b(K,"class","content"),b(S,"class","alert alert-success"),b(y,"class","section-title"),b(E,"class","table-compact table-border m-b-base"),b(B,"class","section-title"),b(q,"class","tabs-header compact combined left"),b(A,"class","tabs-content"),b(w,"class","tabs")},m(e,t){d(e,l,t),o(l,s),o(l,_),o(l,u),d(e,i,t),d(e,p,t),o(p,f),o(f,C),o(f,$),o($,F),o(f,te),d(e,I,t),ve(P,e,t),d(e,L,t),d(e,R,t),d(e,Q,t),d(e,S,t),o(S,N),o(S,le),o(S,K),o(K,T),o(T,se),o(T,z),o(z,G),o(T,ae),d(e,J,t),d(e,y,t),d(e,V,t),d(e,E,t),d(e,X,t),d(e,B,t),d(e,Z,t),d(e,w,t),o(w,q);for(let c=0;c<g.length;c+=1)g[c]&&g[c].m(q,null);o(w,oe),o(w,A);for(let c=0;c<h.length;c+=1)h[c]&&h[c].m(A,null);O=!0},p(e,[t]){var pe,fe;(!O||t&1)&&a!==(a=e[0].name+"")&&j(_,a),(!O||t&1)&&D!==(D=e[0].name+"")&&j(F,D);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(pe=e[0])==null?void 0:pe.name}').confirmEmailChange(
            'TOKEN',
            'YOUR_PASSWORD',
        );
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(fe=e[0])==null?void 0:fe.name}').confirmEmailChange(
          'TOKEN',
          'YOUR_PASSWORD',
        );
    `),P.$set(c),(!O||t&1)&&M!==(M=e[0].name+"")&&j(G,M),t&6&&(W=Y(e[2]),g=ue(g,t,ce,1,e,W,ne,q,we,he,null,_e)),t&6&&(U=Y(e[2]),Oe(),h=ue(h,t,re,1,e,U,ie,A,Re,ke,null,be),Te())},i(e){if(!O){x(P.$$.fragment,e);for(let t=0;t<U.length;t+=1)x(h[t]);O=!0}},o(e){ee(P.$$.fragment,e);for(let t=0;t<h.length;t+=1)ee(h[t]);O=!1},d(e){e&&(m(l),m(i),m(p),m(I),m(L),m(R),m(Q),m(S),m(J),m(y),m(V),m(E),m(X),m(B),m(Z),m(w)),Ce(P,e);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<h.length;t+=1)h[t].d()}}}function De(n,l,s){let a,{collection:_}=l,u=204,i=[];const p=f=>s(1,u=f.code);return n.$$set=f=>{"collection"in f&&s(0,_=f.collection)},s(3,a=ye.getApiExampleUrl(Ee.baseUrl)),s(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,u,i,a,p]}class Me extends $e{constructor(l){super(),Pe(this,l,De,Ue,Se,{collection:0})}}export{Me as default};
