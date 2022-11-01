import{S as Ye,i as Ze,s as xe,O as Xe,e as a,w as b,b as f,c as me,f as m,g as r,h as l,m as he,x as J,P as Ie,Q as et,k as tt,R as lt,n as ot,t as Q,a as W,o as c,d as _e,L as st,C as Je,p as at,r as z,u as nt}from"./index.be8ffbe5.js";import{S as it}from"./SdkTabs.8f55857f.js";function Qe(i,o,s){const n=i.slice();return n[5]=o[s],n}function We(i,o,s){const n=i.slice();return n[5]=o[s],n}function ze(i,o){let s,n=o[5].code+"",k,v,d,p;function h(){return o[4](o[5])}return{key:i,first:null,c(){s=a("button"),k=b(n),v=f(),m(s,"class","tab-item"),z(s,"active",o[1]===o[5].code),this.first=s},m(y,g){r(y,s,g),l(s,k),l(s,v),d||(p=nt(s,"click",h),d=!0)},p(y,g){o=y,g&4&&n!==(n=o[5].code+"")&&J(k,n),g&6&&z(s,"active",o[1]===o[5].code)},d(y){y&&c(s),d=!1,p()}}}function Ge(i,o){let s,n,k,v;return n=new Xe({props:{content:o[5].body}}),{key:i,first:null,c(){s=a("div"),me(n.$$.fragment),k=f(),m(s,"class","tab-item"),z(s,"active",o[1]===o[5].code),this.first=s},m(d,p){r(d,s,p),he(n,s,null),l(s,k),v=!0},p(d,p){o=d;const h={};p&4&&(h.content=o[5].body),n.$set(h),(!v||p&6)&&z(s,"active",o[1]===o[5].code)},i(d){v||(Q(n.$$.fragment,d),v=!0)},o(d){W(n.$$.fragment,d),v=!1},d(d){d&&c(s),_e(n)}}}function rt(i){var He,Le;let o,s,n=i[0].name+"",k,v,d,p,h,y,g,H=i[0].name+"",G,ke,ve,X,Y,w,Z,D,x,C,L,Se,U,E,$e,ee,j=i[0].name+"",te,ye,le,q,oe,M,se,N,ae,T,ne,ge,ie,P,re,we,ce,Ce,_,Te,B,Pe,Oe,Ve,de,Be,fe,De,Ee,qe,pe,Me,ue,R,be,O,F,$=[],Ne=new Map,Re,K,S=[],Fe=new Map,V;w=new it({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[3]}');

        ...

        const authData = await pb.collection('${(He=i[0])==null?void 0:He.name}').confirmVerification('TOKEN');

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[3]}');

        ...

        final authData = await pb.collection('${(Le=i[0])==null?void 0:Le.name}').confirmVerification('TOKEN');

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `}}),B=new Xe({props:{content:"?expand=relField1,relField2.subRelField"}});let I=i[2];const Ke=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=We(i,I,e),u=Ke(t);Ne.set(u,$[e]=ze(u,t))}let A=i[2];const Ae=e=>e[5].code;for(let e=0;e<A.length;e+=1){let t=Qe(i,A,e),u=Ae(t);Fe.set(u,S[e]=Ge(u,t))}return{c(){o=a("h3"),s=b("Confirm verification ("),k=b(n),v=b(")"),d=f(),p=a("div"),h=a("p"),y=b("Confirms "),g=a("strong"),G=b(H),ke=b(" account verification request."),ve=f(),X=a("p"),X.textContent="Returns the refreshed auth data.",Y=f(),me(w.$$.fragment),Z=f(),D=a("h6"),D.textContent="API details",x=f(),C=a("div"),L=a("strong"),L.textContent="POST",Se=f(),U=a("div"),E=a("p"),$e=b("/api/collections/"),ee=a("strong"),te=b(j),ye=b("/confirm-verification"),le=f(),q=a("div"),q.textContent="Body Parameters",oe=f(),M=a("table"),M.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>token</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The token from the verification request email.</td></tr></tbody>`,se=f(),N=a("div"),N.textContent="Query parameters",ae=f(),T=a("table"),ne=a("thead"),ne.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,ge=f(),ie=a("tbody"),P=a("tr"),re=a("td"),re.textContent="expand",we=f(),ce=a("td"),ce.innerHTML='<span class="label">String</span>',Ce=f(),_=a("td"),Te=b(`Auto expand record relations. Ex.:
                `),me(B.$$.fragment),Pe=b(`
                Supports up to 6-levels depth nested relations expansion. `),Oe=a("br"),Ve=b(`
                The expanded relations will be appended to the record under the
                `),de=a("code"),de.textContent="expand",Be=b(" property (eg. "),fe=a("code"),fe.textContent='"expand": {"relField1": {...}, ...}',De=b(`).
                `),Ee=a("br"),qe=b(`
                Only the relations to which the account has permissions to `),pe=a("strong"),pe.textContent="view",Me=b(" will be expanded."),ue=f(),R=a("div"),R.textContent="Responses",be=f(),O=a("div"),F=a("div");for(let e=0;e<$.length;e+=1)$[e].c();Re=f(),K=a("div");for(let e=0;e<S.length;e+=1)S[e].c();m(o,"class","m-b-sm"),m(p,"class","content txt-lg m-b-sm"),m(D,"class","m-b-xs"),m(L,"class","label label-primary"),m(U,"class","content"),m(C,"class","alert alert-success"),m(q,"class","section-title"),m(M,"class","table-compact table-border m-b-base"),m(N,"class","section-title"),m(T,"class","table-compact table-border m-b-base"),m(R,"class","section-title"),m(F,"class","tabs-header compact left"),m(K,"class","tabs-content"),m(O,"class","tabs")},m(e,t){r(e,o,t),l(o,s),l(o,k),l(o,v),r(e,d,t),r(e,p,t),l(p,h),l(h,y),l(h,g),l(g,G),l(h,ke),l(p,ve),l(p,X),r(e,Y,t),he(w,e,t),r(e,Z,t),r(e,D,t),r(e,x,t),r(e,C,t),l(C,L),l(C,Se),l(C,U),l(U,E),l(E,$e),l(E,ee),l(ee,te),l(E,ye),r(e,le,t),r(e,q,t),r(e,oe,t),r(e,M,t),r(e,se,t),r(e,N,t),r(e,ae,t),r(e,T,t),l(T,ne),l(T,ge),l(T,ie),l(ie,P),l(P,re),l(P,we),l(P,ce),l(P,Ce),l(P,_),l(_,Te),he(B,_,null),l(_,Pe),l(_,Oe),l(_,Ve),l(_,de),l(_,Be),l(_,fe),l(_,De),l(_,Ee),l(_,qe),l(_,pe),l(_,Me),r(e,ue,t),r(e,R,t),r(e,be,t),r(e,O,t),l(O,F);for(let u=0;u<$.length;u+=1)$[u].m(F,null);l(O,Re),l(O,K);for(let u=0;u<S.length;u+=1)S[u].m(K,null);V=!0},p(e,[t]){var Ue,je;(!V||t&1)&&n!==(n=e[0].name+"")&&J(k,n),(!V||t&1)&&H!==(H=e[0].name+"")&&J(G,H);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const authData = await pb.collection('${(Ue=e[0])==null?void 0:Ue.name}').confirmVerification('TOKEN');

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final authData = await pb.collection('${(je=e[0])==null?void 0:je.name}').confirmVerification('TOKEN');

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `),w.$set(u),(!V||t&1)&&j!==(j=e[0].name+"")&&J(te,j),t&6&&(I=e[2],$=Ie($,t,Ke,1,e,I,Ne,F,et,ze,null,We)),t&6&&(A=e[2],tt(),S=Ie(S,t,Ae,1,e,A,Fe,K,lt,Ge,null,Qe),ot())},i(e){if(!V){Q(w.$$.fragment,e),Q(B.$$.fragment,e);for(let t=0;t<A.length;t+=1)Q(S[t]);V=!0}},o(e){W(w.$$.fragment,e),W(B.$$.fragment,e);for(let t=0;t<S.length;t+=1)W(S[t]);V=!1},d(e){e&&c(o),e&&c(d),e&&c(p),e&&c(Y),_e(w,e),e&&c(Z),e&&c(D),e&&c(x),e&&c(C),e&&c(le),e&&c(q),e&&c(oe),e&&c(M),e&&c(se),e&&c(N),e&&c(ae),e&&c(T),_e(B),e&&c(ue),e&&c(R),e&&c(be),e&&c(O);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<S.length;t+=1)S[t].d()}}}function ct(i,o,s){let n,{collection:k=new st}=o,v=200,d=[];const p=h=>s(1,v=h.code);return i.$$set=h=>{"collection"in h&&s(0,k=h.collection)},i.$$.update=()=>{i.$$.dirty&1&&s(2,d=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:Je.dummyCollectionRecord(k)},null,2)},{code:400,body:`
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
            `}])},s(3,n=Je.getApiExampleUrl(at.baseUrl)),[k,v,d,n,p]}class pt extends Ye{constructor(o){super(),Ze(this,o,ct,rt,xe,{collection:0})}}export{pt as default};
