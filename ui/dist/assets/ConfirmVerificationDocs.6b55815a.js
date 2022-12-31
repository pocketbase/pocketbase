import{S as we,i as Ce,s as Pe,e as c,w as h,b as v,c as ve,f as b,g as r,h as n,m as he,x as D,O as de,P as Te,k as ge,Q as ye,n as Be,t as Z,a as x,o as f,d as $e,R as qe,C as Oe,p as Se,r as H,u as Ee,N as Ne}from"./index.89a3f554.js";import{S as Ve}from"./SdkTabs.0a6ad1c9.js";function ue(i,l,s){const o=i.slice();return o[5]=l[s],o}function be(i,l,s){const o=i.slice();return o[5]=l[s],o}function _e(i,l){let s,o=l[5].code+"",_,u,a,p;function d(){return l[4](l[5])}return{key:i,first:null,c(){s=c("button"),_=h(o),u=v(),b(s,"class","tab-item"),H(s,"active",l[1]===l[5].code),this.first=s},m(w,C){r(w,s,C),n(s,_),n(s,u),a||(p=Ee(s,"click",d),a=!0)},p(w,C){l=w,C&4&&o!==(o=l[5].code+"")&&D(_,o),C&6&&H(s,"active",l[1]===l[5].code)},d(w){w&&f(s),a=!1,p()}}}function ke(i,l){let s,o,_,u;return o=new Ne({props:{content:l[5].body}}),{key:i,first:null,c(){s=c("div"),ve(o.$$.fragment),_=v(),b(s,"class","tab-item"),H(s,"active",l[1]===l[5].code),this.first=s},m(a,p){r(a,s,p),he(o,s,null),n(s,_),u=!0},p(a,p){l=a;const d={};p&4&&(d.content=l[5].body),o.$set(d),(!u||p&6)&&H(s,"active",l[1]===l[5].code)},i(a){u||(Z(o.$$.fragment,a),u=!0)},o(a){x(o.$$.fragment,a),u=!1},d(a){a&&f(s),$e(o)}}}function Ke(i){var re,fe;let l,s,o=i[0].name+"",_,u,a,p,d,w,C,M=i[0].name+"",I,ee,F,P,L,B,Q,T,A,te,R,q,le,z,U=i[0].name+"",G,se,J,O,W,S,X,E,Y,g,N,$=[],oe=new Map,ie,V,k=[],ne=new Map,y;P=new Ve({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[3]}');

        ...

        await pb.collection('${(re=i[0])==null?void 0:re.name}').confirmVerification('TOKEN');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[3]}');

        ...

        await pb.collection('${(fe=i[0])==null?void 0:fe.name}').confirmVerification('TOKEN');
    `}});let j=i[2];const ae=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=be(i,j,e),m=ae(t);oe.set(m,$[e]=_e(m,t))}let K=i[2];const ce=e=>e[5].code;for(let e=0;e<K.length;e+=1){let t=ue(i,K,e),m=ce(t);ne.set(m,k[e]=ke(m,t))}return{c(){l=c("h3"),s=h("Confirm verification ("),_=h(o),u=h(")"),a=v(),p=c("div"),d=c("p"),w=h("Confirms "),C=c("strong"),I=h(M),ee=h(" account verification request."),F=v(),ve(P.$$.fragment),L=v(),B=c("h6"),B.textContent="API details",Q=v(),T=c("div"),A=c("strong"),A.textContent="POST",te=v(),R=c("div"),q=c("p"),le=h("/api/collections/"),z=c("strong"),G=h(U),se=h("/confirm-verification"),J=v(),O=c("div"),O.textContent="Body Parameters",W=v(),S=c("table"),S.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>token</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The token from the verification request email.</td></tr></tbody>`,X=v(),E=c("div"),E.textContent="Responses",Y=v(),g=c("div"),N=c("div");for(let e=0;e<$.length;e+=1)$[e].c();ie=v(),V=c("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(A,"class","label label-primary"),b(R,"class","content"),b(T,"class","alert alert-success"),b(O,"class","section-title"),b(S,"class","table-compact table-border m-b-base"),b(E,"class","section-title"),b(N,"class","tabs-header compact left"),b(V,"class","tabs-content"),b(g,"class","tabs")},m(e,t){r(e,l,t),n(l,s),n(l,_),n(l,u),r(e,a,t),r(e,p,t),n(p,d),n(d,w),n(d,C),n(C,I),n(d,ee),r(e,F,t),he(P,e,t),r(e,L,t),r(e,B,t),r(e,Q,t),r(e,T,t),n(T,A),n(T,te),n(T,R),n(R,q),n(q,le),n(q,z),n(z,G),n(q,se),r(e,J,t),r(e,O,t),r(e,W,t),r(e,S,t),r(e,X,t),r(e,E,t),r(e,Y,t),r(e,g,t),n(g,N);for(let m=0;m<$.length;m+=1)$[m].m(N,null);n(g,ie),n(g,V);for(let m=0;m<k.length;m+=1)k[m].m(V,null);y=!0},p(e,[t]){var me,pe;(!y||t&1)&&o!==(o=e[0].name+"")&&D(_,o),(!y||t&1)&&M!==(M=e[0].name+"")&&D(I,M);const m={};t&9&&(m.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').confirmVerification('TOKEN');
    `),t&9&&(m.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(pe=e[0])==null?void 0:pe.name}').confirmVerification('TOKEN');
    `),P.$set(m),(!y||t&1)&&U!==(U=e[0].name+"")&&D(G,U),t&6&&(j=e[2],$=de($,t,ae,1,e,j,oe,N,Te,_e,null,be)),t&6&&(K=e[2],ge(),k=de(k,t,ce,1,e,K,ne,V,ye,ke,null,ue),Be())},i(e){if(!y){Z(P.$$.fragment,e);for(let t=0;t<K.length;t+=1)Z(k[t]);y=!0}},o(e){x(P.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);y=!1},d(e){e&&f(l),e&&f(a),e&&f(p),e&&f(F),$e(P,e),e&&f(L),e&&f(B),e&&f(Q),e&&f(T),e&&f(J),e&&f(O),e&&f(W),e&&f(S),e&&f(X),e&&f(E),e&&f(Y),e&&f(g);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function Me(i,l,s){let o,{collection:_=new qe}=l,u=204,a=[];const p=d=>s(1,u=d.code);return i.$$set=d=>{"collection"in d&&s(0,_=d.collection)},s(3,o=Oe.getApiExampleUrl(Se.baseUrl)),s(2,a=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,u,a,o,p]}class Ue extends we{constructor(l){super(),Ce(this,l,Me,Ke,Pe,{collection:0})}}export{Ue as default};
