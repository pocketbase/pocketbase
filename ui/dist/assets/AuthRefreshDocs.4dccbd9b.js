import{S as ze,i as Ue,s as je,M as Ve,e as a,w as k,b as p,c as ae,f as b,g as c,h as o,m as ne,x as re,N as qe,O as xe,k as Ie,P as Je,n as Ke,t as U,a as j,o as d,d as ie,Q as Qe,C as He,p as We,r as x,u as Ge}from"./index.2a48c906.js";import{S as Xe}from"./SdkTabs.9dbb2ca3.js";function Ee(r,l,s){const n=r.slice();return n[5]=l[s],n}function Fe(r,l,s){const n=r.slice();return n[5]=l[s],n}function Le(r,l){let s,n=l[5].code+"",m,_,i,f;function v(){return l[4](l[5])}return{key:r,first:null,c(){s=a("button"),m=k(n),_=p(),b(s,"class","tab-item"),x(s,"active",l[1]===l[5].code),this.first=s},m(g,w){c(g,s,w),o(s,m),o(s,_),i||(f=Ge(s,"click",v),i=!0)},p(g,w){l=g,w&4&&n!==(n=l[5].code+"")&&re(m,n),w&6&&x(s,"active",l[1]===l[5].code)},d(g){g&&d(s),i=!1,f()}}}function Ne(r,l){let s,n,m,_;return n=new Ve({props:{content:l[5].body}}),{key:r,first:null,c(){s=a("div"),ae(n.$$.fragment),m=p(),b(s,"class","tab-item"),x(s,"active",l[1]===l[5].code),this.first=s},m(i,f){c(i,s,f),ne(n,s,null),o(s,m),_=!0},p(i,f){l=i;const v={};f&4&&(v.content=l[5].body),n.$set(v),(!_||f&6)&&x(s,"active",l[1]===l[5].code)},i(i){_||(U(n.$$.fragment,i),_=!0)},o(i){j(n.$$.fragment,i),_=!1},d(i){i&&d(s),ie(n)}}}function Ye(r){var Ae,Be;let l,s,n=r[0].name+"",m,_,i,f,v,g,w,A,I,S,F,ce,L,B,de,J,N=r[0].name+"",K,ue,pe,V,Q,D,W,T,G,fe,X,C,Y,he,Z,be,h,me,R,_e,ke,ve,ee,ge,te,ye,Se,$e,oe,we,le,O,se,P,q,$=[],Te=new Map,Ce,H,y=[],Pe=new Map,M;g=new Xe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${r[3]}');

        ...

        const authData = await pb.collection('${(Ae=r[0])==null?void 0:Ae.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${r[3]}');

        ...

        final authData = await pb.collection('${(Be=r[0])==null?void 0:Be.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);
    `}}),R=new Ve({props:{content:"?expand=relField1,relField2.subRelField"}});let z=r[2];const Re=e=>e[5].code;for(let e=0;e<z.length;e+=1){let t=Fe(r,z,e),u=Re(t);Te.set(u,$[e]=Le(u,t))}let E=r[2];const Me=e=>e[5].code;for(let e=0;e<E.length;e+=1){let t=Ee(r,E,e),u=Me(t);Pe.set(u,y[e]=Ne(u,t))}return{c(){l=a("h3"),s=k("Auth refresh ("),m=k(n),_=k(")"),i=p(),f=a("div"),f.innerHTML=`<p>Returns a new auth response (token and record data) for an
        <strong>already authenticated record</strong>.</p> 
    <p><em>This method is usually called by users on page/screen reload to ensure that the previously stored
            data in <code>pb.authStore</code> is still valid and up-to-date.</em></p>`,v=p(),ae(g.$$.fragment),w=p(),A=a("h6"),A.textContent="API details",I=p(),S=a("div"),F=a("strong"),F.textContent="POST",ce=p(),L=a("div"),B=a("p"),de=k("/api/collections/"),J=a("strong"),K=k(N),ue=k("/auth-refresh"),pe=p(),V=a("p"),V.innerHTML="Requires record <code>Authorization:TOKEN</code> header",Q=p(),D=a("div"),D.textContent="Query parameters",W=p(),T=a("table"),G=a("thead"),G.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,fe=p(),X=a("tbody"),C=a("tr"),Y=a("td"),Y.textContent="expand",he=p(),Z=a("td"),Z.innerHTML='<span class="label">String</span>',be=p(),h=a("td"),me=k(`Auto expand record relations. Ex.:
                `),ae(R.$$.fragment),_e=k(`
                Supports up to 6-levels depth nested relations expansion. `),ke=a("br"),ve=k(`
                The expanded relations will be appended to the record under the
                `),ee=a("code"),ee.textContent="expand",ge=k(" property (eg. "),te=a("code"),te.textContent='"expand": {"relField1": {...}, ...}',ye=k(`).
                `),Se=a("br"),$e=k(`
                Only the relations to which the request user has permissions to `),oe=a("strong"),oe.textContent="view",we=k(" will be expanded."),le=p(),O=a("div"),O.textContent="Responses",se=p(),P=a("div"),q=a("div");for(let e=0;e<$.length;e+=1)$[e].c();Ce=p(),H=a("div");for(let e=0;e<y.length;e+=1)y[e].c();b(l,"class","m-b-sm"),b(f,"class","content txt-lg m-b-sm"),b(A,"class","m-b-xs"),b(F,"class","label label-primary"),b(L,"class","content"),b(V,"class","txt-hint txt-sm txt-right"),b(S,"class","alert alert-success"),b(D,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(O,"class","section-title"),b(q,"class","tabs-header compact left"),b(H,"class","tabs-content"),b(P,"class","tabs")},m(e,t){c(e,l,t),o(l,s),o(l,m),o(l,_),c(e,i,t),c(e,f,t),c(e,v,t),ne(g,e,t),c(e,w,t),c(e,A,t),c(e,I,t),c(e,S,t),o(S,F),o(S,ce),o(S,L),o(L,B),o(B,de),o(B,J),o(J,K),o(B,ue),o(S,pe),o(S,V),c(e,Q,t),c(e,D,t),c(e,W,t),c(e,T,t),o(T,G),o(T,fe),o(T,X),o(X,C),o(C,Y),o(C,he),o(C,Z),o(C,be),o(C,h),o(h,me),ne(R,h,null),o(h,_e),o(h,ke),o(h,ve),o(h,ee),o(h,ge),o(h,te),o(h,ye),o(h,Se),o(h,$e),o(h,oe),o(h,we),c(e,le,t),c(e,O,t),c(e,se,t),c(e,P,t),o(P,q);for(let u=0;u<$.length;u+=1)$[u].m(q,null);o(P,Ce),o(P,H);for(let u=0;u<y.length;u+=1)y[u].m(H,null);M=!0},p(e,[t]){var De,Oe;(!M||t&1)&&n!==(n=e[0].name+"")&&re(m,n);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const authData = await pb.collection('${(De=e[0])==null?void 0:De.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final authData = await pb.collection('${(Oe=e[0])==null?void 0:Oe.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);
    `),g.$set(u),(!M||t&1)&&N!==(N=e[0].name+"")&&re(K,N),t&6&&(z=e[2],$=qe($,t,Re,1,e,z,Te,q,xe,Le,null,Fe)),t&6&&(E=e[2],Ie(),y=qe(y,t,Me,1,e,E,Pe,H,Je,Ne,null,Ee),Ke())},i(e){if(!M){U(g.$$.fragment,e),U(R.$$.fragment,e);for(let t=0;t<E.length;t+=1)U(y[t]);M=!0}},o(e){j(g.$$.fragment,e),j(R.$$.fragment,e);for(let t=0;t<y.length;t+=1)j(y[t]);M=!1},d(e){e&&d(l),e&&d(i),e&&d(f),e&&d(v),ie(g,e),e&&d(w),e&&d(A),e&&d(I),e&&d(S),e&&d(Q),e&&d(D),e&&d(W),e&&d(T),ie(R),e&&d(le),e&&d(O),e&&d(se),e&&d(P);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<y.length;t+=1)y[t].d()}}}function Ze(r,l,s){let n,{collection:m=new Qe}=l,_=200,i=[];const f=v=>s(1,_=v.code);return r.$$set=v=>{"collection"in v&&s(0,m=v.collection)},r.$$.update=()=>{r.$$.dirty&1&&s(2,i=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:He.dummyCollectionRecord(m)},null,2)},{code:401,body:`
                {
                  "code": 401,
                  "message": "The request requires valid record authorization token to be set.",
                  "data": {}
                }
            `},{code:403,body:`
                {
                  "code": 403,
                  "message": "The authorized record model is not allowed to perform this action.",
                  "data": {}
                }
            `},{code:404,body:`
                {
                  "code": 404,
                  "message": "Missing auth record context.",
                  "data": {}
                }
            `}])},s(3,n=He.getApiExampleUrl(We.baseUrl)),[m,_,i,n,f]}class ot extends ze{constructor(l){super(),Ue(this,l,Ze,Ye,je,{collection:0})}}export{ot as default};
