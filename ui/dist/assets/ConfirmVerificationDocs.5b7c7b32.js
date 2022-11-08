import{S as Je,i as We,s as Xe,O as Ge,e as n,w as p,b,c as me,f as m,g as r,h as l,m as ue,x as Q,P as Ue,Q as Ye,k as Ze,R as xe,n as et,t as z,a as G,o as c,d as _e,L as tt,C as lt,p as st,r as J,u as ot}from"./index.a710f1eb.js";import{S as nt}from"./SdkTabs.d25acbcc.js";function je(a,s,o){const i=a.slice();return i[5]=s[o],i}function Ie(a,s,o){const i=a.slice();return i[5]=s[o],i}function Qe(a,s){let o,i=s[5].code+"",h,v,d,u;function _(){return s[4](s[5])}return{key:a,first:null,c(){o=n("button"),h=p(i),v=b(),m(o,"class","tab-item"),J(o,"active",s[1]===s[5].code),this.first=o},m(C,y){r(C,o,y),l(o,h),l(o,v),d||(u=ot(o,"click",_),d=!0)},p(C,y){s=C,y&4&&i!==(i=s[5].code+"")&&Q(h,i),y&6&&J(o,"active",s[1]===s[5].code)},d(C){C&&c(o),d=!1,u()}}}function ze(a,s){let o,i,h,v;return i=new Ge({props:{content:s[5].body}}),{key:a,first:null,c(){o=n("div"),me(i.$$.fragment),h=b(),m(o,"class","tab-item"),J(o,"active",s[1]===s[5].code),this.first=o},m(d,u){r(d,o,u),ue(i,o,null),l(o,h),v=!0},p(d,u){s=d;const _={};u&4&&(_.content=s[5].body),i.$set(_),(!v||u&6)&&J(o,"active",s[1]===s[5].code)},i(d){v||(z(i.$$.fragment,d),v=!0)},o(d){G(i.$$.fragment,d),v=!1},d(d){d&&c(o),_e(i)}}}function it(a){var Le,Ne;let s,o,i=a[0].name+"",h,v,d,u,_,C,y,R=a[0].name+"",W,ke,X,T,Y,E,Z,P,D,ve,U,M,he,x,j=a[0].name+"",ee,$e,te,F,le,V,se,A,oe,g,ne,we,ie,B,ae,Ce,re,ye,k,Te,q,Pe,ge,Be,ce,Se,de,Oe,qe,Ee,fe,Me,pe,H,be,S,K,w=[],Fe=new Map,Ve,L,$=[],Ae=new Map,O;T=new nt({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${a[3]}');

        ...

        await pb.collection('${(Le=a[0])==null?void 0:Le.name}').confirmVerification('TOKEN');
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${a[3]}');

        ...

        await pb.collection('${(Ne=a[0])==null?void 0:Ne.name}').confirmVerification('TOKEN');
    `}}),q=new Ge({props:{content:"?expand=relField1,relField2.subRelField"}});let I=a[2];const He=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=Ie(a,I,e),f=He(t);Fe.set(f,w[e]=Qe(f,t))}let N=a[2];const Ke=e=>e[5].code;for(let e=0;e<N.length;e+=1){let t=je(a,N,e),f=Ke(t);Ae.set(f,$[e]=ze(f,t))}return{c(){s=n("h3"),o=p("Confirm verification ("),h=p(i),v=p(")"),d=b(),u=n("div"),_=n("p"),C=p("Confirms "),y=n("strong"),W=p(R),ke=p(" account verification request."),X=b(),me(T.$$.fragment),Y=b(),E=n("h6"),E.textContent="API details",Z=b(),P=n("div"),D=n("strong"),D.textContent="POST",ve=b(),U=n("div"),M=n("p"),he=p("/api/collections/"),x=n("strong"),ee=p(j),$e=p("/confirm-verification"),te=b(),F=n("div"),F.textContent="Body Parameters",le=b(),V=n("table"),V.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>token</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The token from the verification request email.</td></tr></tbody>`,se=b(),A=n("div"),A.textContent="Query parameters",oe=b(),g=n("table"),ne=n("thead"),ne.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,we=b(),ie=n("tbody"),B=n("tr"),ae=n("td"),ae.textContent="expand",Ce=b(),re=n("td"),re.innerHTML='<span class="label">String</span>',ye=b(),k=n("td"),Te=p(`Auto expand record relations. Ex.:
                `),me(q.$$.fragment),Pe=p(`
                Supports up to 6-levels depth nested relations expansion. `),ge=n("br"),Be=p(`
                The expanded relations will be appended to the record under the
                `),ce=n("code"),ce.textContent="expand",Se=p(" property (eg. "),de=n("code"),de.textContent='"expand": {"relField1": {...}, ...}',Oe=p(`).
                `),qe=n("br"),Ee=p(`
                Only the relations to which the account has permissions to `),fe=n("strong"),fe.textContent="view",Me=p(" will be expanded."),pe=b(),H=n("div"),H.textContent="Responses",be=b(),S=n("div"),K=n("div");for(let e=0;e<w.length;e+=1)w[e].c();Ve=b(),L=n("div");for(let e=0;e<$.length;e+=1)$[e].c();m(s,"class","m-b-sm"),m(u,"class","content txt-lg m-b-sm"),m(E,"class","m-b-xs"),m(D,"class","label label-primary"),m(U,"class","content"),m(P,"class","alert alert-success"),m(F,"class","section-title"),m(V,"class","table-compact table-border m-b-base"),m(A,"class","section-title"),m(g,"class","table-compact table-border m-b-base"),m(H,"class","section-title"),m(K,"class","tabs-header compact left"),m(L,"class","tabs-content"),m(S,"class","tabs")},m(e,t){r(e,s,t),l(s,o),l(s,h),l(s,v),r(e,d,t),r(e,u,t),l(u,_),l(_,C),l(_,y),l(y,W),l(_,ke),r(e,X,t),ue(T,e,t),r(e,Y,t),r(e,E,t),r(e,Z,t),r(e,P,t),l(P,D),l(P,ve),l(P,U),l(U,M),l(M,he),l(M,x),l(x,ee),l(M,$e),r(e,te,t),r(e,F,t),r(e,le,t),r(e,V,t),r(e,se,t),r(e,A,t),r(e,oe,t),r(e,g,t),l(g,ne),l(g,we),l(g,ie),l(ie,B),l(B,ae),l(B,Ce),l(B,re),l(B,ye),l(B,k),l(k,Te),ue(q,k,null),l(k,Pe),l(k,ge),l(k,Be),l(k,ce),l(k,Se),l(k,de),l(k,Oe),l(k,qe),l(k,Ee),l(k,fe),l(k,Me),r(e,pe,t),r(e,H,t),r(e,be,t),r(e,S,t),l(S,K);for(let f=0;f<w.length;f+=1)w[f].m(K,null);l(S,Ve),l(S,L);for(let f=0;f<$.length;f+=1)$[f].m(L,null);O=!0},p(e,[t]){var Re,De;(!O||t&1)&&i!==(i=e[0].name+"")&&Q(h,i),(!O||t&1)&&R!==(R=e[0].name+"")&&Q(W,R);const f={};t&9&&(f.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        await pb.collection('${(Re=e[0])==null?void 0:Re.name}').confirmVerification('TOKEN');
    `),t&9&&(f.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        await pb.collection('${(De=e[0])==null?void 0:De.name}').confirmVerification('TOKEN');
    `),T.$set(f),(!O||t&1)&&j!==(j=e[0].name+"")&&Q(ee,j),t&6&&(I=e[2],w=Ue(w,t,He,1,e,I,Fe,K,Ye,Qe,null,Ie)),t&6&&(N=e[2],Ze(),$=Ue($,t,Ke,1,e,N,Ae,L,xe,ze,null,je),et())},i(e){if(!O){z(T.$$.fragment,e),z(q.$$.fragment,e);for(let t=0;t<N.length;t+=1)z($[t]);O=!0}},o(e){G(T.$$.fragment,e),G(q.$$.fragment,e);for(let t=0;t<$.length;t+=1)G($[t]);O=!1},d(e){e&&c(s),e&&c(d),e&&c(u),e&&c(X),_e(T,e),e&&c(Y),e&&c(E),e&&c(Z),e&&c(P),e&&c(te),e&&c(F),e&&c(le),e&&c(V),e&&c(se),e&&c(A),e&&c(oe),e&&c(g),_e(q),e&&c(pe),e&&c(H),e&&c(be),e&&c(S);for(let t=0;t<w.length;t+=1)w[t].d();for(let t=0;t<$.length;t+=1)$[t].d()}}}function at(a,s,o){let i,{collection:h=new tt}=s,v=204,d=[];const u=_=>o(1,v=_.code);return a.$$set=_=>{"collection"in _&&o(0,h=_.collection)},o(3,i=lt.getApiExampleUrl(st.baseUrl)),o(2,d=[{code:204,body:"null"},{code:400,body:`
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
            `}]),[h,v,d,i,u]}class dt extends Je{constructor(s){super(),We(this,s,at,it,Xe,{collection:0})}}export{dt as default};
