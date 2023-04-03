import{S as ze,i as Ue,s as je,N as Ve,e as a,w as k,b as p,c as ae,f as b,g as d,h as o,m as ne,x as re,O as qe,P as xe,k as Je,Q as Ke,n as Qe,t as U,a as j,o as u,d as ie,T as Ie,C as He,p as We,r as x,u as Ge}from"./index-4eea3e34.js";import{S as Xe}from"./SdkTabs-5d6cc1d4.js";function Ee(r,l,s){const n=r.slice();return n[5]=l[s],n}function Fe(r,l,s){const n=r.slice();return n[5]=l[s],n}function Le(r,l){let s,n=l[5].code+"",m,_,i,f;function v(){return l[4](l[5])}return{key:r,first:null,c(){s=a("button"),m=k(n),_=p(),b(s,"class","tab-item"),x(s,"active",l[1]===l[5].code),this.first=s},m(g,w){d(g,s,w),o(s,m),o(s,_),i||(f=Ge(s,"click",v),i=!0)},p(g,w){l=g,w&4&&n!==(n=l[5].code+"")&&re(m,n),w&6&&x(s,"active",l[1]===l[5].code)},d(g){g&&u(s),i=!1,f()}}}function Ne(r,l){let s,n,m,_;return n=new Ve({props:{content:l[5].body}}),{key:r,first:null,c(){s=a("div"),ae(n.$$.fragment),m=p(),b(s,"class","tab-item"),x(s,"active",l[1]===l[5].code),this.first=s},m(i,f){d(i,s,f),ne(n,s,null),o(s,m),_=!0},p(i,f){l=i;const v={};f&4&&(v.content=l[5].body),n.$set(v),(!_||f&6)&&x(s,"active",l[1]===l[5].code)},i(i){_||(U(n.$$.fragment,i),_=!0)},o(i){j(n.$$.fragment,i),_=!1},d(i){i&&u(s),ie(n)}}}function Ye(r){var Be,Me;let l,s,n=r[0].name+"",m,_,i,f,v,g,w,B,J,$,F,ce,L,M,de,K,N=r[0].name+"",Q,ue,pe,V,I,D,W,T,G,fe,X,C,Y,he,Z,be,h,me,R,_e,ke,ve,ee,ge,te,ye,Se,$e,oe,we,le,O,se,P,q,S=[],Te=new Map,Ce,H,y=[],Pe=new Map,A;g=new Xe({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${r[3]}');

        ...

        const authData = await pb.collection('${(Be=r[0])==null?void 0:Be.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${r[3]}');

        ...

        final authData = await pb.collection('${(Me=r[0])==null?void 0:Me.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);
    `}}),R=new Ve({props:{content:"?expand=relField1,relField2.subRelField"}});let z=r[2];const Re=e=>e[5].code;for(let e=0;e<z.length;e+=1){let t=Fe(r,z,e),c=Re(t);Te.set(c,S[e]=Le(c,t))}let E=r[2];const Ae=e=>e[5].code;for(let e=0;e<E.length;e+=1){let t=Ee(r,E,e),c=Ae(t);Pe.set(c,y[e]=Ne(c,t))}return{c(){l=a("h3"),s=k("Auth refresh ("),m=k(n),_=k(")"),i=p(),f=a("div"),f.innerHTML=`<p>Returns a new auth response (token and record data) for an
        <strong>already authenticated record</strong>.</p> 
    <p><em>This method is usually called by users on page/screen reload to ensure that the previously stored
            data in <code>pb.authStore</code> is still valid and up-to-date.</em></p>`,v=p(),ae(g.$$.fragment),w=p(),B=a("h6"),B.textContent="API details",J=p(),$=a("div"),F=a("strong"),F.textContent="POST",ce=p(),L=a("div"),M=a("p"),de=k("/api/collections/"),K=a("strong"),Q=k(N),ue=k("/auth-refresh"),pe=p(),V=a("p"),V.innerHTML="Requires record <code>Authorization:TOKEN</code> header",I=p(),D=a("div"),D.textContent="Query parameters",W=p(),T=a("table"),G=a("thead"),G.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,fe=p(),X=a("tbody"),C=a("tr"),Y=a("td"),Y.textContent="expand",he=p(),Z=a("td"),Z.innerHTML='<span class="label">String</span>',be=p(),h=a("td"),me=k(`Auto expand record relations. Ex.:
                `),ae(R.$$.fragment),_e=k(`
                Supports up to 6-levels depth nested relations expansion. `),ke=a("br"),ve=k(`
                The expanded relations will be appended to the record under the
                `),ee=a("code"),ee.textContent="expand",ge=k(" property (eg. "),te=a("code"),te.textContent='"expand": {"relField1": {...}, ...}',ye=k(`).
                `),Se=a("br"),$e=k(`
                Only the relations to which the request user has permissions to `),oe=a("strong"),oe.textContent="view",we=k(" will be expanded."),le=p(),O=a("div"),O.textContent="Responses",se=p(),P=a("div"),q=a("div");for(let e=0;e<S.length;e+=1)S[e].c();Ce=p(),H=a("div");for(let e=0;e<y.length;e+=1)y[e].c();b(l,"class","m-b-sm"),b(f,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(F,"class","label label-primary"),b(L,"class","content"),b(V,"class","txt-hint txt-sm txt-right"),b($,"class","alert alert-success"),b(D,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(O,"class","section-title"),b(q,"class","tabs-header compact left"),b(H,"class","tabs-content"),b(P,"class","tabs")},m(e,t){d(e,l,t),o(l,s),o(l,m),o(l,_),d(e,i,t),d(e,f,t),d(e,v,t),ne(g,e,t),d(e,w,t),d(e,B,t),d(e,J,t),d(e,$,t),o($,F),o($,ce),o($,L),o(L,M),o(M,de),o(M,K),o(K,Q),o(M,ue),o($,pe),o($,V),d(e,I,t),d(e,D,t),d(e,W,t),d(e,T,t),o(T,G),o(T,fe),o(T,X),o(X,C),o(C,Y),o(C,he),o(C,Z),o(C,be),o(C,h),o(h,me),ne(R,h,null),o(h,_e),o(h,ke),o(h,ve),o(h,ee),o(h,ge),o(h,te),o(h,ye),o(h,Se),o(h,$e),o(h,oe),o(h,we),d(e,le,t),d(e,O,t),d(e,se,t),d(e,P,t),o(P,q);for(let c=0;c<S.length;c+=1)S[c]&&S[c].m(q,null);o(P,Ce),o(P,H);for(let c=0;c<y.length;c+=1)y[c]&&y[c].m(H,null);A=!0},p(e,[t]){var De,Oe;(!A||t&1)&&n!==(n=e[0].name+"")&&re(m,n);const c={};t&9&&(c.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const authData = await pb.collection('${(De=e[0])==null?void 0:De.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `),t&9&&(c.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final authData = await pb.collection('${(Oe=e[0])==null?void 0:Oe.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);
    `),g.$set(c),(!A||t&1)&&N!==(N=e[0].name+"")&&re(Q,N),t&6&&(z=e[2],S=qe(S,t,Re,1,e,z,Te,q,xe,Le,null,Fe)),t&6&&(E=e[2],Je(),y=qe(y,t,Ae,1,e,E,Pe,H,Ke,Ne,null,Ee),Qe())},i(e){if(!A){U(g.$$.fragment,e),U(R.$$.fragment,e);for(let t=0;t<E.length;t+=1)U(y[t]);A=!0}},o(e){j(g.$$.fragment,e),j(R.$$.fragment,e);for(let t=0;t<y.length;t+=1)j(y[t]);A=!1},d(e){e&&u(l),e&&u(i),e&&u(f),e&&u(v),ie(g,e),e&&u(w),e&&u(B),e&&u(J),e&&u($),e&&u(I),e&&u(D),e&&u(W),e&&u(T),ie(R),e&&u(le),e&&u(O),e&&u(se),e&&u(P);for(let t=0;t<S.length;t+=1)S[t].d();for(let t=0;t<y.length;t+=1)y[t].d()}}}function Ze(r,l,s){let n,{collection:m=new Ie}=l,_=200,i=[];const f=v=>s(1,_=v.code);return r.$$set=v=>{"collection"in v&&s(0,m=v.collection)},r.$$.update=()=>{r.$$.dirty&1&&s(2,i=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:He.dummyCollectionRecord(m)},null,2)},{code:401,body:`
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
