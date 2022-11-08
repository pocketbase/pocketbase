import{S as Ge,i as Je,s as Ve,O as ze,e as a,w as f,b as m,c as be,f as b,g as c,h as l,m as ue,x as I,P as We,Q as Xe,k as Ze,R as xe,n as et,t as Q,a as z,o as r,d as _e,L as tt,C as lt,p as st,r as G,u as nt}from"./index.a710f1eb.js";import{S as at}from"./SdkTabs.d25acbcc.js";function Ye(i,s,n){const o=i.slice();return o[5]=s[n],o}function je(i,s,n){const o=i.slice();return o[5]=s[n],o}function Ie(i,s){let n,o=s[5].code+"",v,h,d,u;function _(){return s[4](s[5])}return{key:i,first:null,c(){n=a("button"),v=f(o),h=m(),b(n,"class","tab-item"),G(n,"active",s[1]===s[5].code),this.first=n},m(C,g){c(C,n,g),l(n,v),l(n,h),d||(u=nt(n,"click",_),d=!0)},p(C,g){s=C,g&4&&o!==(o=s[5].code+"")&&I(v,o),g&6&&G(n,"active",s[1]===s[5].code)},d(C){C&&r(n),d=!1,u()}}}function Qe(i,s){let n,o,v,h;return o=new ze({props:{content:s[5].body}}),{key:i,first:null,c(){n=a("div"),be(o.$$.fragment),v=m(),b(n,"class","tab-item"),G(n,"active",s[1]===s[5].code),this.first=n},m(d,u){c(d,n,u),ue(o,n,null),l(n,v),h=!0},p(d,u){s=d;const _={};u&4&&(_.content=s[5].body),o.$set(_),(!h||u&6)&&G(n,"active",s[1]===s[5].code)},i(d){h||(Q(o.$$.fragment,d),h=!0)},o(d){z(o.$$.fragment,d),h=!1},d(d){d&&r(n),_e(o)}}}function ot(i){var He,Ke;let s,n,o=i[0].name+"",v,h,d,u,_,C,g,L=i[0].name+"",J,ke,V,S,X,A,Z,P,N,he,W,B,ve,x,Y=i[0].name+"",ee,$e,te,q,le,D,se,U,ne,y,ae,we,oe,O,ie,Ce,ce,ge,k,Se,E,Pe,ye,Oe,re,Te,de,Re,Ee,Ae,pe,Be,fe,M,me,T,F,w=[],qe=new Map,De,H,$=[],Ue=new Map,R;S=new at({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[3]}');

        ...

        await pb.collection('${(He=i[0])==null?void 0:He.name}').confirmEmailChange(
            'TOKEN',
            'YOUR_PASSWORD',
        );
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[3]}');

        ...

        await pb.collection('${(Ke=i[0])==null?void 0:Ke.name}').confirmEmailChange(
          'TOKEN',
          'YOUR_PASSWORD',
        );
    `}}),E=new ze({props:{content:"?expand=relField1,relField2.subRelField"}});let j=i[2];const Me=e=>e[5].code;for(let e=0;e<j.length;e+=1){let t=je(i,j,e),p=Me(t);qe.set(p,w[e]=Ie(p,t))}let K=i[2];const Fe=e=>e[5].code;for(let e=0;e<K.length;e+=1){let t=Ye(i,K,e),p=Fe(t);Ue.set(p,$[e]=Qe(p,t))}return{c(){s=a("h3"),n=f("Confirm email change ("),v=f(o),h=f(")"),d=m(),u=a("div"),_=a("p"),C=f("Confirms "),g=a("strong"),J=f(L),ke=f(" email change request."),V=m(),be(S.$$.fragment),X=m(),A=a("h6"),A.textContent="API details",Z=m(),P=a("div"),N=a("strong"),N.textContent="POST",he=m(),W=a("div"),B=a("p"),ve=f("/api/collections/"),x=a("strong"),ee=f(Y),$e=f("/confirm-email-change"),te=m(),q=a("div"),q.textContent="Body Parameters",le=m(),D=a("table"),D.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>token</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The token from the change email request email.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>password</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The account password to confirm the email change.</td></tr></tbody>`,se=m(),U=a("div"),U.textContent="Query parameters",ne=m(),y=a("table"),ae=a("thead"),ae.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,we=m(),oe=a("tbody"),O=a("tr"),ie=a("td"),ie.textContent="expand",Ce=m(),ce=a("td"),ce.innerHTML='<span class="label">String</span>',ge=m(),k=a("td"),Se=f(`Auto expand record relations. Ex.:
                `),be(E.$$.fragment),Pe=f(`
                Supports up to 6-levels depth nested relations expansion. `),ye=a("br"),Oe=f(`
                The expanded relations will be appended to the record under the
                `),re=a("code"),re.textContent="expand",Te=f(" property (eg. "),de=a("code"),de.textContent='"expand": {"relField1": {...}, ...}',Re=f(`).
                `),Ee=a("br"),Ae=f(`
                Only the relations to which the account has permissions to `),pe=a("strong"),pe.textContent="view",Be=f(" will be expanded."),fe=m(),M=a("div"),M.textContent="Responses",me=m(),T=a("div"),F=a("div");for(let e=0;e<w.length;e+=1)w[e].c();De=m(),H=a("div");for(let e=0;e<$.length;e+=1)$[e].c();b(s,"class","m-b-sm"),b(u,"class","content txt-lg m-b-sm"),b(A,"class","m-b-xs"),b(N,"class","label label-primary"),b(W,"class","content"),b(P,"class","alert alert-success"),b(q,"class","section-title"),b(D,"class","table-compact table-border m-b-base"),b(U,"class","section-title"),b(y,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(F,"class","tabs-header compact left"),b(H,"class","tabs-content"),b(T,"class","tabs")},m(e,t){c(e,s,t),l(s,n),l(s,v),l(s,h),c(e,d,t),c(e,u,t),l(u,_),l(_,C),l(_,g),l(g,J),l(_,ke),c(e,V,t),ue(S,e,t),c(e,X,t),c(e,A,t),c(e,Z,t),c(e,P,t),l(P,N),l(P,he),l(P,W),l(W,B),l(B,ve),l(B,x),l(x,ee),l(B,$e),c(e,te,t),c(e,q,t),c(e,le,t),c(e,D,t),c(e,se,t),c(e,U,t),c(e,ne,t),c(e,y,t),l(y,ae),l(y,we),l(y,oe),l(oe,O),l(O,ie),l(O,Ce),l(O,ce),l(O,ge),l(O,k),l(k,Se),ue(E,k,null),l(k,Pe),l(k,ye),l(k,Oe),l(k,re),l(k,Te),l(k,de),l(k,Re),l(k,Ee),l(k,Ae),l(k,pe),l(k,Be),c(e,fe,t),c(e,M,t),c(e,me,t),c(e,T,t),l(T,F);for(let p=0;p<w.length;p+=1)w[p].m(F,null);l(T,De),l(T,H);for(let p=0;p<$.length;p+=1)$[p].m(H,null);R=!0},p(e,[t]){var Le,Ne;(!R||t&1)&&o!==(o=e[0].name+"")&&I(v,o),(!R||t&1)&&L!==(L=e[0].name+"")&&I(J,L);const p={};t&9&&(p.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(Le=e[0])==null?void 0:Le.name}').confirmEmailChange(
            'TOKEN',
            'YOUR_PASSWORD',
        );
    `),t&9&&(p.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(Ne=e[0])==null?void 0:Ne.name}').confirmEmailChange(
          'TOKEN',
          'YOUR_PASSWORD',
        );
    `),S.$set(p),(!R||t&1)&&Y!==(Y=e[0].name+"")&&I(ee,Y),t&6&&(j=e[2],w=We(w,t,Me,1,e,j,qe,F,Xe,Ie,null,je)),t&6&&(K=e[2],Ze(),$=We($,t,Fe,1,e,K,Ue,H,xe,Qe,null,Ye),et())},i(e){if(!R){Q(S.$$.fragment,e),Q(E.$$.fragment,e);for(let t=0;t<K.length;t+=1)Q($[t]);R=!0}},o(e){z(S.$$.fragment,e),z(E.$$.fragment,e);for(let t=0;t<$.length;t+=1)z($[t]);R=!1},d(e){e&&r(s),e&&r(d),e&&r(u),e&&r(V),_e(S,e),e&&r(X),e&&r(A),e&&r(Z),e&&r(P),e&&r(te),e&&r(q),e&&r(le),e&&r(D),e&&r(se),e&&r(U),e&&r(ne),e&&r(y),_e(E),e&&r(fe),e&&r(M),e&&r(me),e&&r(T);for(let t=0;t<w.length;t+=1)w[t].d();for(let t=0;t<$.length;t+=1)$[t].d()}}}function it(i,s,n){let o,{collection:v=new tt}=s,h=204,d=[];const u=_=>n(1,h=_.code);return i.$$set=_=>{"collection"in _&&n(0,v=_.collection)},n(3,o=lt.getApiExampleUrl(st.baseUrl)),n(2,d=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[v,h,d,o,u]}class dt extends Ge{constructor(s){super(),Je(this,s,it,ot,Ve,{collection:0})}}export{dt as default};
