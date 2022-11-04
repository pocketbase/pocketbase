import{S as Pe,i as Te,s as Be,e as c,w as v,b as h,c as Ce,f,g as r,h as n,m as Ee,x as D,P as ve,Q as Se,k as Re,R as Me,n as Ae,t as x,a as ee,o as m,d as ye,L as Ve,C as ze,p as He,r as I,u as Le,O as Oe}from"./index.e29b34a9.js";import{S as Ue}from"./SdkTabs.46703f8b.js";function we(o,l,s){const a=o.slice();return a[5]=l[s],a}function ge(o,l,s){const a=o.slice();return a[5]=l[s],a}function $e(o,l){let s,a=l[5].code+"",_,b,i,p;function u(){return l[4](l[5])}return{key:o,first:null,c(){s=c("button"),_=v(a),b=h(),f(s,"class","tab-item"),I(s,"active",l[1]===l[5].code),this.first=s},m($,q){r($,s,q),n(s,_),n(s,b),i||(p=Le(s,"click",u),i=!0)},p($,q){l=$,q&4&&a!==(a=l[5].code+"")&&D(_,a),q&6&&I(s,"active",l[1]===l[5].code)},d($){$&&m(s),i=!1,p()}}}function qe(o,l){let s,a,_,b;return a=new Oe({props:{content:l[5].body}}),{key:o,first:null,c(){s=c("div"),Ce(a.$$.fragment),_=h(),f(s,"class","tab-item"),I(s,"active",l[1]===l[5].code),this.first=s},m(i,p){r(i,s,p),Ee(a,s,null),n(s,_),b=!0},p(i,p){l=i;const u={};p&4&&(u.content=l[5].body),a.$set(u),(!b||p&6)&&I(s,"active",l[1]===l[5].code)},i(i){b||(x(a.$$.fragment,i),b=!0)},o(i){ee(a.$$.fragment,i),b=!1},d(i){i&&m(s),ye(a)}}}function je(o){var de,pe,ue,fe;let l,s,a=o[0].name+"",_,b,i,p,u,$,q,z=o[0].name+"",F,te,K,C,N,P,Q,w,H,le,L,T,se,G,O=o[0].name+"",J,ae,oe,U,W,B,X,S,Y,R,Z,E,M,g=[],ne=new Map,ie,A,k=[],ce=new Map,y;C=new Ue({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(de=o[0])==null?void 0:de.name}').authViaEmail('test@example.com', '123456');

        await pb.collection('${(pe=o[0])==null?void 0:pe.name}').requestEmailChange('new@example.com');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(ue=o[0])==null?void 0:ue.name}').authViaEmail('test@example.com', '123456');

        await pb.collection('${(fe=o[0])==null?void 0:fe.name}').requestEmailChange('new@example.com');
    `}});let j=o[2];const re=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=ge(o,j,e),d=re(t);ne.set(d,g[e]=$e(d,t))}let V=o[2];const me=e=>e[5].code;for(let e=0;e<V.length;e+=1){let t=we(o,V,e),d=me(t);ce.set(d,k[e]=qe(d,t))}return{c(){l=c("h3"),s=v("Request email change ("),_=v(a),b=v(")"),i=h(),p=c("div"),u=c("p"),$=v("Sends "),q=c("strong"),F=v(z),te=v(" email change request."),K=h(),Ce(C.$$.fragment),N=h(),P=c("h6"),P.textContent="API details",Q=h(),w=c("div"),H=c("strong"),H.textContent="POST",le=h(),L=c("div"),T=c("p"),se=v("/api/collections/"),G=c("strong"),J=v(O),ae=v("/confirm-email-change"),oe=h(),U=c("p"),U.innerHTML="Requires record <code>Authorization:TOKEN</code> header",W=h(),B=c("div"),B.textContent="Body Parameters",X=h(),S=c("table"),S.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>newEmail</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The new email address to send the change email request.</td></tr></tbody>`,Y=h(),R=c("div"),R.textContent="Responses",Z=h(),E=c("div"),M=c("div");for(let e=0;e<g.length;e+=1)g[e].c();ie=h(),A=c("div");for(let e=0;e<k.length;e+=1)k[e].c();f(l,"class","m-b-sm"),f(p,"class","content txt-lg m-b-sm"),f(P,"class","m-b-xs"),f(H,"class","label label-primary"),f(L,"class","content"),f(U,"class","txt-hint txt-sm txt-right"),f(w,"class","alert alert-success"),f(B,"class","section-title"),f(S,"class","table-compact table-border m-b-base"),f(R,"class","section-title"),f(M,"class","tabs-header compact left"),f(A,"class","tabs-content"),f(E,"class","tabs")},m(e,t){r(e,l,t),n(l,s),n(l,_),n(l,b),r(e,i,t),r(e,p,t),n(p,u),n(u,$),n(u,q),n(q,F),n(u,te),r(e,K,t),Ee(C,e,t),r(e,N,t),r(e,P,t),r(e,Q,t),r(e,w,t),n(w,H),n(w,le),n(w,L),n(L,T),n(T,se),n(T,G),n(G,J),n(T,ae),n(w,oe),n(w,U),r(e,W,t),r(e,B,t),r(e,X,t),r(e,S,t),r(e,Y,t),r(e,R,t),r(e,Z,t),r(e,E,t),n(E,M);for(let d=0;d<g.length;d+=1)g[d].m(M,null);n(E,ie),n(E,A);for(let d=0;d<k.length;d+=1)k[d].m(A,null);y=!0},p(e,[t]){var be,_e,he,ke;(!y||t&1)&&a!==(a=e[0].name+"")&&D(_,a),(!y||t&1)&&z!==(z=e[0].name+"")&&D(F,z);const d={};t&9&&(d.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(be=e[0])==null?void 0:be.name}').authViaEmail('test@example.com', '123456');

        await pb.collection('${(_e=e[0])==null?void 0:_e.name}').requestEmailChange('new@example.com');
    `),t&9&&(d.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(he=e[0])==null?void 0:he.name}').authViaEmail('test@example.com', '123456');

        await pb.collection('${(ke=e[0])==null?void 0:ke.name}').requestEmailChange('new@example.com');
    `),C.$set(d),(!y||t&1)&&O!==(O=e[0].name+"")&&D(J,O),t&6&&(j=e[2],g=ve(g,t,re,1,e,j,ne,M,Se,$e,null,ge)),t&6&&(V=e[2],Re(),k=ve(k,t,me,1,e,V,ce,A,Me,qe,null,we),Ae())},i(e){if(!y){x(C.$$.fragment,e);for(let t=0;t<V.length;t+=1)x(k[t]);y=!0}},o(e){ee(C.$$.fragment,e);for(let t=0;t<k.length;t+=1)ee(k[t]);y=!1},d(e){e&&m(l),e&&m(i),e&&m(p),e&&m(K),ye(C,e),e&&m(N),e&&m(P),e&&m(Q),e&&m(w),e&&m(W),e&&m(B),e&&m(X),e&&m(S),e&&m(Y),e&&m(R),e&&m(Z),e&&m(E);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function De(o,l,s){let a,{collection:_=new Ve}=l,b=204,i=[];const p=u=>s(1,b=u.code);return o.$$set=u=>{"collection"in u&&s(0,_=u.collection)},s(3,a=ze.getApiExampleUrl(He.baseUrl)),s(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,b,i,a,p]}class Ke extends Pe{constructor(l){super(),Te(this,l,De,je,Be,{collection:0})}}export{Ke as default};
