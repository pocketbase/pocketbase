import{S as Ge,i as Je,s as Ve,O as ze,e as o,w as f,b,c as me,f as m,g as r,h as l,m as ue,x as Q,P as Le,Q as Xe,k as Ye,R as Ze,n as et,t as x,a as z,o as d,d as _e,L as tt,C as lt,p as st,r as G,u as nt}from"./index.a710f1eb.js";import{S as ot}from"./SdkTabs.d25acbcc.js";function Ue(i,s,n){const a=i.slice();return a[5]=s[n],a}function je(i,s,n){const a=i.slice();return a[5]=s[n],a}function Qe(i,s){let n,a=s[5].code+"",w,v,c,u;function _(){return s[4](s[5])}return{key:i,first:null,c(){n=o("button"),w=f(a),v=b(),m(n,"class","tab-item"),G(n,"active",s[1]===s[5].code),this.first=n},m(P,R){r(P,n,R),l(n,w),l(n,v),c||(u=nt(n,"click",_),c=!0)},p(P,R){s=P,R&4&&a!==(a=s[5].code+"")&&Q(w,a),R&6&&G(n,"active",s[1]===s[5].code)},d(P){P&&d(n),c=!1,u()}}}function xe(i,s){let n,a,w,v;return a=new ze({props:{content:s[5].body}}),{key:i,first:null,c(){n=o("div"),me(a.$$.fragment),w=b(),m(n,"class","tab-item"),G(n,"active",s[1]===s[5].code),this.first=n},m(c,u){r(c,n,u),ue(a,n,null),l(n,w),v=!0},p(c,u){s=c;const _={};u&4&&(_.content=s[5].body),a.$set(_),(!v||u&6)&&G(n,"active",s[1]===s[5].code)},i(c){v||(x(a.$$.fragment,c),v=!0)},o(c){z(a.$$.fragment,c),v=!1},d(c){c&&d(n),_e(a)}}}function at(i){var Be,Ie;let s,n,a=i[0].name+"",w,v,c,u,_,P,R,H=i[0].name+"",J,ke,V,$,X,E,Y,C,K,ve,L,A,we,Z,U=i[0].name+"",ee,he,te,D,le,g,se,M,ne,O,oe,Se,ae,N,ie,Pe,re,Re,k,$e,y,Ce,Oe,Ne,de,Te,ce,We,ye,Ee,pe,Ae,fe,F,be,T,q,S=[],De=new Map,ge,B,h=[],Me=new Map,W;$=new ot({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[3]}');

        ...

        await pb.collection('${(Be=i[0])==null?void 0:Be.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[3]}');

        ...

        await pb.collection('${(Ie=i[0])==null?void 0:Ie.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );
    `}}),y=new ze({props:{content:"?expand=relField1,relField2.subRelField"}});let j=i[2];const Fe=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=je(i,j,e),p=Fe(t);De.set(p,S[e]=Qe(p,t))}let I=i[2];const qe=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=Ue(i,I,e),p=qe(t);Me.set(p,h[e]=xe(p,t))}return{c(){s=o("h3"),n=f("Confirm password reset ("),w=f(a),v=f(")"),c=b(),u=o("div"),_=o("p"),P=f("Confirms "),R=o("strong"),J=f(H),ke=f(" password reset request."),V=b(),me($.$$.fragment),X=b(),E=o("h6"),E.textContent="API details",Y=b(),C=o("div"),K=o("strong"),K.textContent="POST",ve=b(),L=o("div"),A=o("p"),we=f("/api/collections/"),Z=o("strong"),ee=f(U),he=f("/confirm-password-reset"),te=b(),D=o("div"),D.textContent="Body Parameters",le=b(),g=o("table"),g.innerHTML=`<thead><tr><th>Param</th> 
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
            <td>The new password confirmation.</td></tr></tbody>`,se=b(),M=o("div"),M.textContent="Query parameters",ne=b(),O=o("table"),oe=o("thead"),oe.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Se=b(),ae=o("tbody"),N=o("tr"),ie=o("td"),ie.textContent="expand",Pe=b(),re=o("td"),re.innerHTML='<span class="label">String</span>',Re=b(),k=o("td"),$e=f(`Auto expand record relations. Ex.:
                `),me(y.$$.fragment),Ce=f(`
                Supports up to 6-levels depth nested relations expansion. `),Oe=o("br"),Ne=f(`
                The expanded relations will be appended to the record under the
                `),de=o("code"),de.textContent="expand",Te=f(" property (eg. "),ce=o("code"),ce.textContent='"expand": {"relField1": {...}, ...}',We=f(`).
                `),ye=o("br"),Ee=f(`
                Only the relations to which the account has permissions to `),pe=o("strong"),pe.textContent="view",Ae=f(" will be expanded."),fe=b(),F=o("div"),F.textContent="Responses",be=b(),T=o("div"),q=o("div");for(let e=0;e<S.length;e+=1)S[e].c();ge=b(),B=o("div");for(let e=0;e<h.length;e+=1)h[e].c();m(s,"class","m-b-sm"),m(u,"class","content txt-lg m-b-sm"),m(E,"class","m-b-xs"),m(K,"class","label label-primary"),m(L,"class","content"),m(C,"class","alert alert-success"),m(D,"class","section-title"),m(g,"class","table-compact table-border m-b-base"),m(M,"class","section-title"),m(O,"class","table-compact table-border m-b-base"),m(F,"class","section-title"),m(q,"class","tabs-header compact left"),m(B,"class","tabs-content"),m(T,"class","tabs")},m(e,t){r(e,s,t),l(s,n),l(s,w),l(s,v),r(e,c,t),r(e,u,t),l(u,_),l(_,P),l(_,R),l(R,J),l(_,ke),r(e,V,t),ue($,e,t),r(e,X,t),r(e,E,t),r(e,Y,t),r(e,C,t),l(C,K),l(C,ve),l(C,L),l(L,A),l(A,we),l(A,Z),l(Z,ee),l(A,he),r(e,te,t),r(e,D,t),r(e,le,t),r(e,g,t),r(e,se,t),r(e,M,t),r(e,ne,t),r(e,O,t),l(O,oe),l(O,Se),l(O,ae),l(ae,N),l(N,ie),l(N,Pe),l(N,re),l(N,Re),l(N,k),l(k,$e),ue(y,k,null),l(k,Ce),l(k,Oe),l(k,Ne),l(k,de),l(k,Te),l(k,ce),l(k,We),l(k,ye),l(k,Ee),l(k,pe),l(k,Ae),r(e,fe,t),r(e,F,t),r(e,be,t),r(e,T,t),l(T,q);for(let p=0;p<S.length;p+=1)S[p].m(q,null);l(T,ge),l(T,B);for(let p=0;p<h.length;p+=1)h[p].m(B,null);W=!0},p(e,[t]){var He,Ke;(!W||t&1)&&a!==(a=e[0].name+"")&&Q(w,a),(!W||t&1)&&H!==(H=e[0].name+"")&&Q(J,H);const p={};t&9&&(p.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(He=e[0])==null?void 0:He.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );
    `),t&9&&(p.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(Ke=e[0])==null?void 0:Ke.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );
    `),$.$set(p),(!W||t&1)&&U!==(U=e[0].name+"")&&Q(ee,U),t&6&&(j=e[2],S=Le(S,t,Fe,1,e,j,De,q,Xe,Qe,null,je)),t&6&&(I=e[2],Ye(),h=Le(h,t,qe,1,e,I,Me,B,Ze,xe,null,Ue),et())},i(e){if(!W){x($.$$.fragment,e),x(y.$$.fragment,e);for(let t=0;t<I.length;t+=1)x(h[t]);W=!0}},o(e){z($.$$.fragment,e),z(y.$$.fragment,e);for(let t=0;t<h.length;t+=1)z(h[t]);W=!1},d(e){e&&d(s),e&&d(c),e&&d(u),e&&d(V),_e($,e),e&&d(X),e&&d(E),e&&d(Y),e&&d(C),e&&d(te),e&&d(D),e&&d(le),e&&d(g),e&&d(se),e&&d(M),e&&d(ne),e&&d(O),_e(y),e&&d(fe),e&&d(F),e&&d(be),e&&d(T);for(let t=0;t<S.length;t+=1)S[t].d();for(let t=0;t<h.length;t+=1)h[t].d()}}}function it(i,s,n){let a,{collection:w=new tt}=s,v=204,c=[];const u=_=>n(1,v=_.code);return i.$$set=_=>{"collection"in _&&n(0,w=_.collection)},n(3,a=lt.getApiExampleUrl(st.baseUrl)),n(2,c=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[w,v,c,a,u]}class ct extends Ge{constructor(s){super(),Je(this,s,it,at,Ve,{collection:0})}}export{ct as default};
