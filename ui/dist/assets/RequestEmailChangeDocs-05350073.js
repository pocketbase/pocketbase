import{S as Te,i as Ee,s as Be,e as r,w,b as k,c as Pe,f,g as m,h as n,m as Ce,x as N,N as ve,P as Se,k as Me,Q as Re,n as Ae,t as x,a as ee,o as d,d as ye,T as We,C as ze,p as He,r as O,u as Ue,M as je}from"./index-a65ca895.js";import{S as De}from"./SdkTabs-ad912c8f.js";function we(o,l,s){const a=o.slice();return a[5]=l[s],a}function ge(o,l,s){const a=o.slice();return a[5]=l[s],a}function $e(o,l){let s,a=l[5].code+"",_,b,i,p;function u(){return l[4](l[5])}return{key:o,first:null,c(){s=r("button"),_=w(a),b=k(),f(s,"class","tab-item"),O(s,"active",l[1]===l[5].code),this.first=s},m($,q){m($,s,q),n(s,_),n(s,b),i||(p=Ue(s,"click",u),i=!0)},p($,q){l=$,q&4&&a!==(a=l[5].code+"")&&N(_,a),q&6&&O(s,"active",l[1]===l[5].code)},d($){$&&d(s),i=!1,p()}}}function qe(o,l){let s,a,_,b;return a=new je({props:{content:l[5].body}}),{key:o,first:null,c(){s=r("div"),Pe(a.$$.fragment),_=k(),f(s,"class","tab-item"),O(s,"active",l[1]===l[5].code),this.first=s},m(i,p){m(i,s,p),Ce(a,s,null),n(s,_),b=!0},p(i,p){l=i;const u={};p&4&&(u.content=l[5].body),a.$set(u),(!b||p&6)&&O(s,"active",l[1]===l[5].code)},i(i){b||(x(a.$$.fragment,i),b=!0)},o(i){ee(a.$$.fragment,i),b=!1},d(i){i&&d(s),ye(a)}}}function Le(o){var de,pe,ue,fe;let l,s,a=o[0].name+"",_,b,i,p,u,$,q,z=o[0].name+"",F,te,I,P,K,T,Q,g,H,le,U,E,se,G,j=o[0].name+"",J,ae,oe,D,V,B,X,S,Y,M,Z,C,R,v=[],ne=new Map,ie,A,h=[],ce=new Map,y;P=new De({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(de=o[0])==null?void 0:de.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(pe=o[0])==null?void 0:pe.name}').requestEmailChange('new@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(ue=o[0])==null?void 0:ue.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(fe=o[0])==null?void 0:fe.name}').requestEmailChange('new@example.com');
    `}});let L=o[2];const re=e=>e[5].code;for(let e=0;e<L.length;e+=1){let t=ge(o,L,e),c=re(t);ne.set(c,v[e]=$e(c,t))}let W=o[2];const me=e=>e[5].code;for(let e=0;e<W.length;e+=1){let t=we(o,W,e),c=me(t);ce.set(c,h[e]=qe(c,t))}return{c(){l=r("h3"),s=w("Request email change ("),_=w(a),b=w(")"),i=k(),p=r("div"),u=r("p"),$=w("Sends "),q=r("strong"),F=w(z),te=w(" email change request."),I=k(),Pe(P.$$.fragment),K=k(),T=r("h6"),T.textContent="API details",Q=k(),g=r("div"),H=r("strong"),H.textContent="POST",le=k(),U=r("div"),E=r("p"),se=w("/api/collections/"),G=r("strong"),J=w(j),ae=w("/request-email-change"),oe=k(),D=r("p"),D.innerHTML="Requires record <code>Authorization:TOKEN</code> header",V=k(),B=r("div"),B.textContent="Body Parameters",X=k(),S=r("table"),S.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>newEmail</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The new email address to send the change email request.</td></tr></tbody>`,Y=k(),M=r("div"),M.textContent="Responses",Z=k(),C=r("div"),R=r("div");for(let e=0;e<v.length;e+=1)v[e].c();ie=k(),A=r("div");for(let e=0;e<h.length;e+=1)h[e].c();f(l,"class","m-b-sm"),f(p,"class","content txt-lg m-b-sm"),f(T,"class","m-b-xs"),f(H,"class","label label-primary"),f(U,"class","content"),f(D,"class","txt-hint txt-sm txt-right"),f(g,"class","alert alert-success"),f(B,"class","section-title"),f(S,"class","table-compact table-border m-b-base"),f(M,"class","section-title"),f(R,"class","tabs-header compact left"),f(A,"class","tabs-content"),f(C,"class","tabs")},m(e,t){m(e,l,t),n(l,s),n(l,_),n(l,b),m(e,i,t),m(e,p,t),n(p,u),n(u,$),n(u,q),n(q,F),n(u,te),m(e,I,t),Ce(P,e,t),m(e,K,t),m(e,T,t),m(e,Q,t),m(e,g,t),n(g,H),n(g,le),n(g,U),n(U,E),n(E,se),n(E,G),n(G,J),n(E,ae),n(g,oe),n(g,D),m(e,V,t),m(e,B,t),m(e,X,t),m(e,S,t),m(e,Y,t),m(e,M,t),m(e,Z,t),m(e,C,t),n(C,R);for(let c=0;c<v.length;c+=1)v[c]&&v[c].m(R,null);n(C,ie),n(C,A);for(let c=0;c<h.length;c+=1)h[c]&&h[c].m(A,null);y=!0},p(e,[t]){var be,_e,he,ke;(!y||t&1)&&a!==(a=e[0].name+"")&&N(_,a),(!y||t&1)&&z!==(z=e[0].name+"")&&N(F,z);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(_e=e[0])==null?void 0:_e.name}').requestEmailChange('new@example.com');
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(he=e[0])==null?void 0:he.name}').authWithPassword('test@example.com', '1234567890');

        await pb.collection('${(ke=e[0])==null?void 0:ke.name}').requestEmailChange('new@example.com');
    `),P.$set(c),(!y||t&1)&&j!==(j=e[0].name+"")&&N(J,j),t&6&&(L=e[2],v=ve(v,t,re,1,e,L,ne,R,Se,$e,null,ge)),t&6&&(W=e[2],Me(),h=ve(h,t,me,1,e,W,ce,A,Re,qe,null,we),Ae())},i(e){if(!y){x(P.$$.fragment,e);for(let t=0;t<W.length;t+=1)x(h[t]);y=!0}},o(e){ee(P.$$.fragment,e);for(let t=0;t<h.length;t+=1)ee(h[t]);y=!1},d(e){e&&d(l),e&&d(i),e&&d(p),e&&d(I),ye(P,e),e&&d(K),e&&d(T),e&&d(Q),e&&d(g),e&&d(V),e&&d(B),e&&d(X),e&&d(S),e&&d(Y),e&&d(M),e&&d(Z),e&&d(C);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<h.length;t+=1)h[t].d()}}}function Ne(o,l,s){let a,{collection:_=new We}=l,b=204,i=[];const p=u=>s(1,b=u.code);return o.$$set=u=>{"collection"in u&&s(0,_=u.collection)},s(3,a=ze.getApiExampleUrl(He.baseUrl)),s(2,i=[{code:204,body:"null"},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "newEmail": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `},{code:401,body:`
                {
                  "code": 401,
                  "message": "The request requires valid record authorization token to be set.",
                  "data": {}
                }
            `},{code:403,body:`
                {
                  "code": 403,
                  "message": "The authorized record model is not allowed to perform this action.",
                  "data": {}
                }
            `}]),[_,b,i,a,p]}class Ie extends Te{constructor(l){super(),Ee(this,l,Ne,Le,Be,{collection:0})}}export{Ie as default};
