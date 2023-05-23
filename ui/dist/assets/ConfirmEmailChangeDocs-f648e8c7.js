import{S as Ce,i as $e,s as we,e as r,w as g,b as h,c as he,f as b,g as f,h as n,m as ve,x as Y,N as pe,P as Pe,k as Se,Q as Oe,n as Te,t as Z,a as x,o as m,d as ge,T as Re,C as Ee,p as ye,r as j,u as Be,M as qe}from"./index-a65ca895.js";import{S as Ae}from"./SdkTabs-ad912c8f.js";function ue(o,l,s){const a=o.slice();return a[5]=l[s],a}function be(o,l,s){const a=o.slice();return a[5]=l[s],a}function _e(o,l){let s,a=l[5].code+"",_,u,i,d;function p(){return l[4](l[5])}return{key:o,first:null,c(){s=r("button"),_=g(a),u=h(),b(s,"class","tab-item"),j(s,"active",l[1]===l[5].code),this.first=s},m(C,$){f(C,s,$),n(s,_),n(s,u),i||(d=Be(s,"click",p),i=!0)},p(C,$){l=C,$&4&&a!==(a=l[5].code+"")&&Y(_,a),$&6&&j(s,"active",l[1]===l[5].code)},d(C){C&&m(s),i=!1,d()}}}function ke(o,l){let s,a,_,u;return a=new qe({props:{content:l[5].body}}),{key:o,first:null,c(){s=r("div"),he(a.$$.fragment),_=h(),b(s,"class","tab-item"),j(s,"active",l[1]===l[5].code),this.first=s},m(i,d){f(i,s,d),ve(a,s,null),n(s,_),u=!0},p(i,d){l=i;const p={};d&4&&(p.content=l[5].body),a.$set(p),(!u||d&6)&&j(s,"active",l[1]===l[5].code)},i(i){u||(Z(a.$$.fragment,i),u=!0)},o(i){x(a.$$.fragment,i),u=!1},d(i){i&&m(s),ge(a)}}}function Ue(o){var re,fe;let l,s,a=o[0].name+"",_,u,i,d,p,C,$,D=o[0].name+"",H,ee,F,w,I,T,L,P,M,te,N,R,le,Q,K=o[0].name+"",z,se,G,E,J,y,V,B,X,S,q,v=[],ae=new Map,oe,A,k=[],ne=new Map,O;w=new Ae({props:{js:`
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
    `}});let W=o[2];const ie=e=>e[5].code;for(let e=0;e<W.length;e+=1){let t=be(o,W,e),c=ie(t);ae.set(c,v[e]=_e(c,t))}let U=o[2];const ce=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=ue(o,U,e),c=ce(t);ne.set(c,k[e]=ke(c,t))}return{c(){l=r("h3"),s=g("Confirm email change ("),_=g(a),u=g(")"),i=h(),d=r("div"),p=r("p"),C=g("Confirms "),$=r("strong"),H=g(D),ee=g(" email change request."),F=h(),he(w.$$.fragment),I=h(),T=r("h6"),T.textContent="API details",L=h(),P=r("div"),M=r("strong"),M.textContent="POST",te=h(),N=r("div"),R=r("p"),le=g("/api/collections/"),Q=r("strong"),z=g(K),se=g("/confirm-email-change"),G=h(),E=r("div"),E.textContent="Body Parameters",J=h(),y=r("table"),y.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>token</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The token from the change email request email.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>password</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The account password to confirm the email change.</td></tr></tbody>`,V=h(),B=r("div"),B.textContent="Responses",X=h(),S=r("div"),q=r("div");for(let e=0;e<v.length;e+=1)v[e].c();oe=h(),A=r("div");for(let e=0;e<k.length;e+=1)k[e].c();b(l,"class","m-b-sm"),b(d,"class","content txt-lg m-b-sm"),b(T,"class","m-b-xs"),b(M,"class","label label-primary"),b(N,"class","content"),b(P,"class","alert alert-success"),b(E,"class","section-title"),b(y,"class","table-compact table-border m-b-base"),b(B,"class","section-title"),b(q,"class","tabs-header compact left"),b(A,"class","tabs-content"),b(S,"class","tabs")},m(e,t){f(e,l,t),n(l,s),n(l,_),n(l,u),f(e,i,t),f(e,d,t),n(d,p),n(p,C),n(p,$),n($,H),n(p,ee),f(e,F,t),ve(w,e,t),f(e,I,t),f(e,T,t),f(e,L,t),f(e,P,t),n(P,M),n(P,te),n(P,N),n(N,R),n(R,le),n(R,Q),n(Q,z),n(R,se),f(e,G,t),f(e,E,t),f(e,J,t),f(e,y,t),f(e,V,t),f(e,B,t),f(e,X,t),f(e,S,t),n(S,q);for(let c=0;c<v.length;c+=1)v[c]&&v[c].m(q,null);n(S,oe),n(S,A);for(let c=0;c<k.length;c+=1)k[c]&&k[c].m(A,null);O=!0},p(e,[t]){var me,de;(!O||t&1)&&a!==(a=e[0].name+"")&&Y(_,a),(!O||t&1)&&D!==(D=e[0].name+"")&&Y(H,D);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(me=e[0])==null?void 0:me.name}').confirmEmailChange(
            'TOKEN',
            'YOUR_PASSWORD',
        );
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(de=e[0])==null?void 0:de.name}').confirmEmailChange(
          'TOKEN',
          'YOUR_PASSWORD',
        );
    `),w.$set(c),(!O||t&1)&&K!==(K=e[0].name+"")&&Y(z,K),t&6&&(W=e[2],v=pe(v,t,ie,1,e,W,ae,q,Pe,_e,null,be)),t&6&&(U=e[2],Se(),k=pe(k,t,ce,1,e,U,ne,A,Oe,ke,null,ue),Te())},i(e){if(!O){Z(w.$$.fragment,e);for(let t=0;t<U.length;t+=1)Z(k[t]);O=!0}},o(e){x(w.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);O=!1},d(e){e&&m(l),e&&m(i),e&&m(d),e&&m(F),ge(w,e),e&&m(I),e&&m(T),e&&m(L),e&&m(P),e&&m(G),e&&m(E),e&&m(J),e&&m(y),e&&m(V),e&&m(B),e&&m(X),e&&m(S);for(let t=0;t<v.length;t+=1)v[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function De(o,l,s){let a,{collection:_=new Re}=l,u=204,i=[];const d=p=>s(1,u=p.code);return o.$$set=p=>{"collection"in p&&s(0,_=p.collection)},s(3,a=Ee.getApiExampleUrl(ye.baseUrl)),s(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,u,i,a,d]}class Ke extends Ce{constructor(l){super(),$e(this,l,De,Ue,we,{collection:0})}}export{Ke as default};
