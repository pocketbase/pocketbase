import{S as Ne,i as Ue,s as je,O as ze,e as s,w as k,b as p,c as se,f as b,g as c,h as o,m as ne,x as re,P as Oe,Q as Ie,k as Je,R as Ke,n as Qe,t as U,a as j,o as d,d as ie,L as xe,C as Fe,p as We,r as I,u as Ge}from"./index.3b22183a.js";import{S as Xe}from"./SdkTabs.72ecb434.js";function He(r,l,a){const n=r.slice();return n[5]=l[a],n}function Le(r,l,a){const n=r.slice();return n[5]=l[a],n}function Ee(r,l){let a,n=l[5].code+"",m,_,i,f;function v(){return l[4](l[5])}return{key:r,first:null,c(){a=s("button"),m=k(n),_=p(),b(a,"class","tab-item"),I(a,"active",l[1]===l[5].code),this.first=a},m(g,w){c(g,a,w),o(a,m),o(a,_),i||(f=Ge(a,"click",v),i=!0)},p(g,w){l=g,w&4&&n!==(n=l[5].code+"")&&re(m,n),w&6&&I(a,"active",l[1]===l[5].code)},d(g){g&&d(a),i=!1,f()}}}function Ve(r,l){let a,n,m,_;return n=new ze({props:{content:l[5].body}}),{key:r,first:null,c(){a=s("div"),se(n.$$.fragment),m=p(),b(a,"class","tab-item"),I(a,"active",l[1]===l[5].code),this.first=a},m(i,f){c(i,a,f),ne(n,a,null),o(a,m),_=!0},p(i,f){l=i;const v={};f&4&&(v.content=l[5].body),n.$set(v),(!_||f&6)&&I(a,"active",l[1]===l[5].code)},i(i){_||(U(n.$$.fragment,i),_=!0)},o(i){j(n.$$.fragment,i),_=!1},d(i){i&&d(a),ie(n)}}}function Ye(r){var Be,Me;let l,a,n=r[0].name+"",m,_,i,f,v,g,w,B,J,S,L,ce,E,M,de,K,V=r[0].name+"",Q,ue,pe,z,x,q,W,T,G,fe,X,C,Y,he,Z,be,h,me,P,_e,ke,ve,ee,ge,te,ye,Se,$e,oe,we,le,D,ae,R,O,$=[],Te=new Map,Ce,F,y=[],Re=new Map,A;g=new Xe({props:{js:`
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
    `}}),P=new ze({props:{content:"?expand=relField1,relField2.subRelField"}});let N=r[2];const Pe=e=>e[5].code;for(let e=0;e<N.length;e+=1){let t=Le(r,N,e),u=Pe(t);Te.set(u,$[e]=Ee(u,t))}let H=r[2];const Ae=e=>e[5].code;for(let e=0;e<H.length;e+=1){let t=He(r,H,e),u=Ae(t);Re.set(u,y[e]=Ve(u,t))}return{c(){l=s("h3"),a=k("Auth refresh ("),m=k(n),_=k(")"),i=p(),f=s("div"),f.innerHTML=`<p>Returns a new auth response (token and account data) for an
        <strong>already authenticated record</strong>.</p> 
    <p><em>This method is usually called by users on page/screen reload to ensure that the previously stored
            data in <code>pb.authStore</code> is still valid and up-to-date.</em></p>`,v=p(),se(g.$$.fragment),w=p(),B=s("h6"),B.textContent="API details",J=p(),S=s("div"),L=s("strong"),L.textContent="POST",ce=p(),E=s("div"),M=s("p"),de=k("/api/collections/"),K=s("strong"),Q=k(V),ue=k("/auth-refresh"),pe=p(),z=s("p"),z.innerHTML="Requires record <code>Authorization:TOKEN</code> header",x=p(),q=s("div"),q.textContent="Query parameters",W=p(),T=s("table"),G=s("thead"),G.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,fe=p(),X=s("tbody"),C=s("tr"),Y=s("td"),Y.textContent="expand",he=p(),Z=s("td"),Z.innerHTML='<span class="label">String</span>',be=p(),h=s("td"),me=k(`Auto expand record relations. Ex.:
                `),se(P.$$.fragment),_e=k(`
                Supports up to 6-levels depth nested relations expansion. `),ke=s("br"),ve=k(`
                The expanded relations will be appended to the record under the
                `),ee=s("code"),ee.textContent="expand",ge=k(" property (eg. "),te=s("code"),te.textContent='"expand": {"relField1": {...}, ...}',ye=k(`).
                `),Se=s("br"),$e=k(`
                Only the relations to which the account has permissions to `),oe=s("strong"),oe.textContent="view",we=k(" will be expanded."),le=p(),D=s("div"),D.textContent="Responses",ae=p(),R=s("div"),O=s("div");for(let e=0;e<$.length;e+=1)$[e].c();Ce=p(),F=s("div");for(let e=0;e<y.length;e+=1)y[e].c();b(l,"class","m-b-sm"),b(f,"class","content txt-lg m-b-sm"),b(B,"class","m-b-xs"),b(L,"class","label label-primary"),b(E,"class","content"),b(z,"class","txt-hint txt-sm txt-right"),b(S,"class","alert alert-success"),b(q,"class","section-title"),b(T,"class","table-compact table-border m-b-base"),b(D,"class","section-title"),b(O,"class","tabs-header compact left"),b(F,"class","tabs-content"),b(R,"class","tabs")},m(e,t){c(e,l,t),o(l,a),o(l,m),o(l,_),c(e,i,t),c(e,f,t),c(e,v,t),ne(g,e,t),c(e,w,t),c(e,B,t),c(e,J,t),c(e,S,t),o(S,L),o(S,ce),o(S,E),o(E,M),o(M,de),o(M,K),o(K,Q),o(M,ue),o(S,pe),o(S,z),c(e,x,t),c(e,q,t),c(e,W,t),c(e,T,t),o(T,G),o(T,fe),o(T,X),o(X,C),o(C,Y),o(C,he),o(C,Z),o(C,be),o(C,h),o(h,me),ne(P,h,null),o(h,_e),o(h,ke),o(h,ve),o(h,ee),o(h,ge),o(h,te),o(h,ye),o(h,Se),o(h,$e),o(h,oe),o(h,we),c(e,le,t),c(e,D,t),c(e,ae,t),c(e,R,t),o(R,O);for(let u=0;u<$.length;u+=1)$[u].m(O,null);o(R,Ce),o(R,F);for(let u=0;u<y.length;u+=1)y[u].m(F,null);A=!0},p(e,[t]){var qe,De;(!A||t&1)&&n!==(n=e[0].name+"")&&re(m,n);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const authData = await pb.collection('${(qe=e[0])==null?void 0:qe.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final authData = await pb.collection('${(De=e[0])==null?void 0:De.name}').authRefresh();

        // after the above you can also access the refreshed auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);
    `),g.$set(u),(!A||t&1)&&V!==(V=e[0].name+"")&&re(Q,V),t&6&&(N=e[2],$=Oe($,t,Pe,1,e,N,Te,O,Ie,Ee,null,Le)),t&6&&(H=e[2],Je(),y=Oe(y,t,Ae,1,e,H,Re,F,Ke,Ve,null,He),Qe())},i(e){if(!A){U(g.$$.fragment,e),U(P.$$.fragment,e);for(let t=0;t<H.length;t+=1)U(y[t]);A=!0}},o(e){j(g.$$.fragment,e),j(P.$$.fragment,e);for(let t=0;t<y.length;t+=1)j(y[t]);A=!1},d(e){e&&d(l),e&&d(i),e&&d(f),e&&d(v),ie(g,e),e&&d(w),e&&d(B),e&&d(J),e&&d(S),e&&d(x),e&&d(q),e&&d(W),e&&d(T),ie(P),e&&d(le),e&&d(D),e&&d(ae),e&&d(R);for(let t=0;t<$.length;t+=1)$[t].d();for(let t=0;t<y.length;t+=1)y[t].d()}}}function Ze(r,l,a){let n,{collection:m=new xe}=l,_=200,i=[];const f=v=>a(1,_=v.code);return r.$$set=v=>{"collection"in v&&a(0,m=v.collection)},r.$$.update=()=>{r.$$.dirty&1&&a(2,i=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:Fe.dummyCollectionRecord(m)},null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "Failed to authenticate.",
                  "data": {
                    "identity": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `},{code:401,body:`
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
            `}])},a(3,n=Fe.getApiExampleUrl(We.baseUrl)),[m,_,i,n,f]}class ot extends Ne{constructor(l){super(),Ue(this,l,Ze,Ye,je,{collection:0})}}export{ot as default};
