import{S as Se,i as Te,s as Be,O as D,e as r,v as g,b as k,c as ye,f as h,g as f,h as n,m as Ce,w as H,P as ke,Q as qe,k as Re,R as Oe,n as Ae,t as x,a as ee,o as d,d as Pe,C as Ee,A as Ne,q as F,r as Ve,N as Ke}from"./index-D0DO79Dq.js";import{S as Me}from"./SdkTabs-DC6EUYpr.js";function ve(o,l,s){const a=o.slice();return a[5]=l[s],a}function ge(o,l,s){const a=o.slice();return a[5]=l[s],a}function we(o,l){let s,a=l[5].code+"",b,m,i,p;function u(){return l[4](l[5])}return{key:o,first:null,c(){s=r("button"),b=g(a),m=k(),h(s,"class","tab-item"),F(s,"active",l[1]===l[5].code),this.first=s},m(w,$){f(w,s,$),n(s,b),n(s,m),i||(p=Ve(s,"click",u),i=!0)},p(w,$){l=w,$&4&&a!==(a=l[5].code+"")&&H(b,a),$&6&&F(s,"active",l[1]===l[5].code)},d(w){w&&d(s),i=!1,p()}}}function $e(o,l){let s,a,b,m;return a=new Ke({props:{content:l[5].body}}),{key:o,first:null,c(){s=r("div"),ye(a.$$.fragment),b=k(),h(s,"class","tab-item"),F(s,"active",l[1]===l[5].code),this.first=s},m(i,p){f(i,s,p),Ce(a,s,null),n(s,b),m=!0},p(i,p){l=i;const u={};p&4&&(u.content=l[5].body),a.$set(u),(!m||p&6)&&F(s,"active",l[1]===l[5].code)},i(i){m||(x(a.$$.fragment,i),m=!0)},o(i){ee(a.$$.fragment,i),m=!1},d(i){i&&d(s),Pe(a)}}}function Ue(o){var fe,de,pe,ue;let l,s,a=o[0].name+"",b,m,i,p,u,w,$,V=o[0].name+"",I,te,L,y,Q,T,z,C,K,le,M,B,se,G,U=o[0].name+"",J,ae,W,q,X,R,Y,O,Z,P,A,v=[],oe=new Map,ne,E,_=[],ie=new Map,S;y=new Me({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(fe=o[0])==null?void 0:fe.name}').confirmVerification('TOKEN');

        // optionally refresh the previous authStore state with the latest record changes
        await pb.collection('${(de=o[0])==null?void 0:de.name}').authRefresh();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(pe=o[0])==null?void 0:pe.name}').confirmVerification('TOKEN');

        // optionally refresh the previous authStore state with the latest record changes
        await pb.collection('${(ue=o[0])==null?void 0:ue.name}').authRefresh();
    `}});let j=D(o[2]);const ce=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=ge(o,j,e),c=ce(t);oe.set(c,v[e]=we(c,t))}let N=D(o[2]);const re=e=>e[5].code;for(let e=0;e<N.length;e+=1){let t=ve(o,N,e),c=re(t);ie.set(c,_[e]=$e(c,t))}return{c(){l=r("h3"),s=g("Confirm verification ("),b=g(a),m=g(")"),i=k(),p=r("div"),u=r("p"),w=g("Confirms "),$=r("strong"),I=g(V),te=g(" account verification request."),L=k(),ye(y.$$.fragment),Q=k(),T=r("h6"),T.textContent="API details",z=k(),C=r("div"),K=r("strong"),K.textContent="POST",le=k(),M=r("div"),B=r("p"),se=g("/api/collections/"),G=r("strong"),J=g(U),ae=g("/confirm-verification"),W=k(),q=r("div"),q.textContent="Body Parameters",X=k(),R=r("table"),R.innerHTML='<thead><tr><th>Param</th> <th>Type</th> <th width="50%">Description</th></tr></thead> <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> <span>token</span></div></td> <td><span class="label">String</span></td> <td>The token from the verification request email.</td></tr></tbody>',Y=k(),O=r("div"),O.textContent="Responses",Z=k(),P=r("div"),A=r("div");for(let e=0;e<v.length;e+=1)v[e].c();ne=k(),E=r("div");for(let e=0;e<_.length;e+=1)_[e].c();h(l,"class","m-b-sm"),h(p,"class","content txt-lg m-b-sm"),h(T,"class","m-b-xs"),h(K,"class","label label-primary"),h(M,"class","content"),h(C,"class","alert alert-success"),h(q,"class","section-title"),h(R,"class","table-compact table-border m-b-base"),h(O,"class","section-title"),h(A,"class","tabs-header compact combined left"),h(E,"class","tabs-content"),h(P,"class","tabs")},m(e,t){f(e,l,t),n(l,s),n(l,b),n(l,m),f(e,i,t),f(e,p,t),n(p,u),n(u,w),n(u,$),n($,I),n(u,te),f(e,L,t),Ce(y,e,t),f(e,Q,t),f(e,T,t),f(e,z,t),f(e,C,t),n(C,K),n(C,le),n(C,M),n(M,B),n(B,se),n(B,G),n(G,J),n(B,ae),f(e,W,t),f(e,q,t),f(e,X,t),f(e,R,t),f(e,Y,t),f(e,O,t),f(e,Z,t),f(e,P,t),n(P,A);for(let c=0;c<v.length;c+=1)v[c]&&v[c].m(A,null);n(P,ne),n(P,E);for(let c=0;c<_.length;c+=1)_[c]&&_[c].m(E,null);S=!0},p(e,[t]){var me,he,be,_e;(!S||t&1)&&a!==(a=e[0].name+"")&&H(b,a),(!S||t&1)&&V!==(V=e[0].name+"")&&H(I,V);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').confirmVerification('TOKEN');

        // optionally refresh the previous authStore state with the latest record changes
        await pb.collection('${(he=e[0])==null?void 0:he.name}').authRefresh();
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').confirmVerification('TOKEN');

        // optionally refresh the previous authStore state with the latest record changes
        await pb.collection('${(_e=e[0])==null?void 0:_e.name}').authRefresh();
    `),y.$set(c),(!S||t&1)&&U!==(U=e[0].name+"")&&H(J,U),t&6&&(j=D(e[2]),v=ke(v,t,ce,1,e,j,oe,A,qe,we,null,ge)),t&6&&(N=D(e[2]),Re(),_=ke(_,t,re,1,e,N,ie,E,Oe,$e,null,ve),Ae())},i(e){if(!S){x(y.$$.fragment,e);for(let t=0;t<N.length;t+=1)x(_[t]);S=!0}},o(e){ee(y.$$.fragment,e);for(let t=0;t<_.length;t+=1)ee(_[t]);S=!1},d(e){e&&(d(l),d(i),d(p),d(L),d(Q),d(T),d(z),d(C),d(W),d(q),d(X),d(R),d(Y),d(O),d(Z),d(P)),Pe(y,e);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<_.length;t+=1)_[t].d()}}}function je(o,l,s){let a,{collection:b}=l,m=204,i=[];const p=u=>s(1,m=u.code);return o.$$set=u=>{"collection"in u&&s(0,b=u.collection)},s(3,a=Ee.getApiExampleUrl(Ne.baseUrl)),s(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[b,m,i,a,p]}class Fe extends Se{constructor(l){super(),Te(this,l,je,Ue,Be,{collection:0})}}export{Fe as default};
