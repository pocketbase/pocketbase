import{S as Ne,i as Ue,s as je,O as ze,e as a,w as k,b as p,c as ae,f as b,g as c,h as o,m as ne,x as re,P as qe,Q as xe,k as Ie,R as Je,n as Ke,t as U,a as j,o as d,d as ie,L as Qe,C as He,p as We,r as x,u as Ge}from"./index.786ddc4b.js";import{S as Xe}from"./SdkTabs.af9891cd.js";function Le(r,l,s){const n=r.slice();return n[5]=l[s],n}function Ee(r,l,s){const n=r.slice();return n[5]=l[s],n}function Fe(r,l){let s,n=l[5].code+"",m,_,i,f;function v(){return l[4](l[5])}return{key:r,first:null,c(){s=a("button"),m=k(n),_=p(),b(s,"class","tab-item"),x(s,"active",l[1]===l[5].code),this.first=s},m(g,w){c(g,s,w),o(s,m),o(s,_),i||(f=Ge(s,"click",v),i=!0)},p(g,w){l=g,w&4&&n!==(n=l[5].code+"")&&re(m,n),w&6&&x(s,"active",l[1]===l[5].code)},d(g){g&&d(s),i=!1,f()}}}function Ve(r,l){let s,n,m,_;return n=new ze({props:{content:l[5].body}}),{key:r,first:null,c(){s=a("div"),ae(n.$$.fragment),m=p(),b(s,"class","tab-item"),x(s,"active",l[1]===l[5].code),this.first=s},m(i,f){c(i,s,f),ne(n,s,null),o(s,m),_=!0},p(i,f){l=i;const v={};f&4&&(v.content=l[5].body),n.$set(v),(!_||f&6)&&x(s,"active",l[1]===l[5].code)},i(i){_||(U(n.$$.fragment,i),_=!0)},o(i){j(n.$$.fragment,i),_=!1},d(i){i&&d(s),ie(n)}}}function Ye(r){var Be,Me;let l,s,n=r[0].name+"",m,_,i,f,v,g,w,B,I,S,E,ce,F,M,de,J,V=r[0].name+"",K,ue,pe,z,Q,D,W,T,G,fe,X,C,Y,he,Z,be,h,me,P,_e,ke,ve,ee,ge,te,ye,Se,$e,oe,we,le,O,se,R,q,$=[],Te=new Map,Ce,H,y=[],Re=new Map,A;g=new Xe({props:{js:`
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
    `}}),P=new ze({props:{content:"?expand=relField1,relField2.subRelField"}});let N=r[2];const Pe=e=>e[5].code;for(let e=0;e<N.length;e+=1){let t=Ee(r,N,e),u=Pe(t);Te.set(u,$[e]=Fe(u,t))}let L=r[2];const Ae=e=>e[5].code;for(let e=0;e<L.length;e+=1){let t=Le(r,L,e),u=Ae(t);Re.set(u,y[e]=Ve(u,t))}return{c(){l=a("h3"),s=k("Auth refresh ("),m=k(n),_=k(")"),i=p(),f=a("div"),f.innerHTML=`<p>Returns a new auth response (token and record data) for an
        <strong>already authenticated record</strong>.</p> 
    <p><em>This method is usually called by users on page/screen reload to ensure that the previously stored
            data in <code>pb.authStore</code> is still valid and up-to-date.</em></p>`,v=p(),ae(g.$$.fragment),w=p(),B=a("h6"),B.textContent="API details",I=p(),S=a("div"),E=a("strong"),E.textContent="POST",ce=p(),F=a("div"),M=a("p"),de=k("/api/collections/"),J=a("strong"),K=k(V),ue=k("/auth-refresh"),pe=p(),z=a("p"),z.innerHTML="Requires record <code>Authorization:TOKEN</code> header",Q=p(),D=a("div"),D.textContent="Query parameters",W=p(),T=a("table"),G=a("thead"),G.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,fe=p(),X=a("tbody"),C=a("tr"),Y=a("td"),Y.textContent="expand",he=p(),Z=a("td"),Z.innerHTML='<span class="label">String</span>',be=p(),h=a("td"),me=k(`Auto expand record relations. Ex.:
                `),ae(P.$$.fragment),_e=k(`
                Supports up to 6-levels depth nested relations expansion. `),ke=a("br"),ve=k(`
                The expanded relations will be appended to the record under the
                `),ee=a("code"),ee.textContent="expand",ge=k(" property (eg. "),te=a("code"),te.textContent='"expand": {"relField1": {...}, ...}',ye=k(`).
                `),Se=a("br"),$e=k(`
                Only the relations to which the request user has permissions to `),oe=a("strong"),oe.textContent="view",we=k(" will be expanded."),le=p(),O=a("div"),O.textContent="Responses",se=p(),R=a("div"),q=a("div");for(let e=0;e<$.length;e+=1)$[e].c();Ce=p(),H=a("div");for(let e=0;e<y.length;e+=1)y[e].c();b(l,"class","m-b-sm"),b(f,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(E,"class","label label-primary"),b(F,"class","content"),b(z,"class","txt-hint txt-sm txt-right"),b(S,"class","alert alert-success"),b(D,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(O,"class","section-title"),b(q,"class","tabs-header compact left"),b(H,"class","tabs-content"),b(R,"class","tabs")},m(e,t){c(e,l,t),o(l,s),o(l,m),o(l,_),c(e,i,t),c(e,f,t),c(e,v,t),ne(g,e,t),c(e,w,t),c(e,B,t),c(e,I,t),c(e,S,t),o(S,E),o(S,ce),o(S,F),o(F,M),o(M,de),o(M,J),o(J,K),o(M,ue),o(S,pe),o(S,z),c(e,Q,t),c(e,D,t),c(e,W,t),c(e,T,t),o(T,G),o(T,fe),o(T,X),o(X,C),o(C,Y),o(C,he),o(C,Z),o(C,be),o(C,h),o(h,me),ne(P,h,null),o(h,_e),o(h,ke),o(h,ve),o(h,ee),o(h,ge),o(h,te),o(h,ye),o(h,Se),o(h,$e),o(h,oe),o(h,we),c(e,le,t),c(e,O,t),c(e,se,t),c(e,R,t),o(R,q);for(let u=0;u<$.length;u+=1)$[u].m(q,null);o(R,Ce),o(R,H);for(let u=0;u<y.length;u+=1)y[u].m(H,null);A=!0},p(e,[t]){var De,Oe;(!A||t&1)&&n!==(n=e[0].name+"")&&re(m,n);const u={};t&9&&(u.js=`
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
    `),g.$set(u),(!A||t&1)&&V!==(V=e[0].name+"")&&re(K,V),t&6&&(N=e[2],$=qe($,t,Pe,1,e,N,Te,q,xe,Fe,null,Ee)),t&6&&(L=e[2],Ie(),y=qe(y,t,Ae,1,e,L,Re,H,Je,Ve,null,Le),Ke())},i(e){if(!A){U(g.$$.fragment,e),U(P.$$.fragment,e);for(let t=0;t<L.length;t+=1)U(y[t]);A=!0}},o(e){j(g.$$.fragment,e),j(P.$$.fragment,e);for(let t=0;t<y.length;t+=1)j(y[t]);A=!1},d(e){e&&d(l),e&&d(i),e&&d(f),e&&d(v),ie(g,e),e&&d(w),e&&d(B),e&&d(I),e&&d(S),e&&d(Q),e&&d(D),e&&d(W),e&&d(T),ie(P),e&&d(le),e&&d(O),e&&d(se),e&&d(R);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<y.length;t+=1)y[t].d()}}}function Ze(r,l,s){let n,{collection:m=new Qe}=l,_=200,i=[];const f=v=>s(1,_=v.code);return r.$$set=v=>{"collection"in v&&s(0,m=v.collection)},r.$$.update=()=>{r.$$.dirty&1&&s(2,i=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:He.dummyCollectionRecord(m)},null,2)},{code:401,body:`
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
            `}])},s(3,n=He.getApiExampleUrl(We.baseUrl)),[m,_,i,n,f]}class ot extends Ne{constructor(l){super(),Ue(this,l,Ze,Ye,je,{collection:0})}}export{ot as default};
