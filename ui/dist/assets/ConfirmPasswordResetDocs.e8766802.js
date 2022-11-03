import{S as Xe,i as Ye,s as Ze,O as Ge,e as a,w as b,b as p,c as me,f as m,g as r,h as l,m as he,x as j,P as Ue,Q as et,k as tt,R as lt,n as st,t as J,a as Q,o as d,d as _e,L as ot,C as je,p as at,r as x,u as nt}from"./index.86238c7e.js";import{S as it}from"./SdkTabs.c9d3b27f.js";function Je(i,s,o){const n=i.slice();return n[5]=s[o],n}function Qe(i,s,o){const n=i.slice();return n[5]=s[o],n}function xe(i,s){let o,n=s[5].code+"",k,S,c,f;function h(){return s[4](s[5])}return{key:i,first:null,c(){o=a("button"),k=b(n),S=p(),m(o,"class","tab-item"),x(o,"active",s[1]===s[5].code),this.first=o},m(P,R){r(P,o,R),l(o,k),l(o,S),c||(f=nt(o,"click",h),c=!0)},p(P,R){s=P,R&4&&n!==(n=s[5].code+"")&&j(k,n),R&6&&x(o,"active",s[1]===s[5].code)},d(P){P&&d(o),c=!1,f()}}}function ze(i,s){let o,n,k,S;return n=new Ge({props:{content:s[5].body}}),{key:i,first:null,c(){o=a("div"),me(n.$$.fragment),k=p(),m(o,"class","tab-item"),x(o,"active",s[1]===s[5].code),this.first=o},m(c,f){r(c,o,f),he(n,o,null),l(o,k),S=!0},p(c,f){s=c;const h={};f&4&&(h.content=s[5].body),n.$set(h),(!S||f&6)&&x(o,"active",s[1]===s[5].code)},i(c){S||(J(n.$$.fragment,c),S=!0)},o(c){Q(n.$$.fragment,c),S=!1},d(c){c&&d(o),_e(n)}}}function rt(i){var Ke,He;let s,o,n=i[0].name+"",k,S,c,f,h,P,R,K=i[0].name+"",z,ke,Se,G,X,g,Y,W,Z,C,H,ve,L,D,we,ee,V=i[0].name+"",te,Pe,le,E,se,A,oe,M,ae,$,ne,Re,ie,y,re,ge,de,Ce,_,$e,T,ye,Oe,Ne,ce,Te,pe,We,De,Ee,fe,Ae,ue,F,be,O,q,w=[],Me=new Map,Fe,B,v=[],qe=new Map,N;g=new it({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[3]}');

        ...

        const authData = await pb.collection('${(Ke=i[0])==null?void 0:Ke.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[3]}');

        ...

        final authData = await pb.collection('${(He=i[0])==null?void 0:He.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `}}),T=new Ge({props:{content:"?expand=relField1,relField2.subRelField"}});let U=i[2];const Be=e=>e[5].code;for(let e=0;e<U.length;e+=1){let t=Qe(i,U,e),u=Be(t);Me.set(u,w[e]=xe(u,t))}let I=i[2];const Ie=e=>e[5].code;for(let e=0;e<I.length;e+=1){let t=Je(i,I,e),u=Ie(t);qe.set(u,v[e]=ze(u,t))}return{c(){s=a("h3"),o=b("Confirm password reset ("),k=b(n),S=b(")"),c=p(),f=a("div"),h=a("p"),P=b("Confirms "),R=a("strong"),z=b(K),ke=b(" password reset request."),Se=p(),G=a("p"),G.textContent="Returns the refreshed auth data.",X=p(),me(g.$$.fragment),Y=p(),W=a("h6"),W.textContent="API details",Z=p(),C=a("div"),H=a("strong"),H.textContent="POST",ve=p(),L=a("div"),D=a("p"),we=b("/api/collections/"),ee=a("strong"),te=b(V),Pe=b("/confirm-password-reset"),le=p(),E=a("div"),E.textContent="Body Parameters",se=p(),A=a("table"),A.innerHTML=`<thead><tr><th>Param</th> 
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
            <td>The new password confirmation.</td></tr></tbody>`,oe=p(),M=a("div"),M.textContent="Query parameters",ae=p(),$=a("table"),ne=a("thead"),ne.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,Re=p(),ie=a("tbody"),y=a("tr"),re=a("td"),re.textContent="expand",ge=p(),de=a("td"),de.innerHTML='<span class="label">String</span>',Ce=p(),_=a("td"),$e=b(`Auto expand record relations. Ex.:
                `),me(T.$$.fragment),ye=b(`
                Supports up to 6-levels depth nested relations expansion. `),Oe=a("br"),Ne=b(`
                The expanded relations will be appended to the record under the
                `),ce=a("code"),ce.textContent="expand",Te=b(" property (eg. "),pe=a("code"),pe.textContent='"expand": {"relField1": {...}, ...}',We=b(`).
                `),De=a("br"),Ee=b(`
                Only the relations to which the account has permissions to `),fe=a("strong"),fe.textContent="view",Ae=b(" will be expanded."),ue=p(),F=a("div"),F.textContent="Responses",be=p(),O=a("div"),q=a("div");for(let e=0;e<w.length;e+=1)w[e].c();Fe=p(),B=a("div");for(let e=0;e<v.length;e+=1)v[e].c();m(s,"class","m-b-sm"),m(f,"class","content txt-lg m-b-sm"),m(W,"class","m-b-xs"),m(H,"class","label label-primary"),m(L,"class","content"),m(C,"class","alert alert-success"),m(E,"class","section-title"),m(A,"class","table-compact table-border m-b-base"),m(M,"class","section-title"),m($,"class","table-compact table-border m-b-base"),m(F,"class","section-title"),m(q,"class","tabs-header compact left"),m(B,"class","tabs-content"),m(O,"class","tabs")},m(e,t){r(e,s,t),l(s,o),l(s,k),l(s,S),r(e,c,t),r(e,f,t),l(f,h),l(h,P),l(h,R),l(R,z),l(h,ke),l(f,Se),l(f,G),r(e,X,t),he(g,e,t),r(e,Y,t),r(e,W,t),r(e,Z,t),r(e,C,t),l(C,H),l(C,ve),l(C,L),l(L,D),l(D,we),l(D,ee),l(ee,te),l(D,Pe),r(e,le,t),r(e,E,t),r(e,se,t),r(e,A,t),r(e,oe,t),r(e,M,t),r(e,ae,t),r(e,$,t),l($,ne),l($,Re),l($,ie),l(ie,y),l(y,re),l(y,ge),l(y,de),l(y,Ce),l(y,_),l(_,$e),he(T,_,null),l(_,ye),l(_,Oe),l(_,Ne),l(_,ce),l(_,Te),l(_,pe),l(_,We),l(_,De),l(_,Ee),l(_,fe),l(_,Ae),r(e,ue,t),r(e,F,t),r(e,be,t),r(e,O,t),l(O,q);for(let u=0;u<w.length;u+=1)w[u].m(q,null);l(O,Fe),l(O,B);for(let u=0;u<v.length;u+=1)v[u].m(B,null);N=!0},p(e,[t]){var Le,Ve;(!N||t&1)&&n!==(n=e[0].name+"")&&j(k,n),(!N||t&1)&&K!==(K=e[0].name+"")&&j(z,K);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const authData = await pb.collection('${(Le=e[0])==null?void 0:Le.name}').confirmPasswordReset(
            'TOKEN',
            'NEW_PASSWORD',
            'NEW_PASSWORD_CONFIRM',
        );

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final authData = await pb.collection('${(Ve=e[0])==null?void 0:Ve.name}').confirmPasswordReset(
          'TOKEN',
          'NEW_PASSWORD',
          'NEW_PASSWORD_CONFIRM',
        );

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `),g.$set(u),(!N||t&1)&&V!==(V=e[0].name+"")&&j(te,V),t&6&&(U=e[2],w=Ue(w,t,Be,1,e,U,Me,q,et,xe,null,Qe)),t&6&&(I=e[2],tt(),v=Ue(v,t,Ie,1,e,I,qe,B,lt,ze,null,Je),st())},i(e){if(!N){J(g.$$.fragment,e),J(T.$$.fragment,e);for(let t=0;t<I.length;t+=1)J(v[t]);N=!0}},o(e){Q(g.$$.fragment,e),Q(T.$$.fragment,e);for(let t=0;t<v.length;t+=1)Q(v[t]);N=!1},d(e){e&&d(s),e&&d(c),e&&d(f),e&&d(X),_e(g,e),e&&d(Y),e&&d(W),e&&d(Z),e&&d(C),e&&d(le),e&&d(E),e&&d(se),e&&d(A),e&&d(oe),e&&d(M),e&&d(ae),e&&d($),_e(T),e&&d(ue),e&&d(F),e&&d(be),e&&d(O);for(let t=0;t<w.length;t+=1)w[t].d();for(let t=0;t<v.length;t+=1)v[t].d()}}}function dt(i,s,o){let n,{collection:k=new ot}=s,S=200,c=[];const f=h=>o(1,S=h.code);return i.$$set=h=>{"collection"in h&&o(0,k=h.collection)},i.$$.update=()=>{i.$$.dirty&1&&o(2,c=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:je.dummyCollectionRecord(k)},null,2)},{code:400,body:`
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
            `}])},o(3,n=je.getApiExampleUrl(at.baseUrl)),[k,S,c,n,f]}class ft extends Xe{constructor(s){super(),Ye(this,s,dt,rt,Ze,{collection:0})}}export{ft as default};
