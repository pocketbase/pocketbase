import{S as Ge,i as Xe,s as Ze,O as ze,e as s,w as m,b as p,c as be,f as b,g as r,h as l,m as he,x as j,P as Ye,Q as et,k as tt,R as lt,n as ot,t as I,a as J,o as c,d as _e,L as at,C as je,p as st,r as Q,u as nt}from"./index.b110ca95.js";import{S as it}from"./SdkTabs.b01956c7.js";function Ie(i,o,a){const n=i.slice();return n[5]=o[a],n}function Je(i,o,a){const n=i.slice();return n[5]=o[a],n}function Qe(i,o){let a,n=o[5].code+"",k,v,d,u;function h(){return o[4](o[5])}return{key:i,first:null,c(){a=s("button"),k=m(n),v=p(),b(a,"class","tab-item"),Q(a,"active",o[1]===o[5].code),this.first=a},m(C,$){r(C,a,$),l(a,k),l(a,v),d||(u=nt(a,"click",h),d=!0)},p(C,$){o=C,$&4&&n!==(n=o[5].code+"")&&j(k,n),$&6&&Q(a,"active",o[1]===o[5].code)},d(C){C&&c(a),d=!1,u()}}}function xe(i,o){let a,n,k,v;return n=new ze({props:{content:o[5].body}}),{key:i,first:null,c(){a=s("div"),be(n.$$.fragment),k=p(),b(a,"class","tab-item"),Q(a,"active",o[1]===o[5].code),this.first=a},m(d,u){r(d,a,u),he(n,a,null),l(a,k),v=!0},p(d,u){o=d;const h={};u&4&&(h.content=o[5].body),n.$set(h),(!v||u&6)&&Q(a,"active",o[1]===o[5].code)},i(d){v||(I(n.$$.fragment,d),v=!0)},o(d){J(n.$$.fragment,d),v=!1},d(d){d&&c(a),_e(n)}}}function rt(i){var We,He;let o,a,n=i[0].name+"",k,v,d,u,h,C,$,W=i[0].name+"",x,ke,ve,z,G,y,X,D,Z,w,H,Se,L,A,ge,ee,V=i[0].name+"",te,Ce,le,B,oe,q,ae,U,se,O,ne,$e,ie,T,re,ye,ce,we,_,Oe,E,Te,Pe,Re,de,Ee,pe,De,Ae,Be,ue,qe,fe,M,me,P,N,g=[],Ue=new Map,Me,F,S=[],Ne=new Map,R;y=new it({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[3]}');

        ...

        const authData = await pb.collection('${(We=i[0])==null?void 0:We.name}').confirmEmailChange(
            'TOKEN',
            'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[3]}');

        ...

        final authData = await pb.collection('${(He=i[0])==null?void 0:He.name}').confirmEmailChange(
          'TOKEN',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `}}),E=new ze({props:{content:"?expand=relField1,relField2.subRelField"}});let Y=i[2];const Fe=e=>e[5].code;for(let e=0;e<Y.length;e+=1){let t=Je(i,Y,e),f=Fe(t);Ue.set(f,g[e]=Qe(f,t))}let K=i[2];const Ke=e=>e[5].code;for(let e=0;e<K.length;e+=1){let t=Ie(i,K,e),f=Ke(t);Ne.set(f,S[e]=xe(f,t))}return{c(){o=s("h3"),a=m("Confirm email change ("),k=m(n),v=m(")"),d=p(),u=s("div"),h=s("p"),C=m("Confirms "),$=s("strong"),x=m(W),ke=m(" email change request."),ve=p(),z=s("p"),z.textContent="Returns the refreshed auth data.",G=p(),be(y.$$.fragment),X=p(),D=s("h6"),D.textContent="API details",Z=p(),w=s("div"),H=s("strong"),H.textContent="POST",Se=p(),L=s("div"),A=s("p"),ge=m("/api/collections/"),ee=s("strong"),te=m(V),Ce=m("/confirm-email-change"),le=p(),B=s("div"),B.textContent="Body Parameters",oe=p(),q=s("table"),q.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>token</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The token from the change email request email.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>password</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The account password to confirm the email change.</td></tr></tbody>`,ae=p(),U=s("div"),U.textContent="Query parameters",se=p(),O=s("table"),ne=s("thead"),ne.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,$e=p(),ie=s("tbody"),T=s("tr"),re=s("td"),re.textContent="expand",ye=p(),ce=s("td"),ce.innerHTML='<span class="label">String</span>',we=p(),_=s("td"),Oe=m(`Auto expand record relations. Ex.:
                `),be(E.$$.fragment),Te=m(`
                Supports up to 6-levels depth nested relations expansion. `),Pe=s("br"),Re=m(`
                The expanded relations will be appended to the record under the
                `),de=s("code"),de.textContent="expand",Ee=m(" property (eg. "),pe=s("code"),pe.textContent='"expand": {"relField1": {...}, ...}',De=m(`).
                `),Ae=s("br"),Be=m(`
                Only the relations to which the account has permissions to `),ue=s("strong"),ue.textContent="view",qe=m(" will be expanded."),fe=p(),M=s("div"),M.textContent="Responses",me=p(),P=s("div"),N=s("div");for(let e=0;e<g.length;e+=1)g[e].c();Me=p(),F=s("div");for(let e=0;e<S.length;e+=1)S[e].c();b(o,"class","m-b-sm"),b(u,"class","content txt-lg m-b-sm"),b(D,"class","m-b-xs"),b(H,"class","label label-primary"),b(L,"class","content"),b(w,"class","alert alert-success"),b(B,"class","section-title"),b(q,"class","table-compact table-border m-b-base"),b(U,"class","section-title"),b(O,"class","table-compact table-border m-b-base"),b(M,"class","section-title"),b(N,"class","tabs-header compact left"),b(F,"class","tabs-content"),b(P,"class","tabs")},m(e,t){r(e,o,t),l(o,a),l(o,k),l(o,v),r(e,d,t),r(e,u,t),l(u,h),l(h,C),l(h,$),l($,x),l(h,ke),l(u,ve),l(u,z),r(e,G,t),he(y,e,t),r(e,X,t),r(e,D,t),r(e,Z,t),r(e,w,t),l(w,H),l(w,Se),l(w,L),l(L,A),l(A,ge),l(A,ee),l(ee,te),l(A,Ce),r(e,le,t),r(e,B,t),r(e,oe,t),r(e,q,t),r(e,ae,t),r(e,U,t),r(e,se,t),r(e,O,t),l(O,ne),l(O,$e),l(O,ie),l(ie,T),l(T,re),l(T,ye),l(T,ce),l(T,we),l(T,_),l(_,Oe),he(E,_,null),l(_,Te),l(_,Pe),l(_,Re),l(_,de),l(_,Ee),l(_,pe),l(_,De),l(_,Ae),l(_,Be),l(_,ue),l(_,qe),r(e,fe,t),r(e,M,t),r(e,me,t),r(e,P,t),l(P,N);for(let f=0;f<g.length;f+=1)g[f].m(N,null);l(P,Me),l(P,F);for(let f=0;f<S.length;f+=1)S[f].m(F,null);R=!0},p(e,[t]){var Le,Ve;(!R||t&1)&&n!==(n=e[0].name+"")&&j(k,n),(!R||t&1)&&W!==(W=e[0].name+"")&&j(x,W);const f={};t&9&&(f.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const authData = await pb.collection('${(Le=e[0])==null?void 0:Le.name}').confirmEmailChange(
            'TOKEN',
            'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `),t&9&&(f.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final authData = await pb.collection('${(Ve=e[0])==null?void 0:Ve.name}').confirmEmailChange(
          'TOKEN',
          'YOUR_PASSWORD',
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `),y.$set(f),(!R||t&1)&&V!==(V=e[0].name+"")&&j(te,V),t&6&&(Y=e[2],g=Ye(g,t,Fe,1,e,Y,Ue,N,et,Qe,null,Je)),t&6&&(K=e[2],tt(),S=Ye(S,t,Ke,1,e,K,Ne,F,lt,xe,null,Ie),ot())},i(e){if(!R){I(y.$$.fragment,e),I(E.$$.fragment,e);for(let t=0;t<K.length;t+=1)I(S[t]);R=!0}},o(e){J(y.$$.fragment,e),J(E.$$.fragment,e);for(let t=0;t<S.length;t+=1)J(S[t]);R=!1},d(e){e&&c(o),e&&c(d),e&&c(u),e&&c(G),_e(y,e),e&&c(X),e&&c(D),e&&c(Z),e&&c(w),e&&c(le),e&&c(B),e&&c(oe),e&&c(q),e&&c(ae),e&&c(U),e&&c(se),e&&c(O),_e(E),e&&c(fe),e&&c(M),e&&c(me),e&&c(P);for(let t=0;t<g.length;t+=1)g[t].d();for(let t=0;t<S.length;t+=1)S[t].d()}}}function ct(i,o,a){let n,{collection:k=new at}=o,v=200,d=[];const u=h=>a(1,v=h.code);return i.$$set=h=>{"collection"in h&&a(0,k=h.collection)},i.$$.update=()=>{i.$$.dirty&1&&a(2,d=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:je.dummyCollectionRecord(k)},null,2)},{code:400,body:`
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
            `}])},a(3,n=je.getApiExampleUrl(st.baseUrl)),[k,v,d,n,u]}class ut extends Ge{constructor(o){super(),Xe(this,o,ct,rt,Ze,{collection:0})}}export{ut as default};
