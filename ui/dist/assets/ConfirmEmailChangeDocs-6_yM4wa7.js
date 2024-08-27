import{S as Pe,i as Se,s as Oe,O as Y,e as r,v,b as k,c as Ce,f as b,g as d,h as n,m as $e,w as j,P as _e,Q as ye,k as Re,R as Te,n as Ae,t as ee,a as te,o as m,d as we,C as Ee,A as qe,q as H,r as Be,N as Ue}from"./index-D0DO79Dq.js";import{S as De}from"./SdkTabs-DC6EUYpr.js";function he(o,l,s){const a=o.slice();return a[5]=l[s],a}function ke(o,l,s){const a=o.slice();return a[5]=l[s],a}function ge(o,l){let s,a=l[5].code+"",_,u,i,p;function f(){return l[4](l[5])}return{key:o,first:null,c(){s=r("button"),_=v(a),u=k(),b(s,"class","tab-item"),H(s,"active",l[1]===l[5].code),this.first=s},m(C,$){d(C,s,$),n(s,_),n(s,u),i||(p=Be(s,"click",f),i=!0)},p(C,$){l=C,$&4&&a!==(a=l[5].code+"")&&j(_,a),$&6&&H(s,"active",l[1]===l[5].code)},d(C){C&&m(s),i=!1,p()}}}function ve(o,l){let s,a,_,u;return a=new Ue({props:{content:l[5].body}}),{key:o,first:null,c(){s=r("div"),Ce(a.$$.fragment),_=k(),b(s,"class","tab-item"),H(s,"active",l[1]===l[5].code),this.first=s},m(i,p){d(i,s,p),$e(a,s,null),n(s,_),u=!0},p(i,p){l=i;const f={};p&4&&(f.content=l[5].body),a.$set(f),(!u||p&6)&&H(s,"active",l[1]===l[5].code)},i(i){u||(ee(a.$$.fragment,i),u=!0)},o(i){te(a.$$.fragment,i),u=!1},d(i){i&&m(s),we(a)}}}function Ne(o){var pe,fe;let l,s,a=o[0].name+"",_,u,i,p,f,C,$,D=o[0].name+"",F,le,se,I,L,w,Q,y,z,P,N,ae,K,R,ne,G,M=o[0].name+"",J,oe,V,T,X,A,Z,E,x,S,q,g=[],ie=new Map,ce,B,h=[],re=new Map,O;w=new De({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(pe=o[0])==null?void 0:pe.name}').confirmEmailChange(
            'TOKEN',
            'YOUR_PASSWORD',
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(fe=o[0])==null?void 0:fe.name}').confirmEmailChange(
          'TOKEN',
          'YOUR_PASSWORD',
        );
    `}});let W=Y(o[2]);const de=e=>e[5].code;for(let e=0;e<W.length;e+=1){let t=ke(o,W,e),c=de(t);ie.set(c,g[e]=ge(c,t))}let U=Y(o[2]);const me=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=he(o,U,e),c=me(t);re.set(c,h[e]=ve(c,t))}return{c(){l=r("h3"),s=v("Confirm email change ("),_=v(a),u=v(")"),i=k(),p=r("div"),f=r("p"),C=v("Confirms "),$=r("strong"),F=v(D),le=v(" email change request."),se=k(),I=r("p"),I.textContent=`After this request all previously issued tokens for the specific record will be automatically
        invalidated.`,L=k(),Ce(w.$$.fragment),Q=k(),y=r("h6"),y.textContent="API details",z=k(),P=r("div"),N=r("strong"),N.textContent="POST",ae=k(),K=r("div"),R=r("p"),ne=v("/api/collections/"),G=r("strong"),J=v(M),oe=v("/confirm-email-change"),V=k(),T=r("div"),T.textContent="Body Parameters",X=k(),A=r("table"),A.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the change email request email.</td></tr> <tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>password</span></div></td> <td><span class="label">String</span></td> <td>The account password to confirm the email change.</td></tr></tbody>',Z=k(),E=r("div"),E.textContent="Responses",x=k(),S=r("div"),q=r("div");for(let e=0;e<g.length;e+=1)g[e].c();ce=k(),B=r("div");for(let e=0;e<h.length;e+=1)h[e].c();b(l,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(y,"class","m-b-xs"),b(N,"class","label label-primary"),b(K,"class","content"),b(P,"class","alert alert-success"),b(T,"class","section-title"),b(A,"class","table-compact table-border m-b-base"),b(E,"class","section-title"),b(q,"class","tabs-header compact combined left"),b(B,"class","tabs-content"),b(S,"class","tabs")},m(e,t){d(e,l,t),n(l,s),n(l,_),n(l,u),d(e,i,t),d(e,p,t),n(p,f),n(f,C),n(f,$),n($,F),n(f,le),n(p,se),n(p,I),d(e,L,t),$e(w,e,t),d(e,Q,t),d(e,y,t),d(e,z,t),d(e,P,t),n(P,N),n(P,ae),n(P,K),n(K,R),n(R,ne),n(R,G),n(G,J),n(R,oe),d(e,V,t),d(e,T,t),d(e,X,t),d(e,A,t),d(e,Z,t),d(e,E,t),d(e,x,t),d(e,S,t),n(S,q);for(let c=0;c<g.length;c+=1)g[c]&&g[c].m(q,null);n(S,ce),n(S,B);for(let c=0;c<h.length;c+=1)h[c]&&h[c].m(B,null);O=!0},p(e,[t]){var ue,be;(!O||t&1)&&a!==(a=e[0].name+"")&&j(_,a),(!O||t&1)&&D!==(D=e[0].name+"")&&j(F,D);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(ue=e[0])==null?void 0:ue.name}').confirmEmailChange(
            'TOKEN',
            'YOUR_PASSWORD',
        );
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').confirmEmailChange(
          'TOKEN',
          'YOUR_PASSWORD',
        );
    `),w.$set(c),(!O||t&1)&&M!==(M=e[0].name+"")&&j(J,M),t&6&&(W=Y(e[2]),g=_e(g,t,de,1,e,W,ie,q,ye,ge,null,ke)),t&6&&(U=Y(e[2]),Re(),h=_e(h,t,me,1,e,U,re,B,Te,ve,null,he),Ae())},i(e){if(!O){ee(w.$$.fragment,e);for(let t=0;t<U.length;t+=1)ee(h[t]);O=!0}},o(e){te(w.$$.fragment,e);for(let t=0;t<h.length;t+=1)te(h[t]);O=!1},d(e){e&&(m(l),m(i),m(p),m(L),m(Q),m(y),m(z),m(P),m(V),m(T),m(X),m(A),m(Z),m(E),m(x),m(S)),we(w,e);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<h.length;t+=1)h[t].d()}}}function Ke(o,l,s){let a,{collection:_}=l,u=204,i=[];const p=f=>s(1,u=f.code);return o.$$set=f=>{"collection"in f&&s(0,_=f.collection)},s(3,a=Ee.getApiExampleUrl(qe.baseUrl)),s(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,u,i,a,p]}class Ye extends Pe{constructor(l){super(),Se(this,l,Ke,Ne,Oe,{collection:0})}}export{Ye as default};
