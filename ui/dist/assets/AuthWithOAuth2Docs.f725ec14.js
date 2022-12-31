import{S as je,i as He,s as Je,N as We,e as s,w as v,b as p,c as re,f as h,g as r,h as a,m as ce,x as de,O as Ve,P as Ne,k as Qe,Q as ze,n as Ke,t as j,a as H,o as c,d as ue,R as Ye,C as Be,p as Ge,r as J,u as Xe}from"./index.89a3f554.js";import{S as Ze}from"./SdkTabs.0a6ad1c9.js";function Fe(i,l,o){const n=i.slice();return n[5]=l[o],n}function Le(i,l,o){const n=i.slice();return n[5]=l[o],n}function xe(i,l){let o,n=l[5].code+"",m,_,d,b;function g(){return l[4](l[5])}return{key:i,first:null,c(){o=s("button"),m=v(n),_=p(),h(o,"class","tab-item"),J(o,"active",l[1]===l[5].code),this.first=o},m(k,R){r(k,o,R),a(o,m),a(o,_),d||(b=Xe(o,"click",g),d=!0)},p(k,R){l=k,R&4&&n!==(n=l[5].code+"")&&de(m,n),R&6&&J(o,"active",l[1]===l[5].code)},d(k){k&&c(o),d=!1,b()}}}function Me(i,l){let o,n,m,_;return n=new We({props:{content:l[5].body}}),{key:i,first:null,c(){o=s("div"),re(n.$$.fragment),m=p(),h(o,"class","tab-item"),J(o,"active",l[1]===l[5].code),this.first=o},m(d,b){r(d,o,b),ce(n,o,null),a(o,m),_=!0},p(d,b){l=d;const g={};b&4&&(g.content=l[5].body),n.$set(g),(!_||b&6)&&J(o,"active",l[1]===l[5].code)},i(d){_||(j(n.$$.fragment,d),_=!0)},o(d){H(n.$$.fragment,d),_=!1},d(d){d&&c(o),ue(n)}}}function et(i){var qe,Ie;let l,o,n=i[0].name+"",m,_,d,b,g,k,R,C,N,y,L,pe,x,D,he,Q,M=i[0].name+"",z,be,K,q,Y,I,G,P,X,O,Z,fe,ee,$,te,me,ae,_e,f,ve,E,ge,ke,we,le,Se,oe,Re,ye,Oe,se,$e,ne,U,ie,A,V,S=[],Ae=new Map,Ee,B,w=[],Te=new Map,T;k=new Ze({props:{js:`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${i[3]}');

        ...

        const authData = await pb.collection('${(qe=i[0])==null?void 0:qe.name}').authWithOAuth2(
            'google',
            'CODE',
            'VERIFIER',
            'REDIRECT_URL',
            // optional data that will be used for the new account on OAuth2 sign-up
            {
              'name': 'test',
            },
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `,dart:`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${i[3]}');

        ...

        final authData = await pb.collection('${(Ie=i[0])==null?void 0:Ie.name}').authWithOAuth2(
          'google',
          'CODE',
          'VERIFIER',
          'REDIRECT_URL',
          // optional data that will be used for the new account on OAuth2 sign-up
          createData: {
            'name': 'test',
          },
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `}}),E=new We({props:{content:"?expand=relField1,relField2.subRelField"}});let W=i[2];const Ce=e=>e[5].code;for(let e=0;e<W.length;e+=1){let t=Le(i,W,e),u=Ce(t);Ae.set(u,S[e]=xe(u,t))}let F=i[2];const De=e=>e[5].code;for(let e=0;e<F.length;e+=1){let t=Fe(i,F,e),u=De(t);Te.set(u,w[e]=Me(u,t))}return{c(){l=s("h3"),o=v("Auth with OAuth2 ("),m=v(n),_=v(")"),d=p(),b=s("div"),b.innerHTML=`<p>Authenticate with an OAuth2 provider and returns a new auth token and record data.</p> 
    <p>This action usually should be called right after the provider login page redirect.</p> 
    <p>You could also check the
        <a href="https://pocketbase.io/docs/authentication/#web-oauth2-integration" target="_blank" rel="noopener noreferrer">OAuth2 web integration example
        </a>.</p>`,g=p(),re(k.$$.fragment),R=p(),C=s("h6"),C.textContent="API details",N=p(),y=s("div"),L=s("strong"),L.textContent="POST",pe=p(),x=s("div"),D=s("p"),he=v("/api/collections/"),Q=s("strong"),z=v(M),be=v("/auth-with-oauth2"),K=p(),q=s("div"),q.textContent="Body Parameters",Y=p(),I=s("table"),I.innerHTML=`<thead><tr><th>Param</th> 
            <th>Type</th> 
            <th width="50%">Description</th></tr></thead> 
    <tbody><tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>provider</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The name of the OAuth2 client provider (eg. &quot;google&quot;).</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>code</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The authorization code returned from the initial request.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>codeVerifier</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The code verifier sent with the initial request as part of the code_challenge.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-success">Required</span> 
                    <span>redirectUrl</span></div></td> 
            <td><span class="label">String</span></td> 
            <td>The redirect url sent with the initial request.</td></tr> 
        <tr><td><div class="inline-flex"><span class="label label-warning">Optional</span> 
                    <span>createData</span></div></td> 
            <td><span class="label">Object</span></td> 
            <td><p>Optional data that will be used when creating the auth record on OAuth2 sign-up.</p> 
                <p>The created auth record must comply with the same requirements and validations in the
                    regular <strong>create</strong> action.
                    <br/> 
                    <em>The data can only be in <code>json</code>, aka. <code>multipart/form-data</code> and files
                        upload currently are not supported during OAuth2 sign-ups.</em></p></td></tr></tbody>`,G=p(),P=s("div"),P.textContent="Query parameters",X=p(),O=s("table"),Z=s("thead"),Z.innerHTML=`<tr><th>Param</th> 
            <th>Type</th> 
            <th width="60%">Description</th></tr>`,fe=p(),ee=s("tbody"),$=s("tr"),te=s("td"),te.textContent="expand",me=p(),ae=s("td"),ae.innerHTML='<span class="label">String</span>',_e=p(),f=s("td"),ve=v(`Auto expand record relations. Ex.:
                `),re(E.$$.fragment),ge=v(`
                Supports up to 6-levels depth nested relations expansion. `),ke=s("br"),we=v(`
                The expanded relations will be appended to the record under the
                `),le=s("code"),le.textContent="expand",Se=v(" property (eg. "),oe=s("code"),oe.textContent='"expand": {"relField1": {...}, ...}',Re=v(`).
                `),ye=s("br"),Oe=v(`
                Only the relations to which the request user has permissions to `),se=s("strong"),se.textContent="view",$e=v(" will be expanded."),ne=p(),U=s("div"),U.textContent="Responses",ie=p(),A=s("div"),V=s("div");for(let e=0;e<S.length;e+=1)S[e].c();Ee=p(),B=s("div");for(let e=0;e<w.length;e+=1)w[e].c();h(l,"class","m-b-sm"),h(b,"class","content txt-lg m-b-sm"),h(C,"class","m-b-xs"),h(L,"class","label label-primary"),h(x,"class","content"),h(y,"class","alert alert-success"),h(q,"class","section-title"),h(I,"class","table-compact table-border m-b-base"),h(P,"class","section-title"),h(O,"class","table-compact table-border m-b-base"),h(U,"class","section-title"),h(V,"class","tabs-header compact left"),h(B,"class","tabs-content"),h(A,"class","tabs")},m(e,t){r(e,l,t),a(l,o),a(l,m),a(l,_),r(e,d,t),r(e,b,t),r(e,g,t),ce(k,e,t),r(e,R,t),r(e,C,t),r(e,N,t),r(e,y,t),a(y,L),a(y,pe),a(y,x),a(x,D),a(D,he),a(D,Q),a(Q,z),a(D,be),r(e,K,t),r(e,q,t),r(e,Y,t),r(e,I,t),r(e,G,t),r(e,P,t),r(e,X,t),r(e,O,t),a(O,Z),a(O,fe),a(O,ee),a(ee,$),a($,te),a($,me),a($,ae),a($,_e),a($,f),a(f,ve),ce(E,f,null),a(f,ge),a(f,ke),a(f,we),a(f,le),a(f,Se),a(f,oe),a(f,Re),a(f,ye),a(f,Oe),a(f,se),a(f,$e),r(e,ne,t),r(e,U,t),r(e,ie,t),r(e,A,t),a(A,V);for(let u=0;u<S.length;u+=1)S[u].m(V,null);a(A,Ee),a(A,B);for(let u=0;u<w.length;u+=1)w[u].m(B,null);T=!0},p(e,[t]){var Pe,Ue;(!T||t&1)&&n!==(n=e[0].name+"")&&de(m,n);const u={};t&9&&(u.js=`
        import PocketBase from 'pocketbase';

        const pb = new PocketBase('${e[3]}');

        ...

        const authData = await pb.collection('${(Pe=e[0])==null?void 0:Pe.name}').authWithOAuth2(
            'google',
            'CODE',
            'VERIFIER',
            'REDIRECT_URL',
            // optional data that will be used for the new account on OAuth2 sign-up
            {
              'name': 'test',
            },
        );

        // after the above you can also access the auth data from the authStore
        console.log(pb.authStore.isValid);
        console.log(pb.authStore.token);
        console.log(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `),t&9&&(u.dart=`
        import 'package:pocketbase/pocketbase.dart';

        final pb = PocketBase('${e[3]}');

        ...

        final authData = await pb.collection('${(Ue=e[0])==null?void 0:Ue.name}').authWithOAuth2(
          'google',
          'CODE',
          'VERIFIER',
          'REDIRECT_URL',
          // optional data that will be used for the new account on OAuth2 sign-up
          createData: {
            'name': 'test',
          },
        );

        // after the above you can also access the auth data from the authStore
        print(pb.authStore.isValid);
        print(pb.authStore.token);
        print(pb.authStore.model.id);

        // "logout" the last authenticated account
        pb.authStore.clear();
    `),k.$set(u),(!T||t&1)&&M!==(M=e[0].name+"")&&de(z,M),t&6&&(W=e[2],S=Ve(S,t,Ce,1,e,W,Ae,V,Ne,xe,null,Le)),t&6&&(F=e[2],Qe(),w=Ve(w,t,De,1,e,F,Te,B,ze,Me,null,Fe),Ke())},i(e){if(!T){j(k.$$.fragment,e),j(E.$$.fragment,e);for(let t=0;t<F.length;t+=1)j(w[t]);T=!0}},o(e){H(k.$$.fragment,e),H(E.$$.fragment,e);for(let t=0;t<w.length;t+=1)H(w[t]);T=!1},d(e){e&&c(l),e&&c(d),e&&c(b),e&&c(g),ue(k,e),e&&c(R),e&&c(C),e&&c(N),e&&c(y),e&&c(K),e&&c(q),e&&c(Y),e&&c(I),e&&c(G),e&&c(P),e&&c(X),e&&c(O),ue(E),e&&c(ne),e&&c(U),e&&c(ie),e&&c(A);for(let t=0;t<S.length;t+=1)S[t].d();for(let t=0;t<w.length;t+=1)w[t].d()}}}function tt(i,l,o){let n,{collection:m=new Ye}=l,_=200,d=[];const b=g=>o(1,_=g.code);return i.$$set=g=>{"collection"in g&&o(0,m=g.collection)},i.$$.update=()=>{i.$$.dirty&1&&o(2,d=[{code:200,body:JSON.stringify({token:"JWT_TOKEN",record:Be.dummyCollectionRecord(m),meta:{id:"abc123",name:"John Doe",username:"john.doe",email:"test@example.com",avatarUrl:"https://example.com/avatar.png"}},null,2)},{code:400,body:`
                {
                  "code": 400,
                  "message": "An error occurred while submitting the form.",
                  "data": {
                    "provider": {
                      "code": "validation_required",
                      "message": "Missing required value."
                    }
                  }
                }
            `}])},o(3,n=Be.getApiExampleUrl(Ge.baseUrl)),[m,_,d,n,b]}class ot extends je{constructor(l){super(),He(this,l,tt,et,Je,{collection:0})}}export{ot as default};
