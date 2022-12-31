import{S as Se,i as he,s as Re,e as c,w,b as v,c as ve,f as b,g as r,h as n,m as we,x as K,O as me,P as Oe,k as Ne,Q as Ce,n as We,t as Z,a as x,o as d,d as Pe,R as $e,C as Ee,p as Te,r as U,u as ge,N as Ae}from"./index.89a3f554.js";import{S as De}from"./SdkTabs.0a6ad1c9.js";function ue(o,s,l){const a=o.slice();return a[5]=s[l],a}function be(o,s,l){const a=o.slice();return a[5]=s[l],a}function _e(o,s){let l,a=s[5].code+"",_,u,i,p;function m(){return s[4](s[5])}return{key:o,first:null,c(){l=c("button"),_=w(a),u=v(),b(l,"class","tab-item"),U(l,"active",s[1]===s[5].code),this.first=l},m(S,h){r(S,l,h),n(l,_),n(l,u),i||(p=ge(l,"click",m),i=!0)},p(S,h){s=S,h&4&&a!==(a=s[5].code+"")&&K(_,a),h&6&&U(l,"active",s[1]===s[5].code)},d(S){S&&d(l),i=!1,p()}}}function ke(o,s){let l,a,_,u;return a=new Ae({props:{content:s[5].body}}),{key:o,first:null,c(){l=c("div"),ve(a.$$.fragment),_=v(),b(l,"class","tab-item"),U(l,"active",s[1]===s[5].code),this.first=l},m(i,p){r(i,l,p),we(a,l,null),n(l,_),u=!0},p(i,p){s=i;const m={};p&4&&(m.content=s[5].body),a.$set(m),(!u||p&6)&&U(l,"active",s[1]===s[5].code)},i(i){u||(Z(a.$$.fragment,i),u=!0)},o(i){x(a.$$.fragment,i),u=!1},d(i){i&&d(l),Pe(a)}}}function ye(o){var re,de;let s,l,a=o[0].name+"",_,u,i,p,m,S,h,q=o[0].name+"",j,ee,H,R,L,W,Q,O,B,te,M,$,se,z,I=o[0].name+"",G,le,J,E,V,T,X,g,Y,N,A,P=[],ae=new Map,oe,D,k=[],ne=new Map,C;R=new De({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${o[3]}');

        ...

        await pb.collection('${(re=o[0])==null?void 0:re.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${o[3]}');

        ...

        await pb.collection('${(de=o[0])==null?void 0:de.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );
    `}});let F=o[2];const ie=e=>e[5].code;for(let e=0;e<F.length;e+=1){let t=be(o,F,e),f=ie(t);ae.set(f,P[e]=_e(f,t))}let y=o[2];const ce=e=>e[5].code;for(let e=0;e<y.length;e+=1){let t=ue(o,y,e),f=ce(t);ne.set(f,k[e]=ke(f,t))}return{c(){s=c("h3"),l=w("Confirm password reset ("),_=w(a),u=w(")"),i=v(),p=c("div"),m=c("p"),S=w("Confirms "),h=c("strong"),j=w(q),ee=w(" password reset request and sets a new password."),H=v(),ve(R.$$.fragment),L=v(),W=c("h6"),W.textContent="API details",Q=v(),O=c("div"),B=c("strong"),B.textContent="POST",te=v(),M=c("div"),$=c("p"),se=w("/api/collections/"),z=c("strong"),G=w(I),le=w("/confirm-password-reset"),J=v(),E=c("div"),E.textContent="Body Parameters",V=v(),T=c("table"),T.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>token</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The token from the password reset request email.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>password</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The new password to set.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>passwordConfirm</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The new password confirmation.</td></tr></tbody>`,X=v(),g=c("div"),g.textContent="Responses",Y=v(),N=c("div"),A=c("div");for(let e=0;e<P.length;e+=1)P[e].c();oe=v(),D=c("div");for(let e=0;e<k.length;e+=1)k[e].c();b(s,"class","m-b-sm"),b(p,"class","content txt-lg m-b-sm"),b(W,"class","m-b-xs"),b(B,"class","label label-primary"),b(M,"class","content"),b(O,"class","alert alert-success"),b(E,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(g,"class","section-title"),b(A,"class","tabs-header compact left"),b(D,"class","tabs-content"),b(N,"class","tabs")},m(e,t){r(e,s,t),n(s,l),n(s,_),n(s,u),r(e,i,t),r(e,p,t),n(p,m),n(m,S),n(m,h),n(h,j),n(m,ee),r(e,H,t),we(R,e,t),r(e,L,t),r(e,W,t),r(e,Q,t),r(e,O,t),n(O,B),n(O,te),n(O,M),n(M,$),n($,se),n($,z),n(z,G),n($,le),r(e,J,t),r(e,E,t),r(e,V,t),r(e,T,t),r(e,X,t),r(e,g,t),r(e,Y,t),r(e,N,t),n(N,A);for(let f=0;f<P.length;f+=1)P[f].m(A,null);n(N,oe),n(N,D);for(let f=0;f<k.length;f+=1)k[f].m(D,null);C=!0},p(e,[t]){var fe,pe;(!C||t&1)&&a!==(a=e[0].name+"")&&K(_,a),(!C||t&1)&&q!==(q=e[0].name+"")&&K(j,q);const f={};t&9&&(f.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(fe=e[0])==null?void 0:fe.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );
    `),t&9&&(f.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(pe=e[0])==null?void 0:pe.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );
    `),R.$set(f),(!C||t&1)&&I!==(I=e[0].name+"")&&K(G,I),t&6&&(F=e[2],P=me(P,t,ie,1,e,F,ae,A,Oe,_e,null,be)),t&6&&(y=e[2],Ne(),k=me(k,t,ce,1,e,y,ne,D,Ce,ke,null,ue),We())},i(e){if(!C){Z(R.$$.fragment,e);for(let t=0;t<y.length;t+=1)Z(k[t]);C=!0}},o(e){x(R.$$.fragment,e);for(let t=0;t<k.length;t+=1)x(k[t]);C=!1},d(e){e&&d(s),e&&d(i),e&&d(p),e&&d(H),Pe(R,e),e&&d(L),e&&d(W),e&&d(Q),e&&d(O),e&&d(J),e&&d(E),e&&d(V),e&&d(T),e&&d(X),e&&d(g),e&&d(Y),e&&d(N);for(let t=0;t<P.length;t+=1)P[t].d();for(let t=0;t<k.length;t+=1)k[t].d()}}}function qe(o,s,l){let a,{collection:_=new $e}=s,u=204,i=[];const p=m=>l(1,u=m.code);return o.$$set=m=>{"collection"in m&&l(0,_=m.collection)},l(3,a=Ee.getApiExampleUrl(Te.baseUrl)),l(2,i=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[_,u,i,a,p]}class Ie extends Se{constructor(s){super(),he(this,s,qe,ye,Re,{collection:0})}}export{Ie as default};
