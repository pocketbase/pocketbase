import{S as we,i as Ce,s as Pe,e as r,w as $,b as v,c as ve,f as b,g as f,h as n,m as he,x as H,N as de,P as Te,k as ge,Q as ye,n as Be,t as Z,a as x,o as m,d as $e,T as qe,C as Se,p as Ee,r as R,u as Me,M as Ne}from"./index-d0d2ced9.js";import{S as Oe}from"./SdkTabs-3a653059.js";function ue(i,l,s){const o=i.slice();return o[5]=l[s],o}function be(i,l,s){const o=i.slice();return o[5]=l[s],o}function _e(i,l){let s,o=l[5].code+"",_,u,a,p;function d(){return l[4](l[5])}return{key:i,first:null,c(){s=r("button"),_=$(o),u=v(),b(s,"class","tab-item"),R(s,"active",l[1]===l[5].code),this.first=s},m(w,C){f(w,s,C),n(s,_),n(s,u),a||(p=Me(s,"click",d),a=!0)},p(w,C){l=w,C&4&&o!==(o=l[5].code+"")&&H(_,o),C&6&&R(s,"active",l[1]===l[5].code)},d(w){w&&m(s),a=!1,p()}}}function ke(i,l){let s,o,_,u;return o=new Ne({props:{content:l[5].body}}),{key:i,first:null,c(){s=r("div"),ve(o.$$.fragment),_=v(),b(s,"class","tab-item"),R(s,"active",l[1]===l[5].code),this.first=s},m(a,p){f(a,s,p),he(o,s,null),n(s,_),u=!0},p(a,p){l=a;const d={};p&4&&(d.content=l[5].body),o.$set(d),(!u||p&6)&&R(s,"active",l[1]===l[5].code)},i(a){u||(Z(o.$$.fragment,a),u=!0)},o(a){x(o.$$.fragment,a),u=!1},d(a){a&&m(s),$e(o)}}}function Ve(i){var re,fe;let l,s,o=i[0].name+"",_,u,a,p,d,w,C,K=i[0].name+"",F,ee,I,P,L,B,Q,T,A,te,U,q,le,z,j=i[0].name+"",G,se,J,S,W,E,X,M,Y,g,N,h=[],oe=new Map,ie,O,k=[],ne=new Map,y;P=new Oe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[3]}');

        ...

        await pb.collection('${(re=i[0])==null?void 0:re.name}').confirmVerification('TOKEN');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[3]}');

        ...

        await pb.collection('${(fe=i[0])==null?void 0:fe.name}').confirmVerification('TOKEN');
    `}});let D=i[2];const ae=e=>e[5].code;for(let e=0;e<D.length;e+=1){let t=be(i,D,e),c=ae(t);oe.set(c,h[e]=_e(c,t))}let V=i[2];const ce=e=>e[5].code;for(let e=0;e<V.length;e+=1){let t=ue(i,V,e),c=ce(t);ne.set(c,k[e]=ke(c,t))}return{c(){l=r("h3"),s=$("Confirm verification ("),_=$(o),u=$(")"),a=v(),p=r("div"),d=r("p"),w=$("Confirms "),C=r("strong"),F=$(K),ee=$(" account verification request."),I=v(),ve(P.$$.fragment),L=v(),B=r("h6"),B.textContent="API details",Q=v(),T=r("div"),A=r("strong"),A.textContent="POST",te=v(),U=r("div"),q=r("p"),le=$("/api/collections/"),z=r("strong"),G=$(j),se=$("/confirm-verification"),J=v(),S=r("div"),S.textContent="Body Parameters",W=v(),E=r("table"),E.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>token</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The token from the verification request email.</td></tr></tbody>`,X=v(),M=r("div"),M.textContent="Responses",Y=v(),g=r("div"),N=r("div");for(let e=0;e<h.length;e+=1)h[e].c();ie=v(),O=r("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(A,"class","label label-primary"),b(U,"class","content"),b(T,"class","alert alert-success"),b(S,"class","section-title"),b(E,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(N,"class","tabs-header compact left"),b(O,"class","tabs-content"),b(g,"class","tabs")},m(e,t){f(e,l,t),n(l,s),n(l,_),n(l,u),f(e,a,t),f(e,p,t),n(p,d),n(d,w),n(d,C),n(C,F),n(d,ee),f(e,I,t),he(P,e,t),f(e,L,t),f(e,B,t),f(e,Q,t),f(e,T,t),n(T,A),n(T,te),n(T,U),n(U,q),n(q,le),n(q,z),n(z,G),n(q,se),f(e,J,t),f(e,S,t),f(e,W,t),f(e,E,t),f(e,X,t),f(e,M,t),f(e,Y,t),f(e,g,t),n(g,N);for(let c=0;c<h.length;c+=1)h[c]&&h[c].m(N,null);n(g,ie),n(g,O);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(O,null);y=!0},p(e,[t]){var me,pe;(!y||t&1)&&o!==(o=e[0].name+"")&&H(_,o),(!y||t&1)&&K!==(K=e[0].name+"")&&H(F,K);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').confirmVerification('TOKEN');
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(pe=e[0])==null?void 0:pe.name}').confirmVerification('TOKEN');
    `),P.$set(c),(!y||t&1)&&j!==(j=e[0].name+"")&&H(G,j),t&6&&(D=e[2],h=de(h,t,ae,1,e,D,oe,N,Te,_e,null,be)),t&6&&(V=e[2],ge(),k=de(k,t,ce,1,e,V,ne,O,ye,ke,null,ue),Be())},i(e){if(!y){Z(P.$$.fragment,e);for(let t=0;t<V.length;t+=1)Z(k[t]);y=!0}},o(e){x(P.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);y=!1},d(e){e&&m(l),e&&m(a),e&&m(p),e&&m(I),$e(P,e),e&&m(L),e&&m(B),e&&m(Q),e&&m(T),e&&m(J),e&&m(S),e&&m(W),e&&m(E),e&&m(X),e&&m(M),e&&m(Y),e&&m(g);for(let t=0;t<h.length;t+=1)h[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function Ke(i,l,s){let o,{collection:_=new qe}=l,u=204,a=[];const p=d=>s(1,u=d.code);return i.$$set=d=>{"collection"in d&&s(0,_=d.collection)},s(3,o=Se.getApiExampleUrl(Ee.baseUrl)),s(2,a=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,u,a,o,p]}class je extends we{constructor(l){super(),Ce(this,l,Ke,Ve,Pe,{collection:0})}}export{je as default};
