import{S as Ce,i as $e,s as we,e as c,w as v,b as h,c as he,f as b,g as r,h as n,m as ve,x as Y,O as pe,P as Pe,k as Se,Q as Oe,n as Re,t as Z,a as x,o as f,d as ge,R as Te,C as Ee,p as ye,r as j,u as Be,N as qe}from"./index.89a3f554.js";import{S as Ae}from"./SdkTabs.0a6ad1c9.js";function ue(o,l,s){const a=o.slice();return a[5]=l[s],a}function be(o,l,s){const a=o.slice();return a[5]=l[s],a}function _e(o,l){let s,a=l[5].code+"",_,u,i,d;function p(){return l[4](l[5])}return{key:o,first:null,c(){s=c("button"),_=v(a),u=h(),b(s,"class","tab-item"),j(s,"active",l[1]===l[5].code),this.first=s},m(C,$){r(C,s,$),n(s,_),n(s,u),i||(d=Be(s,"click",p),i=!0)},p(C,$){l=C,$&4&&a!==(a=l[5].code+"")&&Y(_,a),$&6&&j(s,"active",l[1]===l[5].code)},d(C){C&&f(s),i=!1,d()}}}function ke(o,l){let s,a,_,u;return a=new qe({props:{content:l[5].body}}),{key:o,first:null,c(){s=c("div"),he(a.$$.fragment),_=h(),b(s,"class","tab-item"),j(s,"active",l[1]===l[5].code),this.first=s},m(i,d){r(i,s,d),ve(a,s,null),n(s,_),u=!0},p(i,d){l=i;const p={};d&4&&(p.content=l[5].body),a.$set(p),(!u||d&6)&&j(s,"active",l[1]===l[5].code)},i(i){u||(Z(a.$$.fragment,i),u=!0)},o(i){x(a.$$.fragment,i),u=!1},d(i){i&&f(s),ge(a)}}}function Ue(o){var re,fe;let l,s,a=o[0].name+"",_,u,i,d,p,C,$,D=o[0].name+"",H,ee,I,w,F,R,L,P,N,te,K,T,le,Q,M=o[0].name+"",z,se,G,E,J,y,V,B,X,S,q,g=[],ae=new Map,oe,A,k=[],ne=new Map,O;w=new Ae({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(re=o[0])==null?void 0:re.name}').confirmEmailChange(
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
    `}});let W=o[2];const ie=e=>e[5].code;for(let e=0;e<W.length;e+=1){let t=be(o,W,e),m=ie(t);ae.set(m,g[e]=_e(m,t))}let U=o[2];const ce=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=ue(o,U,e),m=ce(t);ne.set(m,k[e]=ke(m,t))}return{c(){l=c("h3"),s=v("Confirm email change ("),_=v(a),u=v(")"),i=h(),d=c("div"),p=c("p"),C=v("Confirms "),$=c("strong"),H=v(D),ee=v(" email change request."),I=h(),he(w.$$.fragment),F=h(),R=c("h6"),R.textContent="API details",L=h(),P=c("div"),N=c("strong"),N.textContent="POST",te=h(),K=c("div"),T=c("p"),le=v("/api/collections/"),Q=c("strong"),z=v(M),se=v("/confirm-email-change"),G=h(),E=c("div"),E.textContent="Body Parameters",J=h(),y=c("table"),y.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>token</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The token from the change email request email.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>password</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The account password to confirm the email change.</td></tr></tbody>`,V=h(),B=c("div"),B.textContent="Responses",X=h(),S=c("div"),q=c("div");for(let e=0;e<g.length;e+=1)g[e].c();oe=h(),A=c("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(d,"class","content txt-lg m-b-sm"),b(R,"class","m-b-xs"),b(N,"class","label label-primary"),b(K,"class","content"),b(P,"class","alert alert-success"),b(E,"class","section-title"),b(y,"class","table-compact table-border m-b-base"),b(B,"class","section-title"),b(q,"class","tabs-header compact left"),b(A,"class","tabs-content"),b(S,"class","tabs")},m(e,t){r(e,l,t),n(l,s),n(l,_),n(l,u),r(e,i,t),r(e,d,t),n(d,p),n(p,C),n(p,$),n($,H),n(p,ee),r(e,I,t),ve(w,e,t),r(e,F,t),r(e,R,t),r(e,L,t),r(e,P,t),n(P,N),n(P,te),n(P,K),n(K,T),n(T,le),n(T,Q),n(Q,z),n(T,se),r(e,G,t),r(e,E,t),r(e,J,t),r(e,y,t),r(e,V,t),r(e,B,t),r(e,X,t),r(e,S,t),n(S,q);for(let m=0;m<g.length;m+=1)g[m].m(q,null);n(S,oe),n(S,A);for(let m=0;m<k.length;m+=1)k[m].m(A,null);O=!0},p(e,[t]){var me,de;(!O||t&1)&&a!==(a=e[0].name+"")&&Y(_,a),(!O||t&1)&&D!==(D=e[0].name+"")&&Y(H,D);const m={};t&9&&(m.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').confirmEmailChange(
            'TOKEN',
            'YOUR_PASSWORD',
        );
    `),t&9&&(m.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(de=e[0])==null?void 0:de.name}').confirmEmailChange(
          'TOKEN',
          'YOUR_PASSWORD',
        );
    `),w.$set(m),(!O||t&1)&&M!==(M=e[0].name+"")&&Y(z,M),t&6&&(W=e[2],g=pe(g,t,ie,1,e,W,ae,q,Pe,_e,null,be)),t&6&&(U=e[2],Se(),k=pe(k,t,ce,1,e,U,ne,A,Oe,ke,null,ue),Re())},i(e){if(!O){Z(w.$$.fragment,e);for(let t=0;t<U.length;t+=1)Z(k[t]);O=!0}},o(e){x(w.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);O=!1},d(e){e&&f(l),e&&f(i),e&&f(d),e&&f(I),ge(w,e),e&&f(F),e&&f(R),e&&f(L),e&&f(P),e&&f(G),e&&f(E),e&&f(J),e&&f(y),e&&f(V),e&&f(B),e&&f(X),e&&f(S);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function De(o,l,s){let a,{collection:_=new Te}=l,u=204,i=[];const d=p=>s(1,u=p.code);return o.$$set=p=>{"collection"in p&&s(0,_=p.collection)},s(3,a=Ee.getApiExampleUrl(ye.baseUrl)),s(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,u,i,a,d]}class Me extends Ce{constructor(l){super(),$e(this,l,De,Ue,we,{collection:0})}}export{Me as default};
